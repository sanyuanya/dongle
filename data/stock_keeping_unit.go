package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddSku(tx *sql.Tx, sku *entity.AddSkuRequest) error {
	_, err := tx.Exec("INSERT INTO stock_keeping_unit (snowflake_id, code, name, stock_quantity, virtual_sales, price, status, sorting, item_id, object_name, bucket_name, ext, image_data) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		sku.SnowflakeId,
		sku.Code,
		sku.Name,
		sku.StockQuantity,
		sku.VirtualSales,
		sku.Price,
		sku.Status,
		sku.Sorting,
		sku.ItemId,
		sku.ObjectName,
		sku.BucketName,
		sku.Ext,
		sku.ImageData,
	)
	return err
}

func UpdateSku(tx *sql.Tx, sku *entity.UpdateSkuRequest) error {
	_, err := tx.Exec("UPDATE stock_keeping_unit SET code = $1, name = $2, stock_quantity = $3, virtual_sales = $4, price = $5, status = $6, sorting = $7, item_id = $8, updated_at = $9, object_name = $10, bucket_name = $11, ext = $12, image_data = $13 WHERE snowflake_id = $14",
		sku.Code,
		sku.Name,
		sku.StockQuantity,
		sku.VirtualSales,
		sku.Price,
		sku.Status,
		sku.Sorting,
		sku.ItemId,
		time.Now(),
		sku.ObjectName,
		sku.BucketName,
		sku.Ext,
		sku.ImageData,
		sku.SnowflakeId)
	return err
}

func DeleteSku(tx *sql.Tx, itemId, skuId string) error {
	_, err := tx.Exec("UPDATE stock_keeping_unit SET deleted_at = $1 WHERE item_id = $2 AND snowflake_id = $3 AND deleted_at IS NULL", time.Now(), itemId, skuId)
	return err
}

func GetSkuList(tx *sql.Tx, payload *entity.GetSkuRequest) ([]*entity.Sku, error) {

	baseQuery := `
		SELECT 
			snowflake_id, 
			code, 
			name, 
			stock_quantity, 
			virtual_sales, 
			price, 
			status, 
			sorting, 
			item_id,
			object_name,
			bucket_name,
			ext,
			image_data,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') created_at, 
			TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at
		FROM 
			stock_keeping_unit 
		WHERE deleted_at IS NULL AND item_id = $1`

	paramsIndex := 2
	executeParams := []interface{}{payload.ItemId}

	if payload.Name != "" {
		baseQuery += fmt.Sprintf(" AND name LIKE $%d", paramsIndex)
		executeParams = append(executeParams, "%"+payload.Name+"%")
		paramsIndex++
	}

	if payload.Status != 0 {
		baseQuery += fmt.Sprintf(" AND status = $%d", paramsIndex)
		executeParams = append(executeParams, payload.Status)
		paramsIndex++
	}

	baseQuery += fmt.Sprintf(" ORDER BY sorting DESC, created_at DESC LIMIT $%d OFFSET $%d", paramsIndex, paramsIndex+1)
	executeParams = append(executeParams, payload.PageSize, (payload.Page-1)*payload.PageSize)

	rows, err := tx.Query(baseQuery, executeParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	skuList := make([]*entity.Sku, 0)

	for rows.Next() {
		sku := &entity.Sku{}
		if err := rows.Scan(&sku.SnowflakeId, &sku.Code, &sku.Name, &sku.StockQuantity, &sku.VirtualSales, &sku.Price, &sku.Status, &sku.Sorting, &sku.ItemId, &sku.ObjectName, &sku.BucketName, &sku.Ext, &sku.ImageData, &sku.CreatedAt, &sku.UpdatedAt); err != nil {
			return nil, err
		}
		skuList = append(skuList, sku)
	}

	return skuList, nil

}

func GetSkuCount(tx *sql.Tx, payload *entity.GetSkuRequest) (int64, error) {

	baseQuery := `
		SELECT 
			COUNT(1) 
		FROM 
			stock_keeping_unit 
		WHERE deleted_at IS NULL AND item_id = $1`

	paramsIndex := 2
	executeParams := []interface{}{payload.ItemId}

	if payload.Name != "" {
		baseQuery += fmt.Sprintf(" AND name LIKE $%d", paramsIndex)
		executeParams = append(executeParams, "%"+payload.Name+"%")
		paramsIndex++
	}

	if payload.Status != 0 {
		baseQuery += fmt.Sprintf(" AND status = $%d", paramsIndex)
		executeParams = append(executeParams, payload.Status)
		paramsIndex++
	}

	var total int64

	err := tx.QueryRow(baseQuery, executeParams...).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil

}

func FindBySkuSnowflakeId(tx *sql.Tx, snowflakeId string) (*entity.Sku, error) {
	sku := &entity.Sku{}
	err := tx.QueryRow("SELECT snowflake_id, code, name, stock_quantity, virtual_sales, price, status, sorting, item_id, object_name, bucket_name, ext FROM stock_keeping_unit WHERE snowflake_id = $1 AND deleted_at IS NULL", snowflakeId).Scan(&sku.SnowflakeId, &sku.Code, &sku.Name, &sku.StockQuantity, &sku.VirtualSales, &sku.Price, &sku.Status, &sku.Sorting, &sku.ItemId, &sku.ObjectName, &sku.BucketName, &sku.Ext)
	return sku, err
}

func FindBySkuCode(tx *sql.Tx, code string) (*entity.Sku, error) {
	sku := &entity.Sku{}
	err := tx.QueryRow("SELECT snowflake_id, code, name, stock_quantity, virtual_sales, price, status, sorting, item_id, object_name, bucket_name, ext FROM stock_keeping_unit WHERE code = $1 AND deleted_at IS NULL", code).Scan(&sku.SnowflakeId, &sku.Code, &sku.Name, &sku.StockQuantity, &sku.VirtualSales, &sku.Price, &sku.Status, &sku.Sorting, &sku.ItemId, &sku.ObjectName, &sku.BucketName, &sku.Ext)
	return sku, err
}

func UpdateSkuStatus(tx *sql.Tx, itemId, skuId string, status int64) error {
	_, err := tx.Exec("UPDATE stock_keeping_unit SET status = $1 WHERE item_id = $2 AND snowflake_id = $3 AND deleted_at IS NULL", status, itemId, skuId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("未找到该商品的SKU, 请检查商品ID和SKU ID是否正确")
		}
	}
	return err
}

func UpdateSkuStockQuantity(tx *sql.Tx, itemId, skuId string, quantity int64) error {
	_, err := tx.Exec("UPDATE stock_keeping_unit SET stock_quantity = $1 WHERE item_id = $2 AND snowflake_id = $3 AND deleted_at IS NULL", quantity, itemId, skuId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("未找到该商品的SKU, 请检查商品ID和SKU ID是否正确")
		}
	}
	return err
}
