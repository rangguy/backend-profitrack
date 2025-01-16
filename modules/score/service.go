package score

import (
	"backend-profitrack/helpers"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/final_score"
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"backend-profitrack/modules/report"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreateScoreByMethodIDService(ctx *gin.Context)
	GetAllScoresByMethodIDService(ctx *gin.Context)
	UtilityScoreSMARTService(ctx *gin.Context)
	ScoreOneTimesWeightByMethodIDService(ctx *gin.Context)
	NormalizeScoreMOORAService(ctx *gin.Context)
	CreateFinalScoresSMARTService(ctx *gin.Context)
	CreateFinalScoresMOORAService(ctx *gin.Context)
	CreateDeleteFinalScoreByMethodIDService(ctx *gin.Context)
	DeleteAllScoresSMARTService(ctx *gin.Context)
	DeleteAllScoresMOORAService(ctx *gin.Context)
}

type scoreService struct {
	repository           Repository
	productRepository    product.Repository
	criteriaRepository   criteria.Repository
	methodRepository     method.Repository
	finalScoreRepository final_score.Repository
}

func NewScoreService(repo Repository, productRepo product.Repository, criteriaRepo criteria.Repository, methodRepo method.Repository, finalScoreRepo final_score.Repository) Service {
	return &scoreService{
		repository:           repo,
		productRepository:    productRepo,
		criteriaRepository:   criteriaRepo,
		methodRepository:     methodRepo,
		finalScoreRepository: finalScoreRepo,
	}
}

func (service *scoreService) CreateScoreByMethodIDService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

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

			newScore := Score{
				ProductID:  produk.ID,
				CriteriaID: kriteria.ID,
				MethodID:   methodID,
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

func (service *scoreService) GetAllScoresByMethodIDService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	values, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []Score
	for _, value := range values {
		result = append(result, Score{
			ID:         value.ID,
			Score:      value.Score,
			ScoreOne:   value.ScoreOne,
			ScoreTwo:   value.ScoreTwo,
			ProductID:  value.ProductID,
			CriteriaID: value.CriteriaID,
			MethodID:   value.MethodID,
			CreatedAt:  value.CreatedAt,
			UpdatedAt:  value.UpdatedAt,
		})
	}

	if len(result) == 0 {
		// Ambil data method untuk mendapatkan namanya
		method, err := service.methodRepository.GetMethodByIdRepository(methodID)
		if err != nil {
			response := map[string]string{"error": "Metode tidak ditemukan"}
			helpers.ResponseJSON(ctx, http.StatusNotFound, response)
			return
		}
		response := map[string]string{"message": fmt.Sprintf("data nilai dengan metode %s masih kosong", method.Name)}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	}

	helpers.ResponseJSON(ctx, http.StatusOK, result)
}

