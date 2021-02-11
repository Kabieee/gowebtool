package controller

import (
	"encoding/json"
	"fmt"
	"gowebtool/common"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type SiteController struct{}

type TempData struct {
	Method      string
	IP          string
	ContentType string
	Agent       string
	Error       string
	Body        map[string]interface{}
	Query       url.Values
}

const execPath = "/home/lzfeng/devgo/gowebtool/shell"

var (
	temp *TempData
)

func (s *SiteController) ParseRequest(c *gin.Context) *TempData {
	temp = &TempData{
		Method:      c.Request.Method,
		ContentType: c.ContentType(),
		Agent:       c.Request.UserAgent(),
		IP:          c.ClientIP(),
	}
	body, _ := c.GetRawData()
	var bodyMap = make(map[string]interface{}, 4)
	if len(body) > 0 {
		err := json.Unmarshal(body, &bodyMap)
		if err != nil {
			temp.Error = err.Error()
		}
	}
	temp.Body = bodyMap
	temp.Query = c.Request.URL.Query()
	return temp
}

func (s *SiteController) Index(c *gin.Context) {
	s.ParseRequest(c)
	c.JSON(http.StatusOK, gin.H{
		"Code":    http.StatusOK,
		"Data":    temp,
		"Message": "ok",
	})
}

func (s *SiteController) User(c *gin.Context) {
	s.ParseRequest(c)
	tk, ok := temp.Query["tk"]
	if !ok || len(tk) == 0 || tk[0] != "user666" {
		return
	}
	cmd := fmt.Sprintf("id")
	execCmd := common.ExecCmd(cmd)
	c.JSON(http.StatusOK, gin.H{
		"Code": http.StatusOK,
		"Data": temp,
		"Exec": execCmd,
	})

}
