package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type HostRepo struct{}

type IHostRepo interface {
	Get(opts ...DBOption) (models.Host, error)
	GetList(opts ...DBOption) ([]models.Host, error)
	Page(limit, offset int, opts ...DBOption) (int64, []models.Host, error)
	Create(host *models.Host) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error

	WithByAddr(addr string) DBOption
	WithByUser(user string) DBOption
	WithByPort(port uint) DBOption
	WithByInfo(info string) DBOption

	ListFirewallRecord() ([]models.Firewall, error)
	DeleteFirewallRecordByID(id uint) error
	DeleteFirewallRecord(fType, port, protocol, address, strategy string) error
	SaveFirewallRecord(firewall *models.Firewall) error
}

func NewIHostRepo() IHostRepo {
	return &HostRepo{}
}

func (h *HostRepo) Get(opts ...DBOption) (models.Host, error) {
	var host models.Host
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&host).Error
	return host, err
}

func (h *HostRepo) GetList(opts ...DBOption) ([]models.Host, error) {
	var hosts []models.Host
	db := global.DB.Model(&models.Host{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&hosts).Error
	return hosts, err
}

func (h *HostRepo) Create(host *models.Host) error {
	return global.DB.Create(host).Error
}

func (h *HostRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.Host{}).Where("id = ?", id).Updates(vars).Error
}

func (h *HostRepo) WithByAddr(addr string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("addr = ?", addr)
	}
}

func (h *HostRepo) WithByUser(user string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("user = ?", user)
	}
}

func (h *HostRepo) WithByPort(port uint) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("port = ?", port)
	}
}

func (h *HostRepo) WithByInfo(info string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(info) == 0 {
			return g
		}
		infoStr := "%" + info + "%"
		return g.Where("name LIKE ? OR addr LIKE ?", infoStr, infoStr)
	}
}

func (h *HostRepo) Page(page, size int, opts ...DBOption) (int64, []models.Host, error) {
	var users []models.Host
	db := global.DB.Model(&models.Host{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (h *HostRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.Host{}).Error
}

func (h *HostRepo) ListFirewallRecord() ([]models.Firewall, error) {
	var datas []models.Firewall
	if err := global.DB.Find(&datas).Error; err != nil {
		return datas, nil
	}
	return datas, nil
}

func (h *HostRepo) DeleteFirewallRecordByID(id uint) error {
	return global.DB.Where("id = ?", id).Delete(&models.Firewall{}).Error
}

func (h *HostRepo) SaveFirewallRecord(firewall *models.Firewall) error {
	if firewall.ID != 0 {
		return global.DB.Save(firewall).Error
	}
	var data models.Firewall
	if firewall.Type == "port" {
		_ = global.DB.Where("type = ? AND port = ? AND protocol = ? AND address = ? AND strategy = ?", "port", firewall.Port, firewall.Protocol, firewall.Address, firewall.Strategy).First(&data)
		if data.ID != 0 {
			firewall.ID = data.ID
		}
	} else {
		_ = global.DB.Where("type = ? AND address = ? AND strategy = ?", "address", firewall.Address, firewall.Strategy).First(&data)
		if data.ID != 0 {
			firewall.ID = data.ID
		}
	}
	return global.DB.Save(firewall).Error
}

func (h *HostRepo) DeleteFirewallRecord(fType, port, protocol, address, strategy string) error {
	return global.DB.Where("type = ? AND port = ? AND protocol = ? AND address = ? AND strategy = ?", fType, port, protocol, address, strategy).Delete(&models.Firewall{}).Error
}
