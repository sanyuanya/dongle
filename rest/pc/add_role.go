package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func AddRole(c fiber.Ctx) error {
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

	payload := &entity.AddRoleRequest{}

	err = c.Bind().Body(payload)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	payload.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	err = data.AddRole(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加角色失败: %v", err)})
	}

	for _, permissionId := range payload.PermissionList {
		err = data.AddRolePermission(tx, &entity.AddRolePermissionRequest{
			SnowflakeId:  tools.SnowflakeUseCase.NextVal(),
			RoleId:       payload.SnowflakeId,
			PermissionId: permissionId,
		})
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加角色权限失败: %v", err)})
		}
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加角色成功",
		Result:  struct{}{},
	})

}
