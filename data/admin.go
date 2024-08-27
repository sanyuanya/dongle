package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyuanya/dongle/entity"
	"golang.org/x/crypto/bcrypt"
)

func Login(tx *sql.Tx, auth *entity.LoginRequest) (string, error) {

	// Check if the account exists
	var snowflakeId string
	var password string

	err := tx.QueryRow("SELECT snowflake_id, password FROM admins WHERE account=$1 AND deleted_at IS NULL", auth.Account).Scan(&snowflakeId, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("账号不存在")
		}

		return "", err
	}

	hashPassword, err := HashPassword(auth.Password)

	if err != nil {
		return "", fmt.Errorf("密码 bcrypt 加密失败")
	}
	// Compare the stored hashed password, with the hashed version of the password that was received
	if CheckPasswordHash(auth.Password, hashPassword) {
		return snowflakeId, nil
	}

	return "", fmt.Errorf("密码错误")
}

func SetApiToken(tx *sql.Tx, snowflakeId string, token string) error {
	_, err := tx.Exec("UPDATE admins SET api_token=$1 WHERE snowflake_id=$2 AND deleted_at IS NULL", token, snowflakeId)
	if err != nil {
		return err
	}
	return nil
}

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with its possible plaintext equivalent.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DeleteAdmin(tx *sql.Tx, snowflakeId string) error {
	_, err := tx.Exec(`UPDATE admins SET deleted_at=$1 WHERE snowflake_id=$2 AND deleted_at IS NULL`, time.Now(), snowflakeId)
	if err != nil {
		return err
	}
	return nil
}

func AddAdmin(tx *sql.Tx, admin *entity.AddAdminRequest) error {
	hashPassword, err := HashPassword(admin.Password)

	if err != nil {
		return fmt.Errorf("密码 bcrypt 加密失败")
	}

	_, err = tx.Exec(`INSERT INTO admins (snowflake_id, account, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`, admin.SnowflakeId, admin.Account, hashPassword, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func UpdateAdmin(tx *sql.Tx, admin *entity.UpdateAdminRequest) error {
	_, err := tx.Exec(`UPDATE admins SET account=$1, updated_at=$2 WHERE snowflake_id=$3 AND deleted_at IS NULL`, admin.Account, time.Now(), admin.SnowflakeId)
	if err != nil {
		return err
	}
	return nil
}

func GetAdminList(tx *sql.Tx, req *entity.GetAdminListRequest) ([]*entity.GetAdminListResponse, error) {

	adminList := make([]*entity.GetAdminListResponse, 0)
	baseSQL := `SELECT snowflake_id, account FROM admins WHERE deleted_at IS NULL AND is_hidden = 0 ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := tx.Query(baseSQL, req.PageSize, req.PageSize*(req.Page-1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var admin entity.GetAdminListResponse
		err := rows.Scan(&admin.SnowflakeId, &admin.Account)
		if err != nil {
			return nil, err
		}
		adminList = append(adminList, &admin)
	}

	return adminList, nil
}

func GetAdminTotal(tx *sql.Tx, req *entity.GetAdminListRequest) (int64, error) {
	var total int64
	err := tx.QueryRow(`SELECT COUNT(*) FROM admins WHERE deleted_at IS NULL AND is_hidden = 0`).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAdminByAccount(tx *sql.Tx, account string) (string, error) {

	baseSQL := `
		SELECT 
			snowflake_id
		FROM 
			admins
		WHERE 
			account = $1 AND deleted_at IS NULL
	`
	var snowflakeId string
	err := tx.QueryRow(baseSQL, account).Scan(&snowflakeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", fmt.Errorf("查询失败:%v", err)
	}
	return snowflakeId, nil
}

func FindBySnowflakeIdNotFoundAndAccount(tx *sql.Tx, snowflakeId string, account string) error {
	var count int
	err := tx.QueryRow(`SELECT COUNT(*) FROM admins WHERE account=$1 AND snowflake_id!=$2 AND deleted_at IS NULL`, account, snowflakeId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("账号已存在")
	}
	return nil
}
