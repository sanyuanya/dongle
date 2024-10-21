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
	SnowflakeId     string `json:"snowflakeId"`
	TaskId          string `json:"taskId"`
	OrderId         string `json:"orderId"`
	ThirdOrderId    string `json:"thirdOrderId"`
	OrderNumber     string `json:"orderNumber"`
	EOrder          string `json:"eOrder"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	Status          int8   `json:"status"`
	UserCancelMsg   string `json:"userCancelMsg"`
	SystemCancelMsg string `json:"systemCancelMsg"`
	CourierName     string `json:"courierName"`
	CourierMobile   string `json:"courierMobile"`
	NetTel          string `json:"netTel"`
	NetCode         string `json:"netCode"`
	Weight          string `json:"weight"`
	DefPrice        string `json:"defPrice"`
	Volume          string `json:"volume"`
	ActualWeight    string `json:"actualWeight"`
	PrintTaskId     string `json:"printTaskId"`
	Label           string `json:"label"`
	PickupCode      string `json:"pickupCode"`
}
