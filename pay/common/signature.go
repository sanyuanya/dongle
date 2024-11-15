package common

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"regexp"
)

func Signature(method string, url string, timestamp string, nonceStr string, body string, pri *rsa.PrivateKey) (string, error) {

	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", method, url, timestamp, nonceStr, body)

	// 计算 SHA256 哈希值
	hashed := sha256.Sum256([]byte(signStr))

	signature, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("无法签名: %v", err)
	}

	signatureStr := base64.StdEncoding.EncodeToString(signature)

	mchid := "1682195529"
	serialNo := "1A1EAB972BD01FB2C072DD11996582D1B9F66F5A"
	authorization := fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%s\",serial_no=\"%s\",signature=\"%s\"", mchid, nonceStr, timestamp, serialNo, signatureStr)

	return authorization, nil
}

func PaySign(appId string, timestampString string, nonceString string, packageString string, pri *rsa.PrivateKey) (string, error) {
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n", appId, timestampString, nonceString, packageString)

	// 计算 SHA256 哈希值
	hashed := sha256.Sum256([]byte(signStr))

	signature, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("无法签名: %v", err)
	}

	signatureStr := base64.StdEncoding.EncodeToString(signature)

	return signatureStr, nil
}

func ExtractSignature(authorization string) (string, error) {
	re := regexp.MustCompile(`signature="([^"]+)"`)
	matches := re.FindStringSubmatch(authorization)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法从 authorization 字符串中提取签名")
	}
	return matches[1], nil
}
