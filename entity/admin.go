package entity

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserPageListRequest struct {
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"page_size,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	IsWhite  int64  `json:"is_white,omitempty"`
}

type UserPageListResponse struct {
	SnowflakeId   int64  `json:"snowflake_id"`
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
	IsWhite       int    `json:"is_white"`
}

type AddWhiteRequest struct {
	WhiteList []int64 `json:"white_list"`
}

type ApprovalWithdrawalRequest struct {
	ApprovalList []int64 `json:"approval_list"`
	LifeCycle    int64   `json:"life_cycle"`
	Rejection    string  `json:"rejection"`
}

type AddIncomeExpenseRequest struct {
	SnowflakeId int64  `json:"snowflake_id"`
	Summary     string `json:"summary"`
	Integral    int64  `json:"integral"`
	Shipments   int64  `json:"shipments"`
	UserId      int64  `json:"user_id"`
	Batch       int64  `json:"batch"`
}