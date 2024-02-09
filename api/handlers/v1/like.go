package v1

import (
	"apii_gateway/api/handlers/models"
	pbl "apii_gateway/genproto/like_service"
	"apii_gateway/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"
)

//rpc LikePost(PostLike) returns (Status);
//rpc LikeComment(CommentLike) returns (Status);
//rpc GetLikeOwnersByPostId(GetPostId) returns (Post);
//rpc GetLikeOwnersByCommentId(GetCommentId) returns (Comment);

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

func (h *handlerV1) GetLikeOwnersByPostId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	postId := c.Param("post_id")

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

func (h *handlerV1) GetLikeOwnersByCommentId(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	commentId := c.Param("comment_id")

	response, err := h.serviceManager.LikeService().GetLikeOwnersByCommentId(ctx, &pbl.GetCommentId{
		CommentId: commentId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get like owners by comment id", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
