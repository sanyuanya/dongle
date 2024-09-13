package data

import (
	"database/sql"

	"github.com/sanyuanya/dongle/entity"
)

func GetPermissionList(tx *sql.Tx) ([]*entity.PermissionMenu, error) {
	baseSQL := `
		SELECT
			snowflake_id,
			summary,
			path,
			parent_id,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') created_at,
			TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS') updated_at
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

	permissions := make([]*entity.PermissionMenu, 0)
	for rows.Next() {

		permission := &entity.PermissionMenu{}
		err := rows.Scan(
			&permission.SnowflakeId,
			&permission.Summary,
			&permission.Path,
			&permission.ParentId,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		permission.Children = make([]*entity.PermissionMenu, 0)
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
		&permission.SnowflakeId,
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
			api_path,
			parent_id
		FROM
			permissions
		WHERE
			snowflake_id = $1 AND deleted_at IS NULL
	`

	permissionMenu := &entity.PermissionMenu{}

	err := tx.QueryRow(baseSQL, permissionId).Scan(
		&permissionMenu.SnowflakeId,
		&permissionMenu.Summary,
		&permissionMenu.Path,
		&permissionMenu.ApiPath,
		&permissionMenu.ParentId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	permissionMenu.Children = make([]*entity.PermissionMenu, 0)

	return permissionMenu, nil
}

func GetPermissionListByRoleId(tx *sql.Tx, roleId string) ([]*entity.Permission, error) {
	baseSQL := `
		SELECT
			p.snowflake_id,
			p.summary,
			p.path,
			p.created_at,
			p.updated_at
		FROM
			role_permission rp
		JOIN
			permissions p
		ON
			rp.permission_id = p.snowflake_id AND p.deleted_at IS NULL
		WHERE
			rp.role_id = $1
			AND rp.deleted_at IS NULL
	`

	rows, err := tx.Query(baseSQL, roleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]*entity.Permission, 0)
	for rows.Next() {

		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.SnowflakeId,
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
