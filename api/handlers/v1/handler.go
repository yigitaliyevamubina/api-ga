package v1

import (
	t "apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/queue/producer"
	"apii_gateway/rabbitmq"
	"apii_gateway/services"
	admin "apii_gateway/storage/postgresrepo"
	"apii_gateway/storage/repo"

	"github.com/casbin/casbin/v2"
)

type handlerV1 struct {
	inMemoryStorage repo.InMemoryStorageI
	log             logger.Logger
	serviceManager  services.IServiceManager
	cfg             config.Config
	jwtHandler      t.JWTHandler
	postgres        admin.AdminStorageI
	casbin          *casbin.Enforcer
	producer        producer.KafkaProducer
	rabbit          rabbitmq.Producer
}

type HandlerV1Config struct {
	InMemoryStorage repo.InMemoryStorageI
	Log             logger.Logger
	ServiceManager  services.IServiceManager
	Cfg             config.Config
	JWTHandler      t.JWTHandler
	Postgres        admin.AdminStorageI
	Casbin          *casbin.Enforcer
	Producer        producer.KafkaProducer
	Rabbit          rabbitmq.Producer
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		inMemoryStorage: c.InMemoryStorage,
		log:             c.Log,
		serviceManager:  c.ServiceManager,
		cfg:             c.Cfg,
		jwtHandler:      c.JWTHandler,
		postgres:        c.Postgres,
		casbin:          c.Casbin,
		producer:        c.Producer,
		rabbit:          c.Rabbit,
	}

}
