package firewall

import (
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/firewall/client"
	"os"
)

type FirewallClient interface {
	Name() string // ufw firewalld
	Start() error
	Stop() error
	Restart() error
	Reload() error
	Status() (string, error) // running not running
	Version() (string, error)

	ListPort() ([]client.FireInfo, error)
	ListForward() ([]client.FireInfo, error)
	ListAddress() ([]client.FireInfo, error)

	Port(port client.FireInfo, operation string) error
	RichRules(rule client.FireInfo, operation string) error
	PortForward(info client.Forward, operation string) error

	EnableForward() error
}

func NewFirewallClient() (FirewallClient, error) {
	if _, err := os.Stat("/usr/sbin/firewalld"); err == nil {
		return client.NewFirewalld()
	}
	if _, err := os.Stat("/usr/sbin/ufw"); err == nil {
		return client.NewUfw()
	}
	return nil, buserr.New(constant.ErrFirewall)
}
