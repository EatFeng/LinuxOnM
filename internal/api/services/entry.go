package services

import "LinuxOnM/internal/repositories"

var (
	logRepo    = repositories.NewLogRepository()
	commonRepo = repositories.NewCommonRepository()
)
