package data

import (
	"database/sql"

	"github.com/sanyuanya/dongle/entity"
)

func GetPermissionList(tx *sql.Tx) ([]*entity.Permission, error) {
	baseSQL := `
		SELECT
			snowflake_id,
			summary,
			path,
			created_at,
			updated_at
		FROM
			permissions
		WHERE
			deleted_at IS NULL
	`

	rows, err := tx.Query(baseSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]*entity.Permission, 0)
	for rows.Next() {

		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.SnowflakeID,
			&permission.Summary,
			&permission.Path,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}
