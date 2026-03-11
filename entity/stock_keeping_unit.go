package entity

import "encoding/json"

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

func (a AddSkuRequest) MarshalJSON() ([]byte, error) {
	type Alias AddSkuRequest
	return json.Marshal(&struct {
		ImageData string `json:"image_data"`
		*Alias
	}{
		ImageData: "",
		Alias:     (*Alias)(&a),
	})
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

func (u UpdateSkuRequest) MarshalJSON() ([]byte, error) {
	type Alias UpdateSkuRequest
	return json.Marshal(&struct {
		ImageData string `json:"image_data"`
		*Alias
	}{
		ImageData: "",
		Alias:     (*Alias)(&u),
	})
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
	ActualSales   int64   `json:"actual_sales"`
}

func (s Sku) MarshalJSON() ([]byte, error) {
	type Alias Sku
	return json.Marshal(&struct {
		ImageData string `json:"image_data"`
		*Alias
	}{
		ImageData: "",
		Alias:     (*Alias)(&s),
	})
}

type GetSkuRequest struct {
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
	Name     string `json:"name"`
	Status   int64  `json:"status"`
	ItemId   string `json:"item_id"`
}

type UpdateSkuStatus struct {
	Status int64 `json:"status"`
}
