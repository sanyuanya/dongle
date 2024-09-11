package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func WithdrawalListCount(tx *sql.Tx, page *entity.WithdrawalPageListRequest) (int64, error) {
	baseSQL := `
		SELECT
			COUNT(*)
		FROM
			withdrawals w
		JOIN
			users u
		ON
			w.user_id = u.snowflake_id
		WHERE w.deleted_at IS NULL
	`

	executeParams := []interface{}{}
	paramIndex := 1

	if page.LifeCycle != 0 {
		baseSQL = baseSQL + fmt.Sprintf(" AND w.life_cycle = $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.LifeCycle)
	}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(w.created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND (u.nick LIKE $%d OR u.phone LIKE $%d)", paramIndex, paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	var count int64
	err := tx.QueryRow(baseSQL, executeParams...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询提现列表数量失败: %v", err)
	}

	return count, nil
}

func WithdrawalPageList(tx *sql.Tx, page *entity.WithdrawalPageListRequest) ([]*entity.WithdrawalList, error) {

	baseSQL := `
		SELECT
			w.snowflake_id, 
			w.user_id, 
			w.integral, 
			w.withdrawal_method, 
			w.life_cycle, 
			TO_CHAR(w.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at, 
			TO_CHAR(w.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at, 
			w.rejection,
			u.nick, 
			u.phone,
			w.detail_id,
			w.pay_id,
			w.initiate_time,
			w.update_time,
			w.open_id,
			w.mch_id,
			w.payment_status
		FROM
			withdrawals w
		JOIN
			users u
		ON
			w.user_id = u.snowflake_id AND u.deleted_at IS NULL
		
		WHERE w.deleted_at IS NULL
	`

	executeParams := []interface{}{}
	paramIndex := 1

	if page.LifeCycle != 0 {
		baseSQL = baseSQL + fmt.Sprintf(" AND w.life_cycle = $%d", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.LifeCycle)
	}

	if page.Keyword != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND (u.nick LIKE $%d OR u.phone LIKE $%d)", paramIndex, paramIndex)
		paramIndex++
		executeParams = append(executeParams, "%"+page.Keyword+"%")
	}

	if page.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(w.created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, page.Date)
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY w.created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, page.PageSize, page.PageSize*(page.Page-1))

	rows, err := tx.Query(baseSQL, executeParams...)
	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %v", err)
	}

	withdrawalList := make([]*entity.WithdrawalList, 0)

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
			&withdrawal.DetailId,
			&withdrawal.PayId,
			&withdrawal.InitiateTime,
			&withdrawal.UpdateTime,
			&withdrawal.OpenId,
			&withdrawal.MchId,
			&withdrawal.PaymentStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描提现列表失败: %v", err)
		}
		withdrawalList = append(withdrawalList, withdrawal)
	}

	return withdrawalList, nil
}

