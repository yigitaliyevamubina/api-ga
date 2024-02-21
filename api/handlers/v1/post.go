package v1

import (
	"apii_gateway/api/handlers/models"
	pbp "apii_gateway/genproto/post_service"
	"apii_gateway/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

//rpc CreatePost(ReqPost) returns (RespPost);
//rpc UpdatePost(ReqPost) returns (ReqPost);
//rpc DeletePost(GetPostId) returns (ReqPost);
//rpc GetPostById(GetPostId) returns (RespPost);
//rpc GetPostsByOwnerId(GetOwnerId) returns (OwnerPosts);

// Create Post
// @Router /v1/post/create [post]
// @Security ApiKeyAuth
// @Summary create post
// @Tags Post
// @Description Create Post
// @Accept json
// @Produce json
// @Param PostInfo body models.Post true "Create Post"
// @Success 201 {object} models.RespPost
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		body        models.Post
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	body.Id = uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.PostService().CreatePost(ctx, &pbp.ReqPost{
		Id:       body.Id,
		Title:    body.Title,
		ImageUrl: body.ImageUrl,
		OwnerId:  body.OwnerId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create post", logger.Error(err))
	}

	c.JSON(http.StatusCreated, response)
}

// Update User
// @Router /v1/post/update [put]
// @Security ApiKeyAuth
// @Summary update post
// @Tags Post
// @Description Update Post
// @Accept json
// @Produce json
// @Param PostInfo body models.Post true "Create Post"
// @Success 201 {object} models.Post
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		body        pbp.ReqPost
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	body.Id = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.PostService().UpdatePost(ctx, &pbp.ReqPost{
		Id:       body.Id,
		Title:    body.Title,
		ImageUrl: body.ImageUrl,
		OwnerId:  body.OwnerId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot update post", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete User
// @Router /v1/post/delete/{id} [delete]
// @Security ApiKeyAuth
// @Summary delete post
// @Tags Post
// @Description Delete Post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.Post
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) DeletePost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.PostService().DeletePost(ctx, &pbp.GetPostId{
		PostId: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		h.log.Error("cannot delete post", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get Post By Id
// @Router /v1/post/get/{id} [get]
// @Security ApiKeyAuth
// @Summary get post
// @Tags Post
// @Description Get Post
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.RespPost
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) GetPostById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.PostService().GetPostById(ctx, &pbp.GetPostId{
		PostId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get post by id", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get All Posts By Owner Id
// @Router /v1/post/owner/{id} [get]
// @Security ApiKeyAuth
// @Summary get posts by owner id
// @Tags Post
// @Description Get Posts by owner id
// @Accept json
// @Produce json
// @Param id path string true "owner_id"
// @Success 201 {object} models.Posts
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) GetPostsByOwnerId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ownerId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.PostService().GetPostsByOwnerId(ctx, &pbp.GetOwnerId{
		OwnerId: ownerId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get posts by owner id", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}
