package data

import (
	"database/sql"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddItemImage(tx *sql.Tx, addItemImage *entity.AddItemImage) error {

	_, err := tx.Exec("INSERT INTO item_images (snowflake_id, item_id, type, data, sorting, ext, object_name, bucket_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		addItemImage.SnowflakeId,
		addItemImage.ItemId,
		addItemImage.Type,
		addItemImage.Data,
		addItemImage.Sorting,
		addItemImage.Ext,
		addItemImage.ObjectName,
		addItemImage.BucketName)
	if err != nil {
		return err
	}

	return nil
}

func DeleteItemImage(tx *sql.Tx, item_id string) error {

	_, err := tx.Exec("UPDATE item_images SET deleted_at = $1 WHERE item_id = $2 AND deleted_at IS NULL", time.Now(), item_id)
	if err != nil {
		return err
	}

	return nil
}

func GetItemImageList(tx *sql.Tx, item_id string, t int64) ([]*entity.Picture, error) {
	rows, err := tx.Query("SELECT snowflake_id, item_id, type, sorting, ext, object_name, bucket_name, data FROM item_images WHERE item_id = $1 AND deleted_at IS NULL AND type = $2 ORDER BY sorting ASC, created_at ASC", item_id, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemImages := make([]*entity.Picture, 0)

	for rows.Next() {
		itemImage := &entity.Picture{}
		err := rows.Scan(&itemImage.SnowflakeId, &itemImage.ItemId, &itemImage.Type, &itemImage.Sorting, &itemImage.Ext, &itemImage.ObjectName, &itemImage.BucketName, &itemImage.ImageData)
		if err != nil {
			return nil, err
		}
		itemImages = append(itemImages, itemImage)
	}

	return itemImages, nil
}
