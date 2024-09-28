package pkg

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))
var jwtTtl = time.Hour * 24

func CreateJwtToken(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_auth_id": userId.String(),
			"exp":          time.Now().Add(jwtTtl).Unix(),
		})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJwtToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

// Helper function to return a string or nil
func GetNullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
