package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId uint, email string, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key_change_me"
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
