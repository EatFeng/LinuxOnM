package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/cmd"
	"LinuxOnM/internal/utils/copier"
	"LinuxOnM/internal/utils/firewall"
	fireClient "LinuxOnM/internal/utils/firewall/client"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const confPath = "/etc/sysctl.conf"

type FirewallService struct{}

type IFirewallService interface {
	LoadBaseInfo() (dto.FirewallBaseInfo, error)
	OperateFirewall(operation string) error
	SearchWithPage(search dto.RuleSearch) (int64, interface{}, error)
	OperatePortRule(req dto.PortRuleOperate, reload bool) error
	OperateForwardRule(req dto.ForwardRuleOperate) error
	OperateAddressRule(req dto.AddrRuleOperate, reload bool) error
	UpdateDescription(req dto.UpdateFirewallDescription) error
	UpdatePortRule(req dto.PortRuleUpdate) error
	UpdateAddrRule(req dto.AddrRuleUpdate) error
	BatchOperateRule(req dto.BatchRuleOperate) error
}

func NewIFirewallService() IFirewallService {
	return &FirewallService{}
}

func (u *FirewallService) LoadBaseInfo() (dto.FirewallBaseInfo, error) {
	var baseInfo dto.FirewallBaseInfo
	baseInfo.Status = "not running"
	baseInfo.Version = "-"
	baseInfo.Name = "-"
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return baseInfo, err
	}
	baseInfo.Name = client.Name()

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		baseInfo.PingStatus = u.pingStatus()
	}()
	go func() {
		defer wg.Done()
		baseInfo.Status, _ = client.Status()
	}()
	go func() {
		defer wg.Done()
		baseInfo.Version, _ = client.Version()
	}()
	wg.Wait()
	return baseInfo, nil
}

func (u *FirewallService) pingStatus() string {
	if _, err := os.Stat("/etc/sysctl.conf"); err != nil {
		return constant.StatusNone
	}
	sudo := cmd.SudoHandleCmd()
	command := fmt.Sprintf("%s cat /etc/sysctl.conf | grep net/ipv4/icmp_echo_ignore_all= ", sudo)
	stdout, _ := cmd.Exec(command)
	if stdout == "net/ipv4/icmp_echo_ignore_all=1\n" {
		return constant.StatusEnable
	}
	return constant.StatusDisable
}

func (u *FirewallService) OperateFirewall(operation string) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}
	switch operation {
	case "start":
		if err := client.Start(); err != nil {
			return err
		}
		if err := u.addPortsBeforeStart(client); err != nil {
			_ = client.Stop()
			return err
		}
		_, _ = cmd.Exec("sudo systemctl restart docker")
		return nil
	case "stop":
		if err := client.Stop(); err != nil {
			return err
		}
		_, _ = cmd.Exec("sudo systemctl restart docker")
		return nil
	case "restart":
		if err := client.Restart(); err != nil {
			return err
		}
		_, _ = cmd.Exec("sudo systemctl restart docker")
		return nil
	case "disablePing":
		return u.updatePingStatus("0")
	case "enablePing":
		return u.updatePingStatus("1")
	}
	return fmt.Errorf("not support such operation: %s", operation)
}

func (u *FirewallService) addPortsBeforeStart(client firewall.FirewallClient) error {
	serverPort, err := settingRepo.Get(settingRepo.WithByKey("ServerPort"))
	if err != nil {
		return err
	}
	if err := client.Port(fireClient.FireInfo{Port: serverPort.Value, Protocol: "tcp", Strategy: "accept"}, "add"); err != nil {
		return err
	}
	if err := client.Port(fireClient.FireInfo{Port: "22", Protocol: "tcp", Strategy: "accept"}, "add"); err != nil {
		return err
	}
	if err := client.Port(fireClient.FireInfo{Port: "80", Protocol: "tcp", Strategy: "accept"}, "add"); err != nil {
		return err
	}
	if err := client.Port(fireClient.FireInfo{Port: "443", Protocol: "tcp", Strategy: "accept"}, "add"); err != nil {
		return err
	}

	return client.Reload()
}

