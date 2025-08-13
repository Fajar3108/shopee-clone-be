package controller

import (
	"strings"

	authaction "github.com/Fajar3108/mafi-course-be/internal/action/auth-action"
	authrequest "github.com/Fajar3108/mafi-course-be/internal/request/auth-request"
	"github.com/Fajar3108/mafi-course-be/internal/resource"
	"github.com/Fajar3108/mafi-course-be/internal/service"
	file_storage "github.com/Fajar3108/mafi-course-be/pkg/file-storage"
	"github.com/Fajar3108/mafi-course-be/pkg/helpers"
	"github.com/Fajar3108/mafi-course-be/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		service: service.NewAuthService(),
	}
}

func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	req := &authrequest.LoginRequest{}

	if err := validation.Validate(ctx, req); err != nil {
		return err
	}

	token, refreshToken, tokenExpired, refreshExpired, user, err := ac.service.Login(req)

	if err != nil {
		return err
	}

	avatarUrl := strings.Join([]string{file_storage.GetStorageURL(ctx), user.Avatar}, "/")
	user.Avatar = avatarUrl

	userSessionRequest := &authrequest.UserSessionRequest{
		UserID:         user.ID,
		Token:          token,
		RefreshToken:   refreshToken,
		TokenExpired:   tokenExpired,
		RefreshExpired: refreshExpired,
	}

	_, err = authaction.CreateNewUserSession(userSessionRequest, ctx, ac.service.DB)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	data := resource.NewAuthResource(token, refreshToken, user)

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Login successful",
		data,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (ac *AuthController) Register(ctx *fiber.Ctx) error {
	req := &authrequest.RegisterRequest{}

	if err := validation.Validate(ctx, req); err != nil {
		return err
	}

	if avatar, err := ctx.FormFile("avatar"); err != nil {
		req.Avatar = nil
	} else {
		req.Avatar = avatar
	}

	tokenJWT, refreshToken, tokenExpired, refreshExpired, user, err := ac.service.Register(req, ctx)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userSessionRequest := &authrequest.UserSessionRequest{
		UserID:         user.ID,
		Token:          tokenJWT,
		RefreshToken:   refreshToken,
		TokenExpired:   tokenExpired,
		RefreshExpired: refreshExpired,
	}

	_, err = authaction.CreateNewUserSession(userSessionRequest, ctx, ac.service.DB)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	data := resource.NewAuthResource(tokenJWT, refreshToken, user)

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Registration successful",
		data,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (ac *AuthController) Logout(ctx *fiber.Ctx) error {
	err := ac.service.Logout(ctx)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Logout successful",
		nil,
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}

func (ac *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	req := &authrequest.RefreshTokenRequest{}

	if err := validation.Validate(ctx, req); err != nil {
		return err
	}

	newToken, err := ac.service.RefreshToken(req, ctx)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res := helpers.NewResponseHelper(
		fiber.StatusOK,
		"Token refreshed successfuly",
		map[string]string{
			"token": newToken,
		},
		nil,
		nil,
	)

	return ctx.Status(res.Code).JSON(res)
}
