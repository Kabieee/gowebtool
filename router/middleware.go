package router

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckGitHubToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		githubSign := c.GetHeader("X-Hub-Signature-256")
		giteeSign := c.GetHeader("X-Gitee-Token")
		if githubSign == "" && giteeSign == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Code":    http.StatusUnauthorized,
				"Message": http.StatusText(http.StatusUnauthorized),
				"Data":    "missing sign",
			})
			return
		}
		body, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if json.Valid(body) {
			buf := bytes.NewBuffer(nil)
			_ = json.Compact(buf, body)
			body = buf.Bytes()
		}

		key := "git-push-pwd-666"
		if githubSign != "" {
			hash := hmac.New(sha256.New, []byte(key))
			hash.Write(body)
			hashString := "sha256=" + hex.EncodeToString(hash.Sum(nil))
			if !hmac.Equal([]byte(githubSign), []byte(hashString)) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"Code":    http.StatusUnauthorized,
					"Message": http.StatusText(http.StatusUnauthorized),
					"Data":    fmt.Sprintf("sign not match github %s != %s", githubSign, hashString),
				})
				return
			}
			c.Next()
			return
		}

		if giteeSign != "" {
			if giteeSign != key {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"Code":    http.StatusUnauthorized,
					"Message": http.StatusText(http.StatusUnauthorized),
					"Data":    fmt.Sprintf("sign not match gitee"),
				})
				return
			}
			c.Next()
			return
		}
		c.Abort()
	}

}
