package data

import (
	"fmt"

	"github.com/sanyuanya/dongle/entity"
	"golang.org/x/crypto/bcrypt"
)

func Login(auth *entity.LoginRequest) (int64, error) {

	// Check if the account exists

	var snowflakeId int64
	var password string

	err := db.QueryRow("SELECT snowflake_id, password FROM admins WHERE account=$1", auth.Account).Scan(&snowflakeId, &password)
	if err != nil {
		return 0, err
	}

	hashPassword, err := HashPassword(auth.Password)

	if err != nil {
		return 0, fmt.Errorf("密码 bcrypt 加密失败")
	}
	// Compare the stored hashed password, with the hashed version of the password that was received
	if CheckPasswordHash(auth.Password, hashPassword) {
		return snowflakeId, nil
	}

	return 0, fmt.Errorf("密码错误")
}

func SetApiToken(snowflakeId int64, token string) error {
	_, err := db.Exec("UPDATE admins SET api_token=$1 WHERE snowflake_id=$2", token, snowflakeId)
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