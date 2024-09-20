package entity

type OrderCommodity struct {
	CommodityId string `json:"commodityId"`
	SkuId       string `json:"skuId"`
	Quantity    uint8  `json:"quantity"`
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
	Quantity             uint8   `json:"quantity"`
	ObjectName           string  `json:"objectName"`
	BucketName           string  `json:"bucketName"`
	OrderId              string  `json:"orderId"`
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
}
