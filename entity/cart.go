package entity

type GetCartListRequest struct {
	Page        int64  `json:"page"`
	PageSize    int64  `json:"pageSize"`
	SnowflakeId string `json:"snowflakeId"`
}

type Cart struct {
	SnowflakeId          string  `json:"snowflakeId"`
	CommodityId          string  `json:"commodityId"`
	SkuId                string  `json:"skuId"`
	Quantity             uint64  `json:"quantity"`
	UserId               string  `json:"userId"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
	CommodityName        string  `json:"commodityName"`
	CommodityCode        string  `json:"commodityCode"`
	CommodityDescription string  `json:"commodityDescription"`
	SkuName              string  `json:"skuName"`
	SkuCode              string  `json:"skuCode"`
	SkuPrice             float64 `json:"skuPrice"`
	SkuObjectName        string  `json:"skuObjectName"`
	SkuBucketName        string  `json:"skuBucketName"`
	SkuActualSales       float64 `json:"skuActualSales"`
	SkuStockQuantity     uint64  `json:"skuStockQuantity"`
	Money                float64 `json:"money"`
}

type AddCardRequest struct {
	SnowflakeId string `json:"snowflakeId"`
	CommodityId string `json:"commodityId"`
	SkuId       string `json:"skuId"`
	Quantity    int    `json:"quantity"`
	UserId      string `json:"userId"`
}

type UpdateCardRequest struct {
	SnowflakeId string `json:"snowflakeId"`
	Quantity    int    `json:"quantity"`
	UserId      string `json:"userId"`
}
