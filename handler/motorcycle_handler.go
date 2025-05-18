package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/domain/dto"
	"motorcycleApp/service"
	"motorcycleApp/utils"
	"net/http"
)

type MotorcycleHandler struct {
	Service   *service.MotorcycleService
	Validator *validator.Validate
}

func (h *MotorcycleHandler) ShowCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "add_motorcycle.html", gin.H{
		"isAuthenticated": c.GetBool("isAuthenticated"),
	})
}

func (h *MotorcycleHandler) AddMotorcycle(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") == "application/json" {
		var req dto.MotorcycleAddRequest
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

		userIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.Error{Code: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}
		phoneRaw, _ := c.Get("phone")

		err := h.Service.AddMotorcycle(req, userIDRaw.(uint), phoneRaw.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.Error{Code: http.StatusInternalServerError, Message: "Failed to add motorcycle"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Motorcycle added"})
		return
	}

	var req dto.MotorcycleAddRequest
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
		} else {
			userIDRaw, exists := c.Get("user_id")
			if !exists {
				c.Redirect(http.StatusSeeOther, "/login")
				return
			}
			phoneRaw, _ := c.Get("phone")

			if err := h.Service.AddMotorcycle(req, userIDRaw.(uint), phoneRaw.(string)); err != nil {
				globalErrors = append(globalErrors, "Не удалось добавить мотоцикл")
			} else {
				c.Redirect(http.StatusSeeOther, "/motorcycle/my")
				return
			}
		}
	}

	renderMotorcycleForm(c, req, fieldErrors, globalErrors)
}

func (h *MotorcycleHandler) GetAllMotorcycles(c *gin.Context) {
	motos, err := h.Service.GetAllMotorcycles()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "all_motorcycles.html", gin.H{
			"error":           "Failed to load motorcycles",
			"isAuthenticated": c.GetBool("isAuthenticated"),
		})
		return
	}
	c.HTML(http.StatusOK, "all_motorcycles.html", gin.H{
		"motorcycles":     motos,
		"isAuthenticated": c.GetBool("isAuthenticated"),
	})
}

func (h *MotorcycleHandler) GetUserMotorcycles(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	userID := userIDRaw.(uint)

	motos, err := h.Service.GetUserMotorcycles(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "my_motorcycles.html", gin.H{
			"error":           "Failed to load your motorcycles",
			"isAuthenticated": c.GetBool("isAuthenticated"),
		})
		return
	}
	c.HTML(http.StatusOK, "my_motorcycles.html", gin.H{
		"motorcycles":     motos,
		"isAuthenticated": c.GetBool("isAuthenticated"),
	})
}

func (h *MotorcycleHandler) DeleteMotorcycle(c *gin.Context) {
	id := c.Param("id")
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if err := h.Service.DeleteMotorcycle(id, userIDRaw.(uint)); err != nil {
		c.String(http.StatusInternalServerError, "Не удалось удалить объявление")
		return
	}

	c.Redirect(http.StatusSeeOther, "/motorcycle/my")
}

func renderMotorcycleForm(c *gin.Context, data dto.MotorcycleAddRequest, fieldErrors map[string]string, globalErrors []string) {
	c.HTML(http.StatusOK, "add_motorcycle.html", gin.H{
		"form":            data,
		"fieldErrors":     fieldErrors,
		"errors":          globalErrors,
		"isAuthenticated": c.GetBool("isAuthenticated"),
	})
}
