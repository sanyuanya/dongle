package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func Encrypt(message []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, message, nil)
	if err != nil {
		return nil, fmt.Errorf("加密失败: %v", err)
	}
	return ciphertext, nil
}
