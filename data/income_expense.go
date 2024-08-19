package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddIncomeExpense(addIncomeExpenseRequest *entity.AddIncomeExpenseRequest) error {

	baseSQL := `
		INSERT INTO 
			income_expense (snowflake_id, summary, integral, shipments, user_id, batch, created_at, updated_at, product_id, product_integral)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		addIncomeExpenseRequest.ProductId,
		addIncomeExpenseRequest.ProductIntegral,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetIncomeListBySnowflakeId(tx *sql.Tx, snowflakeId string, page *entity.GetIncomeListRequest) ([]*entity.GetIncomeListResponse, error) {

	baseSQL := `
		SELECT 
			i.snowflake_id, i.user_id, i.summary, i.integral, i.shipments, i.batch, TO_CHAR(i.created_at, 'YYYY-MM-DD') created_at, TO_CHAR(i.updated_at, 'YYYY-MM-DD') updated_at,
		  i.product_integral, p.name
		FROM 
			income_expense i
		JOIN
			product p
		ON
			i.product_id = p.snowflake_id AND p.deleted_at IS NULL
		WHERE 
			user_id = $1 AND deleted_at IS NULL
		`
	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, page.PageSize, page.PageSize*(page.Page-1))

	rows, err := tx.Query(baseSQL, executeParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	incomeList := make([]*entity.GetIncomeListResponse, 0)

	for rows.Next() {
		income := new(entity.GetIncomeListResponse)
		err := rows.Scan(
			&income.SnowflakeId,
			&income.UserId,
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

	if incomeList == nil {
		incomeList = []*entity.GetIncomeListResponse{}
	}
	return incomeList, nil
}

func GetIncomeCountBySnowflakeId(tx *sql.Tx, snowflakeId string, page *entity.GetIncomeListRequest) (int64, error) {

	baseSQL := `
		SELECT
			COUNT(*)
		FROM
			income_expense
		WHERE
			user_id = $1 AND deleted_at IS NULL
		`

	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	var count int64
	err := tx.QueryRow(baseSQL, executeParams...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func UpdateIncomeExpense(tx *sql.Tx, new string, old string) error {

	baseSQL := `
		UPDATE
			income_expense
		SET
			user_id=$1
		WHERE
			user_id=$2 AND deleted_at IS NULL
			`
	_, err := tx.Exec(baseSQL, new, old)

	if err != nil {
		return err
	}

	return nil
}

func IncomeListCount(tx *sql.Tx, page *entity.IncomePageListExpenseRequest) (int64, error) {

	baseSQL := `
		SELECT
			COUNT(*)
		FROM
			income_expense i
		JOIN
			users u
		ON
			i.user_id = u.snowflake_id AND u.deleted_at IS NULL
		JOIN
			product p
		ON
			i.product_id = p.snowflake_id
		WHERE
			i.deleted_at IS NULL
		`
	paramIndex := 1
	executeParams := []interface{}{}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(i.created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND u.phone LIKE $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	if page.UserId != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND i.user_id = $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.UserId)
	}

	var count int64
	err := tx.QueryRow(baseSQL, executeParams...).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func IncomePageList(tx *sql.Tx, page *entity.IncomePageListExpenseRequest) ([]*entity.IncomePageListExpenseResponse, error) {

	baseSQL := `
		SELECT 
			i.snowflake_id, i.user_id, i.summary, i.integral, i.shipments, i.batch, TO_CHAR(i.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at, TO_CHAR(i.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at, u.nick, u.phone,
		  i.product_integral, p.name
		FROM 
			income_expense i
		JOIN 
			users u
		ON
			i.user_id = u.snowflake_id AND u.deleted_at IS NULL
		JOIN
			product p
		ON
			i.product_id = p.snowflake_id
		WHERE 
			i.deleted_at IS NULL
		`
	paramIndex := 1
	executeParams := []interface{}{}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(i.created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND u.phone LIKE $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	if page.UserId != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND i.user_id = $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.UserId)
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY i.created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, page.PageSize, page.PageSize*(page.Page-1))

	rows, err := tx.Query(baseSQL, executeParams...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	incomeList := make([]*entity.IncomePageListExpenseResponse, 0)

	for rows.Next() {
		income := new(entity.IncomePageListExpenseResponse)
		err := rows.Scan(
			&income.SnowflakeId,
			&income.UserId,
			&income.Summary,
			&income.Integral,
			&income.Shipments,
			&income.Batch,
			&income.CreatedAt,
			&income.UpdatedAt,
			&income.Nick,
			&income.Phone,
			&income.ProductIntegral,
			&income.ProductName,
		)
		if err != nil {
			return nil, err
		}
		incomeList = append(incomeList, income)
	}

	if incomeList == nil {
		incomeList = []*entity.IncomePageListExpenseResponse{}
	}

	return incomeList, nil
}
