package controller

import (
	"github.com/logicinu/nest/module/result"
	"net/http"

	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	c.JSON(http.StatusOK, result.GetResultOk())
}

func panicTest(c *gin.Context) {
	panic("panic test")
}
