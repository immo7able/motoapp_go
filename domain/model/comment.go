package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID        uint
	User          User `gorm:"foreignKey:UserID"`
	MotorcycleID  uint
	MotorcycleAdd MotorcycleAdd `gorm:"foreignKey:MotorcycleID"`
	Content       string        `gorm:"type:text;not null"`
}
