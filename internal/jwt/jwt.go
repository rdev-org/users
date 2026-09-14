package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rdev-org/users/internal/config"
)

func GenerateToken(userID int64) (string, error) {
	appConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	jwtKey := []byte(appConfig.JWTSecretKey)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
