package validation

import (
	"strings"

	"github.com/Fajar3108/mafi-course-be/database"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func isUnique(db *gorm.DB, tableName string, fieldName string, value any) bool {
	var count int64

	if err := db.Table(tableName).Where(fieldName+" = ?", value).Count(&count).Error; err != nil || count > 0 {
		return false
	}

	return true
}

func RegisterUniqueValidation(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ".")

	if len(params) != 2 {
		return false
	}

	tableName := params[0]
	fieldName := params[1]

	fieldValue := fl.Field().Interface()

	return isUnique(database.DB(), tableName, fieldName, fieldValue)
}
