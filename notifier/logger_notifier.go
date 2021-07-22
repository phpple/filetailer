package notifier

import (
	"log"
	"strings"
)

type LoggerNotifier struct {
}

func (self *LoggerNotifier) Notify(msgs []string) {
	log.Println(strings.Join(msgs, "\n"))
}
