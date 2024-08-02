package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/rest"
)

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'

	// pc
	app.Post("/api/pc/import", rest.ExcelImport)

	app.Post("/api/pc/login", rest.PcLogin)

	app.Get("/api/pc/userList", rest.UserList)

	app.Post("/api/pc/addWhite", rest.AddWhite)

	app.Get("/api/pc/withdrawalList", rest.WithdrawalList)

	app.Post("/api/pc/approvalWithdrawal", rest.ApprovalWithdrawal)

	//mini program

	app.Post("/api/mini/login", rest.MiniLogin)

	app.Post("/api/mini/updateUserInfo", rest.UpdateUserInfo)

	app.Get("/api/mini/getUserInfo", rest.GetUserInfo)

	app.Post("/api/mini/setUserInfo", rest.SetUserInfo)

	app.Post("/api/mini/applyForWithdrawal", rest.ApplyForWithdrawal)

	app.Get("/api/mini/getWithdrawalList", rest.GetWithdrawalList)

	app.Get("/api/mini/incomeList", rest.GetIncomeList)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
