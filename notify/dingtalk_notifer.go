package notify

import (
	"github.com/blinkbean/dingtalk"
	"log"
	"strings"
	"time"
)

type DingtalkOption struct {
	Tokens  []string `yaml:"tokens"`
	Keyword string   `yaml:"keyword"`
}

var DefaultSendFrequece time.Duration = 2 * time.Second

type DingtalkNotier struct {
	Client *dingtalk.DingTalk
}

type DingtalkMessage struct {
	Title string
	Lines []string
}

var msgQueue chan DingtalkMessage

func init() {
	msgQueue = make(chan DingtalkMessage, 20)
}

func writeMessage(msg DingtalkMessage) {
	msgQueue <- msg
}

func (self DingtalkNotier) fetchQueue() {
	for {
		// 每5秒中从chan t.C 中读取一次
		if msg, ok := <-msgQueue; ok {
			log.Println(msg.Lines)
			msg := strings.Join(msg.Lines, "\n> ###### ")
			log.Println(msg)

			log.Println("dingtalk.send")
			err := self.Client.SendMarkDownMessage("error found", msg)
			if err != nil {
				log.Fatalln("error found:", err)
			}
		}
	}
}

func NewDingtalkNotier(option DingtalkOption) DingtalkNotier {
	log.Println("dingtalk created:", option.Tokens)
	client := dingtalk.InitDingTalk(option.Tokens, option.Keyword)

	notifer := DingtalkNotier{
		client,
	}
	go notifer.fetchQueue()
	return notifer
}

func (self *DingtalkNotier) Notify(msgs []string) {
	writeMessage(DingtalkMessage{
		Title: "error found",
		Lines: msgs,
	})
}
