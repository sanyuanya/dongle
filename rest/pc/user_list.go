package pc

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

	userPageListRequest := &entity.UserPageListRequest{}

	if userPageListRequest.Page, err = strconv.ParseInt(c.Query("page"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if userPageListRequest.PageSize, err = strconv.ParseInt(c.Query("page_size"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	if userPageListRequest.IsWhite, err = strconv.ParseInt(c.Query("is_white"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("is_white 参数错误: %v", err)})
	}

	userPageListRequest.Keyword = c.Query("keyword")

	userList, err := data.GetUserPageList(userPageListRequest)
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取用户列表失败: %v", err)})
	}

	userTotal, err := data.GetUserPageCount(userPageListRequest)
	if err != nil {
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取用户总数失败: %v", err)})
	}

	return c.JSON(tools.Response{
		Code:    0,
		Message: "success",
		Result: map[string]any{
			"data":  userList,
			"total": userTotal,
		},
	})
}
