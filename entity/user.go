package entity

import "time"

type RegisterRequest struct {
	Code   string `json:"code"`
	JsCode string `json:"js_code"`
	Nick   string `json:"nick"`
	Avatar string `json:"avatar"`
}

type SetUserInfoRequest struct {
	Nick        string `json:"nick"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	IDCard      string `json:"id_card"`
	Province    string `json:"province"`
	City        string `json:"city"`
	District    string `json:"district"`
	CompanyName string `json:"company_name"`
	Job         string `json:"job"`
	SnowflakeId string `json:"snowflake_id"`
}

type UserInfo struct {
	OpenID      string `json:"open_id"`
	Nick        string `json:"nick"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	ApiToken    string `json:"api_token"`
	SessionKey  string `json:"session_key"`
	SnowflakeId string `json:"snowflake_id"`
}

type UserDetail struct {
	SnowflakeId        string `json:"snowflake_id"`
	OpenID             string `json:"open_id"`
	Nick               string `json:"nick"`
	Avatar             string `json:"avatar"`
	Phone              string `json:"phone"`
	Integral           int    `json:"integral"`
	Shipments          int    `json:"shipments"`
	Province           string `json:"province"`
	City               string `json:"city"`
	District           string `json:"district"`
	IDCard             string `json:"id_card"`
	CompanyName        string `json:"company_name"`
	Job                string `json:"job"`
	AlipayAccount      string `json:"alipay_account"`
	SessionKey         string `json:"session_key"`
	IsWhite            int64  `json:"is_white"`
	WithdrawablePoints int64  `json:"withdrawable_points"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type ImportUserInfo struct {
	Nick               string    `json:"nick"`
	Phone              string    `json:"phone"`
	Province           string    `json:"province"`
	City               string    `json:"city"`
	Shipments          int64     `json:"shipments"`
	Integral           int64     `json:"integral"`
	SnowflakeId        string    `json:"snowflake_id"`
	ImportdAt          time.Time `json:"importd_at"`
	WithdrawablePoints int64     `json:"withdrawable_points"`
}

type MiniLoginRequest struct {
	JsCode string `json:"js_code"`
}

type RegisterUserRequest struct {
	OpenId      string `json:"open_id"`
	SessionKey  string `json:"session_key"`
	SnowflakeId string `json:"snowflake_id"`
}

type UpdateUserInfoRequest struct {
	Nick   string `json:"nick"`
	Avatar string `json:"avatar"`
	Code   string `json:"code"`
	OpenId string `json:"open_id"`
}

type UserInfoReplace struct {
	Nick        string `json:"nick"`
	SnowflakeId string `json:"snowflake_id"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Shipments   int64  `json:"shipments"`
	Integral    int64  `json:"integral"`
	IsWhite     int64  `json:"is_white"`
	OpenId      string `json:"open_id"`
}

type UpdateUserDetailRequest struct {
	Nick               string `json:"nick"`
	Phone              string `json:"phone"`
	Province           string `json:"province"`
	City               string `json:"city"`
	District           string `json:"district"`
	CompanyName        string `json:"company_name"`
	Job                string `json:"job"`
	SnowflakeId        string `json:"snowflake_id"`
	IsWhite            int64  `json:"is_white"`
	Shipments          int64  `json:"shipments"`
	Integral           int64  `json:"integral"`
	WithdrawablePoints int64  `json:"withdrawable_points"`
}