func (u *FirewallService) updatePingStatus(enable string) error {
	lineBytes, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}
	files := strings.Split(string(lineBytes), "\n")
	var newFiles []string
	hasLine := false
	for _, line := range files {
		if strings.Contains(line, "net/ipv4/icmp_echo_ignore_all") || strings.HasPrefix(line, "net/ipv4/icmp_echo_ignore_all") {
			newFiles = append(newFiles, "net/ipv4/icmp_echo_ignore_all="+enable)
			hasLine = true
		} else {
			newFiles = append(newFiles, line)
		}
	}
	if !hasLine {
		newFiles = append(newFiles, "net/ipv4/icmp_echo_ignore_all="+enable)
	}
	file, err := os.OpenFile(confPath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(newFiles, "\n"))
	if err != nil {
		return err
	}

	sudo := cmd.SudoHandleCmd()
	command := fmt.Sprintf("%s sysctl -p", sudo)
	stdout, err := cmd.Exec(command)
	if err != nil {
		return fmt.Errorf("update ping status failed, err: %v", stdout)
	}

	return nil
}

func (u *FirewallService) SearchWithPage(req dto.RuleSearch) (int64, interface{}, error) {
	var (
		datas     []fireClient.FireInfo
		backDatas []fireClient.FireInfo
	)

	client, err := firewall.NewFirewallClient()
	if err != nil {
		return 0, nil, err
	}

	var rules []fireClient.FireInfo
	switch req.Type {
	case "port":
		rules, err = client.ListPort()
	case "forward":
		rules, err = client.ListForward()
	case "address":
		rules, err = client.ListAddress()
	}
	if err != nil {
		return 0, nil, err
	}

	if len(req.Info) != 0 {
		for _, addr := range rules {
			if strings.Contains(addr.Address, req.Info) ||
				strings.Contains(addr.Port, req.Info) ||
				strings.Contains(addr.TargetPort, req.Info) ||
				strings.Contains(addr.TargetIP, req.Info) {
				datas = append(datas, addr)
			}
		}
	} else {
		datas = rules
	}

	var datasFilterStatus []fireClient.FireInfo
	if len(req.Status) != 0 {
		for _, data := range datas {
			if req.Status == "free" && len(data.UsedStatus) == 0 {
				datasFilterStatus = append(datasFilterStatus, data)
			}
			if req.Status == "used" && len(data.UsedStatus) != 0 {
				datasFilterStatus = append(datasFilterStatus, data)
			}
		}
	} else {
		datasFilterStatus = datas
	}

	var datasFilterStrategy []fireClient.FireInfo
	if len(req.Strategy) != 0 {
		for _, data := range datasFilterStatus {
			if req.Strategy == data.Strategy {
				datasFilterStrategy = append(datasFilterStrategy, data)
			}
		}
	} else {
		datasFilterStrategy = datasFilterStatus
	}

	total, start, end := len(datasFilterStrategy), (req.Page-1)*req.PageSize, req.Page*req.PageSize
	if start > total {
		backDatas = make([]fireClient.FireInfo, 0)
	} else {
		if end >= total {
			end = total
		}
		backDatas = datasFilterStrategy[start:end]
	}

	datasFromDB, _ := hostRepo.ListFirewallRecord()
	for i := 0; i < len(backDatas); i++ {
		for _, des := range datasFromDB {
			if req.Type != des.Type {
				continue
			}
			if backDatas[i].Port == des.Port &&
				req.Type == "port" &&
				backDatas[i].Protocol == des.Protocol &&
				backDatas[i].Strategy == des.Strategy &&
				backDatas[i].Address == des.Address {
				backDatas[i].Description = des.Description
				break
			}
			if req.Type == "address" && backDatas[i].Strategy == des.Strategy && backDatas[i].Address == des.Address {
				backDatas[i].Description = des.Description
				break
			}
		}
	}

	go u.cleanUnUsedData(client)

	return int64(total), backDatas, nil
}

