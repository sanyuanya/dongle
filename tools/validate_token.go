package tools

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func ValidateUserToken(authorization string, identity string) (string, error) {

	token, err := ValidateToken(authorization)

	if err != nil {
		return "", fmt.Errorf("无效的token: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)

	snowflakeId, ok := claims["snowflake_id"].(string)
	if !ok {
		return "", fmt.Errorf("snowflake_id is not a string")
	}

	role, ok := claims["role"].(string)
	_ = role
	if !ok {
		return "", fmt.Errorf("role is not a string")
	}

	// if role != identity {
	// 	return "", fmt.Errorf("未经授权")
	// }

	return snowflakeId, nil
}
