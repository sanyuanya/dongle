package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/tools"
)

func GetRoleList(c fiber.Ctx) error {
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

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	roleList, err := data.GetRoleList(tx)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取角色列表失败: %v", err)})
	}

	for _, role := range roleList {
		permissionList, err := data.GetPermissionListByRoleId(tx, role.SnowflakeID)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("获取角色权限失败: %v", err)})
		}
		role.PermissionList = permissionList
	}

	data.Commit(tx)

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取角色列表成功",
		Result: map[string]any{
			"role_list": roleList,
		},
	})

}
