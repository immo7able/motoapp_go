package service

import (
	"errors"
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
			BrandID:     req.BrandID,
			ModelID:     req.ModelID,
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
	err := s.DB.Preload("Images").
		Preload("MotoModel").
		Preload("Brand").
		Where("author_id = ?", userID).
		Find(&motorcycles).Error
	return motorcycles, err
}

func (s *MotorcycleService) GetAllMotorcycles() ([]model.MotorcycleAdd, error) {
	var motorcycles []model.MotorcycleAdd
	err := s.DB.Preload("Images").
		Preload("MotoModel").
		Preload("Brand").
		Find(&motorcycles).Error
	return motorcycles, err
}

func (s *MotorcycleService) DeleteMotorcycle(id string, userID uint) error {
	return s.DB.
		Where("id = ? AND author_id = ?", id, userID).
		Delete(&model.MotorcycleAdd{}).
		Error
}

func (s *MotorcycleService) GetMotorcycleByIDAndAuthor(id uint, userID uint) (*model.MotorcycleAdd, error) {
	var moto model.MotorcycleAdd
	err := s.DB.Preload("Images").
		Where("id = ? AND author_id = ?", id, userID).
		First(&moto).Error
	if err != nil {
		return nil, err
	}
	return &moto, nil
}

func (s *MotorcycleService) GetMotorcycleByID(id uint) (*model.MotorcycleAdd, error) {
	var moto model.MotorcycleAdd
	err := s.DB.Preload("Images").
		Preload("MotoModel").
		Preload("Brand").
		Where("id = ?", id).
		First(&moto).Error
	if err != nil {
		return nil, err
	}
	return &moto, nil
}

func (s *MotorcycleService) UpdateMotorcycle(ID uint, req dto.MotorcycleUpdateRequest, userID uint, imagePaths []string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var moto model.MotorcycleAdd
		if err := tx.Where("id = ? AND author_id = ?", ID, userID).First(&moto).Error; err != nil {
			return errors.New("мотоцикл не найден или не принадлежит пользователю")
		}

		moto.BrandID = req.BrandID
		moto.ModelID = req.ModelID
		moto.Year = req.Year
		moto.Volume = req.Volume
		moto.Mileage = req.Mileage
		moto.Description = req.Description

		if err := tx.Save(&moto).Error; err != nil {
			return err
		}

		if len(imagePaths) > 0 {
			if err := tx.Where("motorcycle_id = ?", moto.ID).Delete(&model.MotorcycleImage{}).Error; err != nil {
				return err
			}

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

func (s *MotorcycleService) GetCommentsByMotorcycleID(motoID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := s.DB.Preload("User").Where("motorcycle_id = ?", motoID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (s *MotorcycleService) SaveComment(comment model.Comment) error {
	return s.DB.Create(&comment).Error
}

func (s *MotorcycleService) GetCommentByID(id string) (*model.Comment, error) {
	var comment model.Comment
	err := s.DB.First(&comment, id).Error
	return &comment, err
}

func (s *MotorcycleService) DeleteComment(id string) error {
	return s.DB.Delete(&model.Comment{}, id).Error
}
