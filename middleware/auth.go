package middleware

import "github.com/gin-gonic/gin"

func AuthToken(c *gin.Context) {
	name := c.GetHeader("x-forward-user")
	token := c.GetHeader("authentication")

	if name == "" || token == "" {

	}

	c.Next()
}
