package authaction

import (
	"time"

	"github.com/Fajar3108/mafi-course-be/internal/model"
	"github.com/Fajar3108/mafi-course-be/pkg/token"
)

func GenerateAuthToken(user *model.User) (string, string, *time.Time, *time.Time, error) {
	tokenExpired := time.Now().Add(24 * time.Hour)
	jwToken, err := token.GenerateJWT(user, tokenExpired)

	if err != nil {
		return "", "", nil, nil, err
	}

	refreshExpired := time.Now().Add(14 * 24 * time.Hour)
	refreshToken, err := token.GenerateJWT(user, refreshExpired)

	if err != nil {
		return "", "", nil, nil, err
	}

	return jwToken, refreshToken, &tokenExpired, &refreshExpired, nil
}
