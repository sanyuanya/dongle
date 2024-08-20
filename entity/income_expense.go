package entity

type GetIncomeListRequest struct {
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"page_size,omitempty"`
	Date     string `json:"date,omitempty"`
}

type GetIncomeListResponse struct {
	SnowflakeId     string `json:"snowflake_id"`
	Summary         string `json:"summary"`
	Integral        int64  `json:"integral"`
	Shipments       int64  `json:"shipments"`
	Batch           string `json:"batch"`
	UserId          string `json:"user_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	ProductName     string `json:"product_name"`
	ProductIntegral int64  `json:"product_integral"`
}

type IncomePageListExpenseRequest struct {
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"page_size,omitempty"`
	Date     string `json:"date,omitempty"`
	Keyword  string `json:"keyword,omitempty"`
	UserId   string `json:"user_id,omitempty"`
}

type IncomePageListExpenseResponse struct {
	SnowflakeId     string `json:"snowflake_id"`
	UserId          string `json:"user_id"`
	Summary         string `json:"summary"`
	Integral        int64  `json:"integral"`
	Shipments       int64  `json:"shipments"`
	Batch           string `json:"batch"`
	Nick            string `json:"nick"`
	Phone           string `json:"phone"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	ProductName     string `json:"product_name"`
	ProductIntegral int64  `json:"product_integral"`
}

type GetProductGroupListResponse struct {
	ProductName string `json:"product_name"`
	Shipments   int64  `json:"shipments"`
	Integral    int64  `json:"integral"`
	Merge       int64  `json:"merge"`
}

type UpdateIncomeRequest struct {
	SnowflakeId string `json:"snowflake_id"`
	Shipments   int64  `json:"shipments"`
	Integral    int64  `json:"integral"`
}
