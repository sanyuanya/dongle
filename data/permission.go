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

func GetPermission(tx *sql.Tx, snowflakeId string) (*entity.Permission, error) {

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
			deleted_at IS NULL AND snowflake_id = $1
	`

	permission := &entity.Permission{}

	err := tx.QueryRow(baseSQL, snowflakeId).Scan(
		&permission.SnowflakeID,
		&permission.Summary,
		&permission.Path,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return permission, nil
}

func GetPermissionMenu(tx *sql.Tx, permissionId string) (*entity.PermissionMenu, error) {
	baseSQL := `
		SELECT
			snowflake_id,
			summary,
			path,
			api_path
		FROM
			permissions
		WHERE
			snowflake_id = $1
	`

	permissionMenu := &entity.PermissionMenu{}

	err := tx.QueryRow(baseSQL, permissionId).Scan(
		&permissionMenu.SnowflakeID,
		&permissionMenu.Summary,
		&permissionMenu.Path,
		&permissionMenu.ApiPath,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return permissionMenu, nil
}
