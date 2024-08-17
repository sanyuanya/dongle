package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func GetProductAll(tx *sql.Tx) ([]*entity.GetProductListResponse, error) {
	rows, err := tx.Query(`
		SELECT
			snowflake_id,
			name,
			integral,
			created_at,
			updated_at
		FROM	
			product
		WHERE
			deleted_at IS NULL
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var productList []*entity.GetProductListResponse

	for rows.Next() {
		product := &entity.GetProductListResponse{}
		err := rows.Scan(&product.SnowflakeId,
			&product.Name,
			&product.Integral,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	if productList == nil {
		productList = []*entity.GetProductListResponse{}
	}

	return productList, nil
}

func GetProductList(tx *sql.Tx, page *entity.GetProductListRequest) ([]*entity.GetProductListResponse, error) {

	baseSQL := `
		SELECT
			snowflake_id,
			name,
			integral,
			created_at,
			updated_at
		FROM	
			product
		WHERE
			deleted_at IS NULL
			`

	paramIndex := 1
	executeParams := []interface{}{}

	if page.Keyword != "" {
		baseSQL += fmt.Sprintf(" AND name LIKE $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	baseSQL += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	executeParams = append(executeParams, page.PageSize, (page.Page-1)*page.PageSize)

	rows, err := tx.Query(baseSQL, executeParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var productList []*entity.GetProductListResponse

	for rows.Next() {
		product := &entity.GetProductListResponse{}
		err := rows.Scan(&product.SnowflakeId,
			&product.Name,
			&product.Integral,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	if productList == nil {
		productList = []*entity.GetProductListResponse{}
	}

	return productList, nil

}

func GetProductTotal(tx *sql.Tx, page *entity.GetProductListRequest) (int64, error) {
	baseSQL := `
		SELECT
			COUNT(*)
		FROM	
			product
		WHERE
			deleted_at IS NULL
			`

	paramIndex := 1
	executeParams := []interface{}{}

	if page.Keyword != "" {
		baseSQL += fmt.Sprintf(" AND name LIKE $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	var count int64
	err := tx.QueryRow(baseSQL, executeParams...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func AddProduct(tx *sql.Tx, product *entity.AddProductRequest) error {
	_, err := tx.Exec(`
		INSERT INTO product (snowflake_id, name, integral, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)
	`, product.SnowflakeId, product.Name, product.Integral, time.Now(), time.Now())
	return err
}

func UpdateProduct(tx *sql.Tx, product *entity.UpdateProductRequest, snowflakeId string) error {
	result, err := tx.Exec(`
		UPDATE product SET name = $1, integral = $2, updated_at = $3 WHERE snowflake_id = $4
	`, product.Name, product.Integral, time.Now(), snowflakeId)

	if err != nil {
		return fmt.Errorf("update product error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func DeleteProduct(tx *sql.Tx, snowflakeId string) error {
	result, err := tx.Exec(`
		UPDATE product SET deleted_at = $1 WHERE snowflake_id = $2
	`, time.Now(), snowflakeId)

	if err != nil {
		return fmt.Errorf("delete product error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
