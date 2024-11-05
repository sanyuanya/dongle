package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddItem(tx *sql.Tx, addItem *entity.AddItem) error {

	_, err := tx.Exec("INSERT INTO commodity (snowflake_id, name, code, categories_id, status, description, sorting) VALUES ($1, $2, $3, $4, $5, $6, $7)", addItem.SnowflakeId, addItem.Name, addItem.Code, addItem.CategoriesId, addItem.Status, addItem.Description, addItem.Sorting)
	if err != nil {
		return err
	}

	return nil
}

func DeleteItem(tx *sql.Tx, snowflake_id string) error {

	_, err := tx.Exec("UPDATE commodity SET deleted_at = $1 WHERE snowflake_id = $2 AND deleted_at IS NULL", time.Now(), snowflake_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateItem(tx *sql.Tx, updateItem *entity.UpdateItem) error {

	_, err := tx.Exec("UPDATE commodity SET name = $1, code = $2, categories_id = $3, status = $4, description = $5, sorting = $6 WHERE snowflake_id = $7 AND deleted_at IS NULL", updateItem.Name, updateItem.Code, updateItem.CategoriesId, updateItem.Status, updateItem.Description, updateItem.Sorting, updateItem.SnowflakeId)
	if err != nil {
		return err
	}

	return nil
}

func ItemList(tx *sql.Tx, itemPage *entity.ItemPage) ([]*entity.Item, error) {

	baseQuery := `
		SELECT
			snowflake_id,
			name,
			code,
			categories_id,
			status,
			description,
			sorting,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') created_at
		FROM
			commodity
		WHERE
			deleted_at IS NULL
			`

	paramsIndex := 1
	executeParams := []interface{}{}

	if itemPage.Name != "" {
		baseQuery += fmt.Sprintf(" AND name LIKE $%d", paramsIndex)
		executeParams = append(executeParams, "%"+itemPage.Name+"%")
		paramsIndex++
	}

	if itemPage.CategoriesId != "" {
		baseQuery += fmt.Sprintf(" AND categories_id = $%d", paramsIndex)
		executeParams = append(executeParams, itemPage.CategoriesId)
		paramsIndex++
	}

	if itemPage.Status != 0 {
		baseQuery += fmt.Sprintf(" AND status = $%d", paramsIndex)
		executeParams = append(executeParams, itemPage.Status)
		paramsIndex++
	}

	baseQuery += fmt.Sprintf(" ORDER BY sorting DESC, created_at DESC LIMIT $%d OFFSET $%d", paramsIndex, paramsIndex+1)
	executeParams = append(executeParams, itemPage.PageSize, (itemPage.Page-1)*itemPage.PageSize)

	rows, err := tx.Query(baseQuery, executeParams...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemList := make([]*entity.Item, 0)

	for rows.Next() {
		item := &entity.Item{}
		err := rows.Scan(&item.SnowflakeId, &item.Name, &item.Code, &item.CategoriesId, &item.Status, &item.Description, &item.Sorting, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		itemList = append(itemList, item)
	}

	return itemList, nil

}

func ItemListCount(tx *sql.Tx, itemPage *entity.ItemPage) (int64, error) {
	var total int64
	baseQuery := `
		SELECT
			COUNT(1)
		FROM
			commodity
		WHERE
			deleted_at IS NULL
		`

	paramsIndex := 1
	executeParams := []interface{}{}

	if itemPage.Name != "" {
		baseQuery += fmt.Sprintf(" AND name LIKE $%d", paramsIndex)
		executeParams = append(executeParams, "%"+itemPage.Name+"%")
		paramsIndex++
	}

	if itemPage.CategoriesId != "" {
		baseQuery += fmt.Sprintf(" AND categories_id = $%d", paramsIndex)
		executeParams = append(executeParams, itemPage.CategoriesId)
		paramsIndex++
	}

	if itemPage.Status != 0 {
		baseQuery += fmt.Sprintf(" AND status = $%d", paramsIndex)
		executeParams = append(executeParams, itemPage.Status)
		paramsIndex++
	}

	err := tx.QueryRow(baseQuery, executeParams...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func FindByItemCode(tx *sql.Tx, code string) (string, error) {
	var snowflakeId string
	err := tx.QueryRow("SELECT snowflake_id FROM commodity WHERE code = $1 AND deleted_at IS NULL", code).Scan(&snowflakeId)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", err
		}
	}
	return snowflakeId, nil
}

func FindByItemId(tx *sql.Tx, itemId string) (*entity.Item, error) {
	item := &entity.Item{}
	err := tx.QueryRow("SELECT snowflake_id, name, code, categories_id, status, description, created_at FROM commodity WHERE snowflake_id = $1 AND deleted_at IS NULL", itemId).Scan(
		&item.SnowflakeId,
		&item.Name,
		&item.Code,
		&item.CategoriesId,
		&item.Status,
		&item.Description,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func UpdateItemStatus(tx *sql.Tx, payload *entity.UpdateItem) error {
	_, err := tx.Exec("UPDATE commodity SET status = $1 WHERE snowflake_id = $2 AND deleted_at IS NULL", payload.Status, payload.SnowflakeId)
	if err != nil {
		return err
	}
	return nil
}
