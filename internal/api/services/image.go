package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/utils/docker"
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

type ImageService struct{}

type IImageService interface {
	List() ([]dto.Options, error)
	ListAll() ([]dto.ImageInfo, error)
}

func NewIImageService() IImageService {
	return &ImageService{}
}

func (u *ImageService) List() ([]dto.Options, error) {
	var (
		list      []image.Summary
		backDatas []dto.Options
	)
	client, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	list, err = client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, image := range list {
		for _, tag := range image.RepoTags {
			backDatas = append(backDatas, dto.Options{
				Option: tag,
			})
		}
	}
	return backDatas, nil
}

func (u *ImageService) ListAll() ([]dto.ImageInfo, error) {
	var records []dto.ImageInfo
	client, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	list, err := client.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		return nil, err
	}
	containers, _ := client.ContainerList(context.Background(), container.ListOptions{All: true})
	for _, image := range list {
		size := formatFileSize(image.Size)
		records = append(records, dto.ImageInfo{
			ID:        image.ID,
			Tags:      image.RepoTags,
			IsUsed:    checkUsed(image.ID, containers),
			CreatedAt: time.Unix(image.Created, 0),
			Size:      size,
		})
	}
	return records, nil
}

func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func checkUsed(imageID string, containers []types.Container) bool {
	for _, container := range containers {
		if container.ImageID == imageID {
			return true
		}
	}
	return false
}
