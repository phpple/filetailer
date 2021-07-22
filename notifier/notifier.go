package notifier

type Notifier interface {
	// 通知
	Notify(msgs []string)
}

type NotifierConfig struct {
	Key    string                 `yaml:"key"` // key
	Type   string                 `yaml:"type"`
	Option map[string]interface{} `yaml:"option"`
}