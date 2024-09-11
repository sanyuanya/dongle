package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddProductCategories(tx *sql.Tx, addProductCategoriesRequest *entity.AddProductCategoriesRequest) error {
	_, err := tx.Exec("INSERT INTO product_categories (snowflake_id, name, status, sorting) VALUES ($1, $2, $3, $4)", addProductCategoriesRequest.SnowflakeId, addProductCategoriesRequest.Name, addProductCategoriesRequest.Status, addProductCategoriesRequest.Sorting)
	if err != nil {
		return err
	}

	return nil
}

func DeleteProductCategories(tx *sql.Tx, id string) error {
	result, err := tx.Exec("UPDATE product_categories SET deleted_at = $1 WHERE snowflake_id = $2 AND deleted_at IS NULL", time.Now(), id)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("未找到记录")
	}

	return nil
}

func UpdateProductCategories(tx *sql.Tx, updateProductCategoriesRequest *entity.UpdateProductCategoriesRequest) error {
	_, err := tx.Exec("UPDATE product_categories SET name = $1, status = $2, sorting = $3 WHERE snowflake_id = $4 AND deleted_at IS NULL", updateProductCategoriesRequest.Name, updateProductCategoriesRequest.Status, updateProductCategoriesRequest.Sorting, updateProductCategoriesRequest.SnowflakeId)
	if err != nil {
		return err
	}

	return nil
}

func GetProductCategoriesList(tx *sql.Tx, payload *entity.GetProductCategoriesListRequest) ([]entity.ProductCategories, error) {
	productCategories := make([]entity.ProductCategories, 0)

	baseQuery := `
		SELECT 
			snowflake_id, 
			name, 
			status, 
			sorting 
		FROM 
			product_categories
		WHERE deleted_at IS NULL
	`

	paramIndex := 1
	executeParams := []interface{}{}

	if payload.Status != 0 {
		baseQuery += fmt.Sprintf(" AND status = $%d", paramIndex)
		executeParams = append(executeParams, payload.Status)
		paramIndex++
	}

	if payload.Keyword != "" {
		baseQuery += fmt.Sprintf(" AND name LIKE $%d", paramIndex)
		executeParams = append(executeParams, "%"+payload.Keyword+"%")
		paramIndex++
	}

	baseQuery += fmt.Sprintf(" ORDER BY sorting DESC, created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, payload.PageSize, (payload.Page-1)*payload.PageSize)

	rows, err := tx.Query(baseQuery, executeParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var productCategory entity.ProductCategories
		err := rows.Scan(&productCategory.SnowflakeId, &productCategory.Name, &productCategory.Status, &productCategory.Sorting)
		if err != nil {
			return nil, err
		}
		productCategories = append(productCategories, productCategory)
	}

	return productCategories, nil
}

func GetProductCategoriesListCount(tx *sql.Tx, payload *entity.GetProductCategoriesListRequest) (int64, error) {
	var total int64

	baseQuery := `
		SELECT 
			COUNT(*) 
		FROM 
			product_categories
		WHERE deleted_at IS NULL
	`

	paramIndex := 1
	executeParams := []interface{}{}

	if payload.Status != 0 {
		baseQuery += fmt.Sprintf(" AND status = $%d", paramIndex)
		executeParams = append(executeParams, payload.Status)
		paramIndex++
	}

	if payload.Keyword != "" {
		baseQuery += fmt.Sprintf(" AND name LIKE $%d", paramIndex)
		executeParams = append(executeParams, "%"+payload.Keyword+"%")
		paramIndex++
	}

	err := tx.QueryRow(baseQuery, executeParams...).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func FindByProductCategoriesName(tx *sql.Tx, name string) (string, error) {
	var id string
	err := tx.QueryRow("SELECT snowflake_id FROM product_categories WHERE name = $1 AND deleted_at IS NULL", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return id, nil
}

func FindByProductCategoriesId(tx *sql.Tx, id string) (string, error) {
	var name string
	err := tx.QueryRow("SELECT name FROM product_categories WHERE snowflake_id = $1 AND deleted_at IS NULL", id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}
