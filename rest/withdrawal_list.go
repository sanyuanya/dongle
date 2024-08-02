package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func WithdrawalList(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(Resp{
				Code:    50000,
				Message: fmt.Sprintf("%v", err),
				Result:  struct{}{},
			})
		}
	}()

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	withdrawalPageListRequest := &entity.WithdrawalPageListRequest{}

	if withdrawalPageListRequest.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if withdrawalPageListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	if withdrawalPageListRequest.LifeCycle, err = strconv.ParseInt(c.Query("life_cycle", "0"), 10, 64); err != nil {
		panic(fmt.Errorf("is_white 参数错误: %v", err))
	}

	withdrawalPageListRequest.Keyword = c.Query("keyword")

	withdrawalList, err := data.WithdrawalPageList(withdrawalPageListRequest)

	if err != nil {
		panic(fmt.Errorf("获取提现列表失败: %v", err))
	}

	total, err := data.WithdrawalListCount(withdrawalPageListRequest)

	if err != nil {
		panic(fmt.Errorf("获取提现总数失败: %v", err))
	}

	return c.JSON(fiber.Map{
		"code":    20000,
		"message": "success",
		"result": map[string]any{
			"data":  withdrawalList,
			"total": total,
		},
	})

}
