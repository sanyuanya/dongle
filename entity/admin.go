package entity

import "time"

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserPageListRequest struct {
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"page_size,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	IsWhite  int64  `json:"is_white,omitempty"`
	City     string `json:"city,omitempty"`
}

type ExportUserRequest struct {
	Keyword string `json:"keyword"`
	IsWhite int64  `json:"is_white"`
	City    string `json:"city"`
}

type UserPageListResponse struct {
	SnowflakeId        string `json:"snowflake_id"`
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
	IsWhite            int    `json:"is_white"`
	WithdrawablePoints int64  `json:"withdrawable_points"`
}

type SetUpWhiteRequest struct {
	WhiteList []string `json:"white_list"`
	Status    int64    `json:"status"`
}

type ApprovalWithdrawalRequest struct {
	ApprovalList []string `json:"approval_list"`
	LifeCycle    int64    `json:"life_cycle"`
	Rejection    string   `json:"rejection"`
}

type AddIncomeExpenseRequest struct {
	SnowflakeId        string    `json:"snowflake_id"`
	Summary            string    `json:"summary"`
	Integral           int64     `json:"integral"`
	Shipments          int64     `json:"shipments"`
	UserId             string    `json:"user_id"`
	Batch              string    `json:"batch"`
	ProductId          string    `json:"product_id"`
	ProductIntegral    int64     `json:"product_integral"`
	ImportdAt          time.Time `json:"importd_at"`
	WithdrawablePoints int64     `json:"withdrawable_points"`
	Path               string    `json:"path"`
	FileName           string    `json:"file_name"`
}

type GetAdminListRequest struct {
	Page              int64  `json:"page,omitempty"`
	PageSize          int64  `json:"page_size,omitempty"`
	Keyword           string `json:"keyword,omitempty"`
	HiddenSnowflakeId string `json:"hidden_snowflake_id,omitempty"`
}

type GetAdminListResponse struct {
	SnowflakeId string                  `json:"snowflake_id"`
	Account     string                  `json:"account"`
	Role        []*GetAdminRoleResponse `json:"role"`
}

type AddAdminRequest struct {
	SnowflakeId string   `json:"snowflake_id"`
	Account     string   `json:"account"`
	Password    string   `json:"password"`
	RoleList    []string `json:"role"`
}

type UpdateAdminRequest struct {
	SnowflakeId string   `json:"snowflake_id"`
	Account     string   `json:"account"`
	RoleList    []string `json:"role"`
}
