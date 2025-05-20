package service

import (
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
