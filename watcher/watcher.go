package watcher

import (
    "errors"
    "filetailer/notifier"
    "filetailer/rule"
    "filetailer/tailer"
    "log"
)

type Watcher struct {
    Paths    []string    `yaml:"paths"`
    Pattern  string      `yaml:"pattern"`
    Rules    []rule.Rule `yaml:"rules"`
    Notifier string      `yaml:"notifier"` // 对应于NotiferConfig的key字段，没有指定则取GlobalNotifierConfig的Default字段
    notifier notifier.Notifier
    inited   bool
}

func (self *Watcher) Init() {
    for _, rule := range self.Rules {
        rule.Init()
    }
    self.inited = true
}

// 选择通知器
func (self *Watcher) ChooseNotifier(notifiers map[string]notifier.Notifier, defNotifier string) error {
    notifierName := self.Notifier
    if notifierName == "" {
        notifierName = defNotifier
    }
    self.Notifier = notifierName

    if notifier, ok := notifiers[notifierName]; ok {
        self.notifier = notifier
        return nil
    }
    return errors.New("notifier not found")
}

// 开始运行
func (self *Watcher) Run() error {
    if !self.inited {
        self.Init()
    }
    if self.notifier == nil {
        return errors.New("notifier not inited")
    }
    for _, path := range self.Paths {
        fileHandler := tailer.NewFileTailer(path, self.Pattern)
        log.Printf("file:%s, notifier:%s\n", path, self.Notifier)
        go fileHandler.Handle(self.Rules, self.notifier)
    }
    return nil
}
