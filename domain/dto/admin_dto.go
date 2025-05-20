package dto

type CreateBrandRequest struct {
	Name string `form:"Name" json:"name" validate:"required"`
}

type CreateModelRequest struct {
	Name    string `form:"Name" json:"name" validate:"required"`
	BrandID uint   `form:"BrandID" json:"brand_id" validate:"required"`
}
