package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type ImageRepoRepo struct{}

type IImageRepoRepo interface {
	Get(opts ...DBOption) (models.ImageRepo, error)
	Page(limit, offset int, opts ...DBOption) (int64, []models.ImageRepo, error)
	List(opts ...DBOption) ([]models.ImageRepo, error)
	Create(imageRepo *models.ImageRepo) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error
}

func NewIImageRepoRepo() IImageRepoRepo {
	return &ImageRepoRepo{}
}

func (u *ImageRepoRepo) Get(opts ...DBOption) (models.ImageRepo, error) {
	var imageRepo models.ImageRepo
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&imageRepo).Error
	return imageRepo, err
}

func (u *ImageRepoRepo) Page(page, size int, opts ...DBOption) (int64, []models.ImageRepo, error) {
	var ops []models.ImageRepo
	db := global.DB.Model(&models.ImageRepo{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&ops).Error
	return count, ops, err
}

func (u *ImageRepoRepo) List(opts ...DBOption) ([]models.ImageRepo, error) {
	var ops []models.ImageRepo
	db := global.DB.Model(&models.ImageRepo{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Find(&ops).Error
	return ops, err
}

func (u *ImageRepoRepo) Create(imageRepo *models.ImageRepo) error {
	return global.DB.Create(imageRepo).Error
}

func (u *ImageRepoRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.ImageRepo{}).Where("id = ?", id).Updates(vars).Error
}

func (u *ImageRepoRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.ImageRepo{}).Error
}
