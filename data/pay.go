package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/pay"
)

func CreatePay(tx *sql.Tx, totalAmount int, totalNum int, batchResponse *pay.BatchesResponse) error {
	baseSQL := `
		INSERT
		INTO
			pay
			(snowflake_id, name, total_amount, total_num, status, created_at, updated_at, batch_id)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.Exec(baseSQL,
		batchResponse.OutBatchNo,
		"分红奖励",
		totalAmount,
		totalNum,
		batchResponse.BatchStatus,
		time.Now(),
		time.Now(),
		batchResponse.BatchId,
	)

	if err != nil {
		return fmt.Errorf("创建支付记录失败: %v", err)
	}

	return nil
}

func UpdatePay(tx *sql.Tx, transferBatch pay.TransferBatch) error {

	baseSQL := `
		UPDATE pay
		SET batch_id = $1,
			name = $2,
			remark = $3,
			total_amount = $4,
			total_num = $5,
			status = $6,
			created_at = $7,
			updated_at = $8,
			close_reason = $9,
			success_amount = $10,
			success_num = $11,
			fail_amount = $12,
			fail_num = $13,
			create_time = $14,
			update_time = $15
		WHERE snowflake_id = $16
	`
	result, err := tx.Exec(baseSQL,
		transferBatch.BatchId,
		transferBatch.BatchName,
		transferBatch.BatchRemark,
		transferBatch.TotalAmount,
		transferBatch.TotalNum,
		transferBatch.BatchStatus,
		time.Now(),
		time.Now(),
		transferBatch.CloseReason,
		transferBatch.SuccessAmount,
		transferBatch.SuccessNum,
		transferBatch.FailAmount,
		transferBatch.FailNum,
		transferBatch.CreateTime,
		transferBatch.UpdateTime,
		transferBatch.Mchid,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("未找到对应的支付单据")
	}

	return nil
}
