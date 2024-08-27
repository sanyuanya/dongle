package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func GetRoleList(tx *sql.Tx) ([]*entity.Role, error) {

	baseSQL := `

		SELECT
			snowflake_id,
			name,
			created_at,
			updated_at
		FROM
			roles
		WHERE
			deleted_at IS NULL ORDER BY created_at DESC
	`

	rows, err := tx.Query(baseSQL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := make([]*entity.Role, 0)
	for rows.Next() {

		role := &entity.Role{}
		err := rows.Scan(
			&role.SnowflakeID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil

}

func AddRole(tx *sql.Tx, payload *entity.AddRoleRequest) error {

	_, err := tx.Exec(`
		INSERT INTO roles
			(snowflake_id, name, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)
	`, payload.SnowflakeId, payload.Name, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func UpdateRole(tx *sql.Tx, payload *entity.UpdateRoleRequest) error {

	_, err := tx.Exec(`
		UPDATE
			roles
		SET
			name = $1,
			updated_at = $2
		WHERE
			snowflake_id = $3 AND deleted_at IS NULL
	`, payload.Name, time.Now(), payload.SnowflakeId)
	if err != nil {
		return err
	}

	return nil
}

func GetRole(tx *sql.Tx, roleId string) (*entity.Role, error) {

	baseSQL := `
		SELECT
			snowflake_id,
			name,
			created_at,
			updated_at
		FROM
			roles
		WHERE
			snowflake_id = $1 AND deleted_at IS NULL
	`

	role := &entity.Role{}
	err := tx.QueryRow(baseSQL, roleId).Scan(
		&role.SnowflakeID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

func GetRoleByName(tx *sql.Tx, name string) (*entity.Role, error) {

	baseSQL := `
		SELECT
			snowflake_id,
			name,
			created_at,
			updated_at
		FROM
			roles
		WHERE
			name = $1 AND deleted_at IS NULL
	`

	role := &entity.Role{}
	err := tx.QueryRow(baseSQL, name).Scan(
		&role.SnowflakeID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

func DeleteRole(tx *sql.Tx, roleId string) error {
	baseSQL := `
		UPDATE 
			roles
		SET
			deleted_at=$1
		WHERE snowflake_id=$2 AND deleted_at IS NULL
	`

	_, err := tx.Exec(baseSQL, time.Now(), roleId)
	if err != nil {
		return fmt.Errorf("角色删除失败：%v", err)
	}

	return nil
}

func FindBySnowflakeIdNotFoundAndRoleName(tx *sql.Tx, snowflakeId, name string) (*entity.Role, error) {
	baseSQL := `
		SELECT
			snowflake_id,
			name,
			created_at,
			updated_at
		FROM
			roles
		WHERE
			snowflake_id != $1 AND name = $2 AND deleted_at IS NULL
	`

	role := &entity.Role{}
	err := tx.QueryRow(baseSQL, snowflakeId, name).Scan(
		&role.SnowflakeID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}
