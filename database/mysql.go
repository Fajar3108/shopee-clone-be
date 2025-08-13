package database

import (
	"fmt"
	"log"

	"github.com/Fajar3108/mafi-course-be/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func getDsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		viper.GetString(config.DbUser),
		viper.GetString(config.DbPassword),
		viper.GetString(config.DbHost),
		viper.GetString(config.DbPort),
		viper.GetString(config.DbName),
	)
}

func DB() *gorm.DB {
	if db != nil {
		return db
	}

	var err error

	db, err = gorm.Open(mysql.Open(getDsn()))

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return db
}
