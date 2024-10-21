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
	if err := tx.QueryRow("SELECT snowflake_id, order_id, task_id, third_order_id, order_number, e_order, created_at, updated_at, status, user_cancel_msg, system_cancel_msg, courier_name, courier_mobile, net_tel, net_code, weight, def_price, volume, actual_weight, print_task_id, label, pickup_code FROM shipping WHERE order_id = $1", orderId).Scan(
		shippingResponse.SnowflakeId,
		shippingResponse.OrderId,
		shippingResponse.TaskId,
		shippingResponse.ThirdOrderId,
		shippingResponse.OrderNumber,
		shippingResponse.EOrder,
		shippingResponse.CreatedAt,
		shippingResponse.UpdatedAt,
		shippingResponse.Status,
		shippingResponse.UserCancelMsg,
		shippingResponse.SystemCancelMsg,
		shippingResponse.CourierName,
		shippingResponse.CourierMobile,
		shippingResponse.NetTel,
		shippingResponse.NetCode,
		shippingResponse.Weight,
		shippingResponse.DefPrice,
		shippingResponse.Volume,
		shippingResponse.ActualWeight,
		shippingResponse.PrintTaskId,
		shippingResponse.Label,
		shippingResponse.PickupCode,
	); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("查询下单信息失败")
	}

	return shippingResponse, nil
}

func UpdatedShippingByThirdOrderId(tx *sql.Tx, payload *entity.OrderCallbackData) error {

	if _, err := tx.Exec(`UPDATE shipping SET updated_at = $1, status = $2, user_cancel_msg = $3, system_cancel_msg = $4, courier_name = $5, courier_mobile = $6, net_tel = $7, net_code = $8, weight = $9, def_price = $10, volume = $11, actual_weight = $12, print_task_id = $13, label = $14, pickup_code = $15 WHERE third_order_id = $16`,
		time.Now(),
		payload.Status,
		payload.CancelMsg9,
		payload.CancelMsg99,
		payload.CourierName,
		payload.CourierMobile,
		payload.NetTel,
		payload.NetCode,
		payload.Weight,
		payload.DefPrice,
		payload.Volume,
		payload.ActualWeight,
		payload.PrintTaskId,
		payload.Label,
		payload.PickupCode,
		payload.OrderId,
	); err != nil {
		return fmt.Errorf("更新寄件信息失败")
	}
	return nil
}
