package main

import (
	"apii_gateway/api"
	"apii_gateway/config"
	"apii_gateway/pkg/db"
	"apii_gateway/pkg/etc"
	"apii_gateway/pkg/logger"
	"apii_gateway/queue/producer"

	// "apii_gateway/rabbitmq"
	"apii_gateway/services"
	admin "apii_gateway/storage/postgres"
	"apii_gateway/storage/redis"
	"fmt"

	// "github.com/casbin/casbin/v2"
	// gormadapter "github.com/casbin/gorm-adapter/v3"

	rds "github.com/gomodule/redigo/redis"
	// "github.com/streadway/amqp"
)

func main() {
	fmt.Println(etc.HashPassword("superadminpass"))
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

	db, _, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("cannot connect to DB", logger.Error(err))
		panic(err)
	}

	writer, err := producer.NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		log.Fatal("cannot create a kafka producer", logger.Error(err))
	}

	// err = writer.ProduceMessages("producer", []byte("testsfksdmfvkotmovwmtiomvwoirtmvwio"))
	// if err != nil {
	// 	panic(err)
	// }

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// channel, err := conn.Channel()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// rabbit := rabbitmq.NewRabbitMQProducer(channel)
	server := api.New(api.Option{
		InMemory:       redis.NewRedisRepo(&redisPool),
		Cfg:            cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		Postgres:       admin.NewAdminRepo(db),
		Producer:       writer,
		// Rabbit:         rabbit,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		panic(err)
	}


}
