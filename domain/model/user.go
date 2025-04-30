package model

import (
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"type:varchar(255);unique;not null"`
	Login    string `json:"login" gorm:"type:varchar(30);not null"`
	Password string `json:"password" gorm:"not null"`
	Phone    string `json:"phone" gorm:"type:varchar(11);unique;not null"`
	Role     Role   `gorm:"type:varchar(20);not null"`
}
