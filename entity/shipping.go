package entity

type Shipping struct {
	OrderList []string `json:"order_list"`
}

type AddShippingRequest struct {
	SnowflakeId  string
	TaskId       string
	OrderId      string
	ThirdOrderId string
	OrderNumber  string
	EOrder       string
}

type GetShippingResponse struct {
	SnowflakeId  string `json:"snowflakeId"`
	TaskId       string `json:"taskId"`
	OrderId      string `json:"orderId"`
	ThirdOrderId string `json:"thirdOrderId"`
	OrderNumber  string `json:"orderNumber"`
	EOrder       string `json:"eOrder"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
