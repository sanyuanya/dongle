package entity

type SubmitOrderRequest struct {
	OrderCommodity []*OrderCommodity `json:"orderCommodity"`
	AddressId      string            `json:"addressId"`
}

type AddOrder struct {
	SnowflakeId     string  `json:"snowflakeId"`
	Total           float64 `json:"total"`
	Currency        string  `json:"currency"`
	OutTradeNo      string  `json:"outTradeNo"`
	PrepayId        string  `json:"prepayId"`
	AddressId       string  `json:"addressId"`
	Consignee       string  `json:"consignee"`
	PhoneNumber     string  `json:"phoneNumber"`
	Location        string  `json:"location"`
	DetailedAddress string  `json:"detailedAddress"`
	OrderState      int64   `json:"orderState"`
	NonceStr        string  `json:"nonceStr"`
	PaySign         string  `json:"paySign"`
	PayTimestamp    string  `json:"payTimestamp"`
	SignType        string  `json:"signType"`
	ExpirationTime  int64   `json:"expirationTime"`
	UserId          string  `json:"userId"`
	OpenId          string  `json:"openId"`
}

type GetOrderListRequest struct {
	Page       int64  `json:"page"`
	PageSize   int64  `json:"pageSize"`
	Status     int64  `json:"status"`
	Keyword    string `json:"keyword"`
	OutTradeNo string `json:"outTradeNo"`
	OpenId     string `json:"openId"`
}

type GetOrderListResponse struct {
	SnowflakeId    string  `json:"snowflakeId"`
	OutTradeNo     string  `json:"outTradeNo"`
	Total          float64 `json:"total"`
	PayerTotal     float64 `json:"payerTotal"`
	SuccessTime    int64   `json:"successTime"`
	TradeType      string  `json:"tradeType"`
	TradeState     string  `json:"tradeState"`
	OrderState     int64   `json:"orderState"`
	Nick           string  `json:"nick"`
	Phone          string  `json:"phone"`
	ExpirationTime int64   `json:"expirationTime"`
	CreatedAt      string  `json:"createdAt"`
}
