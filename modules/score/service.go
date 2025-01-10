package score

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
	CreateScoreService(ctx *gin.Context)
	GetAllScoreService(ctx *gin.Context)
	NormalizeScoreService(ctx *gin.Context)
	DeleteAllScoreService(ctx *gin.Context)
}

type scoreService struct {
	repository         Repository
	productRepository  product.Repository
	criteriaRepository criteria.Repository
}

func NewScoreService(repo Repository, productRepo product.Repository, criteriaRepo criteria.Repository) Service {
	return &scoreService{
		repository:         repo,
		productRepository:  productRepo,
		criteriaRepository: criteriaRepo,
	}
}

func (service *scoreService) GetAllScoreService(ctx *gin.Context) {
	values, err := service.repository.GetAllScoreRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []Score
	for _, value := range values {
		result = append(result, Score{
			ID:         value.ID,
			Score:      value.Score,
			ProductID:  value.ProductID,
			CriteriaID: value.CriteriaID,
		})
	}

	if result == nil {
		response := map[string]string{"message": "data nilai masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	} else {
		helpers.ResponseJSON(ctx, http.StatusOK, result)
	}
}

func (service *scoreService) CreateScoreService(ctx *gin.Context) {
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	productList, err := service.productRepository.GetAllProductRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	for _, kriteria := range criteriaList {
		for _, produk := range productList {
			var nilai float64
			purchaseCost := float64(produk.PurchaseCost)
			priceSale := float64(produk.PriceSale)
			profit := float64(produk.Profit)

			switch strings.ToLower(kriteria.Name) {
			case strings.ToLower("Return On Investment"):
				if produk.PurchaseCost != 0 {
					nilai = profit / float64(produk.Stock) * 100
				}
			case strings.ToLower("Profit Margin"):
				if produk.PriceSale != 0 {
					nilai = profit / priceSale * 100
				}
			case strings.ToLower("Rasio Efisiensi"):
				if produk.Profit != 0 {
					nilai = purchaseCost / (float64(produk.Stock) * priceSale) * 100
				}
			default:
				response := map[string]string{"error": "kriteria tidak dikenali"}
				helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
				return
			}

			newScore := Score{
				ProductID:  produk.ID,
				CriteriaID: kriteria.ID,
				Score:      nilai,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			err = service.repository.CreateScoreRepository(&newScore)
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

func (service *scoreService) NormalizeScoreService(ctx *gin.Context) {

}

func (service *scoreService) DeleteAllScoreService(ctx *gin.Context) {
	err := service.repository.DeleteAllScoresRepository()
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "semua data nilai berhasil dihapus"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
