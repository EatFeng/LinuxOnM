package process

import (
	"LinuxOnM/internal/utils/cmd"
	"fmt"
	"os/exec"
	"strings"
)

func ExtractFragmentPath(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "FragmentPath=") {
			return strings.TrimSpace(strings.TrimPrefix(line, "FragmentPath="))
		}
	}
	return ""
}

func GetServiceFilePath(serviceName string) (string, error) {
	cmdStr := fmt.Sprintf("systemctl show --property=FragmentPath %s", serviceName)
	output, err := cmd.Execf(cmdStr)
	if err != nil {
		return "", fmt.Errorf("failed to get service file path for service %s: %v", serviceName, err)
	}

	return ExtractFragmentPath(output), nil
}

func IsDockerContainer(serviceName string) bool {
	cmdStr := fmt.Sprintf("docker ps -a --filter name=^/%s$ --format '{{.ID}}'", serviceName)
	output, err := exec.Command("bash", "-c", cmdStr).Output()
	return err == nil && strings.TrimSpace(string(output)) != ""
}

func GetDockerStatus(containerName string) (string, error) {
	cmdStr := fmt.Sprintf("docker inspect --format '{{.State.Status}}' %s", containerName)
	output, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		if strings.Contains(err.Error(), "No such object") {
			return "", nil
		}
		return "", fmt.Errorf("failed to get Docker container status for %s: %v", containerName, err)
	}

	switch strings.TrimSpace(string(output)) {
	case "running":
		return "running", nil
	case "exited":
		return "stopping", nil
	default:
		return "unknown", nil
	}
}

func GetServiceStatus(serviceName string) (string, error) {
	cmdStr := fmt.Sprintf("systemctl is-active %s", serviceName)
	output, err := exec.Command("bash", "-c", cmdStr).Output()
	status := strings.TrimSpace(string(output))

	switch status {
	case "active":
		return "running", nil
	case "inactive":
		return "stopping", nil
	case "unknown":
		return "unknown", nil
	default:
		if err != nil {
			return "", nil
		}
	}
	return "", nil
}
