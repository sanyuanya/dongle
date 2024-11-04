package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/rest/mini"
	"github.com/sanyuanya/dongle/rest/pc"
)

// var (
// 	httpRequestsTotal = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "HTTP 请求的总数量",
// 		},
// 		[]string{"path", "method", "status"},
// 	)
// )

// func init() {
// 	// 注册自定义指标
// 	prometheus.MustRegister(httpRequestsTotal)
// }

func main() {
	// Initialize a new Fiber app
	app := fiber.New()
	app.Use(requestid.New())

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format:     "${pid} ${time} ${locals:requestid} ${status} - ${method} ${path} ${resBody} ${reqHeaders} ${queryParams} ${bytesSent} ${bytesReceived}​\n",
		TimeZone:   "Asia/Shanghai",
		TimeFormat: time.RFC3339Nano,
	}))

	go data.StartTicker()

	// app.Use(adaptor.HTTPMiddleware(middlewares.RecordLog))

	// Define a route for the GET method on the root path '/'
	// pc

	app.Get("/api/pc/productGroup", pc.GetProductGroup)

	app.Get("/api/pc/upload/:fileName", pc.DownloadFile)

	app.Get("/api/pc/tableMarkUp", pc.TableMarkUp)

	app.Get("/api/pc/productGroupList", pc.GetProductGroupList)

	app.Post("/api/pc/import", pc.ExcelImport)

	app.Post("/api/pc/login", pc.PcLogin)

	app.Get("/api/pc/userList", pc.UserList)

	app.Get("/api/pc/exportUser", pc.ExportUser)

	app.Post("/api/pc/addWhite", pc.SetUpWhite)

	app.Get("/api/pc/withdrawalList", pc.WithdrawalList)

	app.Post("/api/pc/approvalWithdrawal", pc.ApprovalWithdrawal)

	app.Post("/api/pc/updateUserInfo", pc.UpdateUserInfo)

	app.Get("/api/pc/incomeList", pc.GetIncomeList)

	app.Post("/api/pc/income/update", pc.UpdateIncome)

	app.Get("/api/pc/batch/:batchId", pc.OutBatchNo)

	app.Get("/api/pc/batch/:batchId/transfer/:transferId", pc.OutTransferNo)

	app.Get("/api/pc/outTradeNo/:outTradeNo/:orderId", pc.OutTradeNo)

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

	app.Get("/api/pc/excel/download", pc.Download)

	app.Get("/api/pc/operationLog/list", pc.GetOperationLogList)

	app.Get("/api/pc/productCategories/list", pc.GetProductCategoriesList)

	app.Post("/api/pc/productCategories/add", pc.AddProductCategories)

	app.Post("/api/pc/productCategories/update/:productCategoriesId", pc.UpdateProductCategories)

	app.Delete("/api/pc/productCategories/delete/:productCategoriesId", pc.DeleteProductCategories)

	app.Get("/api/pc/item/list", pc.GetItemList)

	app.Post("/api/pc/item/add", pc.AddItem)

	app.Post("/api/pc/item/update/:itemId", pc.UpdateItem)

	app.Patch("/api/pc/item/updateStatus/:itemId", pc.UpdateItemStatus)

	app.Delete("/api/pc/item/delete/:itemId", pc.DeleteItem)

	app.Get("/api/pc/item/show/:itemId", pc.ShowItem)

	//Stock Keeping Unit

	app.Get("/api/pc/sku/list/:itemId", pc.GetSkuList)

	app.Post("/api/pc/sku/add/:itemId", pc.AddSku)

	app.Post("/api/pc/sku/update/:itemId/:skuId", pc.UpdateSku)

	app.Delete("/api/pc/sku/delete/:itemId/:skuId", pc.DeleteSku)

	app.Patch("/api/pc/sku/updateStatus/:itemId/:skuId", pc.UpdateSkuStatus)

	app.Post("/api/pc/pay", pc.Pay)

	//mini program

	app.Post("/api/mini/login", mini.MiniLogin)

	app.Post("/api/mini/updateUserInfo", mini.UpdateUserInfo)

	app.Get("/api/mini/getUserInfo", mini.GetUserInfo)

	app.Post("/api/mini/setUserInfo", mini.SetUserInfo)

	app.Post("/api/mini/applyForWithdrawal", mini.ApplyForWithdrawal)

	app.Get("/api/mini/getWithdrawalList", mini.GetWithdrawalList)

	app.Get("/api/mini/incomeList", mini.GetIncomeList)

	app.Get("/api/mini/productGroupList", mini.GetProductGroupList)

	app.Post("/api/mini/address/add", mini.AddAddress)

	app.Get("/api/mini/address/list", mini.GetAddressList)

	app.Post("/api/mini/address/update/:addressId", mini.UpdateAddress)

	app.Delete("/api/mini/address/delete/:addressId", mini.DeleteAddress)

	app.Post("/api/order/submit", mini.Submit)

	app.Get("/api/order/list", mini.GetOrderList)

	app.Get("/api/order/detail/:orderId", mini.GetOrderDetail)

	app.Post("/api/order/cancel/:orderId", mini.CancelOrder)

	app.Get("/api/cart/index", mini.CartIndex)

	app.Post("/api/cart/add", mini.CartAdd)

	app.Post("/api/cart/update/:cartId", mini.CartUpdate)

	app.Post("/api/cart/delete", mini.CartDelete)

	app.Post("/api/jsApi/notify", mini.Notify)

	app.Post("/api/order/logistics", pc.Logistics)

	app.Post("/api/order/label", pc.Label)

	app.Post("/api/order/pollQuery", pc.PollQuery)

	app.Post("/api/order/shipping", pc.Shipping)

	app.Post("/api/order/orderCallback", pc.OrderCallback)

	app.Post("/api/order/confirm/receipt/good", mini.ConfirmReceiptGood)

	app.Post("/api/order/shipping/cancel:orderId", pc.CancelShipping)

	// Metrics route
	// app.Get("/metrics", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	r.Header.Set("Content-Type", "application/vnd.prometheus.metrics; version=0.0.4; charset=utf-8")
	// 	promhttp.Handler().ServeHTTP(w, r)
	// }))

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
