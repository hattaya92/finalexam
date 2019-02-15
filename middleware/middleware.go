package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func LoginMiddleware(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Status code is 401 Unauthorized")
		c.Abort()
		return
	}

	c.Next()
	log.Println("ending middleware")
}
