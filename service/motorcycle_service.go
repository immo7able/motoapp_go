package service

import (
	"gorm.io/gorm"
	"motorcycleApp/domain/dto"
	"motorcycleApp/domain/model"
)

type MotorcycleService struct {
	DB     *gorm.DB
	JWTKey []byte
}

func (s *MotorcycleService) AddMotorcycle(req dto.MotorcycleAddRequest, userID uint, phone string) error {
	moto := model.MotorcycleAdd{
		Brand:       req.Brand,
		MotoModel:   req.MotoModel,
		Year:        req.Year,
		Volume:      req.Volume,
		Mileage:     req.Mileage,
		Description: req.Description,
		Phone:       phone,
		AuthorID:    userID,
	}
	return s.DB.Create(&moto).Error
}

func (s *MotorcycleService) GetUserMotorcycles(userID uint) ([]model.MotorcycleAdd, error) {
	var motorcycles []model.MotorcycleAdd
	err := s.DB.Where("author_id = ?", userID).Find(&motorcycles).Error
	return motorcycles, err
}

func (s *MotorcycleService) GetAllMotorcycles() ([]model.MotorcycleAdd, error) {
	var motorcycles []model.MotorcycleAdd
	err := s.DB.Find(&motorcycles).Error
	return motorcycles, err
}

func (s *MotorcycleService) DeleteMotorcycle(id string, userID uint) error {
	return s.DB.
		Where("id = ? AND author_id = ?", id, userID).
		Delete(&model.MotorcycleAdd{}).
		Error
}
