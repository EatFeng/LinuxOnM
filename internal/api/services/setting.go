package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/cmd"
	"LinuxOnM/internal/utils/common"
	"LinuxOnM/internal/utils/encrypt"
	"LinuxOnM/internal/utils/files"
	"LinuxOnM/internal/utils/firewall"
	fireClient "LinuxOnM/internal/utils/firewall/client"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type SettingService struct{}

type ISettingService interface {
	GetSettingInfo() (*dto.SettingInfo, error)
	LoadInterfaceAddr() ([]string, error)
	Update(key, value string) error
	UpdatePassword(c *gin.Context, old, new string) error
	UpdateProxy(req dto.ProxyUpdate) error
	UpdateBindInfo(req dto.BindInfo) error
	UpdatePort(port uint) error
	HandlePasswordExpired(c *gin.Context, old, new string) error
	UpdateSSL(c *gin.Context, req dto.SSLUpdate) error
	LoadFromCert() (*dto.SSLInfo, error)
}

func NewISettingService() ISettingService {
	return &SettingService{}
}

func (u *SettingService) GetSettingInfo() (*dto.SettingInfo, error) {
	settings, err := settingRepo.GetList()
	if err != nil {
		return nil, constant.ErrRecordNotFound
	}
	settingMap := make(map[string]string)
	for _, setting := range settings {
		settingMap[setting.Key] = setting.Value
	}
	var info dto.SettingInfo
	arr, err := json.Marshal(settingMap)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(arr, &info); err != nil {
		return nil, err
	}
	if info.ProxyPasswdKeep != constant.StatusEnable {
		info.ProxyPasswd = ""
	} else {
		info.ProxyPasswd, _ = encrypt.StringDecrypt(info.ProxyPasswd)
	}

	info.LocalTime = time.Now().Format("2006-01-02 15:04:05 MST -0700")
	return &info, err
}

func (u *SettingService) LoadInterfaceAddr() ([]string, error) {
	addrMap := make(map[string]struct{})
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && ipNet.IP.To16() != nil {
			addrMap[ipNet.IP.String()] = struct{}{}
		}
	}
	var data []string
	for key := range addrMap {
		data = append(data, key)
	}
	return data, nil
}

func (u *SettingService) Update(key, value string) error {
	switch key {
	case "MonitorStatus":
		if value == "enable" && global.MonitorCronID == 0 {
			interval, err := settingRepo.Get(settingRepo.WithByKey("MonitorInterval"))
			if err != nil {
				return err
			}
			if err := StartMonitor(false, interval.Value); err != nil {
				return err
			}
		}
		if value == "disable" && global.MonitorCronID != 0 {
			monitorCancel()
			global.Cron.Remove(cron.EntryID(global.MonitorCronID))
			global.MonitorCronID = 0
		}
	case "MonitorInterval":
		status, err := settingRepo.Get(settingRepo.WithByKey("MonitorStatus"))
		if err != nil {
			return err
		}
		if status.Value == "enable" && global.MonitorCronID != 0 {
			if err := StartMonitor(true, value); err != nil {
				return err
			}
		}
	}

	if err := settingRepo.Update(key, value); err != nil {
		return err
	}

	switch key {
	case "ExpirationDays":
		timeout, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if err := settingRepo.Update("ExpirationTime", time.Now().AddDate(0, 0, timeout).Format(constant.DateTimeLayout)); err != nil {
			return err
		}
	case "BindDomain":
		if len(value) != 0 {
			_ = global.SESSION.Clean()
		}
	case "UserName", "Password":
		_ = global.SESSION.Clean()

	}

	return nil
}

func (u *SettingService) UpdatePassword(c *gin.Context, old, new string) error {
	if err := u.HandlePasswordExpired(c, old, new); err != nil {
		return err
	}
	_ = global.SESSION.Clean()
	return nil
}

func (u *SettingService) UpdateProxy(req dto.ProxyUpdate) error {
	if err := settingRepo.Update("ProxyUrl", req.ProxyUrl); err != nil {
		return err
	}
	if err := settingRepo.Update("ProxyType", req.ProxyType); err != nil {
		return err
	}
	if err := settingRepo.Update("ProxyPort", req.ProxyPort); err != nil {
		return err
	}
	if err := settingRepo.Update("ProxyUser", req.ProxyUser); err != nil {
		return err
	}
	pass, _ := encrypt.StringEncrypt(req.ProxyPasswd)
	if err := settingRepo.Update("ProxyPasswd", pass); err != nil {
		return err
	}
	if err := settingRepo.Update("ProxyPasswdKeep", req.ProxyPasswdKeep); err != nil {
		return err
	}
	return nil
}

