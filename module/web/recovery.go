package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logicinu/nest/module/logger"
)

func recovery() gin.HandlerFunc {

	logger := logger.GetLogger()
	sugar := logger.Sugar()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				sugar.Errorw("PANIC", "error", err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
