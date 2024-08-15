package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func Decrypt(ciphertext string) (string, error) {

	cipherData, _ := base64.StdEncoding.DecodeString(ciphertext)

	fmt.Println(string(cipherData))

	certPath := "/Users/sanyuanya/hjworkspace/go_dev/dongle_new/pay/cert"

	publicFilePath := fmt.Sprintf("%s/apiclient_key.pem", certPath)
	pri, err := ReadPrivateKey(publicFilePath)
	if err != nil {
		return "", fmt.Errorf("无法读取公钥文件: %v", err)
	}
	hash := sha1.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, pri, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("error from decryption: %s", err)
	}
	return string(plaintext), nil
}
