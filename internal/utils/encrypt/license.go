package encrypt

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"sort"
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
