package application

import (
	"fmt"
	"time"

	"github.com/IgorGrieder/Leaky-Bucket/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	UserID string
	jwt.RegisteredClaims
}

type AuthService struct {
	config *config.Config
}

func (auth *AuthService) Authenticate(tokenString string, hashSecret string) (*jwt.Token, error) {
	claims := &JWT{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hashSecret, nil
	})

	return token, err

}

func (auth *AuthService) GenerateToken(userID string) (string, error) {
	claims := &JWT{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secretKey := auth.config.HASH

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error while signing the JWT token: %v", err)
	}

	return signedToken, nil
}
