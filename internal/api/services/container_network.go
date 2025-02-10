package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/utils/docker"
	"context"
	"github.com/docker/docker/api/types/network"
	"sort"
)

func (u *ContainerService) ListNetwork() ([]dto.Options, error) {
	client, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	list, err := client.NetworkList(context.TODO(), network.ListOptions{})
	if err != nil {
		return nil, err
	}
	var datas []dto.Options
	for _, item := range list {
		datas = append(datas, dto.Options{Option: item.Name})
	}
	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Option < datas[j].Option
	})
	return datas, nil
}
