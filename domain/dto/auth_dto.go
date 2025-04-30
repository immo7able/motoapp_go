package dto

type RegisterRequest struct {
	Login    string `binding:"required"`
	Password string `binding:"required,gte=8"`
	Email    string `binding:"required,email"`
	Phone    string `binding:"required,numeric,len=10"`
}

type LoginRequest struct {
	Email    string `binding:"required"`
	Password string `binding:"required,gte=8"`
}