func (service *scoreService) UtilityScoreSMARTService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Step 1: Retrieve all scores for the given methodID
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data score"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Step 2: Retrieve all criteria to get weights and types
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Step 3: Create maps for criteria types
	criteriaTypes := make(map[int]string)
	for _, criterion := range criteriaList {
		criteriaTypes[criterion.ID] = criterion.Type
	}

	// Step 4: Group normalized scores by criteria
	scoreByCriteria := make(map[int][]float64)
	for _, score := range scores {
		scoreByCriteria[score.CriteriaID] = append(scoreByCriteria[score.CriteriaID], score.Score)
	}

	// Step 5: Calculate ScoreOne based on normalized scores
	for _, criteria := range criteriaList {
		scoreValues := scoreByCriteria[criteria.ID]
		if len(scoreValues) == 0 {
			continue
		}

		// Find min and max from normalized scores
		minNorm := scoreValues[0]
		maxNorm := scoreValues[0]
		for _, value := range scoreValues {
			if value < minNorm {
				minNorm = value
			}
			if value > maxNorm {
				maxNorm = value
			}
		}

		// Calculate ScoreOne for each score based on criteria type
		for _, score := range scores {
			if score.CriteriaID == criteria.ID {
				var scoreOne float64

				if maxNorm != minNorm {
					if criteriaTypes[criteria.ID] == "benefit" {
						scoreOne = (score.Score - minNorm) / (maxNorm - minNorm)
					} else if criteriaTypes[criteria.ID] == "cost" {
						scoreOne = (maxNorm - score.Score) / (maxNorm - minNorm)
					}
				}

				// Update score with ScoreOne value
				score.ScoreOne = scoreOne
				err = service.repository.UpdateScoreByMethodIDRepository(methodID, &score)
				if err != nil {
					response := map[string]string{"error": "gagal update score"}
					helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
					return
				}
			}
		}
	}

	response := map[string]string{"message": "Perhitungan utility nilai SMART berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) NormalizeScoreMOORAService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data score"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Kelompokkan score berdasarkan criteria
	scoreByCriteria := make(map[int][]float64)
	for _, score := range scores {
		scoreByCriteria[score.CriteriaID] = append(scoreByCriteria[score.CriteriaID], score.Score)
	}

	// Lakukan normalisasi untuk setiap kriteria
	for _, criteria := range criteriaList {
		scoreValues := scoreByCriteria[criteria.ID]
		if len(scoreValues) == 0 {
			continue
		}

		// Hitung akar dari jumlah kuadrat
		sumOfSquares := 0.0
		for _, value := range scoreValues {
			sumOfSquares += value * value
		}
		sqrtSum := math.Sqrt(sumOfSquares)

		// Normalisasi setiap score untuk kriteria ini menggunakan metode vector normalization
		for _, score := range scores {
			if score.CriteriaID == criteria.ID {
				// Rumus normalisasi MOORA: rij = xij / sqrt(Σ(xij^2))
				normalizedValue := 0.0
				if sqrtSum != 0 {
					normalizedValue = score.Score / sqrtSum
				}

				// Update score dengan nilai yang sudah dinormalisasi
				score.ScoreOne = normalizedValue
				err = service.repository.UpdateScoreByMethodIDRepository(methodID, &score)
				if err != nil {
					response := map[string]string{"error": "gagal update score"}
					helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
					return
				}
			}
		}
	}

	response := map[string]string{"message": "Normalisasi nilai MOORA berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) ScoreOneTimesWeightByMethodIDService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Step 1: Retrieve all scores for the given methodID
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data score"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Step 2: Retrieve all criteria to get weights and types
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Step 3: Create maps for criteria weight
	criteriaWeight := make(map[int]float64)
	for _, criterion := range criteriaList {
		criteriaWeight[criterion.ID] = criterion.Weight
	}

	// Step 4: Group normalized scores by criteria
	scoreOneByCriteria := make(map[int][]float64)
	for _, score := range scores {
		scoreOneByCriteria[score.CriteriaID] = append(scoreOneByCriteria[score.CriteriaID], score.ScoreOne)
	}

	// Step 5: Calculate ScoreTwo based on ScoreOne scores
	for _, criteria := range criteriaList {
		scoreOne := scoreOneByCriteria[criteria.ID]
		if len(scoreOne) == 0 {
			continue
		}

		// Calculate ScoreTwo for each score based on criteria weight
		for _, score := range scores {
			if score.CriteriaID == criteria.ID {
				var scoreTwo float64
				// Mengambil bobot kriteria
				weight := criteria.Weight

				// Mengalikan ScoreOne dengan bobot kriteria
				scoreTwo = score.ScoreOne * weight

				// Update score with ScoreOne value
				score.ScoreTwo = scoreTwo
				err = service.repository.UpdateScoreByMethodIDRepository(methodID, &score)
				if err != nil {
					response := map[string]string{"error": "gagal update score two"}
					helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
					return
				}
			}
		}
	}

	response := map[string]string{"message": "Perhitungan score one x bobot berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) CreateFinalScoresSMARTService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Ambil semua score untuk method ini
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data score"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Buat map untuk menyimpan total score per produk
	productScores := make(map[int]float64)

	// Jumlahkan ScoreTwo untuk setiap produk
	for _, score := range scores {
		productScores[score.ProductID] += score.ScoreTwo
	}

	// Simpan final score untuk setiap produk
	for productID, totalScore := range productScores {
		finalScore := &final_score.FinalScore{
			ProductID:  productID,
			MethodID:   methodID,
			FinalScore: totalScore,
		}

		err = service.repository.CreateFinalScoreByMethodIDRepository(methodID, finalScore)
		if err != nil {
			response := map[string]string{"error": "gagal menyimpan final score"}
			helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
			return
		}
	}

	response := map[string]string{"message": "perhitungan final score SMART berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) CreateFinalScoresMOORAService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Ambil semua score untuk method ini
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data score"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Ambil semua kriteria untuk mendapatkan tipe kriteria (benefit/cost)
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		response := map[string]string{"error": "gagal mengambil data kriteria"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Buat map untuk menyimpan tipe kriteria
	criteriaTypes := make(map[int]string)
	for _, criteria := range criteriaList {
		criteriaTypes[criteria.ID] = criteria.Type
	}

	// Kelompokkan score berdasarkan produk
	type productScore struct {
		benefitSum float64
		costSum    float64
	}
	productScores := make(map[int]*productScore)

	// Hitung jumlah benefit dan cost berdasarkan ScoreTwo
	for _, score := range scores {
		if _, exists := productScores[score.ProductID]; !exists {
			productScores[score.ProductID] = &productScore{0, 0}
		}

		// ScoreTwo sudah merupakan weighted normalized value
		if strings.EqualFold(criteriaTypes[score.CriteriaID], "benefit") {
			productScores[score.ProductID].benefitSum += score.ScoreTwo
		} else if strings.EqualFold(criteriaTypes[score.CriteriaID], "cost") {
			productScores[score.ProductID].costSum += score.ScoreTwo
		}
	}

	// Hitung nilai Yi untuk setiap produk
	for productID, score := range productScores {
		fmt.Printf("Product ID: %d, Total Benefit: %f, Total Cost: %f\n", productID, score.benefitSum, score.costSum)

		// Hitung Yi = (sum of benefit criteria) - (sum of cost criteria)
		yi := score.benefitSum - score.costSum

		finalScore := &final_score.FinalScore{
			ProductID:  productID,
			MethodID:   methodID,
			FinalScore: yi,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = service.repository.CreateFinalScoreByMethodIDRepository(methodID, finalScore)
		if err != nil {
			response := map[string]string{"error": "gagal menyimpan final score"}
			helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
			return
		}
	}

	response := map[string]string{"message": "perhitungan final score MOORA berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) DeleteAllScoresSMARTService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Setelah menyimpan semua final score, baru hapus data score
	err = service.repository.DeleteAllScoresByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "perhitungan final score dan penghapusan data nilai berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) DeleteAllScoresMOORAService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	//Setelah menyimpan semua final score, baru hapus data score
	err = service.repository.DeleteAllScoresByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "perhitungan final score dan penghapusan data nilai berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) CreateDeleteFinalScoreByMethodIDService(ctx *gin.Context) {
	var reportFinalScores report.Report

	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Retrieve all final scores by method ID
	finalScores, err := service.finalScoreRepository.GetAllFinalScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusNotFound, response)
		return
	}

	// Iterate over final scores and create report entries
	for _, score := range finalScores {
		// Extract month and year from CreatedAt
		year, month, _ := score.CreatedAt.Date()
		period := fmt.Sprintf("%s %d", month.String(), year) // e.g., "January 2023"

		// Create a new report entry
		reportFinalScores = report.Report{
			ProductID:  score.ProductID,
			MethodID:   score.MethodID,
			FinalScore: score.FinalScore,
			Period:     period,
		}

		err = service.repository.CreateReportFinalScoreByMethodIDRepository(&reportFinalScores)
		if err != nil {
			response := map[string]string{"error": "gagal memasukkan nilai ke dalam laporan"}
			helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
			return
		}
	}

	err = service.repository.DeleteAllScoresByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus nilai akhir"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	err = service.repository.DeleteFinalScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus nilai akhir"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Berhasil menghapus nilai dan memasukkan ke dalam laporan"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
