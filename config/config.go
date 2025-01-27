package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string //develop, staging, production

	RedisHost string
	RedisPort int

	UserServiceHost string
	UserServicePort int

	PostServiceHost string
	PostServicePort int

	CommentServiceHost string
	CommentServicePort int

	LikeServiceHost string
	LikeServicePort int

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	//context timeout in seconds
	CtxTimeOut int

	LogLevel string
	HTTPPort string

	AccessTokenTimeout  int //minutes
	RefreshTokenTimeout int //hours
	AuthConfigPath      string
	AuthCSVPath         string

	SignInKey string

	RabbitQueue string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":3030"))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redisdb"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "userdb"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "mubina2007"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user-service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 7070))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post-service"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", 8080))

	c.CommentServiceHost = cast.ToString(getOrReturnDefault("COMMENT_SERVICE_HOST", "comment-service"))
	c.CommentServicePort = cast.ToInt(getOrReturnDefault("COMMENT_SERVICE_PORT", 8088))

	c.LikeServiceHost = cast.ToString(getOrReturnDefault("LIKE_SERVICE_HOST", "like-service"))
	c.LikeServicePort = cast.ToInt(getOrReturnDefault("LIKE_SERVICE_PORT", 4040))

	c.AccessTokenTimeout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 500))
	c.RefreshTokenTimeout = cast.ToInt(getOrReturnDefault("REFRESH_TOKEN_TIMEOUT", 3))

	c.CtxTimeOut = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "AAAAKWEFWEKkfkw"))

	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))
	c.AuthCSVPath = cast.ToString(getOrReturnDefault("AUTH_CSV_PATH", "./config/auth.csv"))

	c.RabbitQueue = cast.ToString(getOrReturnDefault("RABBITMQ_QUEUE", "golang"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return defaultValue
	}

	return defaultValue
}
