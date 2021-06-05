package main

import (
	"filetailer/config"
	"filetailer/notify"
	"filetailer/tailer"
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

	log.Println("config:", appConfig)

	notifer := notify.BuildNotifer(appConfig.Notifer.Name, appConfig.Notifer.Option)

	for _, path := range appConfig.File.Paths {
		fileHandler := tailer.NewFileTailer(path, appConfig.File.Pattern)
		go fileHandler.Handle(notifer)
	}
	// 不退出主进程
	for {
	}
}

func getConfig(configFile string) config.AppConfig {
	var appConfig config.AppConfig
	config, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatalln("config file read error:", err)
	}
	err = yaml.Unmarshal(config, &appConfig)
	if err != nil {
		log.Fatalln("yaml parse failed:", err)
	}
	for k, path := range appConfig.File.Paths {
		appConfig.File.Paths[k], err = filepath.Abs(path)
		if err != nil {
			log.Fatalln("file path error:", err)
		}
	}
	return appConfig
}
