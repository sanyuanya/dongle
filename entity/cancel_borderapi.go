package entity

type CancelBorderApiRequest struct {
	CancelMsg string `json:"cancelMsg"`
}

type CancelKOrderApiRequest struct {
	TaskId    string `json:"taskId"`
	OrderId   string `json:"orderId"`
	CancelMsg string `json:"cancelMsg"`
}

type CancelKOrderApiResponse struct {
	Result     bool   `json:"result"`
	ReturnCode string `json:"returnCode"`
	Message    string `json:"message"`
}
