package helpers

import (
	"github.com/gin-gonic/gin"
)

func ResponseJSON(ctx *gin.Context, statusCode int, payload interface{}) {
	ctx.JSON(
		statusCode,
		payload,
	)
}
