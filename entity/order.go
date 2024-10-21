package entity

type SubmitOrderRequest struct {
	OrderCommodity []*OrderCommodity `json:"orderCommodity"`
	AddressId      string            `json:"addressId"`
}

type AddOrder struct {
	SnowflakeId     string `json:"snowflakeId"`
	Total           int64  `json:"total"`
	Currency        string `json:"currency"`
	OutTradeNo      string `json:"outTradeNo"`
	PrepayId        string `json:"prepayId"`
	AddressId       string `json:"addressId"`
	Consignee       string `json:"consignee"`
	PhoneNumber     string `json:"phoneNumber"`
	Location        string `json:"location"`
	DetailedAddress string `json:"detailedAddress"`
	OrderState      int64  `json:"orderState"`
	NonceStr        string `json:"nonceStr"`
	PaySign         string `json:"paySign"`
	PayTimestamp    string `json:"payTimestamp"`
	SignType        string `json:"signType"`
	ExpirationTime  int64  `json:"expirationTime"`
	UserId          string `json:"userId"`
	OpenId          string `json:"openId"`
}

type GetOrderListRequest struct {
	Page       int64  `json:"page"`
	PageSize   int64  `json:"pageSize"`
	Status     int64  `json:"status"`
	Keyword    string `json:"keyword"`
	OutTradeNo string `json:"outTradeNo"`
	OpenId     string `json:"openId"`
	TradeState string `json:"trade_state"`
}

type GetOrderListResponse struct {
	SnowflakeId     string  `json:"snowflakeId"`
	TransactionId   string  `json:"transactionId"`
	AppId           string  `json:"appId"`
	MchId           string  `json:"mchId"`
	TradeType       string  `json:"tradeType"`
	TradeState      string  `json:"tradeState"`
	TradeStateDesc  string  `json:"tradeStateDesc"`
	BankType        string  `json:"bankType"`
	SuccessTime     string  `json:"successTime"`
	OpenId          string  `json:"openId"`
	UserId          string  `json:"userId"`
	Total           float64 `json:"total"`
	PayerTotal      float64 `json:"payerTotal"`
	Currency        string  `json:"currency"`
	PayerCurrency   string  `json:"payerCurrency"`
	OutTradeNo      string  `json:"outTradeNo"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
	PrepayId        string  `json:"prepayId"`
	ExpirationTime  int64   `json:"expirationTime"`
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
	Nick            string  `json:"nick"`
	Phone           string  `json:"phone"`
	OrderCommodity  []*GetOrderCommodityListResponse
	Package         string               `json:"package"`
	OrderShipping   *GetShippingResponse `json:"orderShipping"`
}

type UpdateOrderByOutTradeNo struct {
	Status     int64
	OutTradeNo string
}

type GetOrderByTradeStateResponse struct {
	OutTradeNo  string
	SnowflakeId string
}

type UpdateOrderStatusRequest struct {
	OrderId string `json:"orderId"`
	Status  int64  `json:"status"`
}
