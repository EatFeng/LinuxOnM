package dto

import "time"

type ImageRepoOption struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DownloadUrl string `json:"downloadUrl"`
}

type ImageRepoCreate struct {
	Name        string `json:"name" validate:"required"`
	DownloadUrl string `json:"downloadUrl"`
	Protocol    string `json:"protocol"`
	Username    string `json:"username" validate:"max=256"`
	Password    string `json:"password" validate:"max=256"`
	Auth        bool   `json:"auth"`
}

type ImageRepoUpdate struct {
	ID          uint   `json:"id"`
	DownloadUrl string `json:"downloadUrl"`
	Protocol    string `json:"protocol"`
	Username    string `json:"username" validate:"max=256"`
	Password    string `json:"password" validate:"max=256"`
	Auth        bool   `json:"auth"`
}

type ImageRepoDelete struct {
	Ids []uint `json:"ids" validate:"required"`
}

type ImageRepoInfo struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	DownloadUrl string    `json:"downloadUrl"`
	Protocol    string    `json:"protocol"`
	Username    string    `json:"username"`
	Auth        bool      `json:"auth"`

	Status  string `json:"status"`
	Message string `json:"message"`
}