func (u *SettingService) UpdateBindInfo(req dto.BindInfo) error {
	if err := settingRepo.Update("Ipv6", req.Ipv6); err != nil {
		return err
	}
	if err := settingRepo.Update("BindAddress", req.BindAddress); err != nil {
		return err
	}
	go func() {
		time.Sleep(1 * time.Second)
		_, err := cmd.Exec("sudo systemctl restart myapp_LinuxOnM.service")
		if err != nil {
			global.LOG.Errorf("restart system with new bind info failed, err: %v", err)
		}
	}()

	return nil
}

func (u *SettingService) HandlePasswordExpired(c *gin.Context, old, new string) error {
	setting, err := settingRepo.Get(settingRepo.WithByKey("Password"))
	if err != nil {
		return err
	}
	passwordFromDB, err := encrypt.StringDecrypt(setting.Value)
	if err != nil {
		return err
	}
	if passwordFromDB == old {
		newPassword, err := encrypt.StringEncrypt(new)
		if err != nil {
			return err
		}
		if err := settingRepo.Update("Password", newPassword); err != nil {
			return err
		}

		expiredSetting, err := settingRepo.Get(settingRepo.WithByKey("ExpirationDays"))
		if err != nil {
			return err
		}
		timeout, _ := strconv.Atoi(expiredSetting.Value)
		if err := settingRepo.Update("ExpirationTime", time.Now().AddDate(0, 0, timeout).Format(constant.DateTimeLayout)); err != nil {
			return err
		}
		return nil
	}
	return constant.ErrInitialPassword
}

func (u *SettingService) UpdatePort(port uint) error {
	if common.ScanPort(int(port)) {
		return buserr.WithDetail(constant.ErrPortInUsed, port, nil)
	}
	serverPort, err := settingRepo.Get(settingRepo.WithByKey("ServerPort"))
	if err != nil {
		return err
	}
	portValue, _ := strconv.Atoi(serverPort.Value)
	if err := OperateFirewallPort([]int{portValue}, []int{int(port)}); err != nil {
		global.LOG.Errorf("set system firewall ports failed, err: %v", err)
	}
	if err := settingRepo.Update("ServerPort", strconv.Itoa(int(port))); err != nil {
		return err
	}
	go func() {
		time.Sleep(1 * time.Second)
		_, err := cmd.Exec("sudo systemctl restart myapp_LinuxOnM.service")
		if err != nil {
			global.LOG.Errorf("restart system port failed, err: %v", err)
		}
	}()
	return nil
}

func OperateFirewallPort(oldPorts, newPorts []int) error {
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return err
	}

	for _, port := range oldPorts {
		if err := client.Port(fireClient.FireInfo{Port: strconv.Itoa(port), Protocol: "tcp", Strategy: "accept"}, "remove"); err != nil {
			return err
		}
	}

	for _, port := range newPorts {
		if err := client.Port(fireClient.FireInfo{Port: strconv.Itoa(port), Protocol: "tcp", Strategy: "accept"}, "add"); err != nil {
			return err
		}
	}
	return client.Reload()
}