func (u *FirewallService) cleanUnUsedData(client firewall.FirewallClient) {
	list, _ := client.ListPort()
	addressList, _ := client.ListAddress()
	list = append(list, addressList...)
	if len(list) == 0 {
		return
	}
	records, _ := hostRepo.ListFirewallRecord()
	if len(records) == 0 {
		return
	}
	for _, item := range list {
		for i := 0; i < len(records); i++ {
			if records[i].Port == item.Port && records[i].Protocol == item.Protocol && records[i].Strategy == item.Strategy && records[i].Address == item.Address {
				records = append(records[:i], records[i+1:]...)
			}
		}
	}

	for _, record := range records {
		_ = hostRepo.DeleteFirewallRecordByID(record.ID)
	}
}

func (u *FirewallService) OperatePortRule(req dto.PortRuleOperate, reload bool) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}
	protos := strings.Split(req.Protocol, "/")
	itemAddress := strings.Split(strings.TrimSuffix(req.Address, ","), ",")

	if client.Name() == "ufw" {
		if strings.Contains(req.Port, ",") || strings.Contains(req.Port, "-") {
			for _, proto := range protos {
				for _, addr := range itemAddress {
					if len(addr) == 0 {
						addr = "Anywhere"
					}
					req.Address = addr
					req.Port = strings.ReplaceAll(req.Port, "-", ":")
					req.Protocol = proto
					if err := u.operatePort(client, req); err != nil {
						return err
					}
					req.Port = strings.ReplaceAll(req.Port, ":", "-")
					if err := u.addPortRecord(req); err != nil {
						return err
					}
				}
			}
			return nil
		}
		for _, addr := range itemAddress {
			if len(addr) == 0 {
				addr = "Anywhere"
			}
			if req.Protocol == "tcp/udp" {
				req.Protocol = ""
			}
			req.Address = addr
			if err := u.operatePort(client, req); err != nil {
				return err
			}
			if len(req.Protocol) == 0 {
				req.Protocol = "tcp/udp"
			}
			if err := u.addPortRecord(req); err != nil {
				return err
			}
		}
		return nil
	}

	itemPorts := req.Port
	for _, proto := range protos {
		if strings.Contains(req.Port, "-") {
			for _, addr := range itemAddress {
				req.Protocol = proto
				req.Address = addr
				if err := u.operatePort(client, req); err != nil {
					return err
				}
				if err := u.addPortRecord(req); err != nil {
					return err
				}
			}
		} else {
			ports := strings.Split(itemPorts, ",")
			for _, port := range ports {
				if len(port) == 0 {
					continue
				}
				for _, addr := range itemAddress {
					req.Address = addr
					req.Port = port
					req.Protocol = proto
					if err := u.operatePort(client, req); err != nil {
						return err
					}
					if err := u.addPortRecord(req); err != nil {
						return err
					}
				}
			}
		}
	}

	if reload {
		return client.Reload()
	}
	return nil
}

func (u *FirewallService) OperateForwardRule(req dto.ForwardRuleOperate) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}

	rules, _ := client.ListForward()
	for _, rule := range rules {
		for _, reqRule := range req.Rules {
			if reqRule.Operation == "remove" {
				continue
			}
			if reqRule.TargetIP == "" {
				reqRule.TargetIP = "127.0.0.1"
			}
			if reqRule.Port == rule.Port && reqRule.TargetPort == rule.TargetPort && reqRule.TargetIP == rule.TargetIP {
				return constant.ErrRecordExist
			}
		}
	}

	sort.SliceStable(req.Rules, func(i, j int) bool {
		n1, _ := strconv.Atoi(req.Rules[i].Num)
		n2, _ := strconv.Atoi(req.Rules[j].Num)
		return n1 > n2
	})

	for _, r := range req.Rules {
		for _, p := range strings.Split(r.Protocol, "/") {
			if r.TargetIP == "" {
				r.TargetIP = "127.0.0.1"
			}
			if err = client.PortForward(fireClient.Forward{
				Num:        r.Num,
				Protocol:   p,
				Port:       r.Port,
				TargetIP:   r.TargetIP,
				TargetPort: r.TargetPort,
			}, r.Operation); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *FirewallService) OperateAddressRule(req dto.AddrRuleOperate, reload bool) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}

	var fireInfo fireClient.FireInfo
	if err := copier.Copy(&fireInfo, &req); err != nil {
		return err
	}

	addressList := strings.Split(req.Address, ",")
	for i := 0; i < len(addressList); i++ {
		if len(addressList[i]) == 0 {
			continue
		}
		fireInfo.Address = addressList[i]
		if err := client.RichRules(fireInfo, req.Operation); err != nil {
			return err
		}
		req.Address = addressList[i]
		if err := u.addAddressRecord(req); err != nil {
			return err
		}
	}
	if reload {
		return client.Reload()
	}
	return nil
}

