package pc

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func AddAdmin(c fiber.Ctx) error {

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

	payload := &entity.AddAdminRequest{}

	err = c.Bind().Body(payload)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	re := regexp.MustCompile(`\s`)
	if re.MatchString(payload.Account) {
		panic(tools.CustomError{Code: 40000, Message: "用户名不能包含空格"})
	}

	if len(payload.Password) < 6 {
		panic(tools.CustomError{Code: 40000, Message: "密码长度不能小于 6 位"})
	}

	if re.MatchString(payload.Password) {
		panic(tools.CustomError{Code: 40000, Message: "密码不能包含空格"})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 检查 account 是否已经存在了
	existed, err := data.GetAdminByAccount(tx, payload.Account)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("检查账号是否存在：%v", err)})
	}

	if existed != "" {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 40000, Message: "账号已存在，请勿重复设置"})
	}

	payload.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	err = data.AddAdmin(tx, payload)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加管理员失败: %v", err)})
	}

	for _, roleId := range payload.RoleList {

		// 查询角色是否存在
		role, err := data.GetRole(tx, roleId)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("查询角色失败: %v", err)})
		}

		if role == nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("角色不存在: %v", roleId)})
		}

		addAdminRoleRequest := &entity.AddAdminRoleRequest{
			SnowflakeId: tools.SnowflakeUseCase.NextVal(),
			AdminId:     payload.SnowflakeId,
			RoleId:      roleId,
		}

		err = data.AddAdminRole(tx, addAdminRoleRequest)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("添加用户角色失败: %v", err)})
		}
	}

	data.Commit(tx)

	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加用户成功",
		Result:  struct{}{},
	})

}
