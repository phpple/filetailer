package notify

import (
	"github.com/blinkbean/dingtalk"
	"log"
	"strings"
)

type DingtalkOption struct {
	Tokens  []string `yaml:"tokens"`
	Keyword string   `yaml:"keyword"`
}

type DingtalkNotier struct {
	Client *dingtalk.DingTalk
}

func NewDingtalkNotier(option DingtalkOption) DingtalkNotier {
	log.Println("dingtalk created:", option.Tokens)
	client := dingtalk.InitDingTalk(option.Tokens, option.Keyword)
	return DingtalkNotier{
		client,
	}
}

func (self *DingtalkNotier) Notify(msgs []string) {
	msg := strings.Join(msgs, "\n")
	log.Println(msg)

	log.Println("dingtalk.send")
	err := self.Client.SendTextMessage(msg)
	if err != nil {
		log.Fatalln("error found:", err)
	}
}
