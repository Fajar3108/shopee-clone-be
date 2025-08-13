package validation

import (
	errorhandler "github.com/Fajar3108/mafi-course-be/pkg/error-handler"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	if validate == nil {
		validate = validator.New(validator.WithRequiredStructEnabled())

		english := en.New()
		uni := ut.New(english, english)
		trans, _ = uni.GetTranslator("en")

		if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
			panic(err)
		}
	}
}

func FiberValidationError(err error) error {
	messages := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		messages[err.Field()] = err.Translate(trans)
	}

	if len(messages) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Validation error")
	}

	return errorhandler.NewValidationError("Validation error", messages)
}

func Validate[T any](ctx *fiber.Ctx, request *T) (err error) {
	err = ctx.BodyParser(request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if errs := validate.Struct(request); errs != nil {
		return FiberValidationError(errs)
	}

	return nil
}
