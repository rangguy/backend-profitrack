package score

import (
	"backend-profitrack/middleware"
	"backend-profitrack/modules/criteria"
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
	finalScoreRepo := final_score.NewFinalScoreRepository(db)
	service := NewScoreService(repo, productRepo, criteriaRepo, methodRepo, finalScoreRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())

	//SMART
	api.PUT("/scores/:methodID/SMART", service.UtilityScoreSMARTService, service.ScoreOneTimesWeightByMethodIDService, service.CreateFinalScoresSMARTService)
	api.DELETE("/scores/:methodID/SMART", service.DeleteAllScoresSMARTService)

	//MOORA
	api.PUT("/scores/:methodID/MOORA", service.NormalizeScoreMOORAService, service.ScoreOneTimesWeightByMethodIDService, service.CreateFinalScoresMOORAService)
	api.DELETE("/scores/:methodID/MOORA", service.DeleteAllScoresMOORAService)

	// Delete Final Score And Create Report
	api.DELETE("/final_scores/:methodID", service.CreateDeleteFinalScoreByMethodIDService)
}
