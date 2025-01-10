package method

import (
	"backend-profitrack/helpers"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreateMethodService(ctx *gin.Context)
	GetAllMethodService(ctx *gin.Context)
	GetMethodByIdService(ctx *gin.Context)
	UpdateMethodService(ctx *gin.Context)
	DeleteMethodService(ctx *gin.Context)
}

type methodService struct {
	repository Repository
}

func NewMethodService(repo Repository) Service {
	return &methodService{
		repository: repo,
	}
}

func (service *methodService) GetAllMethodService(ctx *gin.Context) {
	methods, err := service.repository.GetAllMethodRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []ResponseMethod
	for _, method := range methods {
		result = append(result, ResponseMethod{
			ID:   method.ID,
			Name: method.Name,
		})
	}

	if result == nil {
		response := map[string]string{"message": "data metode masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	} else {
		helpers.ResponseJSON(ctx, http.StatusOK, result)
	}
}

func (service *methodService) CreateMethodService(ctx *gin.Context) {
	var newMethod Method
	if err := ctx.ShouldBindJSON(&newMethod); err != nil {
		helpers.ResponseJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if newMethod.Name == "" {
		response := map[string]string{"error": "inputan tidak valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	newMethod.CreatedAt = time.Now()
	newMethod.UpdatedAt = time.Now()

	err := service.repository.CreateMethodRepository(&newMethod)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key score violates unique constraint \"uni_method_name\"") {
			response := map[string]string{"error": "Nama metode sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}
		response := map[string]string{"error": "gagal menambahkan data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, newMethod)
}

func (service *methodService) GetMethodByIdService(ctx *gin.Context) {
	var (
		method Method
		id     = ctx.Param("id")
	)

	methodID, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	method, err = service.repository.GetMethodByIdRepository(methodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Metode dengan ID:%d tidak ditemukan", methodID)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, method)
}

func (service *methodService) UpdateMethodService(ctx *gin.Context) {
	var method Method

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if err = ctx.ShouldBindJSON(&method); err != nil {
		response := map[string]string{"error": "failed to read json"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingMethod, err := service.repository.GetMethodByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Metode dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if method.Name == "" {
		response := map[string]string{"error": "inputan tidak valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if existingMethod.Name == method.Name {
		response := map[string]string{"error": "Masukkan setidaknya satu data baru"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingMethod.Name = method.Name
	existingMethod.UpdatedAt = time.Now()

	err = service.repository.UpdateMethodRepository(existingMethod)
	if err != nil {
		response := map[string]string{"error": "Gagal mengubah data metode"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Data metode berhasil diperbarui"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *methodService) DeleteMethodService(ctx *gin.Context) {
	var method Method
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	method, err = service.repository.GetMethodByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Metode dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	method.ID = id
	err = service.repository.DeleteMethodRepository(method)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus data metode"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": fmt.Sprintf("Metode dengan ID:%d berhasil dihapus", id)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
