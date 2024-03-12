package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	"os"
	"strings"
)

func OAuth2Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is missing"})
			return
		}

		clientID := os.Getenv("CLIENT_ID")
		if clientID == "" {
			log.Fatal("CLIENT_ID not set in environment")
			c.AbortWithStatus(500)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := idtoken.Validate(c, idToken, clientID)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			log.Println(err.Error())
			return
		}
		c.Next()
	}

}
