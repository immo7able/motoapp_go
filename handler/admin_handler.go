package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"motorcycleApp/service"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	Service   *service.AdminService
	Validator *validator.Validate
}

func (h *AdminHandler) BrandsPage(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		brand := c.PostForm("brand")
		if brand != "" {
			_ = h.Service.CreateBrand(brand)
		}
	}

	brands, _ := h.Service.GetAllBrands()

	c.HTML(http.StatusOK, "brands.html", gin.H{
		"brands": brands,
	})
}

func (h *AdminHandler) ModelsPage(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		brandIDStr := c.PostForm("brand_id")
		modelName := c.PostForm("model")

		if brandID, err := strconv.ParseUint(brandIDStr, 10, 64); err == nil && modelName != "" {
			_ = h.Service.CreateModel(uint(brandID), modelName)
		}
	}

	models, _ := h.Service.GetAllModels()
	brands, _ := h.Service.GetAllBrands()

	c.HTML(http.StatusOK, "models.html", gin.H{
		"models": models,
		"brands": brands,
	})
}
