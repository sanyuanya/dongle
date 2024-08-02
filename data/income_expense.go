package data

import (
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddIncomeExpense(addIncomeExpenseRequest *entity.AddIncomeExpenseRequest) error {

	baseSQL := `
		INSERT INTO 
			income_expense (snowflake_id, summary, integral, shipments, user_id, batch, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
	_, err := db.Exec(baseSQL,
		addIncomeExpenseRequest.SnowflakeId,
		addIncomeExpenseRequest.Summary,
		addIncomeExpenseRequest.Integral,
		addIncomeExpenseRequest.Shipments,
		addIncomeExpenseRequest.UserId,
		addIncomeExpenseRequest.Batch,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
