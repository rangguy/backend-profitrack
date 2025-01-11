package final_score

import (
	"backend-profitrack/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Service interface {
	GetAllFinalScoreByMethodID(ctx *gin.Context)
}

type newFinalScoreService struct {
	repository Repository
}

func NewFinalScoreService(repository Repository) Service {
	return &newFinalScoreService{repository}
}

func (service *newFinalScoreService) GetAllFinalScoreByMethodID(ctx *gin.Context) {
	methodID, err := strconv.Atoi(ctx.Param("methodID"))
	if err != nil {
		response := map[string]string{"error": "ID tidak sesuai"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	finalScore, err := service.repository.GetAllFinalScoreByMethodID(methodID)
	if err != nil {
		helpers.ResponseJSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	var result []FinalScore
	for _, score := range finalScore {
		result = append(result, FinalScore{
			ID:         score.ID,
			FinalScore: score.FinalScore,
			ProductID:  score.ProductID,
			MethodID:   score.MethodID,
			CreatedAt:  score.CreatedAt,
			UpdatedAt:  score.UpdatedAt,
		})
	}

	if result == nil {
		response := map[string]string{"message": "data nilai akhir masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusOK, response)
		return
	} else {
		helpers.ResponseJSON(ctx, http.StatusOK, result)
	}
}
