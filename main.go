package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/sanyuanya/dongle/middlewares"
	"github.com/sanyuanya/dongle/rest/mini"
	"github.com/sanyuanya/dongle/rest/pc"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()
	app.Use(recover.New())

	app.Use(adaptor.HTTPMiddleware(middlewares.RecordLog))
	// Define a route for the GET method on the root path '/'
	// pc
	app.Post("/api/pc/import", pc.ExcelImport)

	app.Post("/api/pc/login", pc.PcLogin)

	app.Get("/api/pc/userList", pc.UserList)

	app.Post("/api/pc/addWhite", pc.SetUpWhite)

	app.Get("/api/pc/withdrawalList", pc.WithdrawalList)

	app.Post("/api/pc/approvalWithdrawal", pc.ApprovalWithdrawal)

	app.Post("/api/pc/updateUserInfo", pc.UpdateUserInfo)

	app.Get("/api/pc/incomeList", pc.GetIncomeList)

	app.Get("/api/pc/batch/:batchId", pc.OutBatchNo)

	app.Get("/api/pc/batch/:batchId/transfer/:transferId", pc.OutTransferNo)

	app.Get("/api/pc/permission/list", pc.GetPermissionList)

	app.Get("/api/pc/role/list", pc.GetRoleList)

	app.Post("/api/pc/role/add", pc.AddRole)

	app.Post("/api/pc/role/update/:roleId", pc.UpdateRole)

	app.Delete("/api/pc/role/delete/:roleId", pc.DeleteRole)

	app.Get("/api/pc/admin/list", pc.GetAdminList)

	app.Post("/api/pc/admin/add", pc.AddAdmin)

	app.Post("/api/pc/admin/update/:adminId", pc.UpdateAdmin)

	app.Delete("/api/pc/admin/delete/:adminId", pc.DeleteAdmin)

	app.Get("/api/pc/product/list", pc.GetProductList)

	app.Post("/api/pc/product/add", pc.AddProduct)

	app.Post("/api/pc/product/update/:productId", pc.UpdateProduct)

	app.Delete("/api/pc/product/delete/:productId", pc.DeleteProduct)

	app.Get("/api/pc/download", pc.Download)

	app.Post("/api/pc/pay", pc.Pay)

	//mini program

	app.Post("/api/mini/login", mini.MiniLogin)

	app.Post("/api/mini/updateUserInfo", mini.UpdateUserInfo)

	app.Get("/api/mini/getUserInfo", mini.GetUserInfo)

	app.Post("/api/mini/setUserInfo", mini.SetUserInfo)

	app.Post("/api/mini/applyForWithdrawal", mini.ApplyForWithdrawal)

	app.Get("/api/mini/getWithdrawalList", mini.GetWithdrawalList)

	app.Get("/api/mini/incomeList", mini.GetIncomeList)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
