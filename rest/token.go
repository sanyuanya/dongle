package rest

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Secret key used to sign tokens
var jwtSecret = []byte("46DBF4EC525D3DF58FA18923C8C8E")

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// GenerateToken generates a JWT token
func GenerateToken(snowflakeId int64, role string) (string, error) {

	// Create the Claims
	claims := jwt.MapClaims{
		"snowflake_id": fmt.Sprintf("%d", snowflakeId),
		"role":         role,
		"nbf":          time.Now().Unix() - 1,
		"iat":          time.Now().Unix(),
		"jti":          fmt.Sprintf("%d", time.Now().Unix()),
		"aud":          "dongle",
		"sub":          "dongle",
		"iss":          "dongle",
		"exp":          time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(jwtSecret)
}
