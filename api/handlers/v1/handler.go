package v1

import (
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/services"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
}

type HandlerV1Config struct {
	Log            logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Log,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
	}
}
