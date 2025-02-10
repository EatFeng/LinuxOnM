package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/utils/docker"
	"context"
	"github.com/docker/docker/api/types/image"
)

type ImageService struct{}

type IImageService interface {
	List() ([]dto.Options, error)
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
