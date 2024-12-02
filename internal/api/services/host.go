package services

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/encrypt"
)

type HostService struct{}

type IHostService interface {
	GetHostInfo(id uint) (*models.Host, error)
}

func NewIHostService() IHostService {
	return &HostService{}
}

func (u *HostService) GetHostInfo(id uint) (*models.Host, error) {
	host, err := hostRepo.Get(commonRepo.WithByID(id))
	if err != nil {
		return nil, constant.ErrRecordNotFound
	}
	if len(host.Password) != 0 {
		host.Password, err = encrypt.StringDecrypt(host.Password)
		if err != nil {
			return nil, err
		}
	}
	if len(host.PrivateKey) != 0 {
		host.PrivateKey, err = encrypt.StringDecrypt(host.PrivateKey)
		if err != nil {
			return nil, err
		}
	}
	if len(host.PassPhrase) != 0 {
		host.PassPhrase, err = encrypt.StringDecrypt(host.PassPhrase)
		if err != nil {
			return nil, err
		}
	}
	return &host, err
}
