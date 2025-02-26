package score

import (
	"backend-profitrack/helpers"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/criteria_score"
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
	GetAllScoreByMethodIDService(ctx *gin.Context)
	UtilityScoreSMARTService(ctx *gin.Context)
	ScoreOneTimesWeightByMethodIDService(ctx *gin.Context)
	NormalizeScoreMOORAService(ctx *gin.Context)
	CreateFinalScoresSMARTService(ctx *gin.Context)
	CreateFinalScoresMOORAService(ctx *gin.Context)
	CreateFinalScoreAndReportByMethodIDService(ctx *gin.Context)
}

type scoreService struct {
	repository              Repository
	productRepository       product.Repository
	criteriaRepository      criteria.Repository
	methodRepository        method.Repository
	criteriaScoreRepository criteria_score.Repository
	finalScoreRepository    final_score.Repository
}

func NewScoreService(repo Repository, productRepo product.Repository, criteriaRepo criteria.Repository, methodRepo method.Repository, criteriaScoreRepo criteria_score.Repository, finalScoreRepo final_score.Repository) Service {
	return &scoreService{
		repository:              repo,
		productRepository:       productRepo,
		criteriaRepository:      criteriaRepo,
		methodRepository:        methodRepo,
		criteriaScoreRepository: criteriaScoreRepo,
		finalScoreRepository:    finalScoreRepo,
	}
}

func (service *scoreService) GetAllScoreByMethodIDService(ctx *gin.Context) {
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
	var result []Score
	for _, score := range scores {
		result = append(result, Score{
			ID:         score.ID,
			ProductID:  score.ProductID,
			CriteriaID: score.CriteriaID,
			ScoreOne:   score.ScoreOne,
			ScoreTwo:   score.ScoreTwo,
			MethodID:   score.MethodID,
			CreatedAt:  score.CreatedAt,
			UpdatedAt:  score.UpdatedAt,
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

func (service *scoreService) UtilityScoreSMARTService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	// Step 1: Retrieve all scores for the given methodID
	scores, err := service.criteriaScoreRepository.GetAllCriteriaScoreRepository()
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

				newScore := &Score{
					ProductID:  score.ProductID,
					CriteriaID: score.CriteriaID,
					ScoreOne:   scoreOne,
					MethodID:   methodID,
					CreatedAt:  score.CreatedAt,
					UpdatedAt:  score.UpdatedAt,
				}

				err = service.repository.CreateScoreRepository(newScore)
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

	scores, err := service.criteriaScoreRepository.GetAllCriteriaScoreRepository()
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

				newScore := &Score{
					ProductID:  score.ProductID,
					CriteriaID: score.CriteriaID,
					ScoreOne:   normalizedValue,
					MethodID:   methodID,
					CreatedAt:  score.CreatedAt,
					UpdatedAt:  score.UpdatedAt,
				}

				err = service.repository.CreateScoreRepository(newScore)
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

func (service *scoreService) CreateFinalScoreAndReportByMethodIDService(ctx *gin.Context) {
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

	if len(finalScores) == 0 {
		response := map[string]string{"error": "data nilai akhir masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusNotFound, response)
		return
	}

	// Buat entri laporan utama dahulu
	newReport := report.Report{
		MethodID:  methodID,
		TotalData: len(finalScores),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Simpan laporan utama ke database
	err = service.repository.CreateReportFinalScoreByMethodIDRepository(&newReport)
	if err != nil {
		response := map[string]string{"error": "gagal membuat laporan utama"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Iterate over final scores and create report detail entries
	for _, score := range finalScores {
		// Create a new report detail entry
		reportDetail := report.ReportDetail{
			MethodID:   methodID,
			ProductID:  score.ProductID,
			ReportID:   newReport.ID,
			FinalScore: score.FinalScore,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err = service.repository.CreateReportDetailRepository(&reportDetail)
		if err != nil {
			response := map[string]string{"error": "gagal memasukkan nilai ke dalam detail laporan"}
			helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
			return
		}
	}

	err = service.repository.DeleteAllScoresByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	err = service.repository.DeleteFinalScoreByMethodIDRepository(methodID)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus semua data nilai akhir"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Berhasil membuat laporan dan memasukkan data ke detail laporan"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