func (u *SettingService) UpdateSSL(c *gin.Context, req dto.SSLUpdate) error {
	secretDir := path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret")
	if req.SSL == "disable" {
		if err := settingRepo.Update("SSL", "disable"); err != nil {
			return err
		}
		if err := settingRepo.Update("SSLType", "self"); err != nil {
			return err
		}
		_ = os.Remove(path.Join(secretDir, "server.crt"))
		_ = os.Remove(path.Join(secretDir, "server.key"))
		sID, _ := c.Cookie(constant.SessionName)
		c.SetCookie(constant.SessionName, sID, 0, "", "", false, true)

		go func() {
			_, err := cmd.Exec("sudo systemctl restart myapp_LinuxOnM.service")
			if err != nil {
				global.LOG.Errorf("restart system failed, err: %v", err)
			}
		}()
		return nil
	}
	if _, err := os.Stat(secretDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(secretDir, os.ModePerm); err != nil {
			return err
		}
	}
	if err := settingRepo.Update("SSLType", req.SSLType); err != nil {
		return err
	}
	var (
		secret string
		key    string
	)

	switch req.SSLType {
	case "self":
		if len(req.Domain) == 0 {
			return fmt.Errorf("load domain failed")
		}
		defaultCA, err := certificateRepo.GetFirst(commonRepo.WithByName("LinuxOnM"))
		if err != nil {
			return err
		}
		websiteSSL, err := NewICertificateService().ObtainSSL(request.CertificateObtain{
			ID:        defaultCA.ID,
			KeyType:   "P256",
			Domains:   req.Domain,
			Time:      1,
			Unit:      "year",
			AutoRenew: true,
		})
		if err != nil {
			return err
		}
		secret = websiteSSL.Pem
		key = websiteSSL.PrivateKey
		if err := settingRepo.Update("SSLID", strconv.Itoa(int(websiteSSL.ID))); err != nil {
			return err
		}
	}

	fileOp := files.NewFileOp()
	if err := fileOp.WriteFile(path.Join(secretDir, "server.crt.tmp"), strings.NewReader(secret), 0600); err != nil {
		return err
	}
	if err := fileOp.WriteFile(path.Join(secretDir, "server.key.tmp"), strings.NewReader(key), 0600); err != nil {
		return err
	}
	if err := checkCertValid(); err != nil {
		return err
	}
	if err := fileOp.Rename(path.Join(secretDir, "server.crt.tmp"), path.Join(secretDir, "server.crt")); err != nil {
		return err
	}
	if err := fileOp.Rename(path.Join(secretDir, "server.key.tmp"), path.Join(secretDir, "server.key")); err != nil {
		return err
	}
	if err := settingRepo.Update("SSL", req.SSL); err != nil {
		return err
	}

	sID, _ := c.Cookie(constant.SessionName)
	c.SetCookie(constant.SessionName, sID, 0, "", "", true, true)
	go func() {
		time.Sleep(1 * time.Second)
		_, err := cmd.Exec("sudo systemctl restart myapp_LinuxOnM.service")
		if err != nil {
			global.LOG.Errorf("restart system failed, err: %v", err)
		}
	}()
	return nil
}

func checkCertValid() error {
	certificate, err := os.ReadFile(path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret/server.crt.tmp"))
	if err != nil {
		return err
	}
	key, err := os.ReadFile(path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret/server.key.tmp"))
	if err != nil {
		return err
	}
	if _, err = tls.X509KeyPair(certificate, key); err != nil {
		return err
	}
	certBlock, _ := pem.Decode(certificate)
	if certBlock == nil {
		return err
	}
	if _, err := x509.ParseCertificate(certBlock.Bytes); err != nil {
		return err
	}

	return nil
}

func (u *SettingService) LoadFromCert() (*dto.SSLInfo, error) {
	ssl, err := settingRepo.Get(settingRepo.WithByKey("SSL"))
	if err != nil {
		return nil, err
	}
	if ssl.Value == "disable" {
		return &dto.SSLInfo{}, nil
	}
	sslType, err := settingRepo.Get(settingRepo.WithByKey("SSLType"))
	if err != nil {
		return nil, err
	}
	var data dto.SSLInfo
	switch sslType.Value {
	case "self":
		data, err = loadInfoFromCert()
		if err != nil {
			return nil, err
		}
	}
	return &data, nil
}

func loadInfoFromCert() (dto.SSLInfo, error) {
	var info dto.SSLInfo
	certFile := path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret/server.crt")
	if _, err := os.Stat(certFile); err != nil {
		return info, err
	}
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return info, err
	}
	certBlock, _ := pem.Decode(certData)
	if certBlock == nil {
		return info, err
	}
	certObj, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return info, err
	}
	var domains []string
	if len(certObj.IPAddresses) != 0 {
		for _, ip := range certObj.IPAddresses {
			domains = append(domains, ip.String())
		}
	}
	if len(certObj.DNSNames) != 0 {
		domains = append(domains, certObj.DNSNames...)
	}
	return dto.SSLInfo{
		Domain:   strings.Join(domains, ","),
		Timeout:  certObj.NotAfter.Format(constant.DateTimeLayout),
		RootPath: path.Join(global.CONF.System.BaseDir, "LinuxOnM/secret/server.crt"),
	}, nil
}
