package v1

import (
	t "apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/services"
	"apii_gateway/storage/repo"
)

type handlerV1 struct {
	inMemoryStorage repo.InMemoryStorageI
	log             logger.Logger
	serviceManager  services.IServiceManager
	cfg             config.Config
	jwtHandler      t.JWTHandler
}

type HandlerV1Config struct {
	InMemoryStorage repo.InMemoryStorageI
	Log             logger.Logger
	ServiceManager  services.IServiceManager
	Cfg             config.Config
	JWTHandler      t.JWTHandler
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		inMemoryStorage: c.InMemoryStorage,
		log:             c.Log,
		serviceManager:  c.ServiceManager,
		cfg:             c.Cfg,
		jwtHandler:      c.JWTHandler,
	}
}
