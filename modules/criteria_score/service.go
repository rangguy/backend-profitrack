package criteria_score

import (
	"backend-profitrack/helpers"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"github.com/gin-gonic/gin"
	"net/http"
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
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(criteriaList) == 0 {
		response := map[string]string{"error": "tidak ada data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	productList, err := service.productRepository.GetAllProductRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(productList) == 0 {
		response := map[string]string{"error": "tidak ada data produk"}
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
				response := map[string]string{"error": "gagal menyimpan data nilai untuk produk " + produk.Name}
				helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
				return
			}
		}
	}

	response := map[string]string{"message": "nilai produk berhasil dihitung untuk semua kriteria dan disimpan"}
	helpers.ResponseJSON(ctx, http.StatusCreated, response)
}

func (service *criteriaScoreService) GetAllCriteriaScoreService(ctx *gin.Context) {
	values, err := service.repository.GetAllCriteriaScoreRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
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
		response := map[string]string{"message": "data nilai kriteria masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	}

	helpers.ResponseJSON(ctx, http.StatusOK, result)
}

func (service *criteriaScoreService) UpdateCriteriaScoreService(ctx *gin.Context) {
	return
}

func (service *criteriaScoreService) DeleteCriteriaScoreService(ctx *gin.Context) {
	return
}
