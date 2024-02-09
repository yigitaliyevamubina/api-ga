package api

import (
	v1 "apii_gateway/api/handlers/v1"
	"apii_gateway/config"
	"apii_gateway/pkg/logger"
	"apii_gateway/services"
	"github.com/gin-gonic/gin"
)

// Option Struct
type Option struct {
	Cfg            config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
}

// New -> constructor
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Log:            option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Cfg,
	})

	api := router.Group("/v1")

	//users
	api.POST("/user/create", handlerV1.CreateUser)
	api.GET("/user/:id", handlerV1.GetUserById)
	api.PUT("/user/update/:id", handlerV1.UpdateUser)
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)
	api.GET("/users", handlerV1.GetAllUsers)

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
	api.GET("/like/post/:post_id", handlerV1.GetLikeOwnersByPostId)
	api.GET("/like/comment/:comment_id", handlerV1.GetLikeOwnersByCommentId)

	//url:=ginSwagger.URL("swagger/doc.json")
	//router.GET("/swagger/*any, ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
