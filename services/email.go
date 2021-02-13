package services

import (
	"bytes"
	"fmt"
	"gowebtool/task"
	"regexp"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailServices struct {
	BaseServices
}

func (e *EmailServices) Send(data EmailData) (err error, result interface{}) {
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
		return err, nil
	}
	reg := regexp.MustCompile("\r\n")
	split := reg.Split(buf.String(), -1)
	return nil, split
}
