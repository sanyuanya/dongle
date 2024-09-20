package mini

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func Submit(c fiber.Ctx) error {
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
	userId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}
	payload := new(entity.SubmitOrderRequest)
	if err = c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("参数错误: %v", err)})
	}
	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	address, err := data.FindByAddressSnowflakeId(tx, payload.AddressId, userId)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取地址失败: %v", err)})
	}

	user, err := data.GetUserDetailBySnowflakeID(tx, userId)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取用户失败: %v", err)})
	}

	addOrder := &entity.AddOrder{
		SnowflakeId:     tools.SnowflakeUseCase.NextVal(),
		AddressId:       payload.AddressId,
		Consignee:       address.Consignee,
		PhoneNumber:     address.PhoneNumber,
		Location:        address.Location,
		DetailedAddress: address.DetailedAddress,
		UserId:          userId,
		ExpirationTime:  time.Now().Add(5 * time.Minute).Unix(),
		OutTradeNo:      tools.SnowflakeUseCase.NextVal(),
		OrderState:      1,
		Currency:        "CNY",
		OpenId:          user.OpenID,
	}

	rdb := tools.Redis{}
	if err = rdb.NewClient(); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Redis 客户端: %v", err)})
	}
	defer rdb.Client.Close()

	addOrderCommodityList := []*entity.AddOrderCommodity{}

	for _, commodity := range payload.OrderCommodity {
		if commodity.Quantity <= 0 {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("商品数量不能小于 1")})
		}
		if commodity.SkuId == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("商品 skuId 不能为空")})
		}
		if commodity.CommodityId == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("商品 commodityId 不能为空")})
		}

		if result, err := rdb.DeductStock(commodity.SkuId, commodity.Quantity); err != nil || !result {
			tx.Rollback()
			panic(tools.CustomError{Code: 50002, Message: fmt.Sprintf("库存不足,下单失败")})
		}

		if err = data.UpdateSkuStockQuantity(tx, commodity.CommodityId, commodity.SkuId, commodity.Quantity); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新库存失败: %v", err)})
		}

		comm, err := data.FindByItemId(tx, commodity.CommodityId)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取商品失败: %v", err)})
		}

		sku, err := data.FindBySkuSnowflakeId(tx, commodity.SkuId)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("获取商品失败: %v", err)})
		}

		addOrder.Total += (sku.Price * 100) * float64(commodity.Quantity)

		addOrderCommodity := &entity.AddOrderCommodity{
			SnowflakeId:          tools.SnowflakeUseCase.NextVal(),
			CommodityId:          comm.SnowflakeId,
			CommodityName:        comm.Name,
			CommodityCode:        comm.Code,
			CommodityDescription: comm.Description,
			CategoriesId:         comm.CategoriesId,
			SkuId:                sku.SnowflakeId,
			SkuCode:              sku.Code,
			SkuName:              sku.Name,
			Price:                sku.Price,
			Quantity:             commodity.Quantity,
			ObjectName:           sku.ObjectName,
			BucketName:           sku.BucketName,
			OrderId:              addOrder.SnowflakeId,
		}

		addOrderCommodityList = append(addOrderCommodityList, addOrderCommodity)
	}

	if err = data.AddOrder(tx, addOrder); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加订单失败: %v", err)})
	}

	for _, addOrderCommodity := range addOrderCommodityList {
		if err = data.AddOrderCommodity(tx, addOrderCommodity); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加订单商品失败: %v", err)})
		}
	}

	// jsApiRequest := &pay.JsApiRequest{
	// 	AppId:       "wx370126c8bcf8d00c",
	// 	Mchid:       "1682195529",
	// 	Description: "购买中心",
	// 	OutTradeNo:  addOrder.OutTradeNo,
	// 	Attach:      "",
	// 	Amount: pay.Amount{
	// 		Total:    addOrder.Total,
	// 		Currency: addOrder.Currency,
	// 	},
	// 	Payer: pay.Payer{
	// 		OpenId: user.OpenId,
	// 	},
	// 	Detail: pay.Detail{
	// 		GoodDetail: []*pay.GoodDetail{},
	// 	},
	// 	NotifyUrl: "https://www.weixin.qq.com/wxpay/pay.php",
	// }

	// log.Printf("jsApiRequest: %#+v", jsApiRequest)
	// jsApiResponse, err := pay.JsApi(jsApiRequest)

	// if err != nil {

	// }
	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "创建订单成功",
		Result:  struct{}{},
	})
}