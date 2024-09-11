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

func GetItemList(c fiber.Ctx) error {
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

	payload := &entity.ItemPage{}

	if payload.Page, err = strconv.ParseUint(c.Query("page", "1"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page 参数错误: %v", err)})
	}

	if payload.PageSize, err = strconv.ParseUint(c.Query("page_size", "10"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("page_size 参数错误: %v", err)})
	}

	if payload.Status, err = strconv.ParseUint(c.Query("status", "0"), 10, 64); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("status 参数错误: %v", err)})
	}

	payload.Name = c.Query("name", "")

	payload.CategoriesId = c.Query("categories_id", "")

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	itemList, err := data.ItemList(tx, payload)
	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取商品列表: %v", err)})
	}

	for _, item := range itemList {
		item.Picture, err = data.GetItemImageList(tx, item.SnowflakeId, 1)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取商品图片: %v", err)})
		}

		item.Detail, err = data.GetItemImageList(tx, item.SnowflakeId, 2)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取商品详情: %v", err)})
		}
	}

	itemTotal, err := data.ItemListCount(tx, payload)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取商品总数: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取商品列表成功",
		Result: map[string]any{
			"item_list": itemList,
			"total":     itemTotal,
		},
	})
}

func AddItem(c fiber.Ctx) error {
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

	addItem := &entity.AddItem{}

	addItem.SnowflakeId = tools.SnowflakeUseCase.NextVal()

	if err := c.Bind().Body(addItem); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	if addItem.Status > 2 {
		panic(tools.CustomError{Code: 50001, Message: "无法绑定请求体, status 有效范围 0-2"})
	}

	tx, err := data.Transaction()
	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	if _, err = data.FindByProductCategoriesId(tx, addItem.CategoriesId); err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法找到商品分类: %v", addItem.CategoriesId)})
	}

	if _, err = data.FindByItemCode(tx, addItem.Code); err == nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("商品编码已存在: %v", addItem.Code)})
	}

	if err = data.AddItem(tx, addItem); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法添加商品: %v", err)})
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

	bucketName := "dongle"

	for _, picture := range addItem.Picture {

		imageData := []byte(picture.ImageData)

		hash.Write(imageData)

		sha256Hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		objectName := fmt.Sprintf("%s/%s%s", addItem.Code, sha256Hash, picture.Ext)

		mimeType := mime.TypeByExtension(picture.Ext)
		if mimeType == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v", err)})
		}

		if !strings.HasPrefix(mimeType, "image/") {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v, 必须是图片类型", err)})
		}

		imageBytes, err := base64.StdEncoding.DecodeString(picture.ImageData)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 Base64 编码: %v", err)})
		}

		_, err = m.PutObject(c.Context(), imageBytes, bucketName, objectName, minio.PutObjectOptions{
			ContentType:        mimeType,
			ContentDisposition: "inline",
		})

		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法上传图片: %v", err)})
		}

		addItemImage := &entity.AddItemImage{
			SnowflakeId: tools.SnowflakeUseCase.NextVal(),
			ItemId:      addItem.SnowflakeId,
			Type:        1,
			Data:        picture.ImageData,
			Sorting:     picture.Sorting,
			ObjectName:  objectName,
			BucketName:  bucketName,
			Ext:         picture.Ext,
		}

		if err = data.AddItemImage(tx, addItemImage); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法添加商品图片: %v", err)})
		}
	}

	for _, detail := range addItem.Detail {

		imageData := []byte(detail.ImageData)

		hash.Write(imageData)

		sha256Hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		objectName := fmt.Sprintf("%s/%s%s", addItem.Code, sha256Hash, detail.Ext)

		mimeType := mime.TypeByExtension(detail.Ext)
		if mimeType == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v", err)})
		}

		if !strings.HasPrefix(mimeType, "image/") {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v, 必须是图片类型", err)})
		}

		imageBytes, err := base64.StdEncoding.DecodeString(detail.ImageData)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 Base64 编码: %v", err)})
		}

		_, err = m.PutObject(c.Context(), imageBytes, bucketName, objectName, minio.PutObjectOptions{
			ContentType:        mimeType,
			ContentDisposition: "inline",
		})

		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法上传图片: %v", err)})
		}

		addItemImage := &entity.AddItemImage{
			SnowflakeId: tools.SnowflakeUseCase.NextVal(),
			ItemId:      addItem.SnowflakeId,
			Type:        2,
			Data:        detail.ImageData,
			Sorting:     detail.Sorting,
			ObjectName:  objectName,
			BucketName:  bucketName,
			Ext:         detail.Ext,
		}

		if err = data.AddItemImage(tx, addItemImage); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法添加商品详情: %v", err)})
		}
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "添加商品成功",
		Result:  struct{}{},
	})
}

