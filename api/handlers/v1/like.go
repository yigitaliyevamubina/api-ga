package v1

import (
	"apii_gateway/api/handlers/models"
	pbl "apii_gateway/genproto/like_service"
	"apii_gateway/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"
)

//rpc LikePost(PostLike) returns (Status);
//rpc LikeComment(CommentLike) returns (Status);
//rpc GetLikeOwnersByPostId(GetPostId) returns (Post);
//rpc GetLikeOwnersByCommentId(GetCommentId) returns (Comment);

// @Router /v1/like/post [post]
// @Security ApiKeyAuth
// @Summary like post
// @Tags Like
// @Description Like post
// @Accept json
// @Produce json
// @Param LikeInfo body models.PostLike true "Like post"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/like/post [post]
func (h *handlerV1) LikePost(c *gin.Context) {

	var (
		body        models.PostLike
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.LikeService().LikePost(ctx, &pbl.PostLike{
		PostId:  body.PostId,
		OwnerId: body.UserId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		h.log.Error("cannot create likepost", logger.Error(err))
		return
	}
	c.JSON(http.StatusCreated, response)
}

// @Router /v1/like/comment [post]
// @Security ApiKeyAuth
// @Summary like comment
// @Tags Like
// @Description Like comment
// @Accept json
// @Produce json
// @Param LikeInfo body models.CommentLike true "Like comment"
// @Success 201 {object} models.Status
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) LikeComment(c *gin.Context) {
	var (
		body        models.CommentLike
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		h.log.Error("cannot create like comment", logger.Error(err))
		return
	}

	response, err := h.serviceManager.LikeService().LikeComment(ctx, &pbl.CommentLike{
		CommentId: body.CommentId,
		OwnerId:   body.UserId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		h.log.Error("cannot create like comment", logger.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)

}

// @Router /v1/like/post/{id} [get]
// @Security ApiKeyAuth
// @Summary like owners by post id
// @Tags Like
// @Description Like owners by post id
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.ResponseLikePost
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) GetLikeOwnersByPostId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	postId := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.LikeService().GetLikeOwnersByPostId(ctx, &pbl.GetPostId{
		PostId: postId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		h.log.Error("cannot get like owners by post id", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Router /v1/like/comment/{id} [get]
// @Security ApiKeyAuth
// @Summary like owners by comment id
// @Tags Like
// @Description Like owners by comment id
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} models.ResponseLikeComment
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
func (h *handlerV1) GetLikeOwnersByCommentId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	commentId := c.Param("id")

	response, err := h.serviceManager.LikeService().GetLikeOwnersByCommentId(ctx, &pbl.GetCommentId{
		CommentId: commentId,
	})

	fmt.Println(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get like owners by comment id", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
