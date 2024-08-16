package data

import (
	"database/sql"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddRolePermission(tx *sql.Tx, payload *entity.AddRolePermissionRequest) error {

	_, err := tx.Exec(`
		INSERT INTO role_permission
			(snowflake_id, role_id, permission_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5)
	`, payload.SnowflakeId, payload.RoleId, payload.PermissionId, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil

}

func GetRolePermissionList(tx *sql.Tx, roleId string) ([]string, error) {

	baseSQL := `

		SELECT
			permission_id
		FROM
			role_permissions
		WHERE
			role_id = $1
	`

	rows, err := tx.Query(baseSQL, roleId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	permissionList := make([]string, 0)
	for rows.Next() {

		var permissionId string
		err := rows.Scan(
			&permissionId,
		)
		if err != nil {
			return nil, err
		}
		permissionList = append(permissionList, permissionId)
	}

	return permissionList, nil

}

func DeleteRolePermission(tx *sql.Tx, snowflakeId string) error {

	_, err := tx.Exec(`
		UPDATE role_permission
		SET
			deleted_at = $1
		WHERE
			snowflake_id = $2
	`, time.Now(), snowflakeId)
	if err != nil {
		return err
	}

	return nil

}

func DeleteRolePermissionByRoleId(tx *sql.Tx, roleId string) error {
	_, err := tx.Exec(`
		UPDATE role_permission
		SET
			deleted_at = $1
		WHERE
			role_id = $2
	`, time.Now(), roleId)
	if err != nil {
		return err
	}

	return nil

}
