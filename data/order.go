package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddOrder(tx *sql.Tx, payload *entity.AddOrder) error {
	_, err := tx.Exec(`
		INSERT INTO "order" (
			snowflake_id,
			address_id,
			consignee,
			phone_number,
			location,
			detailed_address,
			open_id,
			user_id,
			expiration_time,
			out_trade_no,
			order_state,
			currency,
			open_id,
			total,
			prepay_id,
			pay_sign,
			pay_timestamp,
			sign_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`, payload.SnowflakeId,
		payload.AddressId,
		payload.Consignee,
		payload.PhoneNumber,
		payload.Location,
		payload.DetailedAddress,
		payload.OpenId,
		payload.UserId,
		payload.ExpirationTime,
		payload.OutTradeNo,
		payload.OrderState,
		payload.Currency,
		payload.OpenId,
		payload.Total,
		payload.PrepayId,
		payload.PaySign,
		payload.PayTimestamp,
		payload.SignType,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetOrderList(tx *sql.Tx, payload *entity.GetOrderListRequest) ([]*entity.GetOrderListResponse, error) {
	baseQuery := `
		SELECT
			o.snowflake_id,
			o.transaction_id,
			o.app_id,
			o.mch_id,
			o.trade_type,
			o.trade_state,
			o.trade_state_desc,
			o.bank_type,
			o.success_time,
			o.open_id,
			o.user_id,
			o.total,
			o.payer_total,
			o.currency,
			o.payer_currency,
			o.out_trade_no,
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			TO_CHAR(o.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at,
			o.prepay_id,
			o.expiration_time,
			o.address_id,
			o.consignee,
			o.phone_number,
			o.location,
			o.detailed_address,
			o.order_state,
			o.nonce_str,
			o.pay_sign,
			o.pay_timestamp,
			o.sign_type,
			u.nick,
			u.phone
		FROM
			"order" o
		LEFT JOIN 
			users u
		ON 
			o.user_id = u.snowflake_id AND u.deleted_at IS NULL
		WHERE
			o.deleted_at IS NULL
			`

	paramsIndex := 1
	executeParams := []any{}
	if payload.Keyword != "" {
		baseQuery += fmt.Sprintf(" AND (u.nick LIKE $%d OR u.phone LIKE $%d)", paramsIndex, paramsIndex+1)
		executeParams = append(executeParams, "%"+payload.Keyword+"%", "%"+payload.Keyword+"%")
		paramsIndex += 2
	}

	if payload.OutTradeNo != "" {
		baseQuery += fmt.Sprintf(" AND o.out_trade_no = $%d", paramsIndex)
		executeParams = append(executeParams, payload.OutTradeNo)
		paramsIndex++
	}

	if payload.OpenId != "" {
		baseQuery += fmt.Sprintf(" AND o.open_id = $%d", paramsIndex)
		executeParams = append(executeParams, payload.OpenId)
		paramsIndex++
	}

	if payload.TradeState != "" {
		baseQuery += fmt.Sprintf(" AND o.trade_state = $%d", paramsIndex)
		executeParams = append(executeParams, payload.TradeState)
		paramsIndex++
	}

	if payload.Status != 0 {
		baseQuery += fmt.Sprintf(" AND o.order_state = $%d", paramsIndex)
		executeParams = append(executeParams, payload.Status)
		paramsIndex++
	}

	baseQuery += fmt.Sprintf(" ORDER BY o.created_at DESC LIMIT $%d OFFSET $%d", paramsIndex, paramsIndex+1)
	executeParams = append(executeParams, payload.PageSize, (payload.Page-1)*payload.PageSize)

	rows, err := tx.Query(baseQuery, executeParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.GetOrderListResponse, 0)
	for rows.Next() {
		order := &entity.GetOrderListResponse{}
		err := rows.Scan(
			&order.SnowflakeId,
			&order.TransactionId,
			&order.AppId,
			&order.MchId,
			&order.TradeType,
			&order.TradeState,
			&order.TradeStateDesc,
			&order.BankType,
			&order.SuccessTime,
			&order.OpenId,
			&order.UserId,
			&order.Total,
			&order.PayerTotal,
			&order.Currency,
			&order.PayerCurrency,
			&order.OutTradeNo,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.PrepayId,
			&order.ExpirationTime,
			&order.AddressId,
			&order.Consignee,
			&order.PhoneNumber,
			&order.Location,
			&order.DetailedAddress,
			&order.OrderState,
			&order.NonceStr,
			&order.PaySign,
			&order.PayTimestamp,
			&order.SignType,
			&order.Nick,
			&order.Phone,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrderCount(tx *sql.Tx, payload *entity.GetOrderListRequest) (uint64, error) {

	baseQuery := `
		SELECT
			COUNT(o.snowflake_id)
		FROM
			"order" o
		LEFT JOIN 
			users u
		ON 
			o.user_id = u.snowflake_id AND u.deleted_at IS NULL
		WHERE
			o.deleted_at IS NULL
		`

	paramsIndex := 1
	executeParams := []any{}
	if payload.Keyword != "" {
		baseQuery += fmt.Sprintf(" AND (u.nick LIKE $%d OR u.phone LIKE $%d)", paramsIndex, paramsIndex+1)
		executeParams = append(executeParams, "%"+payload.Keyword+"%", "%"+payload.Keyword+"%")
		paramsIndex += 2
	}

	if payload.OutTradeNo != "" {
		baseQuery += fmt.Sprintf(" AND o.out_trade_no = $%d", paramsIndex)
		executeParams = append(executeParams, payload.OutTradeNo)
		paramsIndex++
	}

	if payload.OpenId != "" {
		baseQuery += fmt.Sprintf(" AND o.open_id = $%d", paramsIndex)
		executeParams = append(executeParams, payload.OpenId)
		paramsIndex++
	}

	var count uint64
	err := tx.QueryRow(baseQuery, executeParams...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func UpdateOrder(tx *sql.Tx, payload *entity.DecryptResourceResponse) error {
	_, err := tx.Exec(`
		UPDATE
			"order"
		SET
			transaction_id = $1,
			app_id = $2,
			mch_id = $3,
			trade_type = $4,
			trade_state = $5,
			trade_state_desc = $6,
			bank_type = $7,
			success_time = $8,
			open_id = $9,
			total = $10,
			payer_total = $11,
			currency = $12,
			payer_currency = $13,
			updated_at = NOW(),
			order_state = 2
		WHERE
			out_trade_no = $14
	`, payload.TransactionId,
		payload.AppId,
		payload.MchId,
		payload.TradeType,
		payload.TradeState,
		payload.TradeStateDesc,
		payload.BankType,
		payload.SuccessTime,
		payload.Payer.OpenId,
		payload.Amount.Total,
		payload.Amount.PayerTotal,
		payload.Amount.Currency,
		payload.Amount.PayerCurrency,
		payload.OutTradeNo,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetOrderByTradeState() ([]string, error) {

	rows, err := db.Query(`
		SELECT
			out_trade_no
		FROM
			"order"
		WHERE
			trade_state = ''
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTradeNos := make([]string, 0)
	for rows.Next() {
		var outTradeNo string
		err := rows.Scan(&outTradeNo)
		if err != nil {
			return nil, err
		}
		outTradeNos = append(outTradeNos, outTradeNo)
	}
	return outTradeNos, nil
}

func GetOrderExpired() ([]string, error) {

	timestamp := time.Now().Add(-10 * time.Second).Unix()

	rows, err := db.Query(`
		SELECT
			out_trade_no
		FROM
			"order"
		WHERE
			expiration_time <= $1 AND order_state != 99
	`, timestamp)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outTradeNos := make([]string, 0)
	for rows.Next() {
		var outTradeNo string
		err := rows.Scan(&outTradeNo)
		if err != nil {
			return nil, err
		}
		outTradeNos = append(outTradeNos, outTradeNo)
	}
	return outTradeNos, nil
}

func UpdateOrderByOutTradeNo(tx *sql.Tx, payload *entity.UpdateOrderByOutTradeNo) error {
	_, err := tx.Exec(`
		UPDATE
			"order"
		SET
			order_state = $1
		WHERE
			out_trade_no = $2
	`,
		payload.Status,
		payload.OutTradeNo,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetOrderDetail(tx *sql.Tx, snowflakeId string) (*entity.GetOrderListResponse, error) {
	baseQuery := `
	SELECT
		o.snowflake_id,
		o.transaction_id,
		o.app_id,
		o.mch_id,
		o.trade_type,
		o.trade_state,
		o.trade_state_desc,
		o.bank_type,
		o.success_time,
		o.open_id,
		o.user_id,
		o.total,
		o.payer_total,
		o.currency,
		o.payer_currency,
		o.out_trade_no,
		TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
		TO_CHAR(o.updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at,
		o.prepay_id,
		o.expiration_time,
		o.address_id,
		o.consignee,
		o.phone_number,
		o.location,
		o.detailed_address,
		o.order_state,
		o.nonce_str,
		o.pay_sign,
		o.pay_timestamp,
		o.sign_type,
		u.nick,
		u.phone
	FROM
		"order" o
	LEFT JOIN 
		users u
	ON 
		o.user_id = u.snowflake_id AND u.deleted_at IS NULL
	WHERE
		o.deleted_at IS NULL AND o.snowflake_id = $1
		`

	order := &entity.GetOrderListResponse{}

	err := tx.QueryRow(baseQuery, snowflakeId).Scan(&order.SnowflakeId,
		&order.TransactionId,
		&order.AppId,
		&order.MchId,
		&order.TradeType,
		&order.TradeState,
		&order.TradeStateDesc,
		&order.BankType,
		&order.SuccessTime,
		&order.OpenId,
		&order.UserId,
		&order.Total,
		&order.PayerTotal,
		&order.Currency,
		&order.PayerCurrency,
		&order.OutTradeNo,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.PrepayId,
		&order.ExpirationTime,
		&order.AddressId,
		&order.Consignee,
		&order.PhoneNumber,
		&order.Location,
		&order.DetailedAddress,
		&order.OrderState,
		&order.NonceStr,
		&order.PaySign,
		&order.PayTimestamp,
		&order.SignType,
		&order.Nick,
		&order.Phone)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("未找到订单信息：%v", err)
		}
		return nil, fmt.Errorf("查询订单信息失败： %v", err)
	}
	return order, nil
}
