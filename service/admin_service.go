package service

import (
	"errors"
	"gorm.io/gorm"
	"motorcycleApp/domain/model"
)

type AdminService struct {
	DB *gorm.DB
}

func (s *AdminService) CreateBrand(brand string) error {
	return s.DB.Create(&model.MotorcycleBrand{Brand: brand}).Error
}

func (s *AdminService) GetAllBrands() ([]model.MotorcycleBrand, error) {
	var brands []model.MotorcycleBrand
	err := s.DB.Find(&brands).Error
	return brands, err
}

func (s *AdminService) CreateModel(brandID uint, modelName string) error {
	return s.DB.Create(&model.MotorcycleModel{
		MotorcycleBrandID: brandID,
		MotoModel:         modelName,
	}).Error
}

func (s *AdminService) GetAllModels() ([]model.MotorcycleModel, error) {
	var models []model.MotorcycleModel
	err := s.DB.Preload("Brand").Find(&models).Error
	return models, err
}

func (s *AdminService) UpdateBrand(id uint, newName string) error {
	var count int64

	if err := s.DB.Model(&model.MotorcycleBrand{}).
		Where("brand = ? AND id != ?", newName, id).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("duplicate_name")
	}

	return s.DB.Model(&model.MotorcycleBrand{}).
		Where("id = ?", id).
		Update("brand", newName).Error
}

func (s *AdminService) DeleteBrand(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("motorcycle_brand_id = ?", id).Delete(&model.MotorcycleModel{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.MotorcycleBrand{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *AdminService) UpdateModel(id uint, newName string, brandID uint) error {
	return s.DB.Model(&model.MotorcycleModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"moto_model":          newName,
		"motorcycle_brand_id": brandID,
	}).Error
}

func (s *AdminService) DeleteModel(id uint) error {
	return s.DB.Delete(&model.MotorcycleModel{}, id).Error
}
