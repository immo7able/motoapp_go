package handler

import (
	"golang.org/x/crypto/bcrypt"
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
	if roleValue, exists := c.Get("role"); exists {
		if roleValue != nil {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}
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
	if roleValue, exists := c.Get("role"); exists {
		if roleValue != nil {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}
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
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "login.html", gin.H{
		"form":        req,
		"fieldErrors": fieldErrors,
		"error":       globalError,
		"role":        roleValue,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.Set("role", nil)
	c.Set("user_id", nil)
	c.Set("role", nil)
	c.Set("phone", nil)

	c.Redirect(http.StatusFound, "/auth/login")
}

func renderRegisterForm(c *gin.Context, data dto.RegisterRequest, fieldErrors map[string]string, globalErrors []string) {
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "register.html", gin.H{
		"form":        data,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"role":        roleValue,
	})
}

func (h *AuthHandler) ProfilePage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := h.AuthService.GetUserByID(userID.(uint))
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при загрузке профиля")
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User": user,
		"role": user.Role,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		Login       string `form:"login"`
		OldPassword string `form:"old_password"`
		NewPassword string `form:"new_password"`
	}

	var (
		successMsg string
		errorMsg   string
	)

	if err := c.ShouldBind(&input); err != nil {
		errorMsg = "Неверный ввод"
		h.renderProfileForm(c, userID.(uint), input.Login, errorMsg, successMsg)
		return
	}

	user, err := h.AuthService.GetUserByID(userID.(uint))
	if err != nil {
		errorMsg = "Пользователь не найден"
		h.renderProfileForm(c, userID.(uint), input.Login, errorMsg, successMsg)
		return
	}

	if !CheckPassword(input.OldPassword, user.Password) {
		errorMsg = "Неверный старый пароль"
		h.renderProfileForm(c, userID.(uint), input.Login, errorMsg, successMsg)
		return
	}

	user.Login = input.Login
	if input.NewPassword != "" {
		hashed, err := HashPassword(input.NewPassword)
		if err != nil {
			errorMsg = "Ошибка при хешировании пароля"
			h.renderProfileForm(c, userID.(uint), input.Login, errorMsg, successMsg)
			return
		}
		user.Password = hashed
	}

	if err := h.AuthService.UpdateUser(user); err != nil {
		errorMsg = "Ошибка обновления данных"
		h.renderProfileForm(c, userID.(uint), input.Login, errorMsg, successMsg)
		return
	}

	successMsg = "Профиль успешно обновлен"
	h.renderProfileForm(c, userID.(uint), user.Login, "", successMsg)
}

func (h *AuthHandler) renderProfileForm(c *gin.Context, userID uint, login string, errorMsg, successMsg string) {
	user, _ := h.AuthService.GetUserByID(userID)
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"User":       user,
		"role":       user.Role,
		"error":      errorMsg,
		"success":    successMsg,
		"loginValue": login,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
