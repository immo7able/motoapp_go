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

func (s *MotorcycleService) AddMotorcycle(req dto.MotorcycleAddRequest, userID uint, phone string, imagePaths []string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
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
		if err := tx.Create(&moto).Error; err != nil {
			return err
		}

		if len(imagePaths) > 0 {
			var images []model.MotorcycleImage
			for _, path := range imagePaths {
				images = append(images, model.MotorcycleImage{
					MotorcycleID: moto.ID,
					ImagePath:    path,
				})
			}
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *MotorcycleService) GetUserMotorcycles(userID uint) ([]model.MotorcycleAdd, error) {
	var motorcycles []model.MotorcycleAdd
	err := s.DB.Preload("Images").Where("author_id = ?", userID).Find(&motorcycles).Error
	return motorcycles, err
}

func (s *MotorcycleService) GetAllMotorcycles() ([]model.MotorcycleAdd, error) {
	var motorcycles []model.MotorcycleAdd
	err := s.DB.Preload("Images").Find(&motorcycles).Error
	return motorcycles, err
}

func (s *MotorcycleService) DeleteMotorcycle(id string, userID uint) error {
	return s.DB.
		Where("id = ? AND author_id = ?", id, userID).
		Delete(&model.MotorcycleAdd{}).
		Error
}
