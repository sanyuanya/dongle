package wechat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type GetPhoneNumberResp struct {
	Errcode   int       `json:"errcode"`
	Errmsg    string    `json:"errmsg"`
	PhoneInfo PhoneInfo `json:"phone_info"`
}

type PhoneInfo struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

type Watermark struct {
	Timestamp int    `json:"timestamp"`
	Appid     string `json:"appid"`
}

func GetPhoneNumber(code string, access_token string) (*GetPhoneNumberResp, error) {

	baseURL := "https://api.weixin.qq.com/wxa/business/getuserphonenumber"

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("access_token", access_token)
	u.RawQuery = q.Encode()

	payload := map[string]string{
		"code": code,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个自定义的 http.Client，跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Post(u.String(), "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	getPhoneNumberResp := &GetPhoneNumberResp{}
	if err := json.NewDecoder(resp.Body).Decode(&getPhoneNumberResp); err != nil {
		log.Fatal(err)
	}

	return getPhoneNumberResp, nil
}
