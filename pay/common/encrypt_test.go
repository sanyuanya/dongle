package common

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {

	message := []byte("韩杰")

	certPath := "/Users/sanyuanya/hjworkspace/go_dev/dongle_new/pay/cert"

	publicFilePath := fmt.Sprintf("%s/wechatpay_17BDDF6F46451DE2C953B628B76D4458B00CF054.pem", certPath)
	publicKey, err := ReadPublicKey(publicFilePath)
	if err != nil {
		t.Errorf("无法读取公钥文件: %v", err)
	}
	en, err := Encrypt(message, publicKey)
	if err != nil {
		t.Errorf("Encrypt() error = %v", err)
	}

	fmt.Printf("Encrypt() = %v", en)
}
