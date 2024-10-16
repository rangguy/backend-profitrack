package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initiator(router *gin.Engine, db *gorm.DB) {
	repo := NewUserRepository(db)
	service := NewUserService(repo)

	api := router.Group("/api")
	api.POST("/login", service.LoginService)
}
