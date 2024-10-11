package entity

type AddProductCategoriesRequest struct {
	SnowflakeId string `json:"snowflakeId"`
	Name        string `json:"name"`
	Status      uint8  `json:"status"`
	Sorting     uint64 `json:"sorting"`
}

type UpdateProductCategoriesRequest struct {
	Name        string `json:"name"`
	Status      uint8  `json:"status"`
	Sorting     uint64 `json:"sorting"`
	SnowflakeId string `json:"snowflakeId"`
}

type GetProductCategoriesListRequest struct {
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
	Keyword  string `json:"keyword"`
	Status   int64  `json:"status"`
}

type ProductCategories struct {
	SnowflakeId string  `json:"snowflakeId"`
	Name        string  `json:"name"`
	Status      uint8   `json:"status"`
	Sorting     uint64  `json:"sorting"`
	Item        []*Item `json:"item"`
}
