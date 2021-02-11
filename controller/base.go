package controller

import (
	"encoding/json"
	"fmt"
	"gowebtool/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

type Success struct {
	Code    int
	Message string
	Data    interface{}
	Extend  interface{} `json:"Extend,omitempty"`
}

type Fail struct {
	Code      int
	Message   string
	ErrorInfo interface{}
}

func (b *BaseController) Success(c *gin.Context, s *Success) {
	if s.Code == 0 {
		s.Code = 2000
	}
	if s.Message == "" {
		s.Message = "success"
	}
	if s.Data == nil {
		s.Data = ""
	}
	c.AbortWithStatusJSON(http.StatusOK, s)
}

func (b *BaseController) Fail(c *gin.Context, f *Fail, status ...int) {
	if f.Code == 0 {
		f.Code = 5000
	}
	if f.Message == "" {
		f.Message = "fail"
	}
	if f.ErrorInfo == nil {
		f.ErrorInfo = ""
	}
	statusCode := http.StatusOK
	if len(status) > 0 && http.StatusText(status[0]) != "" {
		statusCode = status[0]
	}
	c.AbortWithStatusJSON(statusCode, f)
}

type TempData struct {
	Method      string
	IP          string
	ContentType string
	Agent       string
	Error       string
	//Body        map[string]interface{}
	//Query       url.Values
}

const execPath = "/home/lzfeng/devgo/gowebtool/shell"

var (
	temp *TempData
)

func (b *BaseController) ParseRequest(c *gin.Context) *TempData {
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
	//temp.Body = bodyMap
	//temp.Query = c.Request.URL.Query()
	return temp
}

func (b *BaseController) Index(c *gin.Context) {
	b.ParseRequest(c)
	b.Success(c, &Success{Data: temp})
}

func (b *BaseController) User(c *gin.Context) {
	b.ParseRequest(c)
	tk := c.Query("tk")
	if tk != "user666" {
		b.Success(c, &Success{Data: temp})
		return
	}
	cmd := fmt.Sprintf("id")
	execCmd := common.ExecCmd(cmd)
	b.Success(c, &Success{Data: temp, Extend: execCmd})

}
