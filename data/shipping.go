package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddShipping(tx *sql.Tx, payload *entity.AddShippingRequest) error {
	currentTime := time.Now()
	_, err := tx.Exec("INSERT INTO shipping (snowflake_id, order_id, task_id, third_order_id, order_number, e_order, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		payload.SnowflakeId,
		payload.OrderId,
		payload.TaskId,
		payload.ThirdOrderId,
		payload.OrderNumber,
		payload.EOrder,
		currentTime,
		currentTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetShippingByOrderId(tx *sql.Tx, orderId string) (*entity.GetShippingResponse, error) {

	shippingResponse := &entity.GetShippingResponse{}
	if err := tx.QueryRow("SELECT snowflake_id, order_id, task_id, third_order_id, order_number, e_order, created_at, updated_at FROM shipping WHERE order_id = $1", orderId).Scan(
		shippingResponse.SnowflakeId,
		shippingResponse.OrderId,
		shippingResponse.TaskId,
		shippingResponse.ThirdOrderId,
		shippingResponse.OrderNumber,
		shippingResponse.EOrder,
		shippingResponse.CreatedAt,
		shippingResponse.UpdatedAt,
	); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("查询下单信息失败")
	}

	return shippingResponse, nil
}
