package entity

type AddItem struct {
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	CategoriesId string     `json:"categories_id"`
	Description  string     `json:"description"`
	Picture      []*Picture `json:"picture_list"`
	Detail       []*Picture `json:"detail_list"`
	SnowflakeId  string     `json:"snowflake_id"`
	Status       uint8      `json:"status"`
}

type Picture struct {
	SnowflakeId string `json:"snowflake_id"`
	ItemId      string `json:"item_id"`
	Type        uint8  `json:"type"`
	ImageData   string `json:"image_data"`
	Sorting     uint8  `json:"sorting"`
	Ext         string `json:"ext"`
	ObjectName  string `json:"object_name"`
	BucketName  string `json:"bucket_name"`
}

type UpdateItem struct {
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	CategoriesId string     `json:"categories_id"`
	Description  string     `json:"description"`
	Picture      []*Picture `json:"picture_list"`
	Detail       []*Picture `json:"detail_list"`
	SnowflakeId  string     `json:"snowflake_id"`
	Status       uint8      `json:"status"`
}

type Item struct {
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	CategoriesId string     `json:"categories_id"`
	Description  string     `json:"description"`
	Picture      []*Picture `json:"picture_list"`
	Detail       []*Picture `json:"detail_list"`
	SnowflakeId  string     `json:"snowflake_id"`
	Status       uint8      `json:"status"`
	CreatedAt    string     `json:"created_at"`
	Sku          []*Sku     `json:"sku"`
}

type ItemPage struct {
	Page         uint64 `json:"page"`
	PageSize     uint64 `json:"page_size"`
	Name         string `json:"name"`
	CategoriesId string `json:"categories_id"`
	Status       uint64 `json:"status"`
}
