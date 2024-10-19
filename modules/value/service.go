package value

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"profitrack/helpers"
	"profitrack/modules/criteria"
	"profitrack/modules/product"
	"strings"
	"time"
)

type Service interface {
	CreateValueService(ctx *gin.Context)
	GetAllValueService(ctx *gin.Context)
	DeleteAllValuesService(ctx *gin.Context)
}

type valueService struct {
	repository         Repository
	productRepository  product.Repository
	criteriaRepository criteria.Repository
}

func NewValueService(repo Repository, productRepo product.Repository, criteriaRepo criteria.Repository) Service {
	return &valueService{
		repository:         repo,
		productRepository:  productRepo,
		criteriaRepository: criteriaRepo,
	}
}

func (service *valueService) GetAllValueService(ctx *gin.Context) {
	values, err := service.repository.GetAllValueRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []Value
	for _, value := range values {
		result = append(result, Value{
			ID:         value.ID,
			Value:      value.Value,
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

func (service *valueService) CreateValueService(ctx *gin.Context) {
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
			netProfit := float64(produk.NetProfit)
			purchaseCost := float64(produk.PurchaseCost)
			grossProfit := float64(produk.GrossProfit)
			grossSale := float64(produk.GrossSale)

			switch strings.ToLower(kriteria.Name) {
			case strings.ToLower("Return On Investment"):
				if produk.PurchaseCost != 0 {
					nilai = grossProfit / purchaseCost * 100

				}
			case strings.ToLower("Net Profit Margin"):
				if produk.GrossSale != 0 {
					nilai = netProfit / grossSale * 100

				}
			case strings.ToLower("Gross Profit Margin"):
				if produk.GrossSale != 0 {
					nilai = grossProfit / grossSale * 100

				}
			default:
				response := map[string]string{"error": "kriteria tidak dikenali"}
				helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
				return
			}

			newValue := Value{
				ProductID:  produk.ID,
				CriteriaID: kriteria.ID,
				Value:      nilai,
				CreatedAt:  time.Now(),
				ModifiedAt: time.Now(),
			}

			err = service.repository.CreateValueRepository(&newValue)
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

func (service *valueService) DeleteAllValuesService(ctx *gin.Context) {
	err := service.repository.DeleteAllValuesRepository()
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "semua data nilai berhasil dihapus"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
