package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/logicinu/nest/module/engine"
	"github.com/logicinu/nest/module/web"
	"github.com/logicinu/nest/module/id"
	"github.com/logicinu/nest/module/logger"
	"github.com/logicinu/nest/module/redis"
	"github.com/logicinu/nest/module/setting"
)

var (
	// BuildVersion from git tag
	BuildVersion string
	// BuildTime from make time
	BuildTime string
	// BuildMode from make mode
	BuildMode string
)

// dispay app info
func info() {
	var v bool

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.Parse()

	if v {
		log.Println(fmt.Sprintf("\nBuildVersion: %v\n   BuildTime: %v\n   BuildMode: %v", BuildVersion, BuildTime, BuildMode))
		os.Exit(0)
	}

	if len(BuildMode) == 0 {
		BuildMode = "dev"
	}
}

func main() {
	info()

	setting.InitSetting(BuildMode)

	logger.InitLogger(BuildMode)
	defer logger.GetLogger().Sync()

	id.InitId()
	engine.InitEngineMap(BuildMode)
	redis.InitRedisPool()
	web.InitGinEngine(BuildMode)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := web.GetHttpServer().Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown:", err)
	}
	log.Println("Server exiting")
}
