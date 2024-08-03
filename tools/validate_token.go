package tools

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func ValidateUserToken(authorization string, identity string) (string, error) {

	token, err := ValidateToken(authorization)

	if err != nil {
		return "", fmt.Errorf("未经授权: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)

	snowflakeId, ok := claims["snowflake_id"].(string)
	if !ok {
		return "", fmt.Errorf("snowflake_id is not a string")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role is not a string")
	}

	if role != identity {
		return "", fmt.Errorf("未经授权")
	}
	// snowflakeId, err := strconv.ParseInt(snowflakeIdStr, 10, 64)
	// if err != nil {
	// 	return 0, fmt.Errorf("无法将 snowflake_id 转换为 int64: %v", err)
	// }

	return snowflakeId, nil
}
