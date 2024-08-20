package pc

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func AddProduct(c fiber.Ctx) error {
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

	_, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := &entity.AddProductRequest{}
	err = c.Bind().Body(payload)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	if payload.Name == "" {
		panic(tools.CustomError{Code: 40000, Message: "商品名称不能为空"})
	}

	if payload.Integral <= 0 {
		panic(tools.CustomError{Code: 40000, Message: "产品积分不能小于0"})
	}

	re := regexp.MustCompile(`\s`)
	if re.MatchString(payload.Name) {
		panic(tools.CustomError{Code: 40000, Message: "产品名称不能包含空白字符"})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 判断产品名称是否已经存在
	product, err := data.FindProductByName(tx, payload.Name)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法获取商品: %v", err)})
	}

	if product != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("商品名称已经存在: %v", err)})
	}

	payload.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	err = data.AddProduct(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法添加商品: %v", err)})
	}

	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加商品成功",
		Result:  struct{}{},
	})

}
