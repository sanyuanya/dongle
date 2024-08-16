package data

import (
	"database/sql"
	"fmt"
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
	result, err := tx.Exec(`DELETE FROM admin_role WHERE admin_id=$1`, adminId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("删除管理员角色失败: %v", err)
	}

	if affected == 0 {
		return fmt.Errorf("未找到对应的管理员角色")
	}

	return nil
}

func GetAdminRoleList(tx *sql.Tx, adminId string) ([]*entity.GetAdminRoleResponse, error) {
	baseSQL := `SELECT admin_id, role_id, created_at, updated_at FROM admin_role WHERE admin_id=$1`
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
		)
		if err != nil {
			return nil, err
		}
		adminRoleList = append(adminRoleList, adminRole)
	}

	return adminRoleList, nil
}
