package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func VerifyAuth(tokenString string) error {
	secretKey := os.Getenv("SECRET_KEY")
	trimBearer := strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(trimBearer, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
