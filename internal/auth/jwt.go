package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	UserID string `json:"user_id"`
	PIX    string `json:"pix"`
	jwt.RegisteredClaims
}

func Authenticate(tokenString string, hashSecret string) (*jwt.Token, error) {
	claims := &JWT{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hashSecret, nil
	})

	return token, err

}
