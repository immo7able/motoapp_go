package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"motorcycleApp/domain/dto"
	"motorcycleApp/domain/model"
	"motorcycleApp/service"
	"motorcycleApp/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MotorcycleHandler struct {
	Service      *service.MotorcycleService
	Validator    *validator.Validate
	AdminService *service.AdminService
}

func (h *MotorcycleHandler) ShowCreatePage(c *gin.Context) {
	roleValue, _ := c.Get("role")
	brands, _ := h.AdminService.GetAllBrands()
	models, _ := h.AdminService.GetAllModels()

	c.HTML(http.StatusOK, "add_motorcycle.html", gin.H{
		"role":   roleValue,
		"brands": brands,
		"models": models,
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
				renderMotorcycleForm(c, h, req, fieldErrors, globalErrors)
				return
			}
			files := form.File["Images"]
			if len(files) > 10 {
				fieldErrors["Images"] = "Максимум 10 изображений"
				renderMotorcycleForm(c, h, req, fieldErrors, globalErrors)
				return
			}

			var imagePaths []string
			for _, file := range files {
				path, err := saveUploadedFile(file)
				if err != nil {
					globalErrors = append(globalErrors, fmt.Sprintf("Ошибка при сохранении файла: %s", file.Filename))
					renderMotorcycleForm(c, h, req, fieldErrors, globalErrors)
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

	renderMotorcycleForm(c, h, req, fieldErrors, globalErrors)
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

func renderMotorcycleForm(c *gin.Context, h *MotorcycleHandler, data dto.MotorcycleAddRequest, fieldErrors map[string]string, globalErrors []string) {
	roleValue, _ := c.Get("role")
	brands, _ := h.AdminService.GetAllBrands()
	models, _ := h.AdminService.GetAllModels()
	c.HTML(http.StatusOK, "add_motorcycle.html", gin.H{
		"form":        data,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"brands":      brands,
		"models":      models,
		"role":        roleValue,
	})
}

func (h *MotorcycleHandler) ShowEditPage(c *gin.Context) {
	idParam := c.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid motorcycle ID"})
		return
	}
	userID, _ := c.Get("user_id")
	roleValue, _ := c.Get("role")

	moto, err := h.Service.GetMotorcycleByIDAndAuthor(uint(idUint), userID.(uint))
	if err != nil {
		c.HTML(http.StatusNotFound, "edit_motorcycle.html", gin.H{"errors": "Мотоцикл не найден"})
		return
	}

	brands, _ := h.AdminService.GetAllBrands()
	models, _ := h.AdminService.GetAllModels()

	c.HTML(http.StatusOK, "edit_motorcycle.html", gin.H{
		"form":        moto,
		"brands":      brands,
		"models":      models,
		"fieldErrors": map[string]string{},
		"role":        roleValue,
	})
}

func (h *MotorcycleHandler) EditMotorcycle(c *gin.Context) {
	var req dto.MotorcycleUpdateRequest
	fieldErrors := map[string]string{}
	var globalErrors []string

	if err := c.ShouldBind(&req); err != nil {
		globalErrors = append(globalErrors, "Неверная форма")
		renderEditForm(c, h, req, fieldErrors, globalErrors)
		return
	}
	if err := h.Validator.Struct(req); err != nil {
		parsed := utils.ParseValidationErrors(err)
		for _, fe := range parsed.FieldErrors {
			fieldErrors[fe.Field] = fe.Message
		}
		renderEditForm(c, h, req, fieldErrors, globalErrors)
		return
	}

	form, _ := c.MultipartForm()
	files := form.File["Images"]
	var imagePaths []string
	for _, file := range files {
		path, err := saveUploadedFile(file)
		if err != nil {
			globalErrors = append(globalErrors, "Ошибка при загрузке файла")
		}
		imagePaths = append(imagePaths, path)
	}

	userID, _ := c.Get("user_id")
	idParam := c.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 64)

	err = h.Service.UpdateMotorcycle(uint(idUint), req, userID.(uint), imagePaths)
	if err != nil {
		globalErrors = append(globalErrors, "Ошибка при обновлении")
		renderEditForm(c, h, req, fieldErrors, globalErrors)
		return
	}
	c.Redirect(http.StatusSeeOther, "/motorcycles/my")
}

func renderEditForm(c *gin.Context, h *MotorcycleHandler, req dto.MotorcycleUpdateRequest, fieldErrors map[string]string, globalErrors []string) {
	brands, _ := h.AdminService.GetAllBrands()
	models, _ := h.AdminService.GetAllModels()
	roleValue, _ := c.Get("role")

	c.HTML(http.StatusOK, "edit_motorcycle.html", gin.H{
		"form":        req,
		"brands":      brands,
		"fieldErrors": fieldErrors,
		"errors":      globalErrors,
		"models":      models,
		"role":        roleValue,
	})
}

func (h *MotorcycleHandler) ShowMotorcyclePage(c *gin.Context) {
	idParam := c.Param("id")
	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid motorcycle ID"})
		return
	}

	moto, err := h.Service.GetMotorcycleByID(uint(idUint))
	if err != nil {
		c.HTML(http.StatusNotFound, "motorcycle_detail.html", gin.H{
			"error": "Объявление не найдено",
		})
		return
	}

	roleValue, _ := c.Get("role")
	userId, _ := c.Get("user_id")
	comments, err := h.Service.GetCommentsByMotorcycleID(uint(idUint))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "motorcycle_detail.html", gin.H{
			"error": "Ошибка загрузки комментариев",
		})
		return
	}
	c.HTML(http.StatusOK, "motorcycle_detail.html", gin.H{
		"Ad":       moto,
		"role":     roleValue,
		"userID":   userId,
		"Comments": comments,
	})
}

func (h *MotorcycleHandler) AddComment(c *gin.Context) {
	motoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}
	userID := userIDVal.(uint)

	content := c.PostForm("content")
	if strings.TrimSpace(content) == "" {
		c.Redirect(http.StatusFound, fmt.Sprintf("/motorcycles/%d", motoID))
		return
	}

	comment := model.Comment{
		UserID:       userID,
		MotorcycleID: uint(motoID),
		Content:      content,
	}
	if err := h.Service.SaveComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения комментария"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/motorcycles/%d", motoID))
}

func (h *MotorcycleHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	roleAny, _ := c.Get("role")
	role := roleAny.(string)

	comment, err := h.Service.GetCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Комментарий не найден"})
		return
	}

	if comment.UserID != userId && role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет прав для удаления комментария"})
		return
	}

	err = h.Service.DeleteComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении комментария"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/motorcycles/"+strconv.Itoa(int(comment.MotorcycleID)))
}
