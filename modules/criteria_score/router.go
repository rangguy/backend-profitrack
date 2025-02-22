package criteria_score

import (
	"backend-profitrack/middleware"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewCriteriaScoreRepository(db)
	criteriaRepo := criteria.NewCriteriaRepository(db)
	productRepo := product.NewProductRepository(db)
	service := NewCriteriaScoreService(repo, criteriaRepo, productRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())

	api.GET("/criteria_scores", service.GetAllCriteriaScoreService)
	api.POST("/criteria_scores", service.CreateAllCriteriaScoreService)
	api.PUT("criteria_scores", service.UpdateCriteriaScoreService)
	api.DELETE("/criteria_scores/:id", service.DeleteCriteriaScoreService)
}
