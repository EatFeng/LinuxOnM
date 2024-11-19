package services

import (
	"linux-ops-platform/internal/models"
	"linux-ops-platform/internal/repositories"
)

// GetHomeData 获取首页数据
func GetHomeData() (*models.HomeData, error) {
	// 调用仓库层获取数据
	data, err := repositories.GetHomeData()
	if err != nil {
		return nil, err
	}

	return data, nil
}
