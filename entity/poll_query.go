package entity

type PollQueryResponse struct {
	Message   string                   `json:"message"`
	State     string                   `json:"state"`
	Status    string                   `json:"status"`
	Condition string                   `json:"condition"`
	IsCheck   string                   `json:"ischeck"`
	Com       string                   `json:"com"`
	Nu        string                   `json:"nu"`
	Data      []*PollQueryResponseData `json:"data"`
}

type PollQueryResponseData struct {
	Context    string `json:"context"`
	Time       string `json:"time"`
	FTime      string `json:"ftime"`
	Status     string `json:"status"`
	StatusCode string `json:"statusCode"`
	AreaCode   string `json:"areaCode"`
	AreaName   string `json:"areaName"`
	AreaCenter string `json:"areaCenter"`
	Location   string `json:"location"`
	AreaPinYin string `json:"areaPinYin"`
}

type PollQueryRequest struct {
	Com      string `json:"com"`
	Num      string `json:"num"`
	Phone    string `json:"phone"`
	From     string `json:"from"`
	To       string `json:"to"`
	ResultV2 string `json:"resultv2"`
	Show     string `json:"show"`
	Order    string `json:"order"`
}
