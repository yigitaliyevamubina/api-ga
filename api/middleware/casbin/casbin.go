package casbin

import (
	"apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type CasbinHandler struct {
	cfg      config.Config
	enforcer *casbin.Enforcer
}

func (c *CasbinHandler) CheckCasbinPermission(casbin *casbin.Enforcer, cfg config.Config) gin.HandlerFunc {
	casbHandler := &CasbinHandler{
		cfg:      cfg,
		enforcer: casbin,
	}

	return func(ctx *gin.Context) {
		allowed, err := casbHandler.CheckPermission(ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{
				"error": "permission denied",
				"message": "sorry, this method is not allowed to you.",
			})
			return
		}
	}
}

func (c *CasbinHandler) GetRole(ctx *http.Request) (string, int) {
	token := ctx.Header.Get("authorization")
	if token == "" {
		return "unauthorized", 200
	}

	var cutToken string

	if strings.Contains(token, "Bearer") {
		cutToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		cutToken = token
	}

	claims, err := tokens.ExtractClaim(cutToken, []byte(c.cfg.SignInKey))
	if err != nil {
		return "unauthorized, token is invalid", http.StatusBadRequest
	}

	return cast.ToString(claims["role"]), 0
}

func (c *CasbinHandler) CheckPermission(req *http.Request) (bool, error) {
	role, status := c.GetRole(req)
	if status != 0 {
		return false, fmt.Errorf(role)
	}

	object := req.URL.Path
	action := req.Method

	response, err := c.enforcer.Enforce(role, object, action)
	if err != nil {
		return false, err
	}

	if !response {
		return false, nil
	}

	return true, nil
}
