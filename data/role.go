package data

import (
	"database/sql"
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
			deleted_at IS NULL
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
			name = ?,
			updated_at = ?
		WHERE
			snowflake_id = ?
	`, payload.Name, time.Now(), payload.SnowflakeId)
	if err != nil {
		return err
	}

	return nil
}
