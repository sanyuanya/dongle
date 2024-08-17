package entity

type AddProductRequest struct {
	SnowflakeId string `json:"snowflake_id"`
	Name        string `json:"name"`
	Integral    uint64 `json:"integral"`
}

type GetProductListRequest struct {
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
	Keyword  string `json:"keyword"`
}

type UpdateProductRequest struct {
	Name     string `json:"name"`
	Integral uint64 `json:"integral"`
}

type GetProductListResponse struct {
	SnowflakeId string `json:"snowflake_id"`
	Name        string `json:"name"`
	Integral    int64  `json:"integral"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
