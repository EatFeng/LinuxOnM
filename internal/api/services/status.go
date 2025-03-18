package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/docker"
	"encoding/json"
	"fmt"
	"time"
)

type SystemStatusService struct {
	dashboard *DashboardService
}

func NewSystemStatusService() *SystemStatusService {
	return &SystemStatusService{
		dashboard: NewDashboardService().(*DashboardService),
	}
}

func (s *SystemStatusService) GetCurrentStatus() (*dto.StatusResponse, error) {
	response := &dto.StatusResponse{
		Timestamp: time.Now().Unix(),
		Hostname:  global.CONF.System.BindAddress,
	}

	fmt.Println("before if")

	// 从配置表读取需要采集的指标
	if config, err := settingRepo.Get(settingRepo.WithByKey("metrics_config")); err == nil {
		fmt.Println(config.Value)
		var metricsConfig struct {
			CPU    bool `json:"cpu"`
			Mem    bool `json:"mem"`
			Disk   bool `json:"disk"`
			Net    bool `json:"net"`
			Docker bool `json:"docker"`
		}
		_ = json.Unmarshal([]byte(config.Value), &metricsConfig)

		fmt.Println(metricsConfig)

		// 根据配置采集指标
		current := s.dashboard.LoadCurrentInfo("all", "all")
		if metricsConfig.CPU {
			response.CPU = &dto.CPUStatus{
				UsedPercent: current.CPUUsedPercent,
				Cores:       current.CPUTotal,
			}
		}
		if metricsConfig.Mem {
			response.Memory = &dto.MemoryStatus{
				Total:       current.MemoryTotal,
				UsedPercent: current.MemoryUsedPercent,
			}
		}
		if metricsConfig.Docker {
			if dockerStatus, err := s.getDockerStatus(); err == nil {
				response.Docker = dockerStatus
			}
		}
	}

	return response, nil
}

func (s *SystemStatusService) getDockerStatus() (*dto.DockerStatus, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ListAllContainers()
	if err != nil {
		return nil, err
	}

	status := &dto.DockerStatus{
		Total:      len(containers),
		Containers: make([]dto.ContainerShortInfo, 0),
	}

	for _, c := range containers {
		status.Containers = append(status.Containers, dto.ContainerShortInfo{
			ID:     c.ID[:12],
			Name:   c.Names[0],
			Status: c.Status,
		})
	}

	return status, nil
}
