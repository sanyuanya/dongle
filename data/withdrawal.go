package data

import (
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func WithdrawalListCount(page *entity.WithdrawalPageListRequest) (int64, error) {
	baseSQL := `
		SELECT
			COUNT(*)
		FROM
			withdrawals w
		JOIN
			users u
		ON
			withdrawals.user_id = users.snowflake_id
		
		WHERE deleted_at IS NULL
	`

	executeParams := []interface{}{}
	paramIndex := 1

	if page.LifeCycle != 0 {
		baseSQL = baseSQL + fmt.Sprintf(" AND w.life_cycle=$%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.LifeCycle)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND (u.nick LIKE $%d OR u.phone LIKE $%d)", paramIndex, paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	var count int64
	err := db.QueryRow(baseSQL, executeParams...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询提现列表数量失败: %v", err)
	}

	return count, nil
}

func WithdrawalPageList(page *entity.WithdrawalPageListRequest) ([]*entity.WithdrawalList, error) {

	baseSQL := `
		SELECT
			w.snowflake_id, w.user_id, w.integral, w.withdrawal_method, w.life_cycle, w.created_at, w.updated_at, w.rejection,
			u.nick, u.phone
		FROM
			withdrawals w
		JOIN
			users u
		ON
			withdrawals.user_id = users.snowflake_id
		
		WHERE deleted_at IS NULL

		-- ORDER BY created_at DESC LIMIT 

	`

	executeParams := []interface{}{}
	paramIndex := 1

	if page.LifeCycle != 0 {
		baseSQL = baseSQL + fmt.Sprintf(" AND w.life_cycle=$%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.LifeCycle)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND (u.nick LIKE $%d OR u.phone LIKE $%d)", paramIndex, paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY w.created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, page.PageSize, page.PageSize*(page.Page-1))

	rows, err := db.Query(baseSQL, executeParams...)
	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %v", err)
	}

	var withdrawalList []*entity.WithdrawalList

	for rows.Next() {
		withdrawal := &entity.WithdrawalList{}
		err := rows.Scan(
			&withdrawal.SnowflakeId,
			&withdrawal.UserId,
			&withdrawal.Integral,
			&withdrawal.WithdrawalMethod,
			&withdrawal.LifeCycle,
			&withdrawal.CreatedAt,
			&withdrawal.UpdatedAt,
			&withdrawal.Rejection,
			&withdrawal.Nick,
			&withdrawal.Phone,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描提现列表失败: %v", err)
		}
		withdrawalList = append(withdrawalList, withdrawal)
	}

	return withdrawalList, nil
}

func ApplyForWithdrawal(applyForWithdrawal *entity.ApplyForWithdrawalRequest) error {

	baseSQL := `

		INSERT INTO
			withdrawals
			(snowflake_id, user_id, integral, withdrawal_method, life_cycle, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(baseSQL,
		applyForWithdrawal.SnowflakeId,
		applyForWithdrawal.UserId,
		applyForWithdrawal.Integral,
		applyForWithdrawal.WithdrawalMethod,
		1,
		time.Now(),
		time.Now())
	if err != nil {
		return fmt.Errorf("申请提现失败: %v", err)
	}

	return nil
}

func ApprovalWithdrawal(approvalWithdrawalRequest *entity.ApprovalWithdrawalRequest) error {

	baseSQL := `
		UPDATE
			withdrawals
		SET life_cycle=$1, rejection=$2, updated_at=$3
		WHERE snowflake_id = $4
	`
	for _, snowflakeId := range approvalWithdrawalRequest.ApprovalList {

		_, err := db.Exec(baseSQL, approvalWithdrawalRequest.LifeCycle, approvalWithdrawalRequest.Rejection, time.Now(), snowflakeId)
		if err != nil {
			return fmt.Errorf("审批提现失败: %v", err)
		}
	}

	return nil
}

func GetWithdrawalListBySnowflakeId(snowflakeId int64, getWithdrawal *entity.GetWithdrawalListRequest) ([]*entity.GetWithdrawalListResponse, error) {

	baseSQL := `
		SELECT
			snowflake_id, life_cycle, integral, withdrawal_method, created_at, updated_at, rejection
		FROM
			withdrawals
		WHERE
			user_id=$1 AND deleted_at IS NULL
		
	`
	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if getWithdrawal.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at)>=DATE(%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, getWithdrawal.Date)
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, getWithdrawal.PageSize, getWithdrawal.PageSize*(getWithdrawal.Page-1))

	rows, err := db.Query(baseSQL, executeParams...)
	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %v", err)
	}

	var withdrawalList []*entity.GetWithdrawalListResponse

	for rows.Next() {
		withdrawal := &entity.GetWithdrawalListResponse{}
		err := rows.Scan(
			&withdrawal.SnowflakeId,
			&withdrawal.LifeCycle,
			&withdrawal.Integral,
			&withdrawal.WithdrawalMethod,
			&withdrawal.CreatedAt,
			&withdrawal.UpdatedAt,
			&withdrawal.Rejection,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描提现列表失败: %v", err)
		}
		withdrawalList = append(withdrawalList, withdrawal)
	}

	return withdrawalList, nil
}
