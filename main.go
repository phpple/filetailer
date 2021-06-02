package main

import (
	"awesomeProject/notify"
	"awesomeProject/tailer"
	"flag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configFile := *flag.String("config", "", "config file")
	if configFile == "" {
		pwd, _ := os.Getwd()
		configFile = pwd + "/config.yml"
	}
	appConfig := getConfig(configFile)

	fileHandler := tailer.NewFileTailer(appConfig.File.Path, appConfig.File.Pattern)
	notifer := notify.BuildNotifer(appConfig.Notifer.Name, appConfig.Notifer.Option)
	fileHandler.Handle(notifer)
}

func getConfig(configFile string) AppConfig {
	var appConfig AppConfig
	config, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatalln("config file read error:", err)
	}
	err = yaml.Unmarshal(config, &appConfig)
	if err != nil {
		log.Fatalln("yaml parse failed:", err)
	}
	appConfig.File.Path, err = filepath.Abs(appConfig.File.Path)
	if err != nil {
		log.Fatalln("file path error:", err)
	}
	return appConfig
}