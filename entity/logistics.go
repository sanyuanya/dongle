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
	// Payment          string `json:"payment"`
	// ServiceType      string `json:"serviceType"`
	// Weight           string `json:"weight"`
	Remark string `json:"remark"`
	// DayType string `json:"dayType"`
	// PickupStartTime  string `json:"pickupStartTime"`
	// PickupEndTime    string `json:"pickupEndTime"`
	// ChannelSw        string `json:"channelSw"`
	// ValinsPay        string `json:"valinsPay"`
	// RealName         string `json:"realName"`
	// SendIdCardType   string `json:"sendIdCardType"`
	// SendIdCard       string `json:"sendIdCard"`
	// PasswordSigning  string `json:"passwordSigning"`
	// Op               string `json:"op"`
	// PollCallBackUrl  string `json:"pollCallBackUrl"`
	// Resultv2         string `json:"resultv2"`
	// ReturnType       string `json:"returnType"`
	// Siid             string `json:"siid"`
	// Tempid           string `json:"tempid"`
	// PrintCallBackUrl string `json:"print"`
	// Salt             string `json:"salt"`
	// ThirdOrderId string `json:"thirdOrderId"`
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
