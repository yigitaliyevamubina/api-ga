package v1

import (
	"apii_gateway/api/handlers/models"
	"apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/email"
	pb "apii_gateway/genproto/user_service"
	"apii_gateway/pkg/etc"
	"apii_gateway/pkg/logger"
	"apii_gateway/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

//rpc Create(User) returns (User);
//rpc GetUserById(GetUserId) returns (UserWithPostsAndComments);
//rpc UpdateUser(User) returns (User);
//rpc DeleteUser(GetUserId) returns (User);
//rpc GetAllUsers(GetAllUsersRequest) returns (AllUsers);
//rpc CheckField(Request) returns (Response);

// Register User
// @Summary register user
// @Tags User
// @Description Register a new user with the provided details
// @Accept json
// @Produce json
// @Param UserInfo body models.User true "Register user"
// @Success 201 {object} models.RegisterRespUser
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
// @Router /v1/user/register [post]
func (h *handlerV1) Register(c *gin.Context) {
	var (
		body       models.User
		code       string
		jspMarshal protojson.MarshalOptions
	)
	jspMarshal.UseProtoNames = true

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", logger.Error(err))
		return
	}
	fmt.Println(body)
	body.Id = uuid.NewString()

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	err = body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		h.log.Error("error while validating", logger.Error(err))
		return
	}

	exists, err := h.serviceManager.UserService().CheckField(ctx, &pb.Request{
		Field: "email",
		Data:  body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check email uniqueness", logger.Error(err))
		return
	}
	if exists.Resp {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email is already in use. Try another email.",
		})
		h.log.Error("email is not unique", logger.Error(err))
		return
	}
	code = utils.GenerateCode(5)
	respUser := models.RegisterUserModel{
		Id:        body.Id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Age:       body.Age,
		Gender:    body.Gender,
		Email:     body.Email,
		Password:  body.Password,
		Code:      code,
	}

	userJson, err := json.Marshal(respUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot bind json", logger.Error(err))
		return
	}

	fmt.Println(respUser)

	timeOut := time.Second * 1000

	err = h.inMemoryStorage.SetWithTTL(respUser.Email, string(userJson), int(timeOut.Seconds()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot set with ttl", logger.Error(err))
		return
	}

	from := "mubinayigitaliyeva00@gmail.com"
	password := "iocd vnhb lnvx digm"
	err = email.SendVerificationCode(email.Params{
		From:     from,
		Password: password,
		To:       respUser.Email,
		Message:  fmt.Sprintf("Hi %s,", respUser.FirstName),
		Code:     respUser.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot send a code to an email", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.RegisterRespUser{Status: "a verification code was sent to your email, please check it."})
}

// Verify User
// @Summary verify user
// @Tags User
// @Description Verify a user with code sent to their email
// @Accept json
// @Product json
// @Param email path string true "email"
// @Param code path string true "code"
// @Success 201 {object} models.VerifyResponse
// @Failure 400 string error models.Error
// @Failure 400 string error models.Error
// @Router /v1/user/verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions //default -> lowerCaseName
	jspbMarshal.UseProtoNames = true

	userEmail := c.Param("email")
	code := c.Param("code")

	user, err := redis.Bytes(h.inMemoryStorage.Get(userEmail))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.VerifyResponse{
			Status: "Code is expired, try again.",
		})
		h.log.Error("Code is expired, TTL is over.")
		return
	}

	var respUser models.RegisterUserModel
	if err := json.Unmarshal(user, &respUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot unmarshal user from redis", logger.Error(err))
		fmt.Println(respUser)
		return
	}

	if respUser.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "verification failed",
		})
		h.log.Error("verification failed", logger.Error(err))
		return
	}

	respUser.Password, err = etc.HashPassword(respUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot hash the password", logger.Error(err))
		return
	}

	fmt.Println(respUser)
	h.jwtHandler = tokens.JWTHandler{
		Sub:       respUser.Id,
		Role:      "user",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create access and refresh token", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeOut))
	defer cancel()

	userResp := models.VerifyUserModel{
		Id:        respUser.Id,
		FirstName: respUser.FirstName,
		LastName:  respUser.LastName,
		Age:       respUser.Age,
		Gender:    respUser.Gender,
		Email:     respUser.Email,
		Password:  respUser.Password,
	}

	userResp.RefreshToken = refresh
	userResp.AccessToken = access

	_, err = h.serviceManager.UserService().Create(ctx, &pb.User{
		Id:           userResp.Id,
		FirstName:    userResp.FirstName,
		LastName:     userResp.LastName,
		Age:          userResp.Age,
		Gender:       pb.Gender(userResp.Gender),
		Email:        userResp.Email,
		Password:     userResp.Password,
		RefreshToken: userResp.RefreshToken,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create user", logger.Error(err))
		return
	}

	res := models.UserModel{
		Id:          userResp.Id,
		FirstName:   userResp.FirstName,
		LastName:    userResp.LastName,
		Age:         userResp.Age,
		Gender:      userResp.Gender,
		Email:       userResp.Email,
		Password:    userResp.Password,
		AccessToken: userResp.AccessToken,
	}

	c.JSON(http.StatusOK, res)
}

// Login User
// @Summary login user
// @Tags User
// @Description Login
// @Accept json
// @Produce json
// @Param User body models.LoginRequest true "Login"
// @Success 201 {object} models.LoginResponse
// @Failure 400 string Error models.Error
// @Failure 400 string Error models.Error
// @Router /v1/user/login [post]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
		body        models.LoginRequest
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

	resp, err := h.serviceManager.UserService().GetUserByEmail(ctx, &pb.GetUserEmailReq{
		Email: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get user by email", logger.Error(err))
		return
	}

	if !etc.CompareHashPassword(resp.Password, body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong password",
		})
		h.log.Error("wrong password", logger.Error(err))
		return
	}
	h.jwtHandler = tokens.JWTHandler{
		Sub:       resp.Id,
		Role:      "user",
		SignInKey: h.cfg.SignInKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, _, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create access token", logger.Error(err))
		return
	}

	res := models.LoginResp{
		Id:          resp.Id,
		Email:       resp.Email,
		Password:    resp.Password,
		AccessToken: access,
	}

	c.JSON(http.StatusOK, res)
}

// CreateUser
// @Router /v1/user/create [post]
// @Security ApiKeyAuth
// @Summary create user
// @Tags User
// @Description Create a new user with the provided details
// @Accept json
// @Produce json
// @Param UserInfo body models.User true "Create user"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
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
// @Router /v1/user/{id} [get]
// @Security ApiKeyAuth
// @Summary get user by id
// @Tags User
// @Description Get user
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} models.UserWithPostsAndComments
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
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
// @Router /v1/user/update [put]
// @Security ApiKeyAuth
// @Summary update user
// @Tags User
// @Description Update user
// @Accept json
// @Produce json
// @Param UserInfo body models.User true "Update User"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
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
// @Router /v1/user/delete/{id} [delete]
// @Security ApiKeyAuth
// @Summary delete user
// @Tags User
// @Description Delete user
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 201 {object} models.User
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
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
// @Router /v1/users [get]
// @Security ApiKeyAuth
// @Summary get all users with posts and comments
// @Tags User
// @Description get all users
// @Accept json
// @Produce json
// @Success 201 {object} models.AllUsers
// @Failure 400 string Error models.Error
// @Failure 500 string Error models.Error
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
