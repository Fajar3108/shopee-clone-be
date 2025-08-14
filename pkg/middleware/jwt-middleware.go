package middleware

import (
	"github.com/Fajar3108/mafi-course-be/config"
	"github.com/Fajar3108/mafi-course-be/database"
	"github.com/Fajar3108/mafi-course-be/internal/model"
	"github.com/Fajar3108/mafi-course-be/pkg/helpers"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(viper.GetString(config.JwtSecretKey)),
		},
		ErrorHandler:   errorHandler,
		SuccessHandler: successHandler,
	})
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	res := helpers.NewResponseHelper(
		fiber.StatusUnauthorized,
		"Unauthorized",
		nil,
		nil,
		err.Error(),
	)
	return ctx.Status(res.Code).JSON(res)
}

func successHandler(ctx *fiber.Ctx) error {
	tokenJwt := ctx.Locals("user").(*jwt.Token)

	db := database.DB()

	userSession := &model.UserSession{}
	if err := db.First(userSession, "token = ?", tokenJwt.Raw).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			res := helpers.NewResponseHelper(
				fiber.StatusUnauthorized,
				"Unauthorized: Session not found",
				nil, nil, err.Error(),
			)
			return ctx.Status(res.Code).JSON(res)
		}

		res := helpers.NewResponseHelper(
			fiber.StatusInternalServerError,
			"Database error",
			nil, nil, err.Error(),
		)

		return ctx.Status(res.Code).JSON(res)
	}

	return ctx.Next()
}
