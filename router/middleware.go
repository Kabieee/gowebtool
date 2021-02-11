package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckGitHubToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		githubSign := c.GetHeader("X-Hub-Signature-256")
		giteeSign := c.GetHeader("X-Gitee-Token")
		if githubSign == "" && giteeSign == "" {
			c.AbortWithStatusJSON(403, gin.H{
				"Code":    http.StatusForbidden,
				"Message": http.StatusText(http.StatusForbidden),
			})
			return
		}
		c.Next()
	}

}
