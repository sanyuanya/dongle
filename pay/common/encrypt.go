package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Encrypt(message []byte, pub *rsa.PublicKey) (string, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, message, nil)
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}

	base64Ciphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return base64Ciphertext, nil
}
