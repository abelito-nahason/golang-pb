package helper

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func VerifyAuth(tokenString string) (*UserClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")
	trimBearer := strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(trimBearer, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, nil
	}

}
