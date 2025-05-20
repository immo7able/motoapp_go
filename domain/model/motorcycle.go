package model

import (
	"gorm.io/gorm"
)

type MotorcycleAdd struct {
	gorm.Model
	Brand       string `json:"brand"`
	MotoModel   string `json:"model"`
	Year        uint   `json:"year"`
	Volume      uint   `json:"volume"`
	Mileage     uint   `json:"mileage"`
	Description string `json:"description"`
	Phone       string `json:"phone" gorm:"not null"`
	AuthorID    uint
	Author      User              `gorm:"foreignKey:AuthorID"`
	Images      []MotorcycleImage `gorm:"foreignKey:MotorcycleID"`
}

type MotorcycleImage struct {
	gorm.Model
	MotorcycleID uint
	Motorcycle   MotorcycleAdd `gorm:"foreignKey:MotorcycleID"`
	ImagePath    string
}

type MotorcycleBrand struct {
	gorm.Model
	Brand string `json:"brand" gorm:"not null;unique"`
}

type MotorcycleModel struct {
	gorm.Model
	MotorcycleBrandID uint
	Brand             MotorcycleBrand `gorm:"foreignKey:MotorcycleBrandID"`
	MotoModel         string          `json:"model" gorm:"not null;unique"`
}
