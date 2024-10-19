package product

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
	CreateProductService(ctx *gin.Context)
	GetAllProductService(ctx *gin.Context)
	GetProductByIdService(ctx *gin.Context)
	UpdateProductService(ctx *gin.Context)
	DeleteProductService(ctx *gin.Context)
}

type productService struct {
	repository Repository
}

func NewProductService(repo Repository) Service {
	return &productService{
		repository: repo,
	}
}

func (service *productService) GetAllProductService(ctx *gin.Context) {
	products, err := service.repository.GetAllProductRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []ResponseProduct
	for _, product := range products {
		result = append(result, ResponseProduct{
			ID:           product.ID,
			Name:         product.Name,
			NetProfit:    product.NetProfit,
			GrossProfit:  product.GrossProfit,
			GrossSale:    product.GrossSale,
			PurchaseCost: product.PurchaseCost,
			InitialStock: product.InitialStock,
			FinalStock:   product.FinalStock,
			CategoryID:   product.CategoryID,
		})
	}

	if result == nil {
		response := map[string]string{"message": "data produk masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	} else {
		helpers.ResponseJSON(ctx, http.StatusOK, result)
	}
}

func (service *productService) CreateProductService(ctx *gin.Context) {
	var newProduct Product
	response := map[string]string{}

	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		helpers.ResponseJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if newProduct.Name == "" ||
		newProduct.NetProfit == 0 ||
		newProduct.GrossProfit == 0 ||
		newProduct.GrossSale == 0 ||
		newProduct.PurchaseCost == 0 ||
		newProduct.InitialStock == 0 ||
		newProduct.FinalStock == 0 ||
		newProduct.CategoryID == 0 {
		response = map[string]string{"error": "semua field harus diisi dengan nilai yang valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	newProduct.CreatedAt = time.Now()
	newProduct.ModifiedAt = time.Now()

	err := service.repository.CreateProductRepository(&newProduct)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_products_name\"") {
			response = map[string]string{"message": "nama produk sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}
		response = map[string]string{"error": "gagal menambahkan data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, newProduct)
}

func (service *productService) GetProductByIdService(ctx *gin.Context) {
	var (
		product  Product
		id       = ctx.Param("id")
		response = map[string]string{}
	)

	productID, err := strconv.Atoi(id)
	if err != nil {
		response = map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	product, err = service.repository.GetProductByIdRepository(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = map[string]string{"error": fmt.Sprintf("Produk dengan ID:%d tidak ditemukan", productID)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response = map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (service *productService) UpdateProductService(ctx *gin.Context) {
	var product Product
	response := map[string]string{}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response = map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if err = ctx.ShouldBindJSON(&product); err != nil {
		response = map[string]string{"message": "failed to read json"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingProduct, err := service.repository.GetProductByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = map[string]string{"error": fmt.Sprintf("Produk dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response = map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if product.Name == "" ||
		product.NetProfit == 0 ||
		product.GrossProfit == 0 ||
		product.GrossSale == 0 ||
		product.PurchaseCost == 0 ||
		product.InitialStock == 0 ||
		product.FinalStock == 0 ||
		product.CategoryID == 0 {
		response = map[string]string{"error": "semua field harus diisi dengan nilai yang valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if existingProduct.Name == product.Name && existingProduct.NetProfit == product.NetProfit && existingProduct.GrossProfit == product.GrossProfit && existingProduct.GrossSale == product.GrossSale && existingProduct.PurchaseCost == product.PurchaseCost && existingProduct.InitialStock == product.InitialStock && existingProduct.FinalStock == product.FinalStock && existingProduct.CategoryID == product.CategoryID {
		response = map[string]string{"error": "masukkan minimal satu data yang baru"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingProduct.Name = product.Name
	existingProduct.NetProfit = product.NetProfit
	existingProduct.GrossProfit = product.GrossProfit
	existingProduct.GrossSale = product.GrossSale
	existingProduct.PurchaseCost = product.PurchaseCost
	existingProduct.InitialStock = product.InitialStock
	existingProduct.FinalStock = product.FinalStock
	existingProduct.CategoryID = product.CategoryID
	existingProduct.ModifiedAt = time.Now()

	err = service.repository.UpdateProductRepository(existingProduct)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_products_name\"") {
			response = map[string]string{"message": "nama produk sudah ada"}
			helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
			return
		}

		response = map[string]string{"error": "gagal mengubah data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response = map[string]string{"message": "Data produk berhasil diperbarui"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *productService) DeleteProductService(ctx *gin.Context) {
	var product Product
	response := map[string]string{}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response = map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	product, err = service.repository.GetProductByIdRepository(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response = map[string]string{"error": fmt.Sprintf("Produk dengan ID:%d tidak ditemukan", id)}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response = map[string]string{"error": err.Error()}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	product.ID = id
	err = service.repository.DeleteProductRepository(product)
	if err != nil {
		response = map[string]string{"error": "gagal menghapus data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response = map[string]string{"message": fmt.Sprintf("Produk dengan ID:%d berhasil dihapus", id)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
