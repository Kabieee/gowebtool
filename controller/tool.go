package controller

import (
	"bytes"
	"fmt"
	"gowebtool/task"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en2 "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type ToolController struct {
	BaseController
}

func (t *ToolController) SendEmail(c *gin.Context) {
	type RequestData struct {
		To   string `json:"to" form:"to" binding:"required,email"`
		Body string `json:"body" form:"body" binding:"required"`
	}
	var tran ut.Translator
	if va, ok := binding.Validator.Engine().(*validator.Validate); ok {
		eng := en.New()
		uni := ut.New(eng, eng)
		tran, _ = uni.GetTranslator(eng.Locale())
		en2.RegisterDefaultTranslations(va, tran)

	}
	var data RequestData
	err := c.ShouldBind(&data)
	if err != nil {
		for _, filedError := range err.(validator.ValidationErrors) {
			t.Fail(c, &Fail{ErrorInfo: filedError.Translate(tran)})
			return
		}
	}
	data.Body = strings.Trim(data.Body, " ")
	if data.Body == "" {
		t.Fail(c, &Fail{ErrorInfo: "body不能为空"})
		return
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", "send@email.makemake.in", "Send")
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", "Notify")
	content := fmt.Sprintf("<h3>Notify Time: %s</h3><br><p>%s</p>", time.Now().Format(time.RFC3339), data.Body)
	m.SetBody("text/html", content)

	task.EmailChan <- &task.Email{
		Message: m,
		To:      data.To,
	}

	buf := bytes.NewBuffer(nil)
	_, err = m.WriteTo(buf)
	if err != nil {
		t.Fail(c, &Fail{ErrorInfo: err.Error()})
		return
	}
	reg := regexp.MustCompile("\r\n")
	split := reg.Split(buf.String(), -1)
	t.Success(c, &Success{Data: split})
}