func (u *FirewallService) operatePort(client firewall.FirewallClient, req dto.PortRuleOperate) error {
	var fireInfo fireClient.FireInfo
	if err := copier.Copy(&fireInfo, &req); err != nil {
		return err
	}

	if client.Name() == "ufw" {
		if len(fireInfo.Address) != 0 && !strings.EqualFold(fireInfo.Address, "Anywhere") {
			return client.RichRules(fireInfo, req.Operation)
		}
		return client.Port(fireInfo, req.Operation)
	}

	if len(fireInfo.Address) != 0 || fireInfo.Strategy == "drop" {
		return client.RichRules(fireInfo, req.Operation)
	}
	return client.Port(fireInfo, req.Operation)
}

func (u *FirewallService) addPortRecord(req dto.PortRuleOperate) error {
	if req.Operation == "remove" {
		return hostRepo.DeleteFirewallRecord("port", req.Port, req.Protocol, req.Address, req.Strategy)
	}

	if err := hostRepo.SaveFirewallRecord(&models.Firewall{
		Type:        "port",
		Port:        req.Port,
		Protocol:    req.Protocol,
		Address:     req.Address,
		Strategy:    req.Strategy,
		Description: req.Description,
	}); err != nil {
		return fmt.Errorf("add record %s/%s failed (strategy: %s, address: %s), err: %v", req.Port, req.Protocol, req.Strategy, req.Address, err)
	}

	return nil
}

func (u *FirewallService) addAddressRecord(req dto.AddrRuleOperate) error {
	if req.Operation == "remove" {
		return hostRepo.DeleteFirewallRecord("address", "", "", req.Address, req.Strategy)
	}
	if err := hostRepo.SaveFirewallRecord(&models.Firewall{
		Type:        "address",
		Address:     req.Address,
		Strategy:    req.Strategy,
		Description: req.Description,
	}); err != nil {
		return fmt.Errorf("add record failed (strategy: %s, address: %s), err: %v", req.Strategy, req.Address, err)
	}
	return nil
}

func (u *FirewallService) UpdatePortRule(req dto.PortRuleUpdate) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}
	if err := u.OperatePortRule(req.OldRule, false); err != nil {
		return err
	}
	if err := u.OperatePortRule(req.NewRule, false); err != nil {
		return err
	}
	return client.Reload()
}

func (u *FirewallService) UpdateAddrRule(req dto.AddrRuleUpdate) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}
	if err := u.OperateAddressRule(req.OldRule, false); err != nil {
		return err
	}
	if err := u.OperateAddressRule(req.NewRule, false); err != nil {
		return err
	}
	return client.Reload()
}

func (u *FirewallService) UpdateDescription(req dto.UpdateFirewallDescription) error {
	var afirewall models.Firewall
	if err := copier.Copy(&afirewall, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	return hostRepo.SaveFirewallRecord(&afirewall)
}

func (u *FirewallService) BatchOperateRule(req dto.BatchRuleOperate) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}
	if req.Type == "port" {
		for _, rule := range req.Rules {
			_ = u.OperatePortRule(rule, false)
		}
		return client.Reload()
	}
	for _, rule := range req.Rules {
		itemRule := dto.AddrRuleOperate{Operation: rule.Operation, Address: rule.Address, Strategy: rule.Strategy}
		_ = u.OperateAddressRule(itemRule, false)
	}
	return client.Reload()
}
