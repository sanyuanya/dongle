package entity

type OperationLog struct {
	SnowflakeId             string `json:"snowflake_id"`
	OperationId             string `json:"operation_id"`
	IncomeExpenseId         string `json:"income_expense_id"`
	UserId                  string `json:"user_id"`
	BeforeUpdatingShipments int64  `json:"before_updating_shipments"`
	AfterUpdatingShipments  int64  `json:"after_updating_shipments"`
	Summary                 string `json:"summary"`
	UserName                string `json:"user_name"`
	Phone                   string `json:"phone"`
	ProductName             string `json:"product_name"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
}

type AddOperationLogRequest struct {
	SnowflakeId             string `json:"snowflake_id"`
	OperationId             string `json:"operation_id"`
	IncomeExpenseId         string `json:"income_expense_id"`
	UserId                  string `json:"user_id"`
	BeforeUpdatingShipments int64  `json:"before_updating_shipments"`
	AfterUpdatingShipments  int64  `json:"after_updating_shipments"`
	Summary                 string `json:"summary"`
}

type GetOperationLogListRequest struct {
	Page     int64  `json:"page,omitempty"`
	PageSize int64  `json:"page_size,omitempty"`
	UserId   string `json:"user_id,omitempty"`
}
