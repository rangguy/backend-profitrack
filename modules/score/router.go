package score

import (
	"backend-profitrack/middleware"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/criteria_score"
	"backend-profitrack/modules/final_score"
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewScoreRepository(db)
	productRepo := product.NewProductRepository(db)
	criteriaRepo := criteria.NewCriteriaRepository(db)
	methodRepo := method.NewMethodRepository(db)
	criteriaScoreRepo := criteria_score.NewCriteriaScoreRepository(db)
	finalScoreRepo := final_score.NewFinalScoreRepository(db)
	service := NewScoreService(repo, productRepo, criteriaRepo, methodRepo, criteriaScoreRepo, finalScoreRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/scores/:methodID", service.GetAllScoreByMethodIDService)

	//SMART
	api.POST("/scores/:methodID/SMART", service.UtilityScoreSMARTService, service.ScoreOneTimesWeightByMethodIDService, service.CreateFinalScoresSMARTService)

	//MOORA
	api.POST("/scores/:methodID/MOORA", service.NormalizeScoreMOORAService, service.ScoreOneTimesWeightByMethodIDService, service.CreateFinalScoresMOORAService)

	// Create Final Scores and Report
	api.POST("/final_scores/:methodID", service.CreateFinalScoreAndReportByMethodIDService)
}
