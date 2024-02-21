package main

import (
	"apii_gateway/api"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/services"
	"apii_gateway/storage/redis"
	"fmt"

	rds "github.com/gomodule/redigo/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	redisPool := rds.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (rds.Conn, error) {
			c, err := rds.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	server := api.New(api.Option{
		InMemory:       redis.NewRedisRepo(&redisPool),
		Cfg:            cfg,
		Logger:         log,
		ServiceManager: serviceManager,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		panic(err)
	}
}
