package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddOrderCommodity(tx *sql.Tx, payload *entity.AddOrderCommodity) error {
	baseQuery := `
		INSERT INTO 
			order_commodity(
				snowflake_id,
				commodity_id,
				commodity_name,
				commodity_code,
				categories_id,
				commodity_description,
				sku_id,
				sku_code,
				price,
				object_name,
				bucket_name,
				order_id,
				created_at,
				updated_at,
				quantity,
				sku_name,
				address_id,
				consignee,
				phone_number,
				location,
				detailed_address,
				cart_id
			)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	_, err := tx.Exec(baseQuery,
		payload.SnowflakeId,
		payload.CommodityId,
		payload.CommodityName,
		payload.CommodityCode,
		payload.CategoriesId,
		payload.CommodityDescription,
		payload.SkuId,
		payload.SkuCode,
		payload.Price,
		payload.ObjectName,
		payload.BucketName,
		payload.OrderId,
		time.Now(),
		time.Now(),
		payload.Quantity,
		payload.SkuName,
		payload.AddressId,
		payload.Consignee,
		payload.PhoneNumber,
		payload.Location,
		payload.DetailedAddress,
		payload.CartId,
	)

	if err != nil {
		return fmt.Errorf("创建订单详情失败：%v", err)
	}

	return nil
}

func GetOrderCommodityList(tx *sql.Tx, orderId string) ([]*entity.GetOrderCommodityListResponse, error) {

	baseQuery := `
		SELECT
			oc.snowflake_id,
			oc.commodity_id,
			oc.commodity_name,
			oc.commodity_code,
			oc.categories_id,
			oc.commodity_description,
			oc.sku_id,
			oc.sku_code,
			oc.sku_name,
			oc.price,
			oc.quantity,
			oc.object_name,
			oc.bucket_name,
			oc.order_id,
			TO_CHAR(oc.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			TO_CHAR(oc.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at,
			address_id,
			consignee,
			phone_number,
			location,
			detailed_address
		FROM
			order_commodity oc
		WHERE
			oc.order_id = $1
	`

	rows, err := tx.Query(baseQuery, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*entity.GetOrderCommodityListResponse, 0)
	for rows.Next() {
		var item entity.GetOrderCommodityListResponse
		err = rows.Scan(
			&item.SnowflakeId,
			&item.CommodityId,
			&item.CommodityName,
			&item.CommodityCode,
			&item.CategoriesId,
			&item.CommodityDescription,
			&item.SkuId,
			&item.SkuCode,
			&item.SkuName,
			&item.Price,
			&item.Quantity,
			&item.ObjectName,
			&item.BucketName,
			&item.OrderId,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.AddressId,
			&item.Consignee,
			&item.PhoneNumber,
			&item.Location,
			&item.DetailedAddress,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &item)
	}

	return result, nil
}
