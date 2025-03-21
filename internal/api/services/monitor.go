package services

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type MonitorService struct {
	DiskIO          chan ([]disk.IOCountersStat)
	NetIO           chan ([]net.IOCountersStat)
	AlertStates     map[string]*AlertState
	ThresholdCache  map[string]float64
	IntervalMinutes int
	mu              sync.Mutex
	cacheMu         sync.RWMutex
	lastUpdated     time.Time
}

type AlertState struct {
	Count              int // 连续超过阈值次数
	LastTriggeredCount int // 最近一次触发报警时的计数
}

var monitorCancel context.CancelFunc

type IMonitorService interface {
	Run()
	GetInterval() int
	SetInterval(minutes int)
	saveIODataToDB(ctx context.Context, interval float64)
	saveNetDataToDB(ctx context.Context, interval float64)
}

func NewIMonitorService() IMonitorService {
	return &MonitorService{
		DiskIO:         make(chan []disk.IOCountersStat, 2),
		NetIO:          make(chan []net.IOCountersStat, 2),
		AlertStates:    make(map[string]*AlertState),
		ThresholdCache: make(map[string]float64),
	}
}

func (m *MonitorService) Run() {
	var itemModel models.MonitorBase
	totalPercent, _ := cpu.Percent(3*time.Second, false)
	if len(totalPercent) == 1 {
		itemModel.Cpu = totalPercent[0]
	}
	cpuCount, _ := cpu.Counts(false)

	loadInfo, _ := load.Avg()
	itemModel.CpuLoad1 = loadInfo.Load1
	itemModel.CpuLoad5 = loadInfo.Load5
	itemModel.CpuLoad15 = loadInfo.Load15
	itemModel.LoadUsage = loadInfo.Load1 / (float64(cpuCount*2) * 0.75) * 100

	memoryInfo, _ := mem.VirtualMemory()
	itemModel.Memory = memoryInfo.UsedPercent

	m.checkThresholds(itemModel)

	if err := settingRepo.CreateMonitorBase(itemModel); err != nil {
		global.LOG.Errorf("Insert basic monitoring data failed, err: %v", err)
	}

	m.loadDiskIO()
	m.loadNetIO()

	MonitorStoreDays, err := settingRepo.Get(settingRepo.WithByKey("MonitorStoreDays"))
	if err != nil {
		return
	}
	storeDays, _ := strconv.Atoi(MonitorStoreDays.Value)
	timeForDelete := time.Now().AddDate(0, 0, -storeDays)
	_ = settingRepo.DelMonitorBase(timeForDelete)
	_ = settingRepo.DelMonitorIO(timeForDelete)
	_ = settingRepo.DelMonitorNet(timeForDelete)
}

func (m *MonitorService) loadDiskIO() {
	ioStat, _ := disk.IOCounters()
	var diskIOList []disk.IOCountersStat
	for _, io := range ioStat {
		diskIOList = append(diskIOList, io)
	}
	m.DiskIO <- diskIOList
}

func (m *MonitorService) loadNetIO() {
	netStat, _ := net.IOCounters(true)
	netStatAll, _ := net.IOCounters(false)
	var netList []net.IOCountersStat
	netList = append(netList, netStat...)
	netList = append(netList, netStatAll...)
	m.NetIO <- netList
}

func (m *MonitorService) saveIODataToDB(ctx context.Context, interval float64) {
	defer close(m.DiskIO)
	for {
		select {
		case <-ctx.Done():
			return
		case ioStat := <-m.DiskIO:
			select {
			case <-ctx.Done():
				return
			case ioStat2 := <-m.DiskIO:
				var ioList []models.MonitorIO
				for _, io2 := range ioStat2 {
					for _, io1 := range ioStat {
						if io2.Name == io1.Name {
							var itemIO models.MonitorIO
							itemIO.Name = io1.Name
							if io2.ReadBytes != 0 && io1.ReadBytes != 0 && io2.ReadBytes > io1.ReadBytes {
								itemIO.Read = uint64(float64(io2.ReadBytes-io1.ReadBytes) / interval / 60)
							}
							if io2.WriteBytes != 0 && io1.WriteBytes != 0 && io2.WriteBytes > io1.WriteBytes {
								itemIO.Write = uint64(float64(io2.WriteBytes-io1.WriteBytes) / interval / 60)
							}

							if io2.ReadCount != 0 && io1.ReadCount != 0 && io2.ReadCount > io1.ReadCount {
								itemIO.Count = uint64(float64(io2.ReadCount-io1.ReadCount) / interval / 60)
							}
							writeCount := uint64(0)
							if io2.WriteCount != 0 && io1.WriteCount != 0 && io2.WriteCount > io1.WriteCount {
								writeCount = uint64(float64(io2.WriteCount-io1.WriteCount) / interval * 60)
							}
							if writeCount > itemIO.Count {
								itemIO.Count = writeCount
							}

							if io2.ReadTime != 0 && io1.ReadTime != 0 && io2.ReadTime > io1.ReadTime {
								itemIO.Time = uint64(float64(io2.ReadTime-io1.ReadTime) / interval / 60)
							}
							writeTime := uint64(0)
							if io2.WriteTime != 0 && io1.WriteTime != 0 && io2.WriteTime > io1.WriteTime {
								writeTime = uint64(float64(io2.WriteTime-io1.WriteTime) / interval / 60)
							}
							if writeTime > itemIO.Time {
								itemIO.Time = writeTime
							}
							ioList = append(ioList, itemIO)
							break
						}
					}
				}
				if err := settingRepo.BatchCreateMonitorIO(ioList); err != nil {
					global.LOG.Errorf("Insert io monitoring data failed, err: %v", err)
				}
				m.DiskIO <- ioStat2
			}
		}
	}
}

