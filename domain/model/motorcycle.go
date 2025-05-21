package model

import (
	"gorm.io/gorm"
)

type MotorcycleAdd struct {
	gorm.Model
	BrandID     uint            `json:"brand_id"`
	Brand       MotorcycleBrand `gorm:"foreignKey:BrandID"`
	ModelID     uint            `json:"model_id"`
	MotoModel   MotorcycleModel `gorm:"foreignKey:ModelID"`
	Year        uint            `json:"year"`
	Volume      uint            `json:"volume"`
	Mileage     uint            `json:"mileage"`
	Description string          `json:"description"`
	Phone       string          `json:"phone" gorm:"not null"`
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
	Brand  string            `json:"brand" gorm:"not null;unique"`
	Models []MotorcycleModel `gorm:"foreignKey:MotorcycleBrandID"`
}

type MotorcycleModel struct {
	gorm.Model
	MotorcycleBrandID uint
	Brand             MotorcycleBrand `gorm:"foreignKey:MotorcycleBrandID"`
	MotoModel         string          `json:"model" gorm:"not null;unique"`
}
