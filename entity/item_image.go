package entity

type AddItemImage struct {
	SnowflakeId string `json:"snowflake_id"`
	ItemId      string `json:"item_id"`
	Type        uint8  `json:"type"`
	Data        string `json:"data"`
	Sorting     uint8  `json:"sorting"`
	ObjectName  string `json:"object_name"`
	BucketName  string `json:"bucket_name"`
	Ext         string `json:"ext"`
}
