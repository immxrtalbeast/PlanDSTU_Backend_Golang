package lib

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/immxrtalbeast/plandstu/internal/domain"
)

func NewToken(user *domain.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["role"] = user.Role

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}
