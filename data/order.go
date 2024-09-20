package data

import (
	"database/sql"
	"fmt"

	"github.com/sanyuanya/dongle/entity"
)

func AddOrder(tx *sql.Tx, payload *entity.AddOrder) error {

	_, err := tx.Exec("INSERT INTO `order` (snowflake_id, commodity_id, sku_id, quantity, user_id) VALUES (?, ?, ?, ?, ?)", payload.SnowflakeId, payload.CommodityId, payload.SkuId, payload.Quantity, payload.UserId)
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
			u.phone,
		FROM
			order o
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
			order o
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
