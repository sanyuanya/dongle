package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/rest/mini"
	"github.com/sanyuanya/dongle/rest/pc"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'

	// pc
	app.Post("/api/pc/import", pc.ExcelImport)

	app.Post("/api/pc/login", pc.PcLogin)

	app.Get("/api/pc/userList", pc.UserList)

	app.Post("/api/pc/addWhite", pc.AddWhite)

	app.Get("/api/pc/withdrawalList", pc.WithdrawalList)

	app.Post("/api/pc/approvalWithdrawal", pc.ApprovalWithdrawal)

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
