package useraction

import (
	"github.com/Fajar3108/mafi-course-be/internal/model"
	authrequest "github.com/Fajar3108/mafi-course-be/internal/request/auth-request"
	file_storage "github.com/Fajar3108/mafi-course-be/pkg/file-storage"
	"github.com/Fajar3108/mafi-course-be/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateNewUser(request *authrequest.RegisterRequest, role string, ctx *fiber.Ctx, db *gorm.DB) (*model.User, error) {
	id, err := helpers.GenerateUUID()

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user := &model.User{
		ID:    id,
		Name:  request.Name,
		Email: request.Email,
		Role:  role,
	}

	if request.Avatar != nil {
		avatarPath, err := file_storage.Store(ctx, request.Avatar, "avatars", true)

		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		user.Avatar = avatarPath
	}

	if result := db.Where("email = ?", request.Email).First(user); result.Error == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user.Password = string(hashedPassword)

	result := db.Create(user)

	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error.Error())
	}

	return user, nil
}
