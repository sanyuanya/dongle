package data

import (
	"database/sql"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func AddAdminRole(tx *sql.Tx, payload *entity.AddAdminRoleRequest) error {
	baseSQL := `INSERT INTO admin_role (snowflake_id, admin_id, role_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(baseSQL, payload.SnowflakeId, payload.AdminId, payload.RoleId, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func DeleteAdminRole(tx *sql.Tx, adminId string) error {
	_, err := tx.Exec(`DELETE FROM admin_role WHERE admin_id=$1`, adminId)
	if err != nil {
		return err
	}

	return nil
}

func GetAdminRoleList(tx *sql.Tx, adminId string) ([]*entity.GetAdminRoleResponse, error) {
	baseSQL := `SELECT a.admin_id, a.role_id, a.created_at, a.updated_at, r.name FROM admin_role a JOIN roles r ON a.role_id = r.snowflake_id WHERE a.admin_id=$1 AND a.deleted_at IS NULL`
	rows, err := tx.Query(baseSQL, adminId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adminRoleList := make([]*entity.GetAdminRoleResponse, 0)
	for rows.Next() {
		adminRole := &entity.GetAdminRoleResponse{}
		err := rows.Scan(
			&adminRole.AdminId,
			&adminRole.RoleId,
			&adminRole.CreatedAt,
			&adminRole.UpdatedAt,
			&adminRole.RoleName,
		)
		if err != nil {
			return nil, err
		}
		adminRoleList = append(adminRoleList, adminRole)
	}

	return adminRoleList, nil
}

func GetRoleUsed(tx *sql.Tx, roleId string) (string, error) {

	baseSQL := `
		SELECT 
			admin_id 
		FROM 
			admin_role
		WHERE role_id = $1 AND deleted_at IS NULL
	`
	var adminId string

	err := tx.QueryRow(baseSQL, roleId).Scan(&adminId)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return adminId, nil
}
