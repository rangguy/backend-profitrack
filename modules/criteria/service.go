package criteria

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"profitrack/helpers"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreateCriteriaService(ctx *gin.Context)
	GetAllCriteriaService(ctx *gin.Context)
	GetCriteriaByIdService(ctx *gin.Context)
	UpdateCriteriaService(ctx *gin.Context)
	DeleteCriteriaService(ctx *gin.Context)
}

type criteriaService struct {
	repository Repository
}

func NewCriteriaService(repo Repository) Service {
	return &criteriaService{
		repository: repo,
	}
}

func (service *criteriaService) GetAllCriteriaService(ctx *gin.Context) {
	criterias, err := service.repository.GetAllCriteriaRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []ResponseCriteria
	for _, criteria := range criterias {
		result = append(result, ResponseCriteria{
			ID:     criteria.ID,
			Name:   criteria.Name,
			Weight: criteria.Weight,
			Type:   criteria.Name,
		})
	}

	if result == nil {
		response := map[string]string{"message": "data kriteria masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	} else {
		helpers.ResponseJSON(ctx, http.StatusOK, result)
	}
}

func (service *criteriaService) CreateCriteriaService(ctx *gin.Context) {
	var newCriteria Criteria
	if err := ctx.ShouldBindJSON(&newCriteria); err != nil {
		helpers.ResponseJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if newCriteria.Name == "" || newCriteria.Weight == 0 || newCriteria.Type == "" {
		response := map[string]string{"error": "inputan tidak valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	newCriteria.CreatedAt = time.Now()
	newCriteria.ModifiedAt = time.Now()

	err := service.repository.CreateCriteriaRepository(&newCriteria)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_criteria_name\"") {
			response := map[string]string{"error": "nama kriteria sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}
		response := map[string]string{"error": "gagal menambahkan data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, newCriteria)
}

func (service *criteriaService) GetCriteriaByIdService(ctx *gin.Context) {
	var (
		criteria Criteria
		id       = ctx.Param("id")
	)

	criteriaID, err := strconv.Atoi(id)
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	criteria, err = service.repository.GetCriteriaByIdRepository(criteriaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Kriteria dengan ID:%d tidak ditemukan", criteriaID)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (service *criteriaService) UpdateCriteriaService(ctx *gin.Context) {
	var criteria Criteria

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if err = ctx.ShouldBindJSON(&criteria); err != nil {
		response := map[string]string{"message": "failed to read json"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingCriteria, err := service.repository.GetCriteriaByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Kriteria dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if existingCriteria.Name == criteria.Name {
		response := map[string]string{"error": "masukkan data nama kriteria yang baru"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingCriteria.Name = criteria.Name
	existingCriteria.Weight = criteria.Weight
	existingCriteria.Type = criteria.Type
	existingCriteria.ModifiedAt = time.Now()

	err = service.repository.UpdateCriteriaRepository(existingCriteria)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_criteria_name\"") {
			response := map[string]string{"message": "nama kriteria sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}
		response := map[string]string{"error": "gagal mengubah data kategori"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Data kriteria berhasil diperbarui"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *criteriaService) DeleteCriteriaService(ctx *gin.Context) {
	var criteria Criteria
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	criteria, err = service.repository.GetCriteriaByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]string{"error": fmt.Sprintf("Kriteria dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	criteria.ID = id
	err = service.repository.DeleteCriteriaRepository(criteria)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"error": fmt.Sprintf("Kriteria dengan ID:%d berhasil dihapus", id)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
