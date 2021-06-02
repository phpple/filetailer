package notify

import (
	"log"
	"strings"
)

type LoggerNotifer struct {
}

func (self *LoggerNotifer)Notify(msgs []string)  {
	log.Println(strings.Join(msgs, "\n"))
}
