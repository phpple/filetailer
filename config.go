package main

type FileConfig struct {
	Path    string `yaml:"path"`
	Pattern string `yaml:"pattern"`
}

type NotiferConfig struct {
	Name string `yaml:"name"`
	Option map[string]string `yaml:"option"`
}

type AppConfig struct {
	File FileConfig `yaml:"file"`
	Notifer NotiferConfig `yaml:"notifer"`
}