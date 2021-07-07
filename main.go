package main

import (
	"filetailer/config"
	"filetailer/notify"
	"filetailer/tailer"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
)

func main() {
	//合建chan
	c := make(chan os.Signal)
	//监听所有信号
	signal.Notify(c)

	configFile := *flag.String("config", "", "config file")
	if configFile == "" {
		pwd, _ := os.Getwd()
		configFile = pwd + "/config.yml"
	}
	appConfig := getConfig(configFile)

	defer func() {
		delPid(appConfig.Pid)
	}()
	writePid(appConfig.Pid)

	log.Println("config:", appConfig)

	notifer := notify.BuildNotifer(appConfig.Notifer.Name, appConfig.Notifer.Option)

	for _, path := range appConfig.File.Paths {
		fileHandler := tailer.NewFileTailer(path, appConfig.File.Pattern)
		go fileHandler.Handle(notifer)
	}
	s := <-c
	fmt.Println("filetailer exit", s)
	delPid(appConfig.Pid)
}

// 删除pid文件
func delPid(pidFile string) {
	os.Remove(pidFile)
}

// 写入pid文件
func writePid(pidFile string) {
	pid := os.Getpid()
	var perm fs.FileMode = 0
	os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), perm)
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

	if appConfig.Pid == "" {
		appConfig.Pid = "app.pid"
	}
	return appConfig
}
