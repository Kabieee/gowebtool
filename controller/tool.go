package controller

import (
	"gowebtool/common"
	"gowebtool/services"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

type ToolController struct {
	BaseController
}

func (t *ToolController) SendEmail(c *gin.Context) {
	e := &services.EmailServices{}
	e.Ctx = c
	var data services.EmailData
	err := c.ShouldBind(&data)
	if err != nil {
		for _, filedError := range err.(validator.ValidationErrors) {
			t.Fail(c, &Fail{ErrorInfo: filedError.Translate(common.MyTran)})
			return
		}
	}
	data.Body = strings.Trim(data.Body, " ")
	if data.Body == "" {
		t.Fail(c, &Fail{ErrorInfo: "body不能为空"})
		return
	}
	err, result := e.Send(data)
	if err != nil {
		t.Fail(c, &Fail{ErrorInfo: err.Error()})
		return
	}
	t.Success(c, &Success{Data: result})
}
