package score

import (
	"backend-profitrack/helpers"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreateScoreByMethodIDService(ctx *gin.Context)
	GetAllScoreByMethodIDService(ctx *gin.Context)
	NormalizeScoreByMethodIDService(ctx *gin.Context)
	UtilityScoreByMethodIDService(ctx *gin.Context)
	FinalScoreByMethodIDService(ctx *gin.Context)
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

func (service *scoreService) GetAllScoreByMethodIDService(ctx *gin.Context) {
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
			ID:             value.ID,
			Score:          value.Score,
			NormalizeScore: value.NormalizeScore,
			ScoreOne:       value.ScoreOne,
			ScoreTwo:       value.ScoreTwo,
			ProductID:      value.ProductID,
			CriteriaID:     value.CriteriaID,
			MethodID:       value.MethodID,
			CreatedAt:      value.CreatedAt,
			UpdatedAt:      value.UpdatedAt,
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
			case strings.ToLower("Profit Margin"):
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

			err = service.repository.CreateScoreByMethodIDRepository(methodID, &newScore)
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

func (service *scoreService) NormalizeScoreByMethodIDService(ctx *gin.Context) {
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

	// Ambil semua score untuk method ini
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

		// Cari nilai min dan max
		minScore := scoreValues[0]
		maxScore := scoreValues[0]
		for _, value := range scoreValues {
			if value < minScore {
				minScore = value
			}
			if value > maxScore {
				maxScore = value
			}
		}

		// Normalisasi setiap score untuk kriteria ini
		for _, score := range scores {
			if score.CriteriaID == criteria.ID {
				normalizedValue := 0.0
				if maxScore != minScore {
					normalizedValue = (score.Score - minScore) / (maxScore - minScore)
				}

				// Update score dengan nilai yang sudah dinormalisasi
				score.NormalizeScore = normalizedValue
				err = service.repository.UpdateScoreByMethodIDRepository(methodID, &score)
				if err != nil {
					response := map[string]string{"error": "gagal update score"}
					helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
					return
				}
			}
		}
	}

	response := map[string]string{"message": "Normalisasi score berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) UtilityScoreByMethodIDService(ctx *gin.Context) {
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
	normalizedScoreByCriteria := make(map[int][]float64)
	for _, score := range scores {
		normalizedScoreByCriteria[score.CriteriaID] = append(normalizedScoreByCriteria[score.CriteriaID], score.NormalizeScore)
	}

	// Step 5: Calculate ScoreOne based on normalized scores
	for _, criteria := range criteriaList {
		normalizedValues := normalizedScoreByCriteria[criteria.ID]
		if len(normalizedValues) == 0 {
			continue
		}

		// Find min and max from normalized scores
		minNorm := normalizedValues[0]
		maxNorm := normalizedValues[0]
		for _, value := range normalizedValues {
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
						scoreOne = (score.NormalizeScore - minNorm) / (maxNorm - minNorm)
					} else if criteriaTypes[criteria.ID] == "cost" {
						scoreOne = (maxNorm - score.NormalizeScore) / (maxNorm - minNorm)
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

	response := map[string]string{"message": "Perhitungan utility score berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *scoreService) FinalScoreByMethodIDService(ctx *gin.Context) {
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

	response := map[string]string{"message": "Perhitungan final score berhasil"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
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
