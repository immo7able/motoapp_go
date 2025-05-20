package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"motorcycleApp/domain/dto"
	"motorcycleApp/service"
	"motorcycleApp/utils"
	"net/http"
	"os"
	"path/filepath"
)

type MotorcycleHandler struct {
	Service   *service.MotorcycleService
	Validator *validator.Validate
}

func (h *MotorcycleHandler) ShowCreatePage(c *gin.Context) {
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "add_motorcycle.html", gin.H{
		"role": roleValue,
	})
}

func (h *MotorcycleHandler) AddMotorcycle(c *gin.Context) {
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

			form, err := c.MultipartForm()
			if err != nil {
				globalErrors = append(globalErrors, "Не удалось прочитать файлы")
				renderMotorcycleForm(c, req, fieldErrors, globalErrors)
				return
			}
			files := form.File["Images"]
			if len(files) > 10 {
				fieldErrors["Images"] = "Максимум 10 изображений"
				renderMotorcycleForm(c, req, fieldErrors, globalErrors)
				return
			}

			var imagePaths []string
			for _, file := range files {
				path, err := saveUploadedFile(file)
				if err != nil {
					globalErrors = append(globalErrors, fmt.Sprintf("Ошибка при сохранении файла: %s", file.Filename))
					renderMotorcycleForm(c, req, fieldErrors, globalErrors)
					return
				}
				imagePaths = append(imagePaths, path)
			}

			err = h.Service.AddMotorcycle(req, userIDRaw.(uint), phoneRaw.(string), imagePaths)
			if err != nil {
				globalErrors = append(globalErrors, "Не удалось добавить мотоцикл")
			} else {
				c.Redirect(http.StatusSeeOther, "/motorcycles/my")
				return
			}
		}
	}

	renderMotorcycleForm(c, req, fieldErrors, globalErrors)
}

func saveUploadedFile(file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	savePath := filepath.Join("uploads", filename)
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return "", err
	}

	if err := saveFile(file, savePath); err != nil {
		return "", err
	}
	webPath := filepath.ToSlash(savePath)
	return webPath, nil
}

func saveFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	_, err = io.Copy(out, src)
	return err
}

func (h *MotorcycleHandler) GetAllMotorcycles(c *gin.Context) {
	motos, err := h.Service.GetAllMotorcycles()
	roleValue, _ := c.Get("role")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "all_motorcycles.html", gin.H{
			"error": "Failed to load motorcycles",
			"role":  roleValue,
		})
		return
	}

	c.HTML(http.StatusOK, "all_motorcycles.html", gin.H{
		"motorcycles": motos,
		"role":        roleValue,
	})
}

func (h *MotorcycleHandler) GetUserMotorcycles(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	roleValue, _ := c.Get("role")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	userID := userIDRaw.(uint)

	motos, err := h.Service.GetUserMotorcycles(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "my_motorcycles.html", gin.H{
			"error": "Failed to load your motorcycles",
			"role":  roleValue,
		})
		return
	}
	c.HTML(http.StatusOK, "my_motorcycles.html", gin.H{
		"motorcycles": motos,
		"role":        roleValue,
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

	c.Redirect(http.StatusSeeOther, "/motorcycles/my")
}

func renderMotorcycleForm(c *gin.Context, data dto.MotorcycleAddRequest, fieldErrors map[string]string, globalErrors []string) {
	roleValue, _ := c.Get("role")
	c.HTML(http.StatusOK, "add_motorcycle.html", gin.H{
		"form":        data,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"role":        roleValue,
	})
}
