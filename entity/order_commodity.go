package entity

type OrderCommodity struct {
	CommodityId string `json:"commodityId"`
	SkuId       string `json:"skuId"`
	Quantity    int64  `json:"quantity"`
	AddressId   string `json:"addressId"`
	CartId      string `json:"cartId"`
}

type AddOrderCommodity struct {
	SnowflakeId          string  `json:"snowflakeId"`
	CommodityId          string  `json:"commodityId"`
	CommodityName        string  `json:"commodityName"`
	CommodityCode        string  `json:"commodityCode"`
	CategoriesId         string  `json:"categoriesId"`
	CommodityDescription string  `json:"commodityDescription"`
	SkuId                string  `json:"skuId"`
	SkuCode              string  `json:"skuCode"`
	SkuName              string  `json:"skuName"`
	Price                float64 `json:"price"`
	Quantity             int64   `json:"quantity"`
	ObjectName           string  `json:"objectName"`
	BucketName           string  `json:"bucketName"`
	OrderId              string  `json:"orderId"`
	AddressId            string  `json:"addressId"`
	Consignee            string  `json:"consignee"`
	PhoneNumber          string  `json:"phoneNumber"`
	Location             string  `json:"location"`
	DetailedAddress      string  `json:"detailedAddress"`
	CartId               string  `json:"cartId"`
}

type GetOrderCommodityListResponse struct {
	SnowflakeId          string  `json:"snowflakeId"`
	CommodityId          string  `json:"commodityId"`
	CommodityName        string  `json:"commodityName"`
	CommodityCode        string  `json:"commodityCode"`
	CategoriesId         string  `json:"categoriesId"`
	CommodityDescription string  `json:"commodityDescription"`
	SkuId                string  `json:"skuId"`
	SkuCode              string  `json:"skuCode"`
	SkuName              string  `json:"skuName"`
	Price                float64 `json:"price"`
	Quantity             uint8   `json:"quantity"`
	ObjectName           string  `json:"objectName"`
	BucketName           string  `json:"bucketName"`
	OrderId              string  `json:"orderId"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
	AddressId            string  `json:"addressId"`
	Consignee            string  `json:"consignee"`
	PhoneNumber          string  `json:"phoneNumber"`
	Location             string  `json:"location"`
	DetailedAddress      string  `json:"detailedAddress"`
}

type GetOrderCommodityActualSalesResponse struct {
	SnowflakeId string
	CommodityId string
	SkuId       string
	Quantity    int64
}
