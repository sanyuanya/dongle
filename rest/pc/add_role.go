package pc

import (
	"fmt"
	"log"
	"regexp"

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

	re := regexp.MustCompile(`\s`)
	if re.MatchString(payload.Name) {
		panic(tools.CustomError{Code: 40000, Message: "角色名称不能包含空格"})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	role, err := data.GetRoleByName(tx, payload.Name)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("查询角色失败: %v", err)})
	}

	if role != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50006, Message: "角色已存在，请勿重复设置"})
	}

	payload.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	err = data.AddRole(tx, payload)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加角色失败: %v", err)})
	}

	for _, permissionId := range payload.PermissionList {

		permission, err := data.GetPermission(tx, permissionId)

		log.Printf("%#+v", permission)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("查询权限失败: %v", err)})
		}

		if permission == nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("权限不存在: %v", permissionId)})
		}

		err = data.AddRolePermission(tx, &entity.AddRolePermissionRequest{
			SnowflakeId:  tools.SnowflakeUseCase.NextVal(),
			RoleId:       payload.SnowflakeId,
			PermissionId: permissionId,
		})

		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加角色权限失败: %v", err)})
		}
	}

	data.Commit(tx)

	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加角色成功",
		Result:  struct{}{},
	})

}
