package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type ComposeTemplateRepo struct{}

type IComposeTemplateRepo interface {
	Get(opts ...DBOption) (models.ComposeTemplate, error)
	List(opts ...DBOption) ([]models.ComposeTemplate, error)
	Page(limit, offset int, opts ...DBOption) (int64, []models.ComposeTemplate, error)
	Create(compose *models.ComposeTemplate) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error

	GetRecord(opts ...DBOption) (models.Compose, error)
	CreateRecord(compose *models.Compose) error
	DeleteRecord(opts ...DBOption) error
	ListRecord() ([]models.Compose, error)
}

func NewIComposeTemplateRepo() IComposeTemplateRepo {
	return &ComposeTemplateRepo{}
}

func (u *ComposeTemplateRepo) Get(opts ...DBOption) (models.ComposeTemplate, error) {
	var compose models.ComposeTemplate
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&compose).Error
	return compose, err
}

func (u *ComposeTemplateRepo) Page(page, size int, opts ...DBOption) (int64, []models.ComposeTemplate, error) {
	var users []models.ComposeTemplate
	db := global.DB.Model(&models.ComposeTemplate{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (u *ComposeTemplateRepo) List(opts ...DBOption) ([]models.ComposeTemplate, error) {
	var composes []models.ComposeTemplate
	db := global.DB.Model(&models.ComposeTemplate{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&composes).Error
	return composes, err
}

func (u *ComposeTemplateRepo) Create(compose *models.ComposeTemplate) error {
	return global.DB.Create(compose).Error
}

func (u *ComposeTemplateRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.ComposeTemplate{}).Where("id = ?", id).Updates(vars).Error
}

func (u *ComposeTemplateRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.ComposeTemplate{}).Error
}

func (u *ComposeTemplateRepo) GetRecord(opts ...DBOption) (models.Compose, error) {
	var compose models.Compose
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&compose).Error
	return compose, err
}

func (u *ComposeTemplateRepo) ListRecord() ([]models.Compose, error) {
	var composes []models.Compose
	if err := global.DB.Find(&composes).Error; err != nil {
		return nil, err
	}
	return composes, nil
}

func (u *ComposeTemplateRepo) CreateRecord(compose *models.Compose) error {
	return global.DB.Create(compose).Error
}

func (u *ComposeTemplateRepo) DeleteRecord(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.Compose{}).Error
}
