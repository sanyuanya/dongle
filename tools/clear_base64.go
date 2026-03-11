package tools

import "github.com/sanyuanya/dongle/entity"

// ClearPictureBase64 清空 Picture 结构体中的 base64 图片数据
func ClearPictureBase64(pictures []*entity.Picture) {
	for _, pic := range pictures {
		if pic != nil {
			pic.ImageData = ""
		}
	}
}

// ClearSkuBase64 清空 Sku 结构体中的 base64 图片数据
func ClearSkuBase64(skus []*entity.Sku) {
	for _, sku := range skus {
		if sku != nil {
			sku.ImageData = ""
		}
	}
}

// ClearItemBase64 清空 Item 结构体及其关联数据中的 base64 图片数据
func ClearItemBase64(item *entity.Item) {
	if item == nil {
		return
	}
	ClearPictureBase64(item.Picture)
	ClearPictureBase64(item.Detail)
	ClearSkuBase64(item.Sku)
}

// ClearItemListBase64 清空 Item 列表中的 base64 图片数据
func ClearItemListBase64(items []*entity.Item) {
	for _, item := range items {
		ClearItemBase64(item)
	}
}
