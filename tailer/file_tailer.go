package tailer

import (
	"filetailer/notify"
	"github.com/hpcloud/tail"
	"log"
	"regexp"
	"runtime"
	"time"
)

type FileTailer struct {
	Path   string
	Regexp *regexp.Regexp
	LastFoundTime time.Time
	Poll bool
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
		Poll: runtime.GOOS == "windows",
	}
}

func (self *FileTailer) Handle(notifer notify.Notifer) {
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
	var buffer = make([]string, 0)

	ticker := time.NewTicker(1 * time.Second)
	// 一定要调用Stop()，回收资源
	defer ticker.Stop()
	go func(t *time.Ticker) {
		for {
			// 每5秒中从chan t.C 中读取一次
			<-t.C
			if (time.Now().Unix() - self.LastFoundTime.Unix() > 2 && len(buffer) > 0) {
				log.Println("timer send")
				notifer.Notify(buffer)
				buffer = make([]string, 0)
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

			if len(buffer) > 0 {
				go notifer.Notify(buffer)
				buffer = make([]string, 0)
			}
		}
		self.LastFoundTime = time.Now()
		buffer = append(buffer, msg.Text)
	}
	if len(buffer) > 0 {
		go notifer.Notify(buffer)
	}
}

func (self *FileTailer) isNewLine(line []byte) bool {
	result := self.Regexp.FindIndex(line)
	return result != nil
}
