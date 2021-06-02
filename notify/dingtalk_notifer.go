package notify

import (
	"github.com/blinkbean/dingtalk"
	"log"
	"strings"
)

type DingtalkNotier struct{
	Client *dingtalk.DingTalk
}
func NewDingtalkNotier(token string) DingtalkNotier {
	log.Println("dingtalk created:", token)
	client := dingtalk.InitDingTalk([]string{token}, ".")
	return DingtalkNotier{
		client,
	}
}

func (self *DingtalkNotier)Notify(msgs []string)  {
	msg := strings.Join(msgs, "\n")
	log.Println(msg)

	log.Println("dingtalk.send")
	err := self.Client.SendTextMessage("[ERROR FOUND]\n" + msg)
	if err != nil {
		log.Fatalln("error found:", err)
	}
}
