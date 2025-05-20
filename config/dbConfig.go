package config

import (
	"golang.org/x/crypto/bcrypt"
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

	err = database.AutoMigrate(&model.MotorcycleAdd{}, &model.MotorcycleImage{}, &model.User{})
	if err != nil {
		panic("Failed to migrate database!")
	}
	err = CreateAdminUser(database)
	if err != nil {
		panic("Failed to create admin user!")
	}

	DB = database
}

func CreateAdminUser(db *gorm.DB) error {
	var count int64
	db.Model(&model.User{}).Where("email = ?", "admin@gmail.com").Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Create(&model.User{
		Login:    "admin",
		Email:    "admin@gmail.com",
		Password: string(hashedPassword),
		Phone:    "7777777777",
		Role:     model.RoleAdmin,
	}).Error
}
