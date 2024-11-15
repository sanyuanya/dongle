package pc

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func UpdateAdmin(c fiber.Ctx) error {

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

	payload := &entity.UpdateAdminRequest{}

	payload.SnowflakeId = c.Params("adminId", "")

	err = c.Bind().Body(payload)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	re := regexp.MustCompile(`\s`)
	if re.MatchString(payload.Account) {
		panic(tools.CustomError{Code: 40000, Message: "用户名不能包含空格"})
	}

	for _, char := range payload.Account {
		if unicode.Is(unicode.Han, char) {
			panic(tools.CustomError{Code: 40000, Message: "用户名不能包含中文"})
		}
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	// 查询用户是否存在
	err = data.FindBySnowflakeIdNotFoundAndAccount(tx, payload.SnowflakeId, payload.Account)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("用户已存在: %v", err)})
	}

	err = data.UpdateAdmin(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("更新管理员失败: %v", err)})
	}

	// 删除原有角色
	err = data.DeleteAdminRole(tx, payload.SnowflakeId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50008, Message: fmt.Sprintf("删除原有角色失败: %v", err)})
	}

	// 添加新角色
	for _, roleId := range payload.RoleList {

		// 查询角色是否存在
		role, err := data.GetRole(tx, roleId)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50008, Message: fmt.Sprintf("查询角色失败: %v", err)})
		}

		if role == nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50008, Message: fmt.Sprintf("角色不存在: %v", roleId)})
		}

		addAdminRoleRequest := &entity.AddAdminRoleRequest{
			SnowflakeId: tools.SnowflakeUseCase.NextVal(),
			AdminId:     payload.SnowflakeId,
			RoleId:      roleId,
		}

		err = data.AddAdminRole(tx, addAdminRoleRequest)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50009, Message: fmt.Sprintf("添加用户角色失败: %v", err)})
		}
	}

	tx.Commit()
	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新用户成功",
		Result:  struct{}{},
	})

}
