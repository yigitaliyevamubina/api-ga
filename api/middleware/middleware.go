package middleware

import (
	"apii_gateway/api/handlers/v1/tokens"
	"apii_gateway/config"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type ACL struct {
	Role   string
	Api    string
	Method string
}

func Auth(ctx *gin.Context) {
	if ctx.Request.URL.Path == "/v1/user/register" || ctx.Request.URL.Path == "/v1/user/login" ||
		ctx.Request.URL.Path == "/v1/swagger/swagger/doc.json" || ctx.Request.URL.Path == "/v1/swagger/index.html" ||
		ctx.Request.URL.Path == "/v1/swagger/swagger-ui.css" || ctx.Request.URL.Path == "/v1/swagger/swagger-ui-bundle.js" || ctx.Request.URL.Path == "/v1/swagger/swagger-ui-standalone-preset.js" {
		ctx.Next()
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var trimmedToken string
	suffix := "Bearer "
	if strings.Contains(token, suffix) {
		trimmedToken = strings.TrimPrefix(token, suffix)
	} else {
		trimmedToken = token
	}

	cfg := config.Load()

	_, err := tokens.ExtractClaim(trimmedToken, []byte(cfg.SignInKey))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	ctx.Next()
}

func CheckACL(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unauthorized",
		})
		return
	}

	cfg := config.Load()
	claims, err := tokens.ExtractClaim(token, []byte(cfg.SignInKey))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid token",
		})
		return
	}

	acls, err := ReadCSVFile()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot access ACL store",
		})
		return
	}

	if (cast.ToString(claims["role"]) == "user" && ctx.Request.URL.Path == acls[0].Api) && (ctx.Request.Method == acls[0].Method) {
		ctx.Next()
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "this method is not allowed for you...",
	})
	ctx.Next()
}

func ReadCSVFile() ([]ACL, error) {
	var acls []ACL
	f, err := os.Open("apii_gateway/rbac.csv")
	if err != nil {
		log.Fatal("error opening file:", err.Error())
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		rec, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}

		if len(rec) != 6 {
			log.Fatalf("invalid record: %v", rec)
		}
		acl := ACL{
			Role:   rec[0],
			Api:    rec[1],
			Method: rec[2],
		}

		acls = append(acls, acl)
	}

	return acls, nil
}
