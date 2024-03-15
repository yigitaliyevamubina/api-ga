package api

import (
	// casb "apii_gateway/api/casbin"
	_ "apii_gateway/api/docs"
	v1 "apii_gateway/api/handlers/v1"
	"apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/queue/producer"
	"apii_gateway/rabbitmq"
	"apii_gateway/services"
	"apii_gateway/storage/repo"

	casb "apii_gateway/api/casbin"

	admin "apii_gateway/storage/postgresrepo"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Option Struct
type Option struct {
	InMemory       repo.InMemoryStorageI
	Cfg            config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Postgres       admin.AdminStorageI
	Producer       producer.KafkaProducer
	Rabbit         rabbitmq.Producer
}

// New -> constructor
// @title Welcome to services
// @version 1.0
// @description microservice
// @host localhost:3030
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	// psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	option.Cfg.PostgresHost,
	// 	option.Cfg.PostgresPort,
	// 	option.Cfg.PostgresUser,
	// 	option.Cfg.PostgresPassword,
	// 	option.Cfg.PostgresDatabase)

	// adapter, err := gormadapter.NewAdapter("postgres", psqlString, true)
	// if err != nil {
	// 	option.Logger.Fatal("error while updating new adapter", logger.Error(err))
	// }

	casbinEnforcer, err := casbin.NewEnforcer(option.Cfg.AuthConfigPath, option.Cfg.AuthCSVPath)
	if err != nil {
		option.Logger.Error("cannot create a new enforcer", logger.Error(err))
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}

	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch3", util.KeyMatch3)

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandle := tokens.JWTHandler{
		SignInKey: option.Cfg.SignInKey,
		Log:       option.Logger,
	}

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// channel, err := conn.Channel()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rabbit := rabbitmq.NewRabbitMQProducer(channel)

	handlerV1 := v1.New(&v1.HandlerV1Config{
		InMemoryStorage: option.InMemory,
		Log:             option.Logger,
		ServiceManager:  option.ServiceManager,
		Cfg:             option.Cfg,
		JWTHandler:      jwtHandle,
		Postgres:        option.Postgres,
		Casbin:          casbinEnforcer,
		Producer:        option.Producer,
		// Rabbit:          option.Rabbit,
	})

	api := router.Group("/v1")

	api.Use(casb.NewAuth(casbinEnforcer, option.Cfg))

	//rbac
	api.GET("/rbac/roles", handlerV1.ListAllRoles)              //superadmin
	api.GET("/rbac/policies/:role", handlerV1.ListRolePolicies) //superadmin
	api.POST("/rbac/add/policy", handlerV1.AddPolicyToRole)     //superadmin
	api.DELETE("/rbac/delete/policy", handlerV1.DeletePolicy)   //superadmin

	//users
	api.POST("/user/create", handlerV1.CreateUser)         //admin
	api.POST("/user/register", handlerV1.Register)         //unauthorized
	api.GET("/user/:id", handlerV1.GetUserById)            //user
	api.PUT("/user/update", handlerV1.UpdateUser)          //user
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)   //admin
	api.GET("/users", handlerV1.GetAllUsers)               //admin
	api.GET("/user/verify/:email/:code", handlerV1.Verify) //unauthorized
	api.POST("/user/login", handlerV1.Login)               //unauthorized

	//posts
	api.POST("/post/create", handlerV1.CreatePost)          //user
	api.GET("/post/get/:id", handlerV1.GetPostById)         //user
	api.PUT("/post/update/:id", handlerV1.UpdatePost)       //user
	api.DELETE("/post/delete/:id", handlerV1.DeletePost)    //user
	api.GET("/post/owner/:id", handlerV1.GetPostsByOwnerId) //user

	//comments
	api.POST("/comment/create", handlerV1.CreateComment)             //user
	api.GET("/comment/post/:id", handlerV1.GetAllCommentsByPostId)   //user
	api.GET("/comment/owner/:id", handlerV1.GetAllCommentsByOwnerId) //user

	//likes
	api.POST("/like/post", handlerV1.LikePost)                       //user
	api.POST("/like/comment", handlerV1.LikeComment)                 //user
	api.GET("/like/post/:id", handlerV1.GetLikeOwnersByPostId)       //user
	api.GET("/like/comment/:id", handlerV1.GetLikeOwnersByCommentId) //user

	//admin
	api.POST("/auth/create", handlerV1.CreateAdmin)   //superadmin
	api.DELETE("/auth/delete", handlerV1.DeleteAdmin) //superadmin
	api.POST("/auth/login", handlerV1.LoginAdmin)     //admin

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
