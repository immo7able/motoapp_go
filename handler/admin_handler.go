package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/domain/dto"
	"motorcycleApp/service"
	"motorcycleApp/utils"
)

type AdminHandler struct {
	Service   *service.AdminService
	Validator *validator.Validate
}

func (h *AdminHandler) CreateBrand(c *gin.Context) {
	var req dto.CreateBrandRequest

	if c.Request.Header.Get("Content-Type") == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.Error{Code: http.StatusBadRequest, Message: "Invalid request"})
			return
		}
		if err := h.Validator.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ParseValidationErrors(err))
			return
		}

		if err := h.Service.CreateBrand(req.Name); err != nil {
			c.JSON(http.StatusInternalServerError, dto.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Brand created successfully"})
		return
	}

	fieldErrors := map[string]string{}
	var globalErrors []string

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&req); err != nil {
			globalErrors = append(globalErrors, "Неверная форма")
		} else if err := h.Validator.Struct(req); err != nil {
			for _, fe := range utils.ParseValidationErrors(err).FieldErrors {
				fieldErrors[fe.Field] = fe.Message
			}
		} else if err := h.Service.CreateBrand(req.Name); err != nil {
			globalErrors = append(globalErrors, "Ошибка при создании бренда")
		} else {
			c.Redirect(http.StatusFound, "/admin/brands")
			return
		}
	}
	brands, _ := h.Service.GetAllBrands()
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "brands.html", gin.H{
		"form":        req,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"brands":      brands,
		"role":        roleValue,
	})
}

func (h *AdminHandler) CreateModel(c *gin.Context) {
	var req dto.CreateModelRequest

	if c.Request.Header.Get("Content-Type") == "application/json" {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, dto.Error{Code: http.StatusBadRequest, Message: "Invalid request"})
			return
		}
		if err := h.Validator.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, utils.ParseValidationErrors(err))
			return
		}

		if err := h.Service.CreateModel(req.BrandID, req.Name); err != nil {
			c.JSON(http.StatusInternalServerError, dto.Error{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Model created successfully"})
		return
	}

	fieldErrors := map[string]string{}
	var globalErrors []string

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&req); err != nil {
			globalErrors = append(globalErrors, "Неверная форма")
		} else if err := h.Validator.Struct(req); err != nil {
			for _, fe := range utils.ParseValidationErrors(err).FieldErrors {
				fieldErrors[fe.Field] = fe.Message
			}
		} else if err := h.Service.CreateModel(req.BrandID, req.Name); err != nil {
			globalErrors = append(globalErrors, "Ошибка при создании модели")
		} else {
			c.Redirect(http.StatusFound, "/admin/models")
			return
		}
	}
	models, _ := h.Service.GetAllModels()
	brands, _ := h.Service.GetAllBrands()
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "models.html", gin.H{
		"form":        req,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"models":      models,
		"brands":      brands,
		"role":        roleValue,
	})
}

func (h *AdminHandler) UpdateBrand(c *gin.Context) {
	id := c.Param("id")
	idUint, _ := strconv.ParseUint(id, 10, 64)
	var req dto.UpdateBrandRequest

	fieldErrors := map[string]string{}
	var globalErrors []string

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&req); err != nil || idUint == 0 {
			globalErrors = append(globalErrors, "Неверная форма")
		} else if err := h.Validator.Struct(req); err != nil {
			for _, fe := range utils.ParseValidationErrors(err).FieldErrors {
				fieldErrors[fe.Field] = fe.Message
			}
		} else if err := h.Service.UpdateBrand(uint(idUint), req.Name); err != nil {
			if err.Error() == "duplicate_name" {
				globalErrors = append(globalErrors, "Бренд с таким названием уже существует")
			} else {
				globalErrors = append(globalErrors, "Ошибка при обновлении бренда")
			}
		} else {
			c.Redirect(http.StatusFound, "/admin/brands")
			return
		}
	}

	brands, _ := h.Service.GetAllBrands()
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "brands.html", gin.H{
		"form":        req,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"brands":      brands,
		"role":        roleValue,
	})
}

func (h *AdminHandler) DeleteBrand(c *gin.Context) {
	id := c.Param("id")
	idUint, _ := strconv.ParseUint(id, 10, 64)

	var globalErrors []string

	if err := h.Service.DeleteBrand(uint(idUint)); err != nil {
		globalErrors = append(globalErrors, "Ошибка при удалении бренда")
		brands, _ := h.Service.GetAllBrands()
		roleValue, _ := c.Get("role")
		c.HTML(http.StatusOK, "brands.html", gin.H{
			"errors":      globalErrors,
			"brands":      brands,
			"role":        roleValue,
			"form":        dto.CreateBrandRequest{},
			"fieldErrors": map[string]string{},
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/brands")
}

func (h *AdminHandler) UpdateModel(c *gin.Context) {
	id := c.Param("id")
	idUint, _ := strconv.ParseUint(id, 10, 64)
	var req dto.UpdateModelRequest

	fieldErrors := map[string]string{}
	var globalErrors []string

	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBind(&req); err != nil || idUint == 0 {
			globalErrors = append(globalErrors, "Неверная форма")
		} else if err := h.Validator.Struct(req); err != nil {
			for _, fe := range utils.ParseValidationErrors(err).FieldErrors {
				fieldErrors[fe.Field] = fe.Message
			}
		} else if err := h.Service.UpdateModel(uint(idUint), req.Name, req.BrandID); err != nil {
			if err.Error() == "duplicate_name" {
				globalErrors = append(globalErrors, "Модель с таким названием уже существует")
			} else {
				globalErrors = append(globalErrors, "Ошибка при обновлении модели")
			}
		} else {
			c.Redirect(http.StatusFound, "/admin/models")
			return
		}
	}

	models, _ := h.Service.GetAllModels()
	brands, _ := h.Service.GetAllBrands()
	roleValue, _ := c.Get("role")

	c.HTML(http.StatusOK, "models.html", gin.H{
		"form":        req,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"models":      models,
		"brands":      brands,
		"role":        roleValue,
	})
}

func (h *AdminHandler) DeleteModel(c *gin.Context) {
	id := c.Param("id")
	idUint, _ := strconv.ParseUint(id, 10, 64)

	var globalErrors []string

	if err := h.Service.DeleteModel(uint(idUint)); err != nil {
		globalErrors = append(globalErrors, "Ошибка при удалении модели")
		models, _ := h.Service.GetAllModels()
		brands, _ := h.Service.GetAllBrands()
		roleValue, _ := c.Get("role")

		c.HTML(http.StatusOK, "models.html", gin.H{
			"errors":      globalErrors,
			"models":      models,
			"brands":      brands,
			"role":        roleValue,
			"form":        dto.CreateModelRequest{},
			"fieldErrors": map[string]string{},
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/models")
}
