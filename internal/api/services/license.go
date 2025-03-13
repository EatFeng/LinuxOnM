package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/encrypt"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type LicenseService struct {
	publicKey *rsa.PublicKey
}

type ILicenseService interface {
	ProcessLicenseUpload(licenseData []byte) (*dto.LicenseUploadResponse, error)
}

func NewILicenseService() ILicenseService {
	return &LicenseService{}
}

func (s *LicenseService) ProcessLicenseUpload(licenseData []byte) (*dto.LicenseUploadResponse, error) {
	publicKey, err := encrypt.LoadRSAPublicKey("/home/feng-yite/GolandProjects/LinuxOnm/LinuxOnM/license/public_key/public_key.pem")
	if err != nil {
		return nil, fmt.Errorf("初始化许可证服务失败: %w", err)
	}

	// 解析许可证结构
	var licensePackage struct {
		Data      dto.SignedLicenseData `json:"data"`
		Signature string                `json:"signature"`
	}

	// fmt.Println("解析许可证结构")
	if err := json.Unmarshal(licenseData, &licensePackage); err != nil {
		return nil, errors.Wrap(err, "解析许可证结构失败")
	}

	// 验证签名
	// fmt.Println("验证签名")
	var signedData struct {
		Data      dto.SignedLicenseData `json:"data"`
		Signature string                `json:"signature"`
	}
	if err := json.Unmarshal(licenseData, &signedData); err != nil {
		return nil, errors.Wrap(err, "解析签名数据失败")
	}

	if err := encrypt.VerifyRSASignature(
		publicKey,
		signedData.Data,
		signedData.Signature,
	); err != nil {
		return nil, errors.Wrap(err, "LICENSE_SIGNATURE_INVALID")
	}

	// 将 DTO 转换为数据库模型
	dbModel, err := s.convertToDBModel(signedData.Data)
	if err != nil {
		return nil, errors.Wrap(err, "数据格式转换失败")
	}

	// 验证有效期（使用数据库模型的 IsExpired 方法）
	if dbModel.IsExpired() {
		return nil, errors.New("LICENSE_EXPIRED: 许可证已过期")
	}

	// fmt.Println("存入数据库")
	if err := licenseRepo.Create(dbModel); err != nil { // 使用转换后的模型
		return nil, errors.Wrap(err, "存入数据库失败")
	}

	// 转换回响应格式
	expiryDateStr := dbModel.ExpiryDate.Format(time.RFC3339) // 使用数据库模型的时间
	return &dto.LicenseUploadResponse{
		LicenseID:    dbModel.LicenseID,
		ExpiryDate:   expiryDateStr,
		HardwareHash: dbModel.HardwareHash,
		IsValid:      true,
	}, nil
}

func (s *LicenseService) convertToDBModel(dtoData dto.SignedLicenseData) (*models.License, error) {
	expiryDate, err := time.Parse("2006-01-02", dtoData.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("解析过期时间失败: %w", err)
	}

	// 返回指针
	return &models.License{
		LicenseID:    dtoData.LicenseID,
		ExpiryDate:   expiryDate,
		HardwareHash: dtoData.HardwareHash,
		IssuedAt:     int64(dtoData.IssuedAt),
	}, nil
}
