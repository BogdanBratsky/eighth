package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret []byte, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: secret,
		ttl:    ttl,
	}
}

func (j *JWTManager) GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(j.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := int(claims["user_id"].(float64))
	return userID, nil
}
