package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/copier"
	"LinuxOnM/internal/utils/encrypt"
	"LinuxOnM/internal/utils/ssh"
	"encoding/base64"
)

type HostService struct{}

type IHostService interface {
	GetHostInfo(id uint) (*models.Host, error)
	TestLocalConn(id uint) bool
	TestByInfo(req dto.HostConnTest) bool
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

func (u *HostService) TestLocalConn(id uint) bool {
	var (
		host models.Host
		err  error
	)
	if id == 0 {
		host, err = hostRepo.Get(hostRepo.WithByAddr("127.0.0.1"))
		if err != nil {
			return false
		}
	} else {
		host, err = hostRepo.Get(commonRepo.WithByID(id))
		if err != nil {
			return false
		}
	}
	var connInfo ssh.ConnInfo
	if err := copier.Copy(&connInfo, &host); err != nil {
		return false
	}
	if len(host.Password) != 0 {
		host.Password, err = encrypt.StringDecrypt(host.Password)
		if err != nil {
			return false
		}
		connInfo.Password = host.Password
	}
	if len(host.PrivateKey) != 0 {
		host.PrivateKey, err = encrypt.StringDecrypt(host.PrivateKey)
		if err != nil {
			return false
		}
		connInfo.PrivateKey = []byte(host.PrivateKey)
	}
	if len(host.PassPhrase) != 0 {
		host.PassPhrase, err = encrypt.StringDecrypt(host.PassPhrase)
		if err != nil {
			return false
		}
		connInfo.PassPhrase = []byte(host.PassPhrase)
	}
	client, err := connInfo.NewClient()
	if err != nil {
		return false
	}
	defer client.Close()

	return true
}

func (u *HostService) TestByInfo(req dto.HostConnTest) bool {
	if req.AuthMode == "password" && len(req.Password) != 0 {
		password, err := base64.StdEncoding.DecodeString(req.Password)
		if err != nil {
			return false
		}
		req.Password = string(password)
	}
	if req.AuthMode == "key" && len(req.PrivateKey) != 0 {
		privateKey, err := base64.StdEncoding.DecodeString(req.PrivateKey)
		if err != nil {
			return false
		}
		req.PrivateKey = string(privateKey)
	}
	if len(req.Password) == 0 && len(req.PrivateKey) == 0 {
		host, err := hostRepo.Get(hostRepo.WithByAddr(req.Addr))
		if err != nil {
			return false
		}
		req.Password = host.Password
		req.AuthMode = host.AuthMode
		req.PrivateKey = host.PrivateKey
		req.PassPhrase = host.PassPhrase
	}

	var connInfo ssh.ConnInfo
	_ = copier.Copy(&connInfo, &req)
	connInfo.PrivateKey = []byte(req.PrivateKey)
	if len(req.PassPhrase) != 0 {
		connInfo.PassPhrase = []byte(req.PassPhrase)
	}
	client, err := connInfo.NewClient()
	if err != nil {
		return false
	}
	defer client.Close()
	return true
}
