package engine

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/logicinu/nest/module/setting"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"gopkg.in/natefinch/lumberjack.v2"
)

var engineMap sync.Map

// InitEngineMap init database engine
func InitEngineMap(mode string) {
	setEngine(mode, "default")
}

// GetEngine return engine by key
func GetEngine(keys ...string) *xorm.Engine {
	if len(keys) == 0 {
		keys = append(keys, "default")
	}

	if len(keys) > 1 {
		panic(fmt.Sprintf("keys overflow %v", keys))
	}

	vv, ok := engineMap.Load(keys[0])
	if !ok {
		panic(fmt.Sprintf("get engine %v err", keys))
	}

	return vv.(*xorm.Engine)
}

// GetSession return engine session by engine key
func GetSession(keys ...string) *xorm.Session {
	e := GetEngine(keys...)
	return e.NewSession()
}

// setEngine set engine by engine key
func setEngine(mode, key string) {
	var orm *xorm.Engine

	cfg := setting.GetSetting()

	keyName := "database." + key

	user := cfg.Section(keyName).Key("user").String()
	pass := cfg.Section(keyName).Key("pass").String()
	host := cfg.Section(keyName).Key("host").String()
	port := cfg.Section(keyName).Key("port").String()
	name := cfg.Section(keyName).Key("name").String()

	switch cfg.Section(keyName).Key("type").String() {
	case "postgres":
		link := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, pass, host, port, name)
		x, err := xorm.NewEngine("postgres", link)
		if err != nil {
			log.Println(fmt.Sprintf("postgres:[%v] connection create failed: %v", link, err))
			os.Exit(1)
		} else if err = x.Ping(); err != nil {
			log.Println(fmt.Sprintf("postgres:[%v] connection ping failed: %v", link, err))
			os.Exit(1)
		}
		orm = x
	case "mysql":
		link := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pass, host, port, name)
		x, err := xorm.NewEngine("mysql", link)
		if err != nil {
			log.Println(fmt.Sprintf("mysql:[%v] connection create failed: %v", link, err))
			os.Exit(1)
		} else if err = x.Ping(); err != nil {
			log.Println(fmt.Sprintf("mysql:[%v] connection ping failed: %v", link, err))
			os.Exit(1)
		}
		orm = x
	default:
		log.Println(fmt.Sprintf("Unknown database:%s type", key))
		os.Exit(1)
	}

	//DB.SetLogger(xorm.NewSimpleLogger(logger.GetOrmLogHandle()))

	prefix := cfg.Section(keyName).Key("prefix").String()
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, prefix)
	orm.SetTableMapper(tbMapper)

	switch mode {
	case "test":
		fallthrough
	case "dev":
		orm.ShowSQL(true)
	case "prod":
		filename := cfg.Section(keyName).Key("Filename").MustString(fmt.Sprintf("log/%v.log", keyName))
		maxSize := cfg.Section(keyName).Key("MaxSize").MustInt(100)
		maxBackups := cfg.Section(keyName).Key("MaxBackups").MustInt(15)
		maxAge := cfg.Section(keyName).Key("MaxAge").MustInt(28)
	
		output := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge, // days
		}
		orm.SetLogger(xorm.NewSimpleLogger(output))

		orm.ShowSQL(false)
	}
	level := cfg.Section(keyName).Key("Level").MustString("info")
	orm.Logger().SetLevel(getXormLevel(level))

	maxIdleConns := cfg.Section(keyName).Key("MaxIdleConns").MustInt(10)
	maxOpenConns := cfg.Section(keyName).Key("MaxOpenConns").MustInt(100)

	orm.SetMaxIdleConns(maxIdleConns)
	orm.SetMaxOpenConns(maxOpenConns)

	engineMap.Store(key, orm)
}

// getXormLevel return xorm logger level
func getXormLevel(level string) core.LogLevel {
	switch level {
	case "debug", "DEBUG":
		return core.LOG_DEBUG
	case "info", "INFO", "": // make the zero value useful
		return core.LOG_INFO
	case "warn", "WARN":
		return core.LOG_WARNING
	case "error", "ERROR":
		return core.LOG_ERR
	case "off", "OFF":
		return core.LOG_OFF
	default:
		return core.LOG_UNKNOWN
	}
}