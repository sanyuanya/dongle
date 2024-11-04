package entity

type KOrderApiRequestParam struct {
	Kuaidicom        string `json:"kuaidicom"`
	RecManName       string `json:"recManName"`
	RecManMobile     string `json:"recManMobile"`
	RecManPrintAddr  string `json:"recManPrintAddr"`
	SendManName      string `json:"sendManName"`
	SendManMobile    string `json:"sendManMobile"`
	SendManPrintAddr string `json:"sendManPrintAddr"`
	CallBackUrl      string `json:"callBackUrl"`
	Cargo            string `json:"cargo"`
	Remark           string `json:"remark"`
}

type KOrderApiRequest struct {
	Method string `json:"method"`
	Param  string `json:"param"`
	Sign   string `json:"sign"`
	T      string `json:"t"`
	Key    string `json:"key"`
}

type KOrderApiResponse struct {
	Result     bool                  `json:"result"`
	ReturnCode string                `json:"returnCode"`
	Message    string                `json:"message"`
	Data       KOrderApiDataResponse `json:"data"`
}

type KOrderApiDataResponse struct {
	TaskId    string `json:"taskId"`
	OrderId   string `json:"orderId"`
	Kuaidinum string `json:"kuaidinum"`
	EOrder    string `json:"eOrder"`
}
