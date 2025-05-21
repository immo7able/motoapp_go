package dto

type MotorcycleAddRequest struct {
	BrandID     uint   `form:"BrandID" validate:"required"`
	ModelID     uint   `form:"ModelID" validate:"required"`
	Year        uint   `form:"Year" json:"year" validate:"required,min=1900,max=2100"`
	Volume      uint   `form:"Volume" json:"volume" validate:"required,min=1"`
	Mileage     uint   `form:"Mileage" json:"mileage" validate:"required,min=0"`
	Description string `form:"Description" json:"description"`
}

type MotorcycleUpdateRequest struct {
	BrandID     uint   `form:"BrandID" validate:"required"`
	ModelID     uint   `form:"ModelID" validate:"required"`
	Year        uint   `form:"Year" validate:"required,min=1900,max=2100"`
	Volume      uint   `form:"Volume" validate:"required,min=1"`
	Mileage     uint   `form:"Mileage" validate:"required,min=0"`
	Description string `form:"Description"`
}
