package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	// q.Set("code", code)
	u.RawQuery = q.Encode()

	payload := map[string]string{
		"code": code,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(u.String())
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(payloadBytes))
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