package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type ISnapshotRepo interface {
	Get(opts ...DBOption) (models.Snapshot, error)
	GetList(opts ...DBOption) ([]models.Snapshot, error)
	Create(snap *models.Snapshot) error
	Update(id uint, vars map[string]interface{}) error
	Page(limit, offset int, opts ...DBOption) (int64, []models.Snapshot, error)
	Delete(opts ...DBOption) error

	GetStatus(snapID uint) (models.SnapshotStatus, error)
	GetStatusList(opts ...DBOption) ([]models.SnapshotStatus, error)
	CreateStatus(snap *models.SnapshotStatus) error
	DeleteStatus(snapID uint) error
	UpdateStatus(id uint, vars map[string]interface{}) error
}

func NewISnapshotRepo() ISnapshotRepo {
	return &SnapshotRepo{}
}

type SnapshotRepo struct{}

func (u *SnapshotRepo) Get(opts ...DBOption) (models.Snapshot, error) {
	var Snapshot models.Snapshot
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&Snapshot).Error
	return Snapshot, err
}

func (u *SnapshotRepo) GetList(opts ...DBOption) ([]models.Snapshot, error) {
	var snaps []models.Snapshot
	db := global.DB.Model(&models.Snapshot{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&snaps).Error
	return snaps, err
}

func (u *SnapshotRepo) Page(page, size int, opts ...DBOption) (int64, []models.Snapshot, error) {
	var users []models.Snapshot
	db := global.DB.Model(&models.Snapshot{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (u *SnapshotRepo) Create(Snapshot *models.Snapshot) error {
	return global.DB.Create(Snapshot).Error
}

func (u *SnapshotRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.Snapshot{}).Where("id = ?", id).Updates(vars).Error
}

func (u *SnapshotRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.Snapshot{}).Error
}

func (u *SnapshotRepo) GetStatus(snapID uint) (models.SnapshotStatus, error) {
	var data models.SnapshotStatus
	if err := global.DB.Where("snap_id = ?", snapID).First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (u *SnapshotRepo) GetStatusList(opts ...DBOption) ([]models.SnapshotStatus, error) {
	var status []models.SnapshotStatus
	db := global.DB.Model(&models.SnapshotStatus{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&status).Error
	return status, err
}

func (u *SnapshotRepo) CreateStatus(snap *models.SnapshotStatus) error {
	return global.DB.Create(snap).Error
}

func (u *SnapshotRepo) DeleteStatus(snapID uint) error {
	return global.DB.Where("snap_id = ?", snapID).Delete(&models.SnapshotStatus{}).Error
}

func (u *SnapshotRepo) UpdateStatus(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.SnapshotStatus{}).Where("id = ?", id).Updates(vars).Error
}
