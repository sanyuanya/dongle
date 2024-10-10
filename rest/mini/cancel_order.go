package mini

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay"
	"github.com/sanyuanya/dongle/tools"
)

func CancelOrder(c fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {

			var code int
			var message string

			switch e := err.(type) {
			case tools.CustomError:
				code = e.Code
				message = e.Message
			case error:
				code = 50001
				message = e.Error()
			default:
				code = 50002
				message = fmt.Sprintf("%v", e)
			}

			c.JSON(tools.Response{
				Code:    code,
				Message: message,
				Result:  struct{}{},
			})
		}
	}()

	outTradeNo := c.Params("orderId")
	err := pay.CloseOrder(outTradeNo)
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("关闭订单失败: %v", err)})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	orderListSearch := &entity.GetOrderListRequest{
		OutTradeNo: outTradeNo,
	}

	orderList, err := data.GetOrderList(tx, orderListSearch)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("查询订单失败: %v", err)})
	}

	rdb := tools.Redis{}
	if err = rdb.NewClient(); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Redis 客户端: %v", err)})
	}
	defer rdb.Client.Close()

	for _, order := range orderList {
		commodity, err := data.GetOrderCommodityList(tx, order.SnowflakeId)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("查询订单详情失败: %v", err)})
		}

		for _, comm := range commodity {
			err := data.CancelOrderReturnStock(tx, comm.CommodityId, comm.SkuId, int64(comm.Quantity))
			if err != nil {
				tx.Rollback()
				panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("返回取消订单所购商品的库存: %v", err)})
			}

			if err = rdb.UpdateSkuStock(comm.SkuId, int64(comm.Quantity)); err != nil {
				tx.Rollback()
				log.Printf("redis 更新库存失败: %#+v", err)
				panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("返回取消订单所购商品的库存, redis 更新库存失败: %v", err)})
			}
		}
	}

	updateOrderByOutTradeNo := &entity.UpdateOrderByOutTradeNo{
		Status:     99,
		OutTradeNo: outTradeNo,
	}

	err = data.UpdateOrderByOutTradeNo(tx, updateOrderByOutTradeNo)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新订单状态失败: %v", err)})
	}

	err = tx.Commit()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("数据库发生错误，请联系管理员: %v", err)})
	}
	return c.JSON(tools.Response{
		Code:    0,
		Message: "关闭订单成功",
		Result:  struct{}{},
	})

}
