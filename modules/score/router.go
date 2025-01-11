package score

import (
	"backend-profitrack/middleware"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewScoreRepository(db)
	productRepo := product.NewProductRepository(db)
	criteriaRepo := criteria.NewCriteriaRepository(db)
	service := NewScoreService(repo, productRepo, criteriaRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/scores/:methodID", service.GetAllScoreByMethodIDService)
	api.POST("/scores/:methodID", service.CreateScoreByMethodIDService)
	api.PUT("scores/:methodID", service.NormalizeScoreByMethodIDService, service.UtilityScoreByMethodIDService, service.FinalScoreByMethodIDService)
	api.DELETE("/scores", service.DeleteAllScoreService)
}
