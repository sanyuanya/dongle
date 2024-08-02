package mini

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func GetIncomeList(c fiber.Ctx) error {

	defer func() {
		if err := recover(); err != nil {
			c.JSON(tools.Response{
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

	payload := new(entity.GetIncomeListRequest)

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	payload.Date = c.Query("date", "")

	incomeList, err := data.GetIncomeListBySnowflakeId(snowflakeId, payload)
	if err != nil {
		panic(fmt.Errorf("获取收支列表失败: %v", err))
	}

	total, err := data.GetIncomeCountBySnowflakeId(snowflakeId, payload)
	if err != nil {
		panic(fmt.Errorf("获取收支列表数量失败: %v", err))
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取收支列表成功",
		Result: map[string]any{
			"data":  incomeList,
			"total": total,
		},
	})
}
