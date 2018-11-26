package controller

import(
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine){
	index := router.Group("/")
	{
		index.GET("/test", test)
		index.GET("/panic", panicTest)
	}
}
