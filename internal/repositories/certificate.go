package repositories

import (
	"LinuxOnM/internal/models"
	"context"
)

type CertificateRepo struct {
}

func NewICertificateRepo() ICertificateRepo {
	return &CertificateRepo{}
}

type ICertificateRepo interface {
	Page(page, size int, opts ...DBOption) (int64, []models.WebsiteCA, error)
	GetFirst(opts ...DBOption) (models.WebsiteCA, error)
	List(opts ...DBOption) ([]models.WebsiteCA, error)
	Create(ctx context.Context, ca *models.WebsiteCA) error
	DeleteBy(opts ...DBOption) error
}

func (w CertificateRepo) Page(page, size int, opts ...DBOption) (int64, []models.WebsiteCA, error) {
	var caList []models.WebsiteCA
	db := getDb(opts...).Model(&models.WebsiteCA{})
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&caList).Error
	return count, caList, err
}

func (w CertificateRepo) GetFirst(opts ...DBOption) (models.WebsiteCA, error) {
	var ca models.WebsiteCA
	db := getDb(opts...).Model(&models.WebsiteCA{})
	if err := db.First(&ca).Error; err != nil {
		return ca, err
	}
	return ca, nil
}

func (w CertificateRepo) List(opts ...DBOption) ([]models.WebsiteCA, error) {
	var caList []models.WebsiteCA
	db := getDb(opts...).Model(&models.WebsiteCA{})
	err := db.Find(&caList).Error
	return caList, err
}

func (w CertificateRepo) Create(ctx context.Context, ca *models.WebsiteCA) error {
	return getTx(ctx).Create(ca).Error
}

func (w CertificateRepo) DeleteBy(opts ...DBOption) error {
	return getDb(opts...).Delete(&models.WebsiteCA{}).Error
}