func ApplyForWithdrawal(tx *sql.Tx, applyForWithdrawal *entity.ApplyForWithdrawalRequest) error {

	baseSQL := `

		INSERT INTO
			withdrawals
			(snowflake_id, user_id, integral, withdrawal_method, life_cycle, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(baseSQL,
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

func ApprovalWithdrawal(tx *sql.Tx, snowflakeId string, rejection string, lifeCycle int64) error {

	baseSQL := `
		UPDATE
			withdrawals
		SET life_cycle=$1, rejection=$2, updated_at=$3
		WHERE snowflake_id = $4 AND deleted_at IS NULL
	`

	result, err := tx.Exec(baseSQL, lifeCycle, rejection, time.Now(), snowflakeId)
	if err != nil {
		return fmt.Errorf("审批提现失败: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if affected == 0 {
		return fmt.Errorf("未找到对应的提现记录")
	}

	return nil
}
func UpdateWithdrawalBatchId(tx *sql.Tx, transferDetailList []*pay.TransferDetail, batchResponse *pay.BatchesResponse) error {
	baseSQL := `
		UPDATE
			withdrawals
		SET pay_id=$1
		WHERE snowflake_id = $2 AND deleted_at IS NULL
	`
	for _, transferDetail := range transferDetailList {
		_, err := tx.Exec(baseSQL, batchResponse.OutBatchNo, transferDetail.OutDetailNo)
		if err != nil {
			return fmt.Errorf("更新提现记录失败: %v", err)
		}
	}
	return nil
}

func ComposeBatches(transferDetailList []*pay.TransferDetail) (*pay.BatchesRequest, error) {

	totalAmount := 0
	for _, transferDetail := range transferDetailList {
		totalAmount += transferDetail.TransferAmount
	}

	batchesRequest := &pay.BatchesRequest{
		AppId:              "wx370126c8bcf8d00c",
		OutBatchNo:         tools.SnowflakeUseCase.NextVal(),
		BatchName:          "分红奖励报销单",
		BatchRemark:        "分红奖励报销单",
		TotalAmount:        totalAmount,
		TotalNum:           len(transferDetailList),
		TransferDetailList: transferDetailList,
	}

	return batchesRequest, nil
}

func ComposeTransferDetail(tx *sql.Tx, snowflakeId string) (*pay.TransferDetail, error) {

	baseSQL := `
		SELECT 
			w.integral * 100, u.openid, u.nick
		FROM
			withdrawals w
		JOIN
			users u
		ON
			w.user_id = u.snowflake_id AND u.deleted_at IS NULL
		WHERE	
			w.life_cycle = 3 AND w.snowflake_id = $1 AND w.deleted_at IS NULL
	`

	transferDetail := &pay.TransferDetail{}
	err := tx.QueryRow(baseSQL, snowflakeId).Scan(
		&transferDetail.TransferAmount,
		&transferDetail.OpenId,
		&transferDetail.UserName)

	transferDetail.OutDetailNo = snowflakeId
	transferDetail.TransferRemark = "分红奖励"

	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %v", err)
	}
	return transferDetail, nil
}

func GetWithdrawalBySnowflakeId(tx *sql.Tx, snowflakeId string) (*entity.Withdrawal, error) {

	baseSQL := `
		SELECT
			user_id, integral
		FROM
			withdrawals
		WHERE
			snowflake_id=$1 AND deleted_at IS NULL
	`
	withdrawal := &entity.Withdrawal{}
	err := tx.QueryRow(baseSQL, snowflakeId).Scan(&withdrawal.UserId, &withdrawal.Integral)
	if err != nil {
		return nil, fmt.Errorf("查询提现失败: %v", err)
	}

	return withdrawal, nil
}

func GetWithdrawalByPaymentStatusIsFailAndPaymentStatusIsSuccess() ([]*entity.Order, error) {

	baseSQL := `
		SELECT
			pay_id, snowflake_id
		FROM
			withdrawals
		WHERE
			deleted_at IS NULL AND payment_status NOT IN ('FAIL', 'SUCCESS') AND life_cycle = 3
	`
	rows, err := db.Query(baseSQL)
	if err != nil {
		return nil, fmt.Errorf("查询提现列表失败: %v", err)
	}

	var withdrawalList []*entity.Order

	for rows.Next() {
		withdrawal := &entity.Order{}
		err := rows.Scan(&withdrawal.PayId, &withdrawal.SnowflakeId)
		if err != nil {
			return nil, fmt.Errorf("扫描提现列表失败: %v", err)
		}
		withdrawalList = append(withdrawalList, withdrawal)
	}

	return withdrawalList, nil
}

func GetWithdrawalCountByUserId(tx *sql.Tx, snowflakeId string, getWithdrawal *entity.GetWithdrawalListRequest) (int64, error) {

	baseSQL := `
		SELECT
			COUNT(*)
		FROM
			withdrawals
		WHERE
			user_id=$1 AND deleted_at IS NULL
	`
	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if getWithdrawal.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, getWithdrawal.Date)
	}

	var count int64
	err := tx.QueryRow(baseSQL, executeParams...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询提现列表数量失败: %v", err)
	}

	return count, nil
}

func GetWithdrawalListByUserId(tx *sql.Tx, snowflakeId string, getWithdrawal *entity.GetWithdrawalListRequest) ([]*entity.GetWithdrawalListResponse, error) {

	baseSQL := `
		SELECT
			snowflake_id, life_cycle, integral, withdrawal_method, TO_CHAR(created_at, 'YYYY-MM-DD') created_at, TO_CHAR(updated_at, 'YYYY-MM-DD') updated_at, rejection
		FROM
			withdrawals
		WHERE
			user_id=$1 AND deleted_at IS NULL
		
	`
	paramIndex := 2
	executeParams := []interface{}{snowflakeId}

	if getWithdrawal.Date != "" {
		baseSQL = baseSQL + fmt.Sprintf(" AND DATE(created_at) = DATE($%d)", paramIndex)
		paramIndex++
		executeParams = append(executeParams, getWithdrawal.Date)
	}

	baseSQL = baseSQL + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	executeParams = append(executeParams, getWithdrawal.PageSize, getWithdrawal.PageSize*(getWithdrawal.Page-1))

	rows, err := tx.Query(baseSQL, executeParams...)
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

	// Ensure withdrawalList is not nil
	if withdrawalList == nil {
		withdrawalList = []*entity.GetWithdrawalListResponse{}
	}

	return withdrawalList, nil
}

func UpdateWithdrawalInfoBySnowflakeId(tx *sql.Tx, withdrawal *pay.OutDetailNoResponse) error {

	baseSQL := `
		UPDATE
			withdrawals
		SET
			updated_at = $1,
			detail_id = $2,
			initiate_time = $3,
			update_time = $4,
			open_id = $5,
			mch_id = $6,
			rejection = $7,
			payment_status = $8
		WHERE
			snowflake_id = $9 AND deleted_at IS NULL
	`

	result, err := tx.Exec(baseSQL, time.Now(),
		withdrawal.DetailId,
		withdrawal.InitiateTime,
		withdrawal.UpdateTime,
		withdrawal.OpenId,
		withdrawal.Mchid,
		withdrawal.FailReason,
		withdrawal.DetailStatus,
		withdrawal.OutDetailNo,
	)

	if err != nil {
		return fmt.Errorf("更新提现记录失败: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if affected == 0 {
		return fmt.Errorf("未找到对应的提现记录")
	}

	return nil
}

func UpdateWithdrawalStatusBySnowflakeId(tx *sql.Tx, snowflakeId string, status string) error {

	baseSQL := `
		UPDATE
			withdrawals
		SET
			life_cycle = $1
		WHERE
			snowflake_id = $2 AND deleted_at IS NULL AND life_cycle = 3
	`

	result, err := tx.Exec(baseSQL, 4, snowflakeId)
	if err != nil {
		return fmt.Errorf("更新提现状态失败: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if affected == 0 {
		return fmt.Errorf("未找到对应的提现记录")
	}

	return nil
}