func (m *MonitorService) saveNetDataToDB(ctx context.Context, interval float64) {
	defer close(m.NetIO)
	for {
		select {
		case <-ctx.Done():
			return
		case netStat := <-m.NetIO:
			select {
			case <-ctx.Done():
				return
			case netStat2 := <-m.NetIO:
				var netList []models.MonitorNetwork
				for _, net2 := range netStat2 {
					for _, net1 := range netStat {
						if net2.Name == net1.Name {
							var itemNet models.MonitorNetwork
							itemNet.Name = net1.Name

							if net2.BytesSent != 0 && net1.BytesSent != 0 && net2.BytesSent > net1.BytesSent {
								itemNet.Up = float64(net2.BytesSent-net1.BytesSent) / 1024 / interval / 60
							}
							if net2.BytesRecv != 0 && net1.BytesRecv != 0 && net2.BytesRecv > net1.BytesRecv {
								itemNet.Down = float64(net2.BytesRecv-net1.BytesRecv) / 1024 / interval / 60
							}
							netList = append(netList, itemNet)
							break
						}
					}
				}

				if err := settingRepo.BatchCreateMonitorNet(netList); err != nil {
					global.LOG.Errorf("Insert network monitoring data failed, err: %v", err)
				}
				m.NetIO <- netStat2
			}
		}
	}
}

func StartMonitor(removeBefore bool, interval string) error {
	if removeBefore {
		monitorCancel()
		global.Cron.Remove(cron.EntryID(global.MonitorCronID))
	}
	intervalItem, err := strconv.Atoi(interval)
	if err != nil {
		return err
	}

	service := NewIMonitorService()
	service.SetInterval(intervalItem)
	ctx, cancel := context.WithCancel(context.Background())
	monitorCancel = cancel
	monitorID, err := global.Cron.AddJob(fmt.Sprintf("@every %sm", interval), service)
	if err != nil {
		return err
	}

	service.Run()

	go service.saveIODataToDB(ctx, float64(intervalItem))
	go service.saveNetDataToDB(ctx, float64(intervalItem))

	global.MonitorCronID = monitorID
	return nil
}

func (m *MonitorService) checkThresholds(data models.MonitorBase) {
	m.refreshThresholds()

	m.checkSingleThreshold("CPU", data.Cpu)

	m.checkSingleThreshold("Memory", data.Memory)
}

func (m *MonitorService) checkSingleThreshold(metric string, currentValue float64) {
	m.cacheMu.RLock()
	threshold, exists := m.ThresholdCache[metric+"Threshold"]
	m.cacheMu.RUnlock()

	if !exists {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.AlertStates[metric]
	if !exists {
		state = &AlertState{Count: 0, LastTriggeredCount: 0}
		m.AlertStates[metric] = state
	}

	if currentValue > threshold {
		state.Count++
		if state.Count-state.LastTriggeredCount >= 3 {
			duration := state.Count * m.GetInterval() * 60
			m.triggerAlert(metric, currentValue, duration)
			state.LastTriggeredCount = state.Count
		}
	} else {
		state.Count = 0
		state.LastTriggeredCount = 0
	}
}

func (m *MonitorService) refreshThresholds() {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	if time.Since(m.lastUpdated) < 30*time.Second {
		return
	}

	keys := []string{"CPUThreshold", "MemoryThreshold"}

	for _, key := range keys {
		setting, err := settingRepo.Get(settingRepo.WithByKey(key))
		if err == nil {
			if value, err := strconv.ParseFloat(setting.Value, 64); err == nil {
				m.ThresholdCache[key] = value
			}
		}
	}
	m.lastUpdated = time.Now()
}

func (m *MonitorService) triggerAlert(metric string, value float64, duration int) {
	msg := fmt.Sprintf("%s使用率为%.2f%%，超过阈值达%d秒", metric, value, duration)
	global.LOG.Info(msg)

	notificationService := NewNotificationService()
	notificationService.SendAlert(metric, value, duration)
}

func (m *MonitorService) GetInterval() int {
	return m.IntervalMinutes
}

func (m *MonitorService) SetInterval(minutes int) {
	m.IntervalMinutes = minutes
}
