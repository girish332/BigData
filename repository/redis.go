package repository

import (
	"github.com/gin-gonic/gin"
)

type RedisRepo interface {
	Ping(c *gin.Context) error
	Get(c *gin.Context, key string) (string, error)
	Set(c *gin.Context, key string, value string) error
	Delete(c *gin.Context, key string) error
	Keys(c *gin.Context, pattern string) ([]string, error)
}