func UpdateItem(c fiber.Ctx) error {
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

	updateItem := &entity.UpdateItem{}

	updateItem.SnowflakeId = c.Params("itemId", "")

	if err = c.Bind().Body(updateItem); err != nil {
		panic(tools.CustomError{Code: 40000, Message: fmt.Sprintf("无法绑定请求体: %v", err)})
	}

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	itemId, err := data.FindByItemCode(tx, updateItem.Code)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法找到商品: %v", err)})
	}

	if itemId != updateItem.SnowflakeId {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("商品编码已存在: %v", updateItem.Code)})
	}

	if err = data.UpdateItem(tx, updateItem); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法更新商品: %v", err)})
	}

	if err = data.DeleteItemImage(tx, updateItem.SnowflakeId); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法删除商品图片: %v", err)})
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

	bucketName := "dongle"

	for _, picture := range updateItem.Picture {

		imageData := []byte(picture.ImageData)

		hash.Write(imageData)

		sha256Hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		objectName := fmt.Sprintf("%s/%s%s", updateItem.Code, sha256Hash, picture.Ext)

		mimeType := mime.TypeByExtension(picture.Ext)
		if mimeType == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v", err)})
		}

		if !strings.HasPrefix(mimeType, "image/") {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v, 必须是图片类型", err)})
		}

		imageBytes, err := base64.StdEncoding.DecodeString(picture.ImageData)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 Base64 编码: %v", err)})
		}

		_, err = m.PutObject(c.Context(), imageBytes, bucketName, objectName, minio.PutObjectOptions{
			ContentType:        mimeType,
			ContentDisposition: "inline",
		})

		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法上传图片: %v", err)})
		}

		addItemImage := &entity.AddItemImage{
			SnowflakeId: tools.SnowflakeUseCase.NextVal(),
			ItemId:      updateItem.SnowflakeId,
			Type:        1,
			Data:        picture.ImageData,
			Sorting:     picture.Sorting,
			ObjectName:  objectName,
			BucketName:  bucketName,
			Ext:         picture.Ext,
		}

		if err = data.AddItemImage(tx, addItemImage); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法添加商品图片: %v", err)})
		}
	}

	for _, detail := range updateItem.Detail {

		imageData := []byte(detail.ImageData)

		hash.Write(imageData)

		sha256Hash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		objectName := fmt.Sprintf("%s/%s%s", updateItem.Code, sha256Hash, detail.Ext)

		mimeType := mime.TypeByExtension(detail.Ext)
		if mimeType == "" {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v", err)})
		}

		if !strings.HasPrefix(mimeType, "image/") {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 MIME 类型: %v, 必须是图片类型", err)})
		}

		imageBytes, err := base64.StdEncoding.DecodeString(detail.ImageData)
		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无效的 Base64 编码: %v", err)})
		}

		_, err = m.PutObject(c.Context(), imageBytes, bucketName, objectName, minio.PutObjectOptions{
			ContentType:        mimeType,
			ContentDisposition: "inline",
		})

		if err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法上传图片: %v", err)})
		}

		addItemImage := &entity.AddItemImage{
			SnowflakeId: tools.SnowflakeUseCase.NextVal(),
			ItemId:      updateItem.SnowflakeId,
			Type:        2,
			Data:        detail.ImageData,
			Sorting:     detail.Sorting,
			ObjectName:  objectName,
			BucketName:  bucketName,
			Ext:         detail.Ext,
		}

		if err = data.AddItemImage(tx, addItemImage); err != nil {
			tx.Rollback()
			panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法添加商品详情: %v", err)})
		}
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "更新商品成功",
		Result:  struct{}{},
	})
}

func DeleteItem(c fiber.Ctx) error {
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

	snowflakeId := c.Params("itemId", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	if err = data.DeleteItem(tx, snowflakeId); err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法删除商品: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "删除商品成功",
		Result:  struct{}{},
	})
}

func ShowItem(c fiber.Ctx) error {
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

	snowflakeId := c.Params("itemId", "")

	tx, err := data.Transaction()

	if err != nil {
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法开启事务: %v", err)})
	}

	item, err := data.FindByItemId(tx, snowflakeId)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: "无法找到商品"})
	}

	item.Picture, err = data.GetItemImageList(tx, snowflakeId, 1)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取商品图片: %v", err)})
	}

	item.Detail, err = data.GetItemImageList(tx, snowflakeId, 2)

	if err != nil {
		tx.Rollback()
		panic(tools.CustomError{Code: 50001, Message: fmt.Sprintf("无法获取商品详情: %v", err)})
	}

	tx.Commit()

	return c.JSON(tools.Response{
		Code:    0,
		Message: "获取商品成功",
		Result:  item,
	})
}
