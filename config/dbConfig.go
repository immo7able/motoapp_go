package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"motorcycleApp/domain/model"
)

var DB *gorm.DB

func ConnectDatabase(databaseUrl string) {
	database, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&model.MotorcycleAdd{}, &model.User{})
	if err != nil {
		panic("Failed to migrate database!")
	}

	DB = database
}
