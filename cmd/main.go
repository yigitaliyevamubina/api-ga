package main

import (
	"apii_gateway/api"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/services"
	"apii_gateway/storage/redis"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"

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


	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "mubina2007", "postgres")


	_, err = gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		log.Fatal("error while updating new adapter", logger.Error(err))
		return
	}

	enforcer, err := casbin.NewEnforcer(cfg.AuthConfigPath, cfg.AuthCSVPath)
	if err != nil {
		log.Error("cannot create a new enforcer", logger.Error(err))
		return
	}


	server := api.New(api.Option{
		InMemory:       redis.NewRedisRepo(&redisPool),
		Cfg:            cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		CasbinEnforser: enforcer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		panic(err)
	}
}
