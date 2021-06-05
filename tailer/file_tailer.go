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
	buffer []string
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
	self.buffer = make([]string, 0)

	ticker := time.NewTicker(1 * time.Second)
	// 一定要调用Stop()，回收资源
	defer ticker.Stop()
	go func(t *time.Ticker) {
		for {
			// 每5秒中从chan t.C 中读取一次
			<-t.C
			if (time.Now().Unix() - self.LastFoundTime.Unix() > 2 && len(self.buffer) > 0) {
				log.Println("timer send")
				self.sendBuffers(notifer)
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

			self.sendBuffers(notifer)
		}
		self.LastFoundTime = time.Now()
		self.buffer = append(self.buffer, msg.Text)
	}
	self.sendBuffers(notifer)
}

func (self *FileTailer)sendBuffers(notifer notify.Notifer)  {
	if len(self.buffer) == 0 {
		return
	}
	self.buffer = append(self.buffer, "file:" + self.Path)
	notifer.Notify(self.buffer)
	self.buffer = make([]string, 0)
}

func (self *FileTailer) isNewLine(line []byte) bool {
	result := self.Regexp.FindIndex(line)
	return result != nil
}
