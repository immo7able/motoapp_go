package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"motorcycleApp/domain/dto"
	"motorcycleApp/domain/model"
	"time"
)

type AuthService struct {
	DB     *gorm.DB
	JWTKey []byte
}

func (s *AuthService) RegisterUser(req dto.RegisterRequest) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Login:    req.Login,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     model.RoleUser,
	}

	return s.DB.Create(&user).Error
}

func (s *AuthService) LoginUser(req dto.LoginRequest) (string, error) {
	var user model.User
	err := s.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return "", errors.New("invalid login")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"login":   user.Login,
		"email":   user.Email,
		"phone":   user.Phone,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(s.JWTKey)
}

func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateUser(user *model.User) error {
	return s.DB.Save(user).Error
}
