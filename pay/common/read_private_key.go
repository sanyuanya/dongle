package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func ReadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	block, _ := pem.Decode(privateKeyData)

	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("无法解码 PEM 块，或者块类型不是 PRIVATE KEY")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("无法解析私钥: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥不是 RSA 私钥")
	}
	return rsaPrivateKey, nil
}
