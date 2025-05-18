package dto

type MotorcycleAddRequest struct {
	Brand       string `form:"Brand" json:"brand" validate:"required"`
	MotoModel   string `form:"Model" json:"model" validate:"required"`
	Year        uint   `form:"Year" json:"year" validate:"required"`
	Volume      uint   `form:"Volume" json:"volume" validate:"required"`
	Mileage     uint   `form:"Mileage" json:"mileage" validate:"required"`
	Description string `form:"Description" json:"description"`
}
