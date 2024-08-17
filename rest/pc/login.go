package pc

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

// mysecretpassword
func PcLogin(c fiber.Ctx) error {

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

	loginRequest := new(entity.LoginRequest)

	err := c.Bind().Body(loginRequest)
	if err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50006, Message: fmt.Sprintf("开始事务失败: %v", err)})
	}

	snowflakeId, err := data.Login(tx, loginRequest)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("登录失败: %v", err)})
	}

	token, err := tools.GenerateToken(snowflakeId, "admin")

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50004, Message: fmt.Sprintf("生成token失败: %v", err)})
	}

	err = data.SetApiToken(tx, snowflakeId, token)

	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50005, Message: fmt.Sprintf("设置token失败 : %v", err)})
	}

	adminRole, err := data.GetAdminRoleList(tx, snowflakeId)
	if err != nil {
		data.Rollback(tx)
		panic(tools.CustomError{Code: 50007, Message: fmt.Sprintf("获取管理员角色失败: %v", err)})
	}

	permissionList := make([]string, 0)
	// 循环角色查询权限
	for _, role := range adminRole {
		rolePermission, err := data.GetRolePermissionList(tx, role.RoleId)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50008, Message: fmt.Sprintf("获取角色权限失败: %v", err)})
		}
		permissionList = append(permissionList, rolePermission...)
	}

	// 循环权限查询菜单
	menuList := make([]*entity.PermissionMenu, 0)
	for _, permission := range permissionList {
		menu, err := data.GetPermissionMenu(tx, permission)
		if err != nil {
			data.Rollback(tx)
			panic(tools.CustomError{Code: 50009, Message: fmt.Sprintf("获取权限菜单失败: %v", err)})
		}
		menuList = append(menuList, menu)
	}

	data.Commit(tx)
	c.Response().Header.Set("Authorization", token)
	return c.JSON(tools.Response{
		Code:    0,
		Message: "登录成功",
		Result: map[string]any{
			"menu_list": menuList,
		},
	})
}
