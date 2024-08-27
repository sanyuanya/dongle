package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func UpdateRole(c fiber.Ctx) error {
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

	payload := &entity.UpdateRoleRequest{}

	err = c.Bind().Body(payload)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}
	payload.SnowflakeId = c.Params("roleId", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	err = data.FindBySnowflakeIdNotFoundAndRoleName(tx, payload.SnowflakeId, payload.Name)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("角色已存在: %v", err)})
	}

	err = data.UpdateRole(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("更新角色失败: %v", err)})
	}

	// 同步角色权限 把原来的权限删除，再添加新的权限
	err = data.DeleteRolePermission(tx, payload.SnowflakeId)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("删除角色权限失败: %v", err)})
	}

	for _, permissionId := range payload.PermissionList {

		permission, err := data.GetPermission(tx, permissionId)

		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("查询权限失败: %v", err)})
		}

		if permission == nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("权限不存在: %v", permissionId)})
		}

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
		Message: "更新角色成功",
		Result:  struct{}{},
	})

}
