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

// func IdFromToken(tokenString string, secret string) (uuid.UUID, error) {
// 	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(secret), nil
// 	})

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if id, exists := claims["uid"].(float64); exists {
// 			return id, nil
// 		}
// 	}

// 	return uuid.Nil, err
// }
