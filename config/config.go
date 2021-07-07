package config

type FileConfig struct {
	Paths   []string `yaml:"paths"`
	Pattern string   `yaml:"pattern"`
}

type NotiferConfig struct {
	Name   string                 `yaml:"name"`
	Option map[string]interface{} `yaml:"option"`
}

type AppConfig struct {
	File    FileConfig    `yaml:"file"`
	Notifer NotiferConfig `yaml:"notifer"`
	Pid     string        `yaml:"pid"`
}
