package encrypt

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// RSA签名验证
func VerifyRSASignature(pubKey *rsa.PublicKey, data interface{}, signature string) error {
	// 序列化数据
	canonical, err := canonicalJSON(data)
	if err != nil {
		return err
	}

	// 调试输出（关键！）
	// log.Printf("规范化的JSON: %s", canonical) // 确保与Python生成的完全一致

	// 计算哈希
	hashed := sha256.Sum256(canonical)
	// log.Printf("SHA256哈希: %x", hashed) // 调试输出

	// 解码签名
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	// log.Printf("签名字节: %x", sigBytes)        // 调试输出
	// log.Printf("签名字节长度: %d", len(sigBytes)) // 必须等于 256

	// 验证签名
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sigBytes)
}

// 加载RSA公钥
func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	// 读取文件
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取公钥文件失败: %w", err)
	}

	// 解码PEM块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("无效的PEM格式")
	}

	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析PKIX公钥失败: %w", err)
	}

	// 类型断言
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("非RSA公钥类型")
	}

	// 检查关键字段
	if rsaPub.N == nil || rsaPub.E == 0 {
		return nil, errors.New("RSA公钥关键字段缺失")
	}

	return rsaPub, nil
}

func canonicalJSON(data interface{}) ([]byte, error) {
	// 将结构体转换为 map 以进行排序
	var m map[string]interface{}
	b, _ := json.Marshal(data)
	json.Unmarshal(b, &m)

	// 按键名排序
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 生成紧凑格式 JSON
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// 序列化键
		keyBytes, _ := json.Marshal(k)
		buf.Write(keyBytes)
		buf.WriteByte(':')
		// 序列化值
		valBytes, _ := json.Marshal(m[k])
		buf.Write(valBytes)
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}

// GenerateHardwareHash 生成硬件哈希值
func GenerateHardwareHash() (string, error) {
	// 1. 获取主板信息
	baseboard, err := getBaseboardInfo()
	if err != nil {
		baseboard = "NULL"
	}
	if baseboard == "" {
		baseboard = "NULL"
	}

	// 2. 获取 CPU 信息
	cpuID, err := getCPUInfo()
	if err != nil {
		cpuID = "NULL"
	}
	if cpuID == "" {
		cpuID = "NULL"
	}

	// 3. 获取磁盘 UUID
	diskUUID, err := getDiskUUID()
	if err != nil {
		diskUUID = "NULL"
	}
	if diskUUID == "" {
		diskUUID = "NULL"
	}

	// 4. 获取 MAC 地址
	macAddress, err := getMACAddress()
	if err != nil {
		macAddress = "NULL"
	}
	if macAddress == "" {
		macAddress = "NULL"
	}

	fmt.Printf("原始信息：%s|%s|%s|%s\n", baseboard, cpuID, diskUUID, macAddress)
	// 组合硬件信息
	hardwareInfo := strings.Join([]string{baseboard, cpuID, diskUUID, macAddress}, "|")

	// 生成 SHA-256 哈希
	hash := sha256.Sum256([]byte(hardwareInfo))
	return hex.EncodeToString(hash[:]), nil
}

// getBaseboardInfo 获取主板信息
func getBaseboardInfo() (string, error) {
	// 尝试通过 dmidecode 获取主板序列号
	cmd := exec.Command("dmidecode", "-t", "baseboard")
	output, err := cmd.Output()
	if err != nil {
		// 如果失败，尝试读取 /sys/class/dmi/id/product_uuid
		uuid, err := exec.Command("cat", "/sys/class/dmi/id/product_uuid").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(uuid)), nil
	}

	// 解析 dmidecode 输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Serial Number") {
			parts := strings.Split(line, ": ")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", nil
}

// getCPUInfo 获取 CPU 信息
func getCPUInfo() (string, error) {
	// 尝试通过 dmidecode 获取 CPU ID
	cmd := exec.Command("dmidecode", "-t", "processor")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "ID") {
				parts := strings.Split(line, ": ")
				if len(parts) > 1 {
					// 移除所有空格和制表符
					cpuID := strings.ReplaceAll(parts[1], " ", "")
					cpuID = strings.ReplaceAll(cpuID, "\t", "")
					return cpuID, nil
				}
			}
		}
	}

	// 如果失败，通过 /proc/cpuinfo 获取 CPU 型号并生成哈希
	cpuModel, err := exec.Command("grep", "-m1", "model name", "/proc/cpuinfo").Output()
	if err != nil {
		return "NULL", nil
	}
	hash := sha256.Sum256(cpuModel)
	return hex.EncodeToString(hash[:]), nil
}

// getDiskUUID 获取磁盘 UUID（严格匹配根分区且返回 NULL）
func getDiskUUID() (string, error) {
	// 1. 尝试通过 lsblk 获取根分区的 UUID
	cmd := exec.Command("lsblk", "-o", "UUID,MOUNTPOINT", "-d", "-n", "-l")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[1] == "/" {
				uuid := strings.TrimSpace(fields[0])
				if uuid != "" {
					return uuid, nil
				}
			}
		}
	}

	// 2. 如果失败，通过 blkid 获取根设备的 UUID
	cmd = exec.Command("sh", "-c", "lsblk -o MOUNTPOINT,PKNAME -n -l | awk '$1 == \"/\" {print $2}'")
	output, err = cmd.Output()
	if err != nil || len(output) == 0 {
		return "NULL", nil // 显式返回 NULL
	}
	device := strings.TrimSpace(string(output))
	if device == "" {
		return "NULL", nil
	}

	cmd = exec.Command("blkid", "-s", "UUID", "-o", "value", "/dev/"+device)
	output, err = cmd.Output()
	if err != nil {
		return "NULL", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// getMACAddress 获取 MAC 地址
func getMACAddress() (string, error) {
	// 通过 ip 命令获取第一个物理网卡的 MAC 地址
	cmd := exec.Command("ip", "link", "show")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 解析 ip 命令输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ether") && !strings.Contains(line, "lo") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return strings.ReplaceAll(parts[1], ":", ""), nil
			}
		}
	}
	return "", nil
}
