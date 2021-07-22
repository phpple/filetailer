package notifier

import (
    "github.com/blinkbean/dingtalk"
    "log"
    "strings"
    "time"
)

type DingtalkOption struct {
    Tokens  []string `yaml:"tokens"`
    Keyword string   `yaml:"keyword"`
    MaxLine int      `yaml:"maxline"`
    MaxChar int      `yaml:"maxchar"`
}

var DefaultSendFrequece time.Duration = 2 * time.Second

type DingtalkNotifier struct {
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

func (self DingtalkNotifier) fetchQueue() {
    for {
        // 每5秒中从chan t.C 中读取一次
        if msg, ok := <-msgQueue; ok {
            var message string
            lineLen := len(msg.Lines)
            if self.Option.MaxLine > 0 && lineLen > self.Option.MaxLine {
                // 取前面几条，以及最后一条
                lines := make([]string, self.Option.MaxLine+1)
                lines[0] = msg.Lines[0]

                for i := 1; i < self.Option.MaxLine; i++ {
                    if self.Option.MaxChar > 0 && len(msg.Lines[i]) > self.Option.MaxChar {
                        lines[i] = msg.Lines[i][0:self.Option.MaxChar] + "..."
                    } else {
                        lines[i] = msg.Lines[i]
                    }
                }
                lines = append(lines, "...")
                lines = append(lines, msg.Lines[len(msg.Lines)-1])

                message = strings.Join(lines, "\n> ###### ")
            } else {
                lines := make([]string, lineLen)
                lines[0] = msg.Lines[0]
                for i := 1; i < lineLen-1; i++ {
                    if self.Option.MaxChar > 0 && len(msg.Lines[i]) > self.Option.MaxChar {
                        lines[i] = msg.Lines[i][0:self.Option.MaxChar] + "..."
                    } else {
                        lines[i] = msg.Lines[i]
                    }
                }
                lines[lineLen-1] = msg.Lines[lineLen-1]
                message = strings.Join(lines, "\n> ###### ")
            }

            log.Println(message)

            err := self.Client.SendMarkDownMessage("error found", message)
            if err != nil {
                log.Fatalln("error found:", err)
            }
        }
    }
}

func NewDingtalkNotier(option DingtalkOption) DingtalkNotifier {
    log.Println("dingtalk created:", option)
    client := dingtalk.InitDingTalk(option.Tokens, option.Keyword)

    notifer := DingtalkNotifier{
        client,
        option,
    }
    go notifer.fetchQueue()
    return notifer
}

func (self *DingtalkNotifier) Notify(msgs []string) {
    writeMessage(DingtalkMessage{
        Title: "error found",
        Lines: msgs,
    })
}
