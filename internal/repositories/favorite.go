package repositories

import (
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type FavoriteRepo struct{}

type IFavoriteRepo interface {
	GetFirst(opts ...DBOption) (models.Favorite, error)
	WithByPath(path string) DBOption
	Delete(opts ...DBOption) error
}

func NewIFavoriteRepo() IFavoriteRepo {
	return &FavoriteRepo{}
}

func (f *FavoriteRepo) GetFirst(opts ...DBOption) (models.Favorite, error) {
	var favorite models.Favorite
	db := getDb(opts...).Model(&models.Favorite{})
	if err := db.First(&favorite).Error; err != nil {
		return favorite, err
	}
	return favorite, nil
}

func (f *FavoriteRepo) WithByPath(path string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("path = ?", path)
	}
}

func (f *FavoriteRepo) Delete(opts ...DBOption) error {
	db := getDb(opts...).Model(&models.Favorite{})
	return db.Delete(&models.Favorite{}).Error
}
