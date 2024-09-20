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
			o.out_trade_no,
			o.total,
			o.payer_total,
			o.success_time,
			o.trade_type,
			o.trade_state,
			o.order_state,
			u.nick,
			u.phone,
			u.expiration_time,
			TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS') created_at
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
		baseQuery += fmt.Sprintf(" AND (u.nick LIKE $%d OR phone LIKE $%d)", paramsIndex, paramsIndex+1)
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
		err := rows.Scan(&order.SnowflakeId, &order.OutTradeNo, &order.Total, &order.PayerTotal, &order.SuccessTime, &order.TradeType, &order.TradeState, &order.OrderState, &order.Nick, &order.Phone, &order.ExpirationTime, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

