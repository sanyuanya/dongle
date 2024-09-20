package data

import (
	"database/sql"

	"github.com/sanyuanya/dongle/entity"
)

func AddOrderCommodity(tx *sql.Tx, payload *entity.AddOrderCommodity) error {
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
			TO_CHAR(oc.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at
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
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &item)
	}

	return result, nil
}
