package entity

type OrderCallback struct {
	Kuaidicom string             `json:"kuaidicom"`
	Kuaidinum string             `json:"kuaidinum"`
	Status    string             `json:"status"`
	Message   string             `json:"message"`
	Data      *OrderCallbackData `json:"data"`
}

type OrderCallbackData struct {
	OrderId       string       `json:"orderId"`
	Status        int8         `json:"status"`
	CancelMsg9    string       `json:"cancelMsg9"`
	CancelMsg99   string       `json:"cancelMsg99"`
	CourierName   string       `json:"courierName"`
	CourierMobile string       `json:"courierMobile"`
	NetTel        string       `json:"netTel"`
	NetCode       string       `json:"netCode"`
	Weight        string       `json:"weight"`
	DefPrice      string       `json:"defPrice"`
	Freight       string       `json:"freight"`
	Volume        string       `json:"volume"`
	ActualWeight  string       `json:"actualWeight"`
	FeeDetails    []*FeeDetail `json:"feeDetails"`
	PrintTaskId   string       `json:"printTaskId"`
	Label         string       `json:"label"`
	PickupCode    string       `json:"pickupCode"`
}

type FeeDetail struct {
	FeeType   string `json:"feeType"`
	FeeDesc   string `json:"feeDesc"`
	Amount    string `json:"amount"`
	PayStatus string `json:"payStatus"`
}
