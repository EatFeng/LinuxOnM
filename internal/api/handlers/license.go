package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/buserr"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HandleLicenseUpload 处理许可证上传
// @Summary 上传许可证文件
// @Tags 许可证管理
// @Accept multipart/form-data
// @Param file formData file true "许可证文件"
// @Param password formData string false "解密密码"
// @Success 200 {object} dto.LicenseUploadResponse
// @Router /api/license/upload [post]
func (b *BaseApi) HandleLicenseUpload(c *gin.Context) {
	// 获取上传文件
	fileHeader, err := c.FormFile("file")
	fmt.Println("获取上传文件")
	if err != nil {
		helper.HandleBusinessError(c, err)
		return
	}

	// 验证文件类型
	fmt.Println("验证文件类型")
	if !strings.HasSuffix(fileHeader.Filename, ".lic") {
		helper.HandleBusinessError(c, buserr.New("LICENSE_INVALID_TYPE"))
		return
	}

	// 读取文件内容
	fmt.Println("读取文件内容")
	file, err := fileHeader.Open()
	if err != nil {
		helper.HandleBusinessError(c, err)
		return
	}
	defer file.Close()

	licenseData, err := io.ReadAll(file)
	if err != nil {
		helper.HandleBusinessError(c, err)
		return
	}

	// 调用服务处理
	fmt.Println("调用服务处理")
	result, err := licenseService.ProcessLicenseUpload(licenseData)
	if err != nil {
		helper.HandleBusinessError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
