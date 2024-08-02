package data

import (
	"fmt"
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

func GetIncomeListBySnowflakeId(snowflakeId int64, page *entity.GetIncomeListRequest) ([]*entity.GetIncomeListResponse, error) {

	baseSQL := `
		SELECT 
			snowflake_id, user_id, summary, integral, shipments, batch, created_at, updated_at
		FROM 
			income_expense
		WHERE 
			user_id = $1
		`
	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at)>=DATE(%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, page.PageSize, page.PageSize*(page.Page-1))

	rows, err := db.Query(baseSQL, executeParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	incomeList := make([]*entity.GetIncomeListResponse, 0)

	for rows.Next() {
		income := new(entity.GetIncomeListResponse)
		err := rows.Scan(
			&income.SnowflakeId,
			&income.Summary,
			&income.Integral,
			&income.Shipments,
			&income.Batch,
			&income.CreatedAt,
			&income.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		incomeList = append(incomeList, income)
	}

	return incomeList, nil
}

func GetIncomeCountBySnowflakeId(snowflakeId int64, page *entity.GetIncomeListRequest) (int64, error) {

	baseSQL := `
		SELECT
			COUNT(*)
		FROM
			income_expense
		WHERE
			user_id = $1
		`

	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at)>=DATE(%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	var count int64
	err := db.QueryRow(baseSQL, executeParams...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
