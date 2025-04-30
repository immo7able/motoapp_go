package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/domain/dto"
	"motorcycleApp/service"
	"motorcycleApp/utils"
)

type AuthHandler struct {
	AuthService *service.AuthService
	Validator   *validator.Validate
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ParseValidationErrors(err))
		return
	}

	if err := h.AuthService.RegisterUser(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
		return
	}

	token, err := h.AuthService.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Error{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func (h *AuthHandler) RegisterForm(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"errors": []string{"Неверная форма"},
		})
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		parsed := utils.ParseValidationErrors(err)
		var messages []string
		for _, fieldErr := range parsed.FieldErrors {
			messages = append(messages, fieldErr.Field+": "+fieldErr.Message)
		}
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"errors": messages,
		})
		return
	}

	if err := h.AuthService.RegisterUser(req); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"errors": []string{err.Error()},
		})
		return
	}

	c.Redirect(http.StatusFound, "/auth/login")
}

func (h *AuthHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) LoginForm(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"error": "Неверная форма",
		})
		return
	}

	token, err := h.AuthService.LoginUser(req)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Неверные данные",
		})
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"message": "Успешный вход: " + token,
	})
}
