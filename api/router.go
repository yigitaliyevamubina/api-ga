package api

import (
	_ "apii_gateway/api/docs"
	v1 "apii_gateway/api/handlers/v1"
	"apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/api/middleware"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/services"
	"apii_gateway/storage/repo"

	// casb "apii-gateway/middleware/casbin"

	"github.com/casbin/casbin/v2"
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
	CasbinEnforser *casbin.Enforcer
}

// New -> constructor
// @title Welcome to services
// @version 1.0
// @description In this swagger documentation you can test all of your microservices
// @host localhost:5555
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandle := tokens.JWTHandler{
		SignInKey: option.Cfg.SignInKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		InMemoryStorage: option.InMemory,
		Log:             option.Logger,
		ServiceManager:  option.ServiceManager,
		Cfg:             option.Cfg,
		JWTHandler:      jwtHandle,
		CasbinEnforcer:  option.CasbinEnforser,
	})

	api := router.Group("/v1")

	api.Use(middleware.Auth)
	//users
	api.POST("/user/create", handlerV1.CreateUser)
	api.POST("/user/register", handlerV1.Register)
	api.GET("/user/:id", handlerV1.GetUserById)
	api.PUT("/user/update", handlerV1.UpdateUser)
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)
	api.GET("/users", handlerV1.GetAllUsers)
	api.GET("/user/verify/:email/:code", handlerV1.Verify)
	api.POST("/user/login", handlerV1.Login)

	//posts
	api.POST("/post/create", handlerV1.CreatePost)
	api.GET("/post/get/:id", handlerV1.GetPostById)
	api.PUT("/post/update/:id", handlerV1.UpdatePost)
	api.DELETE("/post/delete/:id", handlerV1.DeletePost)
	api.GET("/post/owner/:id", handlerV1.GetPostsByOwnerId)

	//comments
	api.POST("/comment/create", handlerV1.CreateComment)
	api.GET("/comment/post/:id", handlerV1.GetAllCommentsByPostId)
	api.GET("/comment/owner/:id", handlerV1.GetAllCommentsByOwnerId)

	//likes
	api.POST("/like/post", handlerV1.LikePost)
	api.POST("/like/comment", handlerV1.LikeComment)
	api.GET("/like/post/:id", handlerV1.GetLikeOwnersByPostId)
	api.GET("/like/comment/:id", handlerV1.GetLikeOwnersByCommentId)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
