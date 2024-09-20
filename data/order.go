package data

import (
	"database/sql"

	"github.com/sanyuanya/dongle/entity"
)

func AddOrder(tx *sql.Tx, payload *entity.SubmitOrderRequest) error {

	_, err := tx.Exec("INSERT INTO `order` (snowflake_id, commodity_id, sku_id, quantity, user_id) VALUES (?, ?, ?, ?, ?)", payload.SnowflakeId, payload.CommodityId, payload.SkuId, payload.Quantity, payload.UserId)
	if err != nil {
		return err
	}

	return nil
}
