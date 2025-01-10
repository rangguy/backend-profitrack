package category

import (
	"backend-profitrack/helpers"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreateCategoryService(ctx *gin.Context)
	GetAllCategoryService(ctx *gin.Context)
	GetCategoryByIdService(ctx *gin.Context)
	UpdateCategoryService(ctx *gin.Context)
	DeleteCategoryService(ctx *gin.Context)
}

type categoryService struct {
	repository Repository
}

func NewCategoryService(repo Repository) Service {
	return &categoryService{
		repository: repo,
	}
}

func (service *categoryService) GetAllCategoryService(ctx *gin.Context) {
	categories, err := service.repository.GetAllCategoryRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []map[string]interface{}
	for _, category := range categories {
		result = append(result, map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	if result == nil {
		response := map[string]string{"message": "Data kategori masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	} else {
		helpers.ResponseJSON(ctx, http.StatusOK, result)
	}
}

func (service *categoryService) CreateCategoryService(ctx *gin.Context) {
	var newCategory Category
	if err := ctx.ShouldBindJSON(&newCategory); err != nil {
		helpers.ResponseJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	newCategory.Name = strings.TrimSpace(newCategory.Name)

	if newCategory.Name == "" {
		response := map[string]string{"error": "Tolong masukan nama kategori"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if match, _ := regexp.MatchString(`[a-zA-Z]`, newCategory.Name); !match {
		response := map[string]string{"error": "Nama kategori harus mengandung setidaknya satu huruf"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	newCategory.CreatedAt = time.Now()
	newCategory.UpdatedAt = time.Now()

	err := service.repository.CreateCategoryRepository(&newCategory)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key score violates unique constraint \"uni_categories_name\"") {
			response := map[string]string{"error": "Nama kategori sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}
		response := map[string]string{"error": "gagal menambahkan data kategori"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, newCategory)
}

func (service *categoryService) GetCategoryByIdService(ctx *gin.Context) {
	var (
		category Category
		id       = ctx.Param("id")
	)

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	category, err = service.repository.GetCategoryByIdRepository(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Kategori dengan ID:%d tidak ditemukan", categoryID)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (service *categoryService) UpdateCategoryService(ctx *gin.Context) {
	var category Category

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if err = ctx.ShouldBindJSON(&category); err != nil {
		response := map[string]string{"error": "failed to read json"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingCategory, err := service.repository.GetCategoryByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Kategori dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if category.Name == "" {
		response := map[string]string{"error": "Masukkan nama kategori"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if strings.ToLower(existingCategory.Name) == strings.ToLower(category.Name) {
		response := map[string]string{"error": "Masukkan data nama kategori yang baru"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingCategory.Name = category.Name
	existingCategory.UpdatedAt = time.Now()

	err = service.repository.UpdateCategoryRepository(existingCategory)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key score violates unique constraint \"uni_categories_name\"") {
			response := map[string]string{"error": "Nama kategori sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}
		response := map[string]string{"error": "gagal mengubah data kategori"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Data kategori berhasil diperbarui"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *categoryService) DeleteCategoryService(ctx *gin.Context) {
	var category Category
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	category, err = service.repository.GetCategoryByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Kategori dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	category.ID = id
	err = service.repository.DeleteCategoryRepository(category)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			response := map[string]string{"error": "Kategori tidak dapat dihapus karena masih digunakan."}
			helpers.ResponseJSON(ctx, http.StatusConflict, response)
			return
		}

		response := map[string]string{"error": "Gagal menghapus data kategori"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": fmt.Sprintf("Kategori dengan ID:%d berhasil dihapus", id)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
