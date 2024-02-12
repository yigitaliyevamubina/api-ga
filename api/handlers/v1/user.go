package v1

import (
	"apii_gateway/api/handlers/models"
	pb "apii_gateway/genproto/user_service"
	"apii_gateway/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

//rpc Create(User) returns (User);
//rpc GetUserById(GetUserId) returns (UserWithPostsAndComments);
//rpc UpdateUser(User) returns (User);
//rpc DeleteUser(GetUserId) returns (User);
//rpc GetAllUsers(GetAllUsersRequest) returns (AllUsers);

// CreateUser
// @Summary create user
// @Tags User
// @Description Create a new user with the provided details
// @Accept json
// @Produce json
// @Param UserInfo body models.User true "Create user"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/user/create [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.User
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

	response, err := h.serviceManager.UserService().Create(ctx, &pb.User{
		Id:        body.Id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Age:       body.Age,
		Gender:    pb.Gender(body.Gender),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create user", logger.Error(err))
	}

	c.JSON(http.StatusCreated, response)
}

// Get User By Id
// @Summary get user by id
// @Tags User
// @Description Get user
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} models.UserWithPostsAndComments
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/user/{id} [get]
func (h *handlerV1) GetUserById(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()
	response, err := h.serviceManager.UserService().GetUserById(ctx, &pb.GetUserId{
		UserId: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update User

// @Summary update user
// @Tags User
// @Description Update user
// @Accept json
// @Produce json
// @Param UserInfo body models.User true "Update User"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/user/update [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pb.User
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

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &pb.User{
		Id:        body.Id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Age:       body.Age,
		Gender:    pb.Gender(body.Gender),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot update user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete User
// @Summary delete user
// @Tags User
// @Description Delete user
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/user/delete/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.UserService().DeleteUser(ctx, &pb.GetUserId{
		UserId: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		h.log.Error("cannot delete user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// Get All Users with Posts and Comments
// @Summary get all users with posts and comments
// @Tags User
// @Description get all users
// @Accept json
// @Produce json
// @Success 201 {object} models.AllUsers
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/users [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	response, err := h.serviceManager.UserService().GetAllUsers(ctx, &pb.GetAllUsersRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get all users", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
