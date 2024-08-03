package wechat

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type getAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func GetAccessToken() (*getAccessTokenResp, error) {
	baseURL := "https://api.weixin.qq.com/cgi-bin/token"

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("appid", appid)
	q.Set("secret", secret)
	q.Set("grant_type", "client_credential")
	u.RawQuery = q.Encode()

	// 创建一个自定义的 http.Client，跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	getAccessTokenResp := &getAccessTokenResp{}
	if err := json.NewDecoder(resp.Body).Decode(&getAccessTokenResp); err != nil {
		log.Fatal(err)
	}

	return getAccessTokenResp, nil
}
