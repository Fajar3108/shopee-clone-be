package token

import (
	"errors"
	"strings"
	"time"

	"github.com/Fajar3108/mafi-course-be/config"
	"github.com/Fajar3108/mafi-course-be/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type UserClaims struct {
	User *model.User `json:"user"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *model.User, expirationTime time.Time) (string, error) {
	claims := &UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(viper.GetString(config.JwtSecretKey)))
}

func ParseJWT(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}

	jwToken, err := jwt.ParseWithClaims(
		strings.TrimPrefix(tokenStr, "Bearer "),
		claims,
		func(tkn *jwt.Token) (any, error) {
			return []byte(viper.GetString(config.JwtSecretKey)), nil
		})

	if err != nil {
		return nil, err
	}

	if !jwToken.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
