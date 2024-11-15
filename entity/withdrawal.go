package entity

type ApplyForWithdrawalRequest struct {
	Integral         int64  `json:"integral"`
	AlipayAccount    string `json:"alipay_account"`
	WithdrawalMethod string `json:"withdrawal_method"`
	SnowflakeId      string `json:"snowflake_id"`
	UserId           string `json:"user_id"`
}

type WithdrawalPageListRequest struct {
	Page      int64  `json:"page,omitempty"`
	PageSize  int64  `json:"page_size,omitempty"`
	Keyword   string `json:"keyword,omitempty"`
	LifeCycle int64  `json:"life_cycle,omitempty"`
	Date      string `json:"date,omitempty"`
}

type WithdrawalList struct {
	SnowflakeId      string `json:"snowflake_id"`
	UserId           string `json:"user_id"`
	Nick             string `json:"nick"`
	Phone            string `json:"phone"`
	Integral         int64  `json:"integral"`
	WithdrawalMethod string `json:"withdrawal_method"`
	LifeCycle        int    `json:"life_cycle"`
	Rejection        string `json:"rejection"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	DetailId         string `json:"detail_id"`
	PayId            string `json:"pay_id"`
	InitiateTime     string `json:"initiate_time"`
	UpdateTime       string `json:"update_time"`
	OpenId           string `json:"open_id"`
	MchId            string `json:"mch_id"`
	PaymentStatus    string `json:"payment_status"`
}

type GetWithdrawalListRequest struct {
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"page_size,omitempty"`
	Date     string `json:"date,omitempty"`
}

type GetWithdrawalListResponse struct {
	SnowflakeId      string `json:"snowflake_id"`
	LifeCycle        int    `json:"life_cycle"`
	Integral         int64  `json:"integral"`
	WithdrawalMethod string `json:"withdrawal_method"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	Rejection        string `json:"rejection"`
}

type Withdrawal struct {
	UserId   string `json:"user_id"`
	Integral int64  `json:"integral"`
}

type Order struct {
	PayId       string `json:"pay_id"`
	SnowflakeId string `json:"snowflake_id"`
}
