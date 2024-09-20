package data

import (
	"database/sql"

	"github.com/sanyuanya/dongle/entity"
)

func AddCart(tx *sql.Tx, payload *entity.AddCardRequest) error {

	_, err := tx.Exec("INSERT INTO cart (snowflake_id, commodity_id, sku_id, quantity, user_id) VALUES ($1, $2, $3, $4, $5)", payload.SnowflakeId, payload.CommodityId, payload.SkuId, payload.Quantity, payload.UserId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCart(tx *sql.Tx, payload *entity.UpdateCardRequest) error {

	_, err := tx.Exec("UPDATE cart SET quantity = $1 WHERE snowflake_id = $2 AND user_id = $3", payload.Quantity, payload.SnowflakeId, payload.UserId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCart(tx *sql.Tx, snowflakeId, userId string) error {

	_, err := tx.Exec("UPDATE cart SET deleted_at = now() WHERE snowflake_id = $1 AND user_id = $2 AND deleted_at IS NULL", snowflakeId, userId)
	if err != nil {
		return err
	}

	return nil
}

func GetCartList(tx *sql.Tx, payload *entity.GetCartListRequest) ([]entity.Cart, error) {
	var carts []entity.Cart

	baseQuery := `
		SELECT
			cart.snowflake_id,
			cart.commodity_id,
			cart.sku_id,
			cart.quantity,
			cart.user_id,
			TO_CHAR(cart.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			TO_CHAR(cart.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at,
			commodity.name AS commodity_name,
			commodity.code AS commodity_code,
			commodity.description AS commodity_description,
			sku.name AS sku_name,
			sku.code AS sku_code,
			sku.price AS sku_price,
			sku.object_name AS sku_object_name,
			sku.bucket_name AS sku_bucket_name,
			sku.actual_sales AS sku_actual_sales,
			sku.stock_quantity AS sku_stock_quantity,
			sku.price * cart.quantity AS money
		FROM 
			cart cart
		JOIN 
			commodity commodity ON cart.commodity_id = commodity.snowflake_id AND commodity.deleted_at IS NULL
		JOIN
			stock_keeping_unit sku ON cart.sku_id = sku.snowflake_id AND sku.item_id = cart.commodity_id AND sku.deleted_at IS NULL
		WHERE 
			cart.user_id = $1 AND cart.deleted_at IS NULL
		ORDER BY cart.created_at DESC
		LIMIT $2 OFFSET $3
			`
	rows, err := tx.Query(baseQuery, payload.SnowflakeId, payload.PageSize, (payload.Page-1)*payload.PageSize)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cart entity.Cart
		err := rows.Scan(
			&cart.SnowflakeId,
			&cart.CommodityId,
			&cart.SkuId,
			&cart.Quantity,
			&cart.UserId,
			&cart.CreatedAt,
			&cart.UpdatedAt,
			&cart.CommodityName,
			&cart.CommodityCode,
			&cart.CommodityDescription,
			&cart.SkuName,
			&cart.SkuCode,
			&cart.SkuPrice,
			&cart.SkuObjectName,
			&cart.SkuBucketName,
			&cart.SkuActualSales,
			&cart.SkuStockQuantity,
			&cart.Money,
		)
		if err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

func CartListTotal(tx *sql.Tx, payload *entity.GetCartListRequest) (uint64, error) {
	var total uint64

	baseQuery := `
		SELECT
			COUNT(1)
		FROM 
			cart cart
		JOIN 
			commodity commodity ON cart.commodity_id = commodity.snowflake_id AND commodity.deleted_at IS NULL
		JOIN
			stock_keeping_unit sku ON cart.sku_id = sku.snowflake_id AND sku.item_id = cart.commodity_id AND sku.deleted_at IS NULL
		WHERE 
			cart.user_id = $1 AND cart.deleted_at IS NULL
		`
	err := tx.QueryRow(baseQuery, payload.SnowflakeId).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func FindByCartSnowflakeId(tx *sql.Tx, snowflakeId string) (string, error) {
	var id string
	err := tx.QueryRow("SELECT snowflake_id FROM cart WHERE snowflake_id = $1 AND deleted_at IS NULL", snowflakeId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return id, nil
}
