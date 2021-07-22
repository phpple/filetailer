package config

import (
    "filetailer/notifier"
    "filetailer/watcher"
)

type GlobalNotifierConfig struct {
    Default string                    `yaml:"default"`
    List    []*notifier.NotifierConfig `yaml:"list"`
}

type AppConfig struct {
    Watchers  []*watcher.Watcher    `yaml:"watchers"`
    Notifiers GlobalNotifierConfig `yaml:"notifiers"`
    Pid       string                  `yaml:"pid"`
}
