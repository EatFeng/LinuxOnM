package services

import (
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/cmd"
	process2 "LinuxOnM/internal/utils/process"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"runtime"
)

type ProcessService struct{}

type IProcessService interface {
	KillProcess(req request.ProcessReq) error
	GetProcessContent(req request.ProcessRequest) (string, error)
	StartProcess(req request.ProcessRequest) error
	StopProcess(req request.ProcessRequest) error
	EnableProcess(req request.ProcessRequest) error
	DisableProcess(req request.ProcessRequest) error
	StatusProcess(req request.ProcessRequest) (string, error)
	CreateProcess(req request.ProcessCreate) error
}

func NewIProcessService() IProcessService {
	return &ProcessService{}
}

func (p *ProcessService) KillProcess(req request.ProcessReq) error {
	proc, err := process.NewProcess(req.PID)
	if err != nil {
		return err
	}
	if err := proc.Kill(); err != nil {
		return err
	}
	return nil
}

func (p *ProcessService) GetProcessContent(req request.ProcessRequest) (string, error) {
	if runtime.GOOS != "linux" {
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	if !cmd.Which("systemctl") {
		return "", fmt.Errorf("systemctl command not found, make sure the system is Linux and has systemd installed")
	}

	isDocker := process2.IsDockerContainer(req.Name)
	if isDocker {
		return "", buserr.New(constant.ErrProcessIsDocker)
	}

	filePath, err := process2.GetServiceFilePath(req.Name)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read service file %s: %v", filePath, err)
	}
	return string(content), nil
}

func (p *ProcessService) StartProcess(req request.ProcessRequest) error {
	isDocker := process2.IsDockerContainer(req.Name)
	if isDocker {
		return buserr.New(constant.ErrProcessIsDocker)
	}
	cmdStr := fmt.Sprintf("systemctl start %s", req.Name)
	_, err := cmd.Execf(cmdStr)
	if err != nil {
		return buserr.New(constant.ErrStartService)
	}
	return nil
}

func (p *ProcessService) StopProcess(req request.ProcessRequest) error {
	isDocker := process2.IsDockerContainer(req.Name)
	if isDocker {
		return buserr.New(constant.ErrProcessIsDocker)
	}

	cmdStr := fmt.Sprintf("systemctl stop %s", req.Name)
	_, err := cmd.Execf(cmdStr)
	if err != nil {
		return buserr.New(constant.ErrStopService)
	}
	return nil
}

func (p *ProcessService) EnableProcess(req request.ProcessRequest) error {
	isDocker := process2.IsDockerContainer(req.Name)
	if isDocker {
		return buserr.New(constant.ErrProcessIsDocker)
	}
	cmdStr := fmt.Sprintf("systemctl enable %s", req.Name)
	_, err := cmd.Execf(cmdStr)
	if err != nil {
		return buserr.New(constant.ErrEnableService)
	}
	return nil
}

func (p *ProcessService) DisableProcess(req request.ProcessRequest) error {
	isDocker := process2.IsDockerContainer(req.Name)
	if isDocker {
		return buserr.New(constant.ErrProcessIsDocker)
	}
	cmdStr := fmt.Sprintf("systemctl disable %s", req.Name)
	_, err := cmd.Execf(cmdStr)
	if err != nil {
		return buserr.New(constant.ErrDisableService)
	}
	return nil
}

func (p *ProcessService) StatusProcess(req request.ProcessRequest) (string, error) {
	isDocker := process2.IsDockerContainer(req.Name)
	if isDocker {
		dockerStatus, err := process2.GetDockerStatus(req.Name)
		if err != nil {
			return "unknown", nil
		}
		return dockerStatus, nil
	}
	status, err := process2.GetServiceStatus(req.Name)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (p *ProcessService) CreateProcess(req request.ProcessCreate) error {
	serviceFilePath := fmt.Sprintf("/etc/systemd/system/%s.service", req.Name)
	err := os.WriteFile(serviceFilePath, []byte(req.Content), 0644)
	if err != nil {
		return err
	}
	cmdStr := "systemctl daemon-reload"
	_, err = cmd.Execf(cmdStr)
	if err != nil {
		return buserr.New(constant.ErrReloadDaemon)
	}
	if err = p.StartProcess(request.ProcessRequest{Name: req.Name}); err != nil {
		return buserr.New(constant.ErrStartService)
	}
	if err = p.EnableProcess(request.ProcessRequest{Name: req.Name}); err != nil {
		return buserr.New(constant.ErrEnableService)
	}
	return nil
}
