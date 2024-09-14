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

type JsApiResponse struct {
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PrepayId  string `json:"prepay_id"`
	PaySign   string `json:"paySign"`
}

type JsApiRequest struct {
	AppId       string `json:"appid"`
	Mchid       string `json:"mchid"`
	Description string `json:"description"`
	OutTradeNo  string `json:"out_trade_no"`
	Attach      string `json:"attach"`
	Amount      Amount `json:"amount"`
	Payer       Payer  `json:"payer"`
	Detail      Detail `json:"detail"`
	NotifyUrl   string `json:"notify_url"`
}

type Amount struct {
	Total    uint64 `json:"total"`
	Currency string `json:"currency"`
}

type Payer struct {
	OpenId string `json:"openid"`
}

type Detail struct {
	GoodDetail []*GoodDetail `json:"goods_detail"`
}

type GoodDetail struct {
	MerchantGoodsId string `json:"merchant_goods_id"`
	GoodsName       string `json:"goods_name"`
	Quantity        uint64 `json:"quantity"`
	UnitPrice       uint64 `json:"unit_price"`
}

func JsApi(payInfo *JsApiRequest) (*JsApiResponse, error) {
	host := "https://api.mch.weixin.qq.com"

	path := "/v3/pay/transactions/jsapi"

	url := fmt.Sprintf("%s%s", host, path)

	method := http.MethodPost

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

	payloadBytes, err := json.Marshal(payInfo)
	if err != nil {
		return nil, fmt.Errorf("无法序列化请求体: %v", err)
	}

	authorization, err := common.Signature(method, path, timestamp, nonceStr, string(payloadBytes), privateKey)
	if err != nil {
		return nil, fmt.Errorf("无法生成签名: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))

	if err != nil {
		return nil, fmt.Errorf("无法创建请求: %v", err)
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
		return nil, fmt.Errorf("无法发送请求: %v", err)
	}
	defer resp.Body.Close()

	jsApiResponse := JsApiResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&jsApiResponse); err != nil {
		return nil, fmt.Errorf("无法解码响应体: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("无法生成预支付交易会话: %v", resp)
	}

	if jsApiResponse.PrepayId == "" {
		return nil, fmt.Errorf("无法生成预支付交易会话: %v", jsApiResponse)
	}

	if jsApiResponse.PaySign, err = common.ExtractSignature(authorization); err != nil {
		return nil, fmt.Errorf("无法提取签名: %v", err)
	}

	jsApiResponse.TimeStamp = timestamp
	jsApiResponse.NonceStr = nonceStr
	jsApiResponse.SignType = "RSA"
	jsApiResponse.Package = fmt.Sprintf("prepay_id=%s", jsApiResponse.PrepayId)
	return &jsApiResponse, nil
}
