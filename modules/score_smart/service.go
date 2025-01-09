package score_smart

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
	CreateScoreSmartService(ctx *gin.Context)
	GetAllScoreSmartService(ctx *gin.Context)
	NormalizeScoreSmartService(ctx *gin.Context)
	DeleteAllScoreSmartService(ctx *gin.Context)
}

type ScoreSmartService struct {
	repository         Repository
	productRepository  product.Repository
	criteriaRepository criteria.Repository
}

func NewScoreSmartService(repo Repository, productRepo product.Repository, criteriaRepo criteria.Repository) Service {
	return &ScoreSmartService{
		repository:         repo,
		productRepository:  productRepo,
		criteriaRepository: criteriaRepo,
	}
}

func (service *ScoreSmartService) GetAllScoreSmartService(ctx *gin.Context) {
	values, err := service.repository.GetAllScoreSmartRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []ScoreSmart
	for _, value := range values {
		result = append(result, ScoreSmart{
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

func (service *ScoreSmartService) CreateScoreSmartService(ctx *gin.Context) {
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

			newScoreSmart := ScoreSmart{
				ProductID:  produk.ID,
				CriteriaID: kriteria.ID,
				Score:      nilai,
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			}

			err = service.repository.CreateScoreSmartRepository(&newScoreSmart)
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

func (service *ScoreSmartService) NormalizeScoreSmartService(ctx *gin.Context) {

}

func (service *ScoreSmartService) DeleteAllScoreSmartService(ctx *gin.Context) {
	err := service.repository.DeleteAllScoreSmartsRepository()
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "semua data nilai berhasil dihapus"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
