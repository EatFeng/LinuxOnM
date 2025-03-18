package dto

type StatusResponse struct {
	Timestamp int64         `json:"timestamp"`
	Hostname  string        `json:"hostname"`
	CPU       *CPUStatus    `json:"cpu,omitempty"`
	Memory    *MemoryStatus `json:"memory,omitempty"`
	Docker    *DockerStatus `json:"docker,omitempty"`
}

type CPUStatus struct {
	UsedPercent float64 `json:"used_percent"`
	Cores       int     `json:"cores"`
}

type MemoryStatus struct {
	Total       uint64  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

type DockerStatus struct {
	Total      int                  `json:"total"`
	Containers []ContainerShortInfo `json:"containers"`
}

type ContainerShortInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
