package dto

type RegisterRequest struct {
	Login    string `form:"Login" json:"login" validate:"required"`
	Password string `form:"Password" json:"password" validate:"required,gte=8"`
	Email    string `form:"Email" json:"email" validate:"required,email"`
	Phone    string `form:"Phone" json:"phone" validate:"required,numeric,len=10"`
}

type LoginRequest struct {
	Email    string `form:"Email" json:"email" validate:"required,email"`
	Password string `form:"Password" json:"password" validate:"required,gte=8"`
}
