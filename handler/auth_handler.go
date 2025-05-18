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
	if c.Request.Header.Get("Content-Type") == "application/json" {
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
		return
	}

	var req dto.RegisterRequest
	fieldErrors := map[string]string{}
	var globalErrors []string

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&req); err != nil {
			globalErrors = append(globalErrors, "Неверная форма")
		} else if err := h.Validator.Struct(req); err != nil {
			parsed := utils.ParseValidationErrors(err)
			for _, fe := range parsed.FieldErrors {
				fieldErrors[fe.Field] = fe.Message
			}
		} else if err := h.AuthService.RegisterUser(req); err != nil {
			globalErrors = append(globalErrors, "User already registered")
		} else {
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}
	}

	renderRegisterForm(c, req, fieldErrors, globalErrors)
}

func (h *AuthHandler) Login(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") == "application/json" {
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

		c.SetCookie("token", token, 3600*24, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	var req dto.LoginRequest
	fieldErrors := map[string]string{}
	var globalError string

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&req); err != nil {
			globalError = "Неверная форма"
		} else if err := h.Validator.Struct(req); err != nil {
			parsed := utils.ParseValidationErrors(err)
			for _, fe := range parsed.FieldErrors {
				fieldErrors[fe.Field] = fe.Message
			}
		} else if token, err := h.AuthService.LoginUser(req); err != nil {
			globalError = "Неверные данные"
		} else {
			c.SetCookie("token", token, 3600*24, "/", "", false, true)

			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"form":            req,
		"fieldErrors":     fieldErrors,
		"error":           globalError,
		"isAuthenticated": c.GetBool("isAuthenticated"),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.Set("isAuthenticated", false)
	c.Set("user_id", nil)
	c.Set("role", nil)
	c.Set("phone", nil)

	c.Redirect(http.StatusFound, "/auth/login")
}

func renderRegisterForm(c *gin.Context, data dto.RegisterRequest, fieldErrors map[string]string, globalErrors []string) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"form":            data,
		"fieldErrors":     fieldErrors,
		"errors":          globalErrors,
		"isAuthenticated": c.GetBool("isAuthenticated"),
	})
}
