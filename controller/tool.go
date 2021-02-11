package controller

import (
	"bytes"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type ToolController struct {
	BaseController
}

func (t *ToolController) SendEmail(c *gin.Context) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "send@email.makemake.in", "Send Message")
	m.SetHeader("To", "345263950@qq.com")
	m.SetHeader("Subject", "haha666")
	m.SetBody("text/plain", "test message 123")
	g := gomail.NewDialer("smtp.yandex.com", 465, "send@email.makemake.in", "Lzf129126")

	err := g.DialAndSend(m)
	if err != nil {
		t.Fail(c, &Fail{ErrorInfo: err})
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = m.WriteTo(buf)
	if err != nil {
		t.Fail(c, &Fail{ErrorInfo: err})
		return
	}
	t.Success(c, &Success{Data: buf.String()})
}
