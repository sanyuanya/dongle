package entity

type AddSkuRequest struct {
	SnowflakeId   string  `json:"snowflake_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	StockQuantity int64   `json:"stock_quantity"`
	VirtualSales  int64   `json:"virtual_sales"`
	Price         float64 `json:"price"`
	Status        int64   `json:"status"`
	Sorting       int64   `json:"sorting"`
	ItemId        string  `json:"item_id"`
	ImageData     string  `json:"image_data"`
	ObjectName    string  `json:"object_name"`
	BucketName    string  `json:"bucket_name"`
	Ext           string  `json:"ext"`
}

type UpdateSkuRequest struct {
	SnowflakeId   string  `json:"snowflake_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	StockQuantity int64   `json:"stock_quantity"`
	VirtualSales  int64   `json:"virtual_sales"`
	Price         float64 `json:"price"`
	Status        int64   `json:"status"`
	Sorting       int64   `json:"sorting"`
	ItemId        string  `json:"item_id"`
	ImageData     string  `json:"image_data"`
	ObjectName    string  `json:"object_name"`
	BucketName    string  `json:"bucket_name"`
	Ext           string  `json:"ext"`
}

type ShowSkuResponse struct {
	SnowflakeId   string  `json:"snowflake_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	StockQuantity int64   `json:"stock_quantity"`
	VirtualSales  int64   `json:"virtual_sales"`
	Price         float64 `json:"price"`
	Status        int64   `json:"status"`
	Sorting       int64   `json:"sorting"`
	ItemId        string  `json:"item_id"`
}

type Sku struct {
	SnowflakeId   string  `json:"snowflake_id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	StockQuantity int64   `json:"stock_quantity"`
	VirtualSales  int64   `json:"virtual_sales"`
	Price         float64 `json:"price"`
	Status        int64   `json:"status"`
	Sorting       int64   `json:"sorting"`
	ItemId        string  `json:"item_id"`
	ObjectName    string  `json:"object_name"`
	BucketName    string  `json:"bucket_name"`
	Ext           string  `json:"ext"`
	ImageData     string  `json:"image_data"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type GetSkuRequest struct {
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
	Name     string `json:"name"`
	Status   int64  `json:"status"`
	ItemId   string `json:"item_id"`
}
