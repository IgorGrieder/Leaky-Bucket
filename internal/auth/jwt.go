package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	UserID string `json:"user_id"`
	PIX    string `json:"pix"`
	jwt.RegisteredClaims
}

func Authenticate(tokenString string, hashSecret string) (*jwt.Token, error) {
	claims := &JWT{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hashSecret, nil
	})

	return token, err

}

func GenerateToken(userID, pix string, secretKey []byte) (string, error) {
	// Create the claims
	claims := &JWT{
		UserID: userID,
		PIX:    pix,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	return signedToken, err
}
