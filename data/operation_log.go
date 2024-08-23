package data

import (
	"database/sql"
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
