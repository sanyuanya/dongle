package entity

type AddAddress struct {
	SnowflakeId     string `json:"snowflake_id"`
	Location        string `json:"location"`
	PhoneNumber     string `json:"phone_number"`
	IsDefault       uint8  `json:"is_default"`
	Consignee       string `json:"consignee"`
	Longitude       int64  `json:"longitude"`
	Latitude        int64  `json:"latitude"`
	DetailedAddress string `json:"detailed_address"`
	UserId          string `json:"user_id"`
}

type AddressList struct {
	SnowflakeId     string `json:"snowflake_id"`
	Location        string `json:"location"`
	PhoneNumber     string `json:"phone_number"`
	IsDefault       uint8  `json:"is_default"`
	Consignee       string `json:"consignee"`
	Longitude       int64  `json:"longitude"`
	Latitude        int64  `json:"latitude"`
	DetailedAddress string `json:"detailed_address"`
}

type UpdateAddress struct {
	SnowflakeId     string `json:"snowflake_id"`
	Location        string `json:"location"`
	PhoneNumber     string `json:"phone_number"`
	IsDefault       uint8  `json:"is_default"`
	Consignee       string `json:"consignee"`
	Longitude       int64  `json:"longitude"`
	Latitude        int64  `json:"latitude"`
	DetailedAddress string `json:"detailed_address"`
	UserId          string `json:"user_id"`
}
