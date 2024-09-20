package pay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay/common"
)

type OutTradeNoRequest struct {
	MchId string `json:"mchid"`
}

func OutTradeNo(outTradeNo string) (*entity.DecryptResourceResponse, error) {
	host := "https://api.mch.weixin.qq.com"

	path := fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s?mchid=%s", outTradeNo, "1682195529")

	url := fmt.Sprintf("%s%s", host, path)

	method := http.MethodGet

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	nonceStr, err := common.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("无法生成随机字符串: %v", err)
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
		return nil, fmt.Errorf("无法读取私钥文件: %v", err)
	}

	authorization, err := common.Signature(method, path, timestamp, nonceStr, string(""), privateKey)
	if err != nil {
		return nil, fmt.Errorf("无法生成签名: %v", err)
	}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("无法发送请求: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: %d", resp.StatusCode)
	}

	// 打印响应详细信息

	outTradeNoResponse := &entity.DecryptResourceResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&outTradeNoResponse); err != nil {
		return nil, fmt.Errorf("无法解析响应体: %v", err)
	}

	return outTradeNoResponse, nil
}
