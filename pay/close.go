package pay

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sanyuanya/dongle/pay/common"
)

type CloseOrderRequest struct {
	MchId string `json:"mchid"`
}

func CloseOrder(outTradeNo string) error {
	host := "https://api.mch.weixin.qq.com"

	path := fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s/close", outTradeNo)

	url := fmt.Sprintf("%s%s", host, path)

	method := http.MethodPost

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	nonceStr, err := common.GenerateRandomString(32)
	if err != nil {
		return fmt.Errorf("无法生成随机字符串: %v", err)
	}
	env := os.Getenv("ENVIRONMENT")

	certPath := ""
	switch env {
	case "production":
		certPath = "/cert"
	default:
		certPath = "/Users/sanyuanya/hjworkspace/go_dev/dongle_new/pay/cert"
	}

	privateFilePath := fmt.Sprintf("%s/apiclient_key.pem", certPath)
	privateKey, err := common.ReadPrivateKey(privateFilePath)
	if err != nil {
		return fmt.Errorf("无法读取私钥文件: %v", err)
	}

	closeOrderRequest := &CloseOrderRequest{
		MchId: "1682195529",
	}
	
	payloadBytes, err := json.Marshal(closeOrderRequest)
	if err != nil {
		return fmt.Errorf("无法序列化请求体: %v", err)
	}

	authorization, err := common.Signature(method, path, timestamp, nonceStr, string(payloadBytes), privateKey)
	if err != nil {
		return fmt.Errorf("无法生成签名: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))

	if err != nil {
		return fmt.Errorf("无法创建请求: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("无法发送请求: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
