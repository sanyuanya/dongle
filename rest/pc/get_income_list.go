package pc

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

	_, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	if err != nil {
		panic(fmt.Errorf("未经授权: %v", err))
	}

	payload := &entity.IncomePageListExpenseRequest{}

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	payload.Date = c.Query("date")
	payload.Keyword = c.Query("keyword")

	incomeList, err := data.IncomePageList(payload)

	if err != nil {
		panic(fmt.Errorf("获取收入列表失败: %v", err))
	}

	total, err := data.IncomeListCount(payload)

	if err != nil {
		panic(fmt.Errorf("获取收入总数失败: %v", err))
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result: map[string]interface{}{
			"data":  incomeList,
			"total": total,
		},
	})

}
