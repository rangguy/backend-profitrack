package report

import (
	"backend-profitrack/middleware"
	"backend-profitrack/modules/method"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewReportRepository(db)
	methodRepo := method.NewMethodRepository(db)
	service := NewReportService(repo, methodRepo)

	api := router.Group("/api")
	api.Use(middleware.LoggingMiddleware())
	api.Use(middleware.JWTMiddleware())
	api.GET("/reports", service.GetAllReportsService)
	api.GET("/reports/count", service.CountReportsService)
	api.GET("/reports/:id", service.GetDetailReportService)
	api.GET("/reports/export/:id", service.ExportPDFService)
	api.DELETE("/reports/:id", service.DeleteDetailReportService, service.DeleteAllReportService)
}
