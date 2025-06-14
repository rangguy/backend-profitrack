package product

import (
	"backend-profitrack/helpers"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CountProductsService(ctx *gin.Context)
	CreateProductService(ctx *gin.Context)
	GetAllProductService(ctx *gin.Context)
	GetProductByIdService(ctx *gin.Context)
	UpdateProductService(ctx *gin.Context)
	DeleteProductService(ctx *gin.Context)
	ImportExcelService(ctx *gin.Context)
	ExportExcelService(ctx *gin.Context)
}

type productService struct {
	repository Repository
}

func NewProductService(repo Repository) Service {
	return &productService{
		repository: repo,
	}
}

func (service *productService) CountProductsService(ctx *gin.Context) {
	result, err := service.repository.CountProductsRepository()
	if err != nil {
		response := map[string]string{"error": "gagal menghitung jumlah data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]int{"count": int(result)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
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
			PriceSale:    product.PriceSale,
			PurchaseCost: product.PurchaseCost,
			Profit:       product.Profit,
			Unit:         product.Unit,
			Stock:        product.Stock,
			Sold:         product.Sold,
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
		newProduct.PurchaseCost == 0 ||
		newProduct.PriceSale == 0 ||
		newProduct.Unit == "" ||
		newProduct.Stock == 0 ||
		newProduct.Sold == 0 {
		response = map[string]string{"error": "semua field harus diisi dengan nilai yang valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	newProduct.Profit = newProduct.PriceSale - newProduct.PurchaseCost
	newProduct.CreatedAt = time.Now()
	newProduct.UpdatedAt = time.Now()

	err := service.repository.CreateProductRepository(&newProduct)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"uni_products_name\"") {
			response = map[string]string{"error": "nama produk sudah ada"}
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

	ctx.JSON(http.StatusOK, ResponseProduct{
		ID:           product.ID,
		Name:         product.Name,
		PurchaseCost: product.PurchaseCost,
		PriceSale:    product.PriceSale,
		Profit:       product.Profit,
		Unit:         product.Unit,
		Stock:        product.Stock,
		Sold:         product.Sold,
	})
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
		product.PriceSale == 0 ||
		product.PurchaseCost == 0 ||
		product.Unit == "" ||
		product.Stock == 0 ||
		product.Sold == 0 {
		response = map[string]string{"error": "semua field harus diisi dengan nilai yang valid"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	if existingProduct.Name == product.Name && existingProduct.PriceSale == product.PriceSale && existingProduct.PurchaseCost == product.PurchaseCost && existingProduct.Profit == product.Profit && existingProduct.Unit == product.Unit && existingProduct.Stock == product.Stock {
		response = map[string]string{"error": "masukkan minimal satu data yang baru"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	existingProduct.Name = product.Name
	existingProduct.PriceSale = product.PriceSale
	existingProduct.PurchaseCost = product.PurchaseCost
	existingProduct.Profit = product.PriceSale - product.PurchaseCost
	existingProduct.Unit = product.Unit
	existingProduct.Stock = product.Stock
	existingProduct.Sold = product.Sold
	existingProduct.UpdatedAt = time.Now()

	err = service.repository.UpdateProductRepository(&existingProduct)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key score violates unique constraint \"uni_products_name\"") {
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
	err = service.repository.DeleteProductRepository(&product)
	if err != nil {
		response = map[string]string{"error": "gagal menghapus data produk"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response = map[string]string{"message": fmt.Sprintf("Produk dengan ID:%d berhasil dihapus", id)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *productService) ImportExcelService(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	if filepath.Ext(file.Filename) != ".xlsx" {
		helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{"error": "Format file harus xlsx"})
		return
	}

	// Simpan file sementara
	tempFile := fmt.Sprintf("temp/%d-%s", time.Now().Unix(), file.Filename)
	if err = ctx.SaveUploadedFile(file, tempFile); err != nil {
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}
	defer os.Remove(tempFile)

	// Baca file Excel
	xlsx, err := excelize.OpenFile(tempFile)
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal membaca file Excel"})
		return
	}
	defer xlsx.Close()

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal membaca sheet"})
		return
	}

	var products []Product
	for i, row := range rows {
		if i == 0 {
			continue // Skip header row
		}

		// Periksa jumlah kolom yang dibutuhkan (6 kolom sesuai format)
		if len(row) < 6 {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Format tidak valid pada baris %d: jumlah kolom kurang dari 6", i+1),
			})
			return
		}

		// Fungsi helper untuk membersihkan string angka
		cleanNumber := func(str string) string {
			str = strings.ReplaceAll(str, ".", "")
			str = strings.ReplaceAll(str, ",", "")
			str = strings.ReplaceAll(str, " ", "")
			return str
		}

		// Bersihkan dan konversi data numerik
		cleanPurchaseCost := cleanNumber(row[1]) // harga beli di kolom 2
		cleanPriceSale := cleanNumber(row[2])    // harga jual di kolom 3
		cleanStock := cleanNumber(row[4])        // stok di kolom 5
		cleanSold := cleanNumber(row[5])         // stok terjual di kolom 6

		// Konversi ke integer dengan validasi
		purchaseCost, err := strconv.Atoi(cleanPurchaseCost)
		if err != nil {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Format harga beli tidak valid pada baris %d: %s", i+1, row[1]),
			})
			return
		}

		priceSale, err := strconv.Atoi(cleanPriceSale)
		if err != nil {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Format harga jual tidak valid pada baris %d: %s", i+1, row[2]),
			})
			return
		}

		stock, err := strconv.Atoi(cleanStock)
		if err != nil {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Format stok tidak valid pada baris %d: %s", i+1, row[4]),
			})
			return
		}

		sold, err := strconv.Atoi(cleanSold)
		if err != nil {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Format stok terjual tidak valid pada baris %d: %s", i+1, row[5]),
			})
			return
		}

		// Validasi angka negatif
		if purchaseCost < 0 || priceSale < 0 || stock < 0 || sold < 0 {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Nilai tidak boleh negatif pada baris %d", i+1),
			})
			return
		}

		product := Product{
			Name:         strings.TrimSpace(row[0]),
			PurchaseCost: purchaseCost,
			PriceSale:    priceSale,
			Profit:       priceSale - purchaseCost,
			Unit:         strings.TrimSpace(row[3]),
			Stock:        stock,
			Sold:         sold,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		products = append(products, product)
	}

	// Bulk insert
	if err = service.repository.BulkCreateProductRepository(products); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			helpers.ResponseJSON(ctx, http.StatusBadRequest, gin.H{"error": "Beberapa produk sudah ada (duplikat nama produk)"})
			return
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	helpers.ResponseJSON(ctx, http.StatusOK, gin.H{"message": fmt.Sprintf("Berhasil import %d produk", len(products))})
}

func (service *productService) ExportExcelService(ctx *gin.Context) {
	products, err := service.repository.GetAllProductRepository()
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}

	// Buat file Excel baru
	f := excelize.NewFile()
	defer f.Close()

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		// Handle error
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal membuat style"})
		return
	}

	// Buat header
	headers := []string{"Nama Produk", "Harga Beli", "Harga Jual", "Keuntungan", "Satuan", "Stok", "Stok terjual"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue("Sheet1", cell, header)
		f.SetCellStyle("Sheet1", cell, cell, style)
	}

	// Isi data
	for i, product := range products {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), product.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), product.PurchaseCost)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), product.PriceSale)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), product.Profit)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), product.Unit)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), product.Stock)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), product.Sold)
	}

	// Format nama file dengan tanggal
	currentTime := time.Now()
	fileName := fmt.Sprintf("data-produk-%s.xlsx", currentTime.Format("02-01-2006"))

	// Set response header
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	if err := f.Write(ctx.Writer); err != nil {
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal membuat file Excel"})
		return
	}
}
