package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/copier"
	"LinuxOnM/internal/utils/encrypt"
	"LinuxOnM/internal/utils/ssh"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
)

type HostService struct{}

type IHostService interface {
	GetHostInfo(id uint) (*models.Host, error)
	TestLocalConn(id uint) bool
	TestByInfo(req dto.HostConnTest) bool
	Create(hostDto dto.HostOperate) (*dto.HostInfo, error)
	Update(id uint, upMap map[string]interface{}) error
	Delete(id []uint) error
	SearchForTree(search dto.SearchForTree) ([]dto.HostTree, error)
	SearchWithPage(search dto.SearchHostWithPage) (int64, interface{}, error)

	EncryptHost(itemVal string) (string, error)
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

func (u *HostService) Create(req dto.HostOperate) (*dto.HostInfo, error) {
	var err error
	if len(req.Password) != 0 && req.AuthMode == "password" {
		req.Password, err = u.EncryptHost(req.Password)
		if err != nil {
			return nil, err
		}
		req.PrivateKey = ""
		req.PassPhrase = ""
	}
	if len(req.PrivateKey) != 0 && req.AuthMode == "key" {
		req.PrivateKey, err = u.EncryptHost(req.PrivateKey)
		if err != nil {
			return nil, err
		}
		if len(req.PassPhrase) != 0 {
			req.PassPhrase, err = encrypt.StringEncrypt(req.PassPhrase)
			if err != nil {
				return nil, err
			}
		}
		req.Password = ""
	}
	var host models.Host
	if err := copier.Copy(&host, &req); err != nil {
		return nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	if req.GroupID == 0 {
		group, err := groupRepo.Get(groupRepo.WithByHostDefault())
		if err != nil {
			return nil, errors.New("get default group failed")
		}
		host.GroupID = group.ID
		req.GroupID = group.ID
	}
	var sameHostID uint
	if req.Addr == "127.0.0.1" {
		hostSame, _ := hostRepo.Get(hostRepo.WithByAddr(req.Addr))
		sameHostID = hostSame.ID
	} else {
		hostSame, _ := hostRepo.Get(hostRepo.WithByAddr(req.Addr), hostRepo.WithByUser(req.User), hostRepo.WithByPort(req.Port))
		sameHostID = hostSame.ID
	}
	if sameHostID != 0 {
		host.ID = sameHostID
		upMap := make(map[string]interface{})
		upMap["name"] = req.Name
		upMap["group_id"] = req.GroupID
		upMap["addr"] = req.Addr
		upMap["port"] = req.Port
		upMap["user"] = req.User
		upMap["auth_mode"] = req.AuthMode
		upMap["password"] = req.Password
		upMap["private_key"] = req.PrivateKey
		upMap["pass_phrase"] = req.PassPhrase
		upMap["remember_password"] = req.RememberPassword
		upMap["description"] = req.Description
		if err := hostRepo.Update(sameHostID, upMap); err != nil {
			return nil, err
		}
		var hostinfo dto.HostInfo
		if err := copier.Copy(&hostinfo, &host); err != nil {
			return nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		return &hostinfo, nil
	}

	if err := hostRepo.Create(&host); err != nil {
		return nil, err
	}
	var hostinfo dto.HostInfo
	if err := copier.Copy(&hostinfo, &host); err != nil {
		return nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	return &hostinfo, nil
}

func (u *HostService) EncryptHost(itemVal string) (string, error) {
	privateKey, err := base64.StdEncoding.DecodeString(itemVal)
	if err != nil {
		return "", err
	}
	keyItem, err := encrypt.StringEncrypt(string(privateKey))
	return keyItem, err
}

func (u *HostService) SearchForTree(search dto.SearchForTree) ([]dto.HostTree, error) {
	hosts, err := hostRepo.GetList(hostRepo.WithByInfo(search.Info))
	if err != nil {
		return nil, err
	}
	groups, err := groupRepo.GetList(commonRepo.WithByType("host"))
	if err != nil {
		return nil, err
	}
	var datas []dto.HostTree
	for _, group := range groups {
		var data dto.HostTree
		data.ID = group.ID + 10000
		data.Label = group.Name
		for _, host := range hosts {
			label := fmt.Sprintf("%s@%s:%d", host.User, host.Addr, host.Port)
			if len(host.Name) != 0 {
				label = fmt.Sprintf("%s - %s@%s:%d", host.Name, host.User, host.Addr, host.Port)
			}
			if host.GroupID == group.ID {
				data.Children = append(data.Children, dto.TreeChild{ID: host.ID, Label: label})
			}
		}
		if len(data.Children) != 0 {
			datas = append(datas, data)
		}
	}
	return datas, err
}

func (u *HostService) SearchWithPage(search dto.SearchHostWithPage) (int64, interface{}, error) {
	total, hosts, err := hostRepo.Page(search.Page, search.PageSize, hostRepo.WithByInfo(search.Info), commonRepo.WithByGroupID(search.GroupID))
	if err != nil {
		return 0, nil, err
	}
	var dtoHosts []dto.HostInfo
	for _, host := range hosts {
		var item dto.HostInfo
		if err := copier.Copy(&item, &host); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		group, _ := groupRepo.Get(commonRepo.WithByID(host.GroupID))
		item.GroupBelong = group.Name
		if !item.RememberPassword {
			item.Password = ""
			item.PrivateKey = ""
			item.PassPhrase = ""
		} else {
			if len(host.Password) != 0 {
				item.Password, err = encrypt.StringDecrypt(host.Password)
				if err != nil {
					return 0, nil, err
				}
			}
			if len(host.PrivateKey) != 0 {
				item.PrivateKey, err = encrypt.StringDecrypt(host.PrivateKey)
				if err != nil {
					return 0, nil, err
				}
			}
			if len(host.PassPhrase) != 0 {
				item.PassPhrase, err = encrypt.StringDecrypt(host.PassPhrase)
				if err != nil {
					return 0, nil, err
				}
			}
		}
		dtoHosts = append(dtoHosts, item)
	}
	return total, dtoHosts, err
}

func (u *HostService) Update(id uint, upMap map[string]interface{}) error {
	return hostRepo.Update(id, upMap)
}

func (u *HostService) Delete(ids []uint) error {
	hosts, _ := hostRepo.GetList(commonRepo.WithIDsIn(ids))
	for _, host := range hosts {
		if host.ID == 0 {
			return constant.ErrRecordNotFound
		}
		if host.Addr == "127.0.0.1" {
			return errors.New("the local connection information cannot be deleted!")
		}
	}
	return hostRepo.Delete(commonRepo.WithIDsIn(ids))
}
