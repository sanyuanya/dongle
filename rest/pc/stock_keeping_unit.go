package pc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"mime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/minio/minio-go/v7"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/tools"
)

func GetSkuList(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := &entity.GetSkuRequest{}

	if payload.Page, err = strconv.ParseInt(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if payload.PageSize, err = strconv.ParseInt(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	if payload.Status, err = strconv.ParseInt(c.Query("status", "0"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("status 参数错误: %v", err)})
	}

	payload.ItemId = c.Params("itemId", "")

	payload.Name = c.Query("name", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	skuList, err := data.GetSkuList(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50002, Message: fmt.Sprintf("无法获取SKU列表: %v", err)})
	}

	total, err := data.GetSkuCount(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法获取SKU数量: %v", err)})

	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取SKU列表成功",
		Result: map[string]interface{}{
			"sku_list": skuList,
			"total":    total,
		},
	})
}

func AddSku(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := &entity.AddSkuRequest{}

	if err = c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("请求参数错误: %v", err)})
	}

	if payload.Status != 0 && payload.Status > 2 {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("status 参数错误: %v", err)})
	}

	payload.SnowflakeId = tools.SnowflakeUseCase.NextVal()
	payload.ItemId = c.Params("itemId", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	_, err = data.FindByItemId(tx, payload.ItemId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50002, Message: "商品不存在"})
	}

	if _, err = data.FindBySkuCode(tx, payload.Code); err == nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: "SKU已存在"})
	}

	m := tools.Minio{
		Config: &tools.MinioConfig{
			Endpoint:        "218.11.1.36:9000",
			AccessKeyID:     "EvQqTpffmcfUD91VhnHZ",
			SecretAccessKey: "qAjVfSTMGWS57MQs5z9m4j0Xyr4y17U8dsXOmrmr",
			UseSSL:          false,
		},
	}

	if err = m.NewClient(); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Minio 客户端: %v", err)})
	}

	if err = m.MakeBucket(c.Context(), "dongle"); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Minio 存储桶: %v", err)})
	}

	hash := sha256.New()

	payload.BucketName = "dongle"

	imageData := []byte(payload.ImageData)

	hash.Write(imageData)

	sha256Hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	payload.ObjectName = fmt.Sprintf("%s/%s/%s%s", payload.ItemId, payload.Code, sha256Hash, payload.Ext)

	mimeType := mime.TypeByExtension(payload.Ext)
	if mimeType == "" {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v", err)})
	}

	if !strings.HasPrefix(mimeType, "image/") {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v, 必须是图片类型", err)})
	}

	imageBytes, err := base64.StdEncoding.DecodeString(payload.ImageData)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 Base64 编码: %v", err)})
	}

	_, err = m.PutObject(c.Context(), imageBytes, payload.BucketName, payload.ObjectName, minio.PutObjectOptions{
		ContentType:        mimeType,
		ContentDisposition: "inline",
	})

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法上传图片: %v", err)})
	}

	err = data.AddSku(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法添加SKU: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加SKU成功",
		Result:  struct{}{},
	})
}

func UpdateSku(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	payload := &entity.UpdateSkuRequest{}

	if err = c.Bind().Body(payload); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("请求参数错误: %v", err)})
	}

	if payload.Status != 0 && payload.Status > 2 {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("status 参数错误: %v", err)})
	}

	payload.ItemId = c.Params("itemId", "")
	payload.SnowflakeId = c.Params("skuId", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	_, err = data.FindByItemId(tx, payload.ItemId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50002, Message: "商品不存在"})
	}

	sku, err := data.FindBySkuSnowflakeId(tx, payload.SnowflakeId)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: "SKU不存在"})
	}

	if sku.SnowflakeId != payload.SnowflakeId {
		tx.Rollback()
		panic(tools.CustomError{Code: 50004, Message: "SKU已存在"})
	}

	m := tools.Minio{
		Config: &tools.MinioConfig{
			Endpoint:        "218.11.1.36:9000",
			AccessKeyID:     "EvQqTpffmcfUD91VhnHZ",
			SecretAccessKey: "qAjVfSTMGWS57MQs5z9m4j0Xyr4y17U8dsXOmrmr",
			UseSSL:          false,
		},
	}

	if err = m.NewClient(); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Minio 客户端: %v", err)})
	}

	if err = m.MakeBucket(c.Context(), "dongle"); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法创建 Minio 存储桶: %v", err)})
	}

	hash := sha256.New()

	payload.BucketName = "dongle"

	imageData := []byte(payload.ImageData)

	hash.Write(imageData)

	sha256Hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	payload.ObjectName = fmt.Sprintf("%s/%s/%s%s", payload.ItemId, payload.Code, sha256Hash, payload.Ext)

	mimeType := mime.TypeByExtension(payload.Ext)
	if mimeType == "" {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v", err)})
	}

	if !strings.HasPrefix(mimeType, "image/") {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v, 必须是图片类型", err)})
	}

	imageBytes, err := base64.StdEncoding.DecodeString(payload.ImageData)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 Base64 编码: %v", err)})
	}

	_, err = m.PutObject(c.Context(), imageBytes, payload.BucketName, payload.ObjectName, minio.PutObjectOptions{
		ContentType:        mimeType,
		ContentDisposition: "inline",
	})

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法上传图片: %v", err)})
	}

	err = data.UpdateSku(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50003, Message: fmt.Sprintf("无法更新SKU: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新SKU成功",
		Result:  struct{}{},
	})
}

func DeleteSku(c fiber.Ctx) error {
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

	snowflakeId, err := tools.ValidateUserToken(c.Get("Authorization"), "admin")
	_ = snowflakeId
	if err != nil {
		panic(tools.CustomError{Code: 50000, Message: fmt.Sprintf("未经授权: %v", err)})
	}

	itemId := c.Params("itemId", "")
	skuId := c.Params("skuId", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	if err = data.DeleteSku(tx, itemId, skuId); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50002, Message: fmt.Sprintf("无法删除SKU: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "删除SKU成功",
		Result:  struct{}{},
	})
}
