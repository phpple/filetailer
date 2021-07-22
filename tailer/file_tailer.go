package tailer

import (
    "filetailer/notifier"
    "filetailer/rule"
    "fmt"
    "github.com/hpcloud/tail"
    "log"
    "regexp"
    "runtime"
    "strings"
    "time"
)

type FileTailer struct {
    Path          string
    Regexp        *regexp.Regexp
    LastFoundTime time.Time
    Poll          bool
    buffer        []string
}

func NewFileTailer(path string, pattern string) *FileTailer {
    if pattern == "" {
        log.Fatalln("pattern is empty")
        return nil
    }
    timestampRegexp := regexp.MustCompile(pattern)

    log.Println("os:", runtime.GOOS)
    return &FileTailer{
        Path:   path,
        Regexp: timestampRegexp,
        Poll:   runtime.GOOS == "windows",
    }
}

// 处理信息
func (self *FileTailer) Handle(rules []rule.Rule, notifer notifier.Notifier) {
    tails, err := tail.TailFile(self.Path, tail.Config{
        ReOpen:    false,
        Follow:    true,
        MustExist: true,
        Poll:      self.Poll,
        Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
    })
    if err != nil {
        log.Fatalln("tail file err:", err)
        return
    }

    var msg *tail.Line
    var ok bool
    self.buffer = make([]string, 0)

    ticker := time.NewTicker(1 * time.Second)
    // 一定要调用Stop()，回收资源
    defer ticker.Stop()
    go func(t *time.Ticker) {
        for {
            // 每5秒中从chan t.C 中读取一次
            <-t.C
            if time.Now().Unix()-self.LastFoundTime.Unix() > 2 && len(self.buffer) > 0 {
                log.Println("timer send")
                self.sendBuffers(notifer, rules)
            }
        }
    }(ticker)

    for true {
        msg, ok = <-tails.Lines
        if !ok {
            log.Printf("tail file close reopen, filename:%s\n", tails.Filename)
            time.Sleep(100 * time.Millisecond)
            continue
        }

        if self.isNewLine([]byte(msg.Text)) {
            log.Print("found new line")
            self.sendBuffers(notifer, rules)
        }
        self.LastFoundTime = time.Now()
        line := strings.TrimSpace(msg.Text)
        if line != "" {
            self.buffer = append(self.buffer, line)
        }
    }
    self.sendBuffers(notifer, rules)
}

func (self *FileTailer) sendBuffers(notifer notifier.Notifier, rules []rule.Rule) {
    if len(self.buffer) == 0 {
        return
    }

    pass := false
    for i, line := range self.buffer {
        // 通过规则校验，决定是否将其放到buffer
        for _, rule := range rules {
            _, matched := rule.Match(line)
            // 第一行做代码格式化
            if matched && i == 0 {
                if line != "" && rule.Msg != "" {
                    line = fmt.Sprintf(rule.Msg, line)
                }
                self.buffer[i] = line
            }
            if matched {
                pass = true
                break
            }
        }
    }
    if (pass) {
        self.buffer = append(self.buffer, "file:"+self.Path)
        notifer.Notify(self.buffer)
    }
    self.buffer = make([]string, 0)
}

func (self *FileTailer) isNewLine(line []byte) bool {
    result := self.Regexp.FindIndex(line)
    return result != nil
}
