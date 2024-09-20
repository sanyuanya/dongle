package entity

type SubmitOrderRequest struct {
	CommodityId string `json:"commodityId"`
	SkuId       string `json:"skuId"`
	Quantity    uint8  `json:"quantity"`
	AddressId   string `json:"addressId"`
}
