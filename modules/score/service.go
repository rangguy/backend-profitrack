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
	CalculateSMARTService(ctx *gin.Context)
	CalculateMOORAService(ctx *gin.Context)
	CreateReportByMethodIDService(ctx *gin.Context)
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

	// Step 2: Format response data
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

func (service *scoreService) CalculateSMARTService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	startTime := time.Now()

	// Lakukan perhitungan utility score
	err = service.utilityScoreSMART(methodID)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"process": "Perhitungan utility nilai SMART",
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Lakukan perhitungan score times weight
	err = service.scoreOneTimesWeightByMethodID(methodID)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"process": "Perhitungan score one x bobot",
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Lakukan perhitungan final score
	err = service.createFinalScoresSMART(methodID)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"process": "Perhitungan final score SMART",
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Hitung selisih waktu
	endTime := time.Now()
	processingTime := endTime.Sub(startTime)

	// Kirim response dengan waktu proses
	response := map[string]interface{}{
		"message":        "Normalisasi nilai utility, utility x bobot, dan perhitungan skor akhir SMART berhasil",
		"processingTime": processingTime.String(),
	}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

// CalculateMOORAService - Public endpoint untuk perhitungan MOORA
func (service *scoreService) CalculateMOORAService(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	startTime := time.Now()

	// Lakukan normalisasi score MOORA
	err = service.normalizeScoreMOORA(methodID)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"process": "Normalisasi nilai MOORA",
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Lakukan perhitungan score times weight
	err = service.scoreOneTimesWeightByMethodID(methodID)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"process": "Perhitungan score one x bobot",
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Lakukan perhitungan final score
	err = service.createFinalScoresMOORA(methodID)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"process": "Perhitungan final score MOORA",
		}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// Hitung selisih waktu
	endTime := time.Now()
	processingTime := endTime.Sub(startTime)

	// Kirim response dengan waktu proses
	response := map[string]interface{}{
		"message":        "Normalisasi nilai normalisasi, normalisasi x bobot, dan perhitungan skor akhir MOORA berhasil",
		"processingTime": processingTime.String(),
	}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) CreateReportByMethodIDService(ctx *gin.Context) {
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
		ReportCode: fmt.Sprintf("LAP-%d-%d", methodID, time.Now().Unix()),
		MethodID:   methodID,
		TotalData:  len(finalScores),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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

func (service *scoreService) utilityScoreSMART(methodID int) error {
	// Step 1: Retrieve all criteria to get weights and types
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		return fmt.Errorf("gagal mengambil data kriteria: %v", err)
	}

	// Step 2: Retrieve all scores
	scores, err := service.criteriaScoreRepository.GetAllCriteriaScoreRepository()
	if err != nil {
		return fmt.Errorf("gagal mengambil data nilai kriteria: %v", err)
	}

	if len(scores) == 0 {
		return fmt.Errorf("data nilai kriteria masih kosong")
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
					return fmt.Errorf("gagal update score: %v", err)
				}
			}
		}
	}

	return nil
}

func (service *scoreService) normalizeScoreMOORA(methodID int) error {
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		return fmt.Errorf("gagal mengambil data kriteria: %v", err)
	}

	scores, err := service.criteriaScoreRepository.GetAllCriteriaScoreRepository()
	if err != nil {
		return fmt.Errorf("gagal mengambil data nilai kriteria: %v", err)
	}

	if len(scores) == 0 {
		return fmt.Errorf("data nilai kriteria masih kosong")
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
					return fmt.Errorf("gagal update score: %v", err)
				}
			}
		}
	}

	return nil
}

func (service *scoreService) scoreOneTimesWeightByMethodID(methodID int) error {
	// Step 1: Retrieve all scores for the given methodID
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		return fmt.Errorf("gagal mengambil data nilai: %v", err)
	}

	if len(scores) == 0 {
		return fmt.Errorf("data nilai masih kosong")
	}

	// Step 2: Retrieve all criteria to get weights and types
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		return fmt.Errorf("gagal mengambil data kriteria: %v", err)
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
					return fmt.Errorf("gagal update score two: %v", err)
				}
			}
		}
	}

	return nil
}

func (service *scoreService) createFinalScoresSMART(methodID int) error {
	// Ambil semua score untuk method ini
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		return fmt.Errorf("gagal mengambil data nilai: %v", err)
	}

	if len(scores) == 0 {
		return fmt.Errorf("data nilai masih kosong")
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
			return fmt.Errorf("gagal menyimpan final score: %v", err)
		}
	}

	return nil
}

func (service *scoreService) createFinalScoresMOORA(methodID int) error {
	// Ambil semua score untuk method ini
	scores, err := service.repository.GetAllScoreByMethodIDRepository(methodID)
	if err != nil {
		return fmt.Errorf("gagal mengambil data nilai kriteria: %v", err)
	}

	// Ambil semua kriteria untuk mendapatkan tipe kriteria (benefit/cost)
	criteriaList, err := service.criteriaRepository.GetAllCriteriaRepository()
	if err != nil {
		return fmt.Errorf("gagal mengambil data kriteria: %v", err)
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
			return fmt.Errorf("gagal menyimpan final score: %v", err)
		}
	}

	return nil
}
