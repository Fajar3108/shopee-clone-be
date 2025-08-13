package authaction

import (
	"github.com/Fajar3108/mafi-course-be/internal/model"
	authrequest "github.com/Fajar3108/mafi-course-be/internal/request/auth-request"
	"github.com/Fajar3108/mafi-course-be/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateNewUserSession(request *authrequest.UserSessionRequest, ctx *fiber.Ctx, db *gorm.DB) (*model.UserSession, error) {
	id, err := helpers.GenerateUUID()

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userSession := &model.UserSession{
		ID:             id,
		UserID:         request.UserID,
		Token:          request.Token,
		RefreshToken:   request.RefreshToken,
		TokenExpired:   request.TokenExpired,
		RefreshExpired: request.RefreshExpired,
	}

	result := db.Create(userSession)

	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return userSession, nil
}
