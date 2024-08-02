package rest

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func UserList(c fiber.Ctx) error {

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

	userPageListRequest := &entity.UserPageListRequest{}

	if userPageListRequest.Page, err = strconv.ParseInt(c.Query("page"), 10, 64); err != nil {
		panic(fmt.Errorf("page 参数错误: %v", err))
	}

	if userPageListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size"), 10, 64); err != nil {
		panic(fmt.Errorf("page_size 参数错误: %v", err))
	}

	if userPageListRequest.IsWhite, err = strconv.ParseInt(c.Query("is_white"), 10, 64); err != nil {
		panic(fmt.Errorf("is_white 参数错误: %v", err))
	}

	userPageListRequest.Keyword = c.Query("keyword")

	userList, err := data.GetUserPageList(userPageListRequest)
	if err != nil {
		panic(fmt.Errorf("获取用户列表失败: %v", err))
	}

	userTotal, err := data.GetUserPageCount(userPageListRequest)
	if err != nil {
		panic(fmt.Errorf("获取用户总数失败: %v", err))
	}

	return c.JSON(Resp{
		Code:    0,
		Message: "success",
		Result: map[string]any{
			"data":  userList,
			"total": userTotal,
		},
	})
}
