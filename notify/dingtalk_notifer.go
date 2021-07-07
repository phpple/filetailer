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
	MaxLine int      `yaml:"max_line"`
}

var DefaultSendFrequece time.Duration = 2 * time.Second

type DingtalkNotier struct {
	Client *dingtalk.DingTalk
	Option DingtalkOption
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

			var message string
			if self.Option.MaxLine > 0 && len(msg.Lines) > self.Option.MaxLine {
				// 取前面几条，以及最后一条
				var lines = msg.Lines[0 : self.Option.MaxLine-1]
				lines = append(lines, "...")
				lines = append(lines, msg.Lines[len(msg.Lines)-1])

				message = strings.Join(lines, "\n> ###### ")
			} else {
				message = strings.Join(msg.Lines, "\n> ###### ")
			}

			log.Println(message)

			err := self.Client.SendMarkDownMessage("error found", message)
			if err != nil {
				log.Fatalln("error found:", err)
			}
		}
	}
}

func NewDingtalkNotier(option DingtalkOption) DingtalkNotier {
	log.Println("dingtalk created:", option)
	client := dingtalk.InitDingTalk(option.Tokens, option.Keyword)

	notifer := DingtalkNotier{
		client,
		option,
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
