package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func ReadPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取公钥文件: %v", err)
	}

	block, _ := pem.Decode(publicKeyData)

	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("无法解码 PEM 块，或者块类型不是 CERTIFICATE")
	}

	publicKey, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析公钥: %v", err)
	}

	rsaPublicKey, ok := publicKey.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥不是 RSA 公钥")
	}
	return rsaPublicKey, nil
}
