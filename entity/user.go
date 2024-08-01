package entity

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
	SnowflakeId int64  `json:"snowflake_id"`
}

type UserInfo struct {
	OpenID      string `json:"open_id"`
	Nick        string `json:"nick"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	ApiToken    string `json:"api_token"`
	SessionKey  string `json:"session_key"`
	SnowflakeId int64  `json:"snowflake_id"`
}

type UserDetail struct {
	SnowflakeId   int64  `json:"snowflake_id"`
	OpenID        string `json:"open_id"`
	Nick          string `json:"nick"`
	Avatar        string `json:"avatar"`
	Phone         string `json:"phone"`
	Integral      int    `json:"integral"`
	Shipments     int    `json:"shipments"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	IDCard        string `json:"id_card"`
	CompanyName   string `json:"company_name"`
	Job           string `json:"job"`
	AlipayAccount string `json:"alipay_account"`
	SessionKey    string `json:"session_key"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ImportUserInfo struct {
	Nick        string `json:"nick"`
	Phone       string `json:"phone"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Shipments   int64  `json:"shipments"`
	Integral    int64  `json:"integral"`
	SnowflakeId int64  `json:"snowflake_id"`
}
