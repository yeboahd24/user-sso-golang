package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yeboahd24/user-sso/config"
)

func GenerateJWT(user interface{}, cfg *config.Config) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(cfg.JWT.ExpiresInHrs)).Unix()

	return token.SignedString([]byte(cfg.JWT.Secret))
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
