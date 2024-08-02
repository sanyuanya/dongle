package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func GetWithdrawalList(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "user")
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	getWithdrawalListRequest := new(entity.GetWithdrawalListRequest)

	if getWithdrawalListRequest.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if getWithdrawalListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	getWithdrawalListRequest.Date = c.Query("date", "")

	withdrawalList, err := data.GetWithdrawalListBySnowflakeId(snowflakeId, getWithdrawalListRequest)

	if err != nil {
		panic(fmt.Errorf("获取提现列表失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    20000,
		Message: "获取提现列表成功",
		Result:  withdrawalList,
	})

}
