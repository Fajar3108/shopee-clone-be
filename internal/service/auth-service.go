package service

import (
	"time"

	"github.com/Fajar3108/mafi-course-be/database"
	authaction "github.com/Fajar3108/mafi-course-be/internal/action/auth-action"
	useraction "github.com/Fajar3108/mafi-course-be/internal/action/user-action"
	"github.com/Fajar3108/mafi-course-be/internal/model"
	authrequest "github.com/Fajar3108/mafi-course-be/internal/request/auth-request"
	"github.com/Fajar3108/mafi-course-be/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		DB: database.DB(),
	}
}

func (as *AuthService) Login(request *authrequest.LoginRequest) (jwToken string, refreshToken string, tokenExpired *time.Time, refreshExpired *time.Time, user *model.User, err error) {
	user = &model.User{}
	result := as.DB.Where("email = ?", request.Email).First(user)

	if result.Error != nil {
		return "", "", nil, nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Your email is not registered")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return "", "", nil, nil, nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	jwToken, refreshToken, tokenExpired, refreshExpired, err = authaction.GenerateAuthToken(user)

	if err != nil {
		return "", "", nil, nil, nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return jwToken, refreshToken, tokenExpired, refreshExpired, user, nil
}

func (as *AuthService) Register(request *authrequest.RegisterRequest, ctx *fiber.Ctx) (string, string, *time.Time, *time.Time, *model.User, error) {
	user, fiberError := useraction.CreateNewUser(request, "user", ctx, as.DB)

	if fiberError != nil {
		return "", "", nil, nil, nil, fiberError
	}

	go authaction.SendWelcomeEmail(user)

	jwToken, refreshToken, tokenExpired, refreshExpired, err := authaction.GenerateAuthToken(user)

	if err != nil {
		return "", "", nil, nil, nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return jwToken, refreshToken, tokenExpired, refreshExpired, user, nil
}

func (as *AuthService) Logout(ctx *fiber.Ctx) error {
	tokenJwt := ctx.Locals("user").(*jwt.Token)

	if err := as.DB.Where("token = ?", tokenJwt.Raw).Delete(&model.UserSession{}).Error; err != nil {
		return err
	}

	return nil
}

func (as *AuthService) RefreshToken(request *authrequest.RefreshTokenRequest, ctx *fiber.Ctx) (string, error) {
	userSession := &model.UserSession{}

	if err := as.DB.Where("refresh_token = ?", request.RefreshToken).First(userSession).Error; err != nil {
		return "", err
	}

	user := &model.User{}
	if err := as.DB.Where("id = ?", userSession.UserID).First(user).Error; err != nil {
		return "", nil
	}

	newExpiration := time.Now().Add(24 * time.Hour)
	jwToken, err := token.GenerateJWT(user, newExpiration)

	if err != nil {
		return "", nil
	}

	userSession.Token = jwToken
	userSession.TokenExpired = &newExpiration
	as.DB.Save(userSession)

	return jwToken, nil
}
