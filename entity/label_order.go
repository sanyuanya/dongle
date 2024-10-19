package entity

type LabelOrderRequest struct {
	PrintType           string  `json:"printType"`
	PartnerId           string  `json:"partnerId"`
	PartnerKey          string  `json:"partnerKey"`
	Kuaidicom           string  `json:"kuaidicom"`
	RecMan              RecMan  `json:"recMan"`
	SendMan             SendMan `json:"sendMan"`
	Cargo               string  `json:"cargo"`
	Count               int64   `json:"count"`
	PayType             string  `json:"payType"`
	ExpType             string  `json:"expType"`
	Remark              string  `json:"remark"`
	TempId              string  `json:"tempId"`
	NeedChild           string  `json:"needChild"`
	NeedBack            string  `json:"needBack"`
	OrderId             string  `json:"orderId"`
	NeedDesensitization bool    `json:"needDesensitization"`
	NeedLogo            bool    `json:"needLogo"`
}

type RecMan struct {
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Tel       string `json:"tel"`
	PrintAddr string `json:"printAddr"`
	Company   string `json:"company"`
}

type SendMan struct {
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Tel       string `json:"tel"`
	PrintAddr string `json:"printAddr"`
	Company   string `json:"company"`
}

type LabelOrderResponse struct {
	Success bool                   `json:"success"`
	Code    int64                  `json:"code"`
	Message string                 `json:"message"`
	Data    LabelOrderResponseData `json:"data"`
}

type LabelOrderResponseData struct {
	TaskId          string `json:"taskId"`
	Kuaidinum       string `json:"kuaidinum"`
	ChildNum        string `json:"childNum"`
	ReturnNum       string `json:"returnNum"`
	Label           string `json:"label"`
	Bulkpen         string `json:"bulkpen"`
	OrgCode         string `json:"orgCode"`
	OrgName         string `json:"orgName"`
	DestCode        string `json:"destCode"`
	DestName        string `json:"destName"`
	OrgSortingCode  string `json:"orgSortingCode"`
	OrgSortingName  string `json:"orgSortingName"`
	DestSortingCode string `json:"destSortingCode"`
	DestSortingName string `json:"destSortingName"`
	OrgExtra        string `json:"orgExtra"`
	DestExtra       string `json:"destExtra"`
	PkgCode         string `json:"pkgCode"`
	PkgName         string `json:"pkgName"`
	Road            string `json:"road"`
	QrCode          string `json:"qrCode"`
	KdComOrderNum   string `json:"kdComOrderNum"`
	ExpressCode     string `json:"expressCode"`
	ExpressName     string `json:"expressName"`
}
