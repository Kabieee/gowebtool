package task

import (
	"fmt"

	"github.com/fatih/color"

	"gopkg.in/gomail.v2"
)

type EmailTask struct {
	//EmailChan chan *Email
}

var EmailChan chan *Email

type Email struct {
	Message *gomail.Message
	To      string
}

func (e *EmailTask) Run() {
	//defer sy.Done()
	EmailChan = make(chan *Email, 4)
	fmt.Println(color.GreenString("%s", "Start Email"))
	e.send()
}

func (e *EmailTask) send() {
	g := gomail.NewDialer("smtp.yandex.com", 465, "send@email.makemake.in", "Lzf129126")
	for email := range EmailChan {
		fmt.Println(color.BlueString("start send to %s", email.To))
		err := g.DialAndSend(email.Message)
		if err != nil {
			fmt.Printf("send email failed: %s", err)
			continue
		}
		fmt.Println(color.BlueString("success send to %s", email.To))
		//fmt.Printf("\ntask id: %d \n", i)
		//time.Sleep(time.Second * 3)
		//buf := bytes.NewBuffer(nil)
		//_, err := email.Message.WriteTo(buf)
		//fmt.Println(buf.String())
	}
}
