package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func Code2Session(jsCode string) (*Code2SessionResp, error) {

	baseURL := "https://api.weixin.qq.com/sns/jscode2session"

	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("appid", appid)
	q.Set("secret", secret)
	q.Set("js_code", jsCode)
	q.Set("grant_type", "authorization_code")

	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	code2SessionResp := Code2SessionResp{}
	if err := json.NewDecoder(resp.Body).Decode(&code2SessionResp); err != nil {
		return nil, err
	}
	return &code2SessionResp, nil
}
