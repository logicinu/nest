package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/logicinu/nest/module/logger"
	"github.com/logicinu/nest/module/setting"
	"github.com/logicinu/nest/controller"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/cors"
)

var ginEngine *gin.Engine
var httpServer *http.Server

// InitGinEngine init gin
func InitGinEngine(mode string) {
	switch mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(ginzap.Ginzap(logger.GetLogger(), time.RFC3339, true))
	router.Use(recovery())
	router.Use(cors.Default())

	if mode == "dev" || mode == "test" {
		pprof.Register(router)
	}

	controller.InitRouter(router)

	cfg := setting.GetSetting()
	host := cfg.Section("web").Key("Host").MustString("127.0.0.1")
	port := cfg.Section("web").Key("Port").MustString("8080")
	readTimeout := cfg.Section("web").Key("ReadTimeout").MustInt(10)
	writeTimeout := cfg.Section("web").Key("WriteTimeout").MustInt(10)
	maxHeaderBytes := cfg.Section("web").Key("MaxHeaderBytes").MustInt(1)

	server := &http.Server{
		Addr:           fmt.Sprintf("%v:%v", host, port),
		Handler:        router,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: maxHeaderBytes << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(fmt.Sprintf("gin ListenAndServe error: %v", err))
		}
	}()

	ginEngine = router
	httpServer = server
}

// GetGinEngine return gin engine
func GetGinEngine() *gin.Engine {
	return ginEngine
}

// GetHttpServer return http server
func GetHttpServer() *http.Server {
	return httpServer
}
