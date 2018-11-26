package redis

import (
	"fmt"
	"time"

	"github.com/logicinu/nest/module/setting"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

// InitRedisPool init redis pool
func InitRedisPool() {
	cfg := setting.GetSetting()

	network := cfg.Section("redis").Key("Network").MustString("tcp")
	host := cfg.Section("redis").Key("Host").MustString("127.0.0.1")
	port := cfg.Section("redis").Key("Port").MustString("6397")
	password := cfg.Section("redis").Key("Password").MustString("")
	db := cfg.Section("redis").Key("DB").MustString("nest")
	maxIdle := cfg.Section("redis").Key("MaxIdle").MustInt(3)
	maxActive := cfg.Section("redis").Key("MaxActive").MustInt(5)
	idleTimeout := cfg.Section("redis").Key("IdleTimeout").MustInt(240)
	wait := cfg.Section("redis").Key("Wait").MustBool(false)

	redisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Wait:        wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, fmt.Sprintf("%v:%v", host, port))
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err := c.Do("PING"); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// GetRedisPool return redis pool
func GetRedisPool() *redis.Pool {
	return redisPool
}
