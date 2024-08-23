package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddOperationLog(tx *sql.Tx, body *entity.AddOperationLogRequest) error {

	sql := `INSERT INTO 
					operation_log (snowflake_id, operation_id, income_expense_id, user_id, before_updating_shipments, after_updating_shipments, summary, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := tx.Exec(sql, body.SnowflakeId, body.OperationId, body.IncomeExpenseId, body.UserId, body.BeforeUpdatingShipments, body.AfterUpdatingShipments, body.Summary, time.Now(), time.Now())
	return err
}

func GetOperationLogList(tx *sql.Tx, body *entity.GetOperationLogListRequest) ([]*entity.OperationLog, error) {

	sql := `SELECT 
			o.snowflake_id, 
			o.operation_id,
			o.income_expense_id, 
			o.user_id, 
			o.before_updating_shipments, 
			o.after_updating_shipments, 
			o.summary, 
			o.created_at, 
			o.updated_at,
			u.nick,
			u.phone,
			p.name
		FROM operation_log o
		JOIN income_expense ie ON o.income_expense_id = ie.snowflake_id
		JOIN users u ON o.user_id = u.snowflake_id
		JOIN product p ON ie.product_id = p.snowflake_id
		WHERE o.deleted_at IS NULL
	`
	paramIndex := 1
	executeParams := []interface{}{}
	if body.UserId != "" {
		sql += fmt.Sprintf(" AND o.user_id = $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, body.UserId)
	}

	sql += fmt.Sprintf(" ORDER BY o.created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	executeParams = append(executeParams, body.PageSize, (body.Page-1)*body.PageSize)

	rows, err := tx.Query(sql, executeParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*entity.OperationLog, 0)

	for rows.Next() {
		item := new(entity.OperationLog)

		err = rows.Scan(&item.SnowflakeId, &item.OperationId, &item.IncomeExpenseId, &item.UserId, &item.BeforeUpdatingShipments, &item.AfterUpdatingShipments, &item.Summary, &item.CreatedAt, &item.UpdatedAt, &item.UserName, &item.Phone, &item.ProductName)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func GetOperationLogCount(tx *sql.Tx, body *entity.GetOperationLogListRequest) (int64, error) {
	sql := `SELECT COUNT(*) FROM operation_log WHERE deleted_at IS NULL`

	paramIndex := 1
	executeParams := []interface{}{}
	if body.UserId != "" {
		sql += fmt.Sprintf(" AND user_id = $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, body.UserId)
	}

	var count int64
	err := tx.QueryRow(sql, executeParams...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
