package main

import (
	"log"

	"github.com/Fajar3108/mafi-course-be/config"
	"github.com/Fajar3108/mafi-course-be/database"
	"github.com/Fajar3108/mafi-course-be/internal/model"
)

func main() {
	config.InitConfig()

	err := database.DB().AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.UserSession{},
	)

	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
}
