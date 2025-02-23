package criteria_score

import (
	"backend-profitrack/helpers"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
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
	CreateAllCriteriaScoreService(ctx *gin.Context)
	GetAllCriteriaScoreService(ctx *gin.Context)
	UpdateCriteriaScoreService(ctx *gin.Context)
	DeleteCriteriaScoreService(ctx *gin.Context)
}

type criteriaScoreService struct {
	repository         Repository
	criteriaRepository criteria.Repository
	productRepository  product.Repository
}

func NewCriteriaScoreService(repo Repository, criteriaRepo criteria.Repository, productRepo product.Repository) Service {
	return &criteriaScoreService{
		repo,
		criteriaRepo,
		productRepo,
	}
}

func (service *criteriaScoreService) CreateAllCriteriaScoreService(ctx *gin.Context) {
	var response map[string]string
	existingScores, err := service.repository.GetAllCriteriaScoreRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data nilai kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(existingScores) != 0 {
		response = map[string]string{"error": "data nilai kriteria telah ada"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(criteriaList) == 0 {
		response = map[string]string{"error": "tidak ada data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	productList, err := service.productRepository.GetAllProductRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(productList) == 0 {
		response = map[string]string{"error": "tidak ada data produk"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	for _, produk := range productList {
		for _, kriteria := range criteriaList {
			var nilai float64
			purchaseCost := float64(produk.PurchaseCost)
			priceSale := float64(produk.PriceSale)
			profit := float64(produk.Profit)
			stock := float64(produk.Stock)
			sold := float64(produk.Sold)

			switch strings.ToLower(kriteria.Name) {
			case strings.ToLower("Return On Investment"):
				if produk.PurchaseCost != 0 {
					nilai = (profit * sold) / (purchaseCost * stock)
				}
			case strings.ToLower("Net Profit Margin"):
				if produk.PriceSale != 0 {
					nilai = (profit * sold) / (priceSale * sold)
				}
			case strings.ToLower("Rasio Efisiensi"):
				if produk.Profit != 0 {
					nilai = (purchaseCost * sold) / (sold * priceSale)
				}
			default:
				response := map[string]string{"error": "kriteria tidak dikenali"}
				helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
				return
			}

			newScore := CriteriaScore{
				ProductID:  produk.ID,
				CriteriaID: kriteria.ID,
				Score:      nilai,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			err = service.repository.CreateCriteriaScoreRepository(&newScore)
			if err != nil {
				response = map[string]string{"error": "gagal menyimpan data nilai untuk produk " + produk.Name}
				helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
				return
			}
		}
	}

	response = map[string]string{"message": "nilai produk berhasil dihitung untuk semua kriteria dan disimpan"}
	helpers.ResponseJSON(ctx, http.StatusCreated, response)
}

func (service *criteriaScoreService) GetAllCriteriaScoreService(ctx *gin.Context) {
	var response map[string]string
	values, err := service.repository.GetAllCriteriaScoreRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data nilai kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	var result []CriteriaScore
	for _, value := range values {
		result = append(result, CriteriaScore{
			ID:         value.ID,
			Score:      value.Score,
			ProductID:  value.ProductID,
			CriteriaID: value.CriteriaID,
			CreatedAt:  value.CreatedAt,
			UpdatedAt:  value.UpdatedAt,
		})
	}

	if len(result) == 0 {
		response = map[string]string{"message": "data nilai kriteria masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	}

	helpers.ResponseJSON(ctx, http.StatusOK, result)
}

func (service *criteriaScoreService) UpdateCriteriaScoreService(ctx *gin.Context) {
	var response map[string]string
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	productList, err := service.productRepository.GetAllProductRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	existingScores, err := service.repository.GetAllCriteriaScoreRepository()
	if err != nil {
		response = map[string]string{"error": "gagal mengambil data nilai kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	scoreMap := make(map[int]map[int]bool) // productID -> criteriaID -> exists
	for _, score := range existingScores {
		if _, exists := scoreMap[score.ProductID]; !exists {
			scoreMap[score.ProductID] = make(map[int]bool)
		}
		scoreMap[score.ProductID][score.CriteriaID] = true
	}

	// Pastikan semua produk memiliki nilai, jika tidak maka buat nilai baru
	for _, product := range productList {
		for _, criteria := range criteriaList {
			if !scoreMap[product.ID][criteria.ID] {
				var nilai float64
				purchaseCost := float64(product.PurchaseCost)
				priceSale := float64(product.PriceSale)
				profit := float64(product.Profit)
				stock := float64(product.Stock)
				sold := float64(product.Sold)

				switch strings.ToLower(criteria.Name) {
				case "return on investment":
					if product.PurchaseCost != 0 {
						nilai = (profit * sold) / (purchaseCost * stock)
					}
				case "net profit margin":
					if product.PriceSale != 0 {
						nilai = (profit * sold) / (priceSale * sold)
					}
				case "rasio efisiensi":
					if product.Profit != 0 {
						nilai = (purchaseCost * sold) / (sold * priceSale)
					}
				}

				newScore := CriteriaScore{
					ProductID:  product.ID,
					CriteriaID: criteria.ID,
					Score:      nilai,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				err = service.repository.CreateCriteriaScoreRepository(&newScore)
				if err != nil {
					response = map[string]string{"error": "gagal menyimpan data nilai untuk produk " + product.Name}
					helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
					return
				}
			}
		}
	}

	// Update nilai yang sudah ada
	for _, score := range existingScores {
		var nilai float64
		produk, _ := service.productRepository.GetProductByIdRepository(score.ProductID)
		kriteria, _ := service.criteriaRepository.GetCriteriaByIdRepository(score.CriteriaID)

		switch strings.ToLower(kriteria.Name) {
		case "return on investment":
			if produk.PurchaseCost != 0 {
				nilai = (float64(produk.Profit) * float64(produk.Sold)) / (float64(produk.PurchaseCost) * float64(produk.Stock))
			}
		case "net profit margin":
			if produk.PriceSale != 0 {
				nilai = (float64(produk.Profit) * float64(produk.Sold)) / (float64(produk.PriceSale) * float64(produk.Sold))
			}
		case "rasio efisiensi":
			if produk.Profit != 0 {
				nilai = (float64(produk.PurchaseCost) * float64(produk.Sold)) / (float64(produk.Sold) * float64(produk.PriceSale))
			}
		}

		score.Score = nilai
		score.UpdatedAt = time.Now()
		err = service.repository.UpdateCriteriaScoreRepository(&score)
		if err != nil {
			response = map[string]string{"error": "gagal memperbarui data nilai"}
			helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
			return
		}
	}

	response = map[string]string{"message": "nilai produk berhasil diperbarui atau ditambahkan jika belum ada"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *criteriaScoreService) DeleteCriteriaScoreService(ctx *gin.Context) {
	var response map[string]string
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response = map[string]string{"error": "produk ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	_, err = service.repository.GetCriteriaScoreByProductIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = map[string]string{"error": fmt.Sprintf("Nilai Kriteria dengan produk ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response = map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	product, err := service.productRepository.GetProductByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = map[string]string{"error": fmt.Sprintf("Data produk ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response = map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	err = service.repository.DeleteCriteriaScoreRepository(id)
	if err != nil {
		response = map[string]string{"error": "gagal menghapus data nilai kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response = map[string]string{"message": fmt.Sprintf("data nilai kriteria produk:%s berhasil dihapus", product.Name)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
