package main

import (
    "filetailer/config"
    "filetailer/notifier"
    "flag"
    "gopkg.in/yaml.v3"
    "io/fs"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

func main() {
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

    log.Printf("config:%q", appConfig)

    // 初始化通知器
    var notifiers map[string]notifier.Notifier = make(map[string]notifier.Notifier)
    for _, notifierConf := range appConfig.Notifiers.List {
        notifiers[notifierConf.Key] = notifier.BuildNotifier(notifierConf.Type, notifierConf.Option)
    }

    // 绑定通知器
    for _, watcher := range appConfig.Watchers {
        err := watcher.ChooseNotifier(notifiers, appConfig.Notifiers.Default)
        if err != nil {
            panic(err)
        }
        watcher.Init()
    }

    // 运行观察器
    for _, watcher := range appConfig.Watchers {
        err := watcher.Run()
        if err != nil {
            panic(err)
        }
    }

    for {
        time.Sleep(10 * time.Second)
    }
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

    for _, watcher := range appConfig.Watchers {
        for k, path := range watcher.Paths {
            watcher.Paths[k], err = filepath.Abs(path)
            if err != nil {
                log.Fatalln("file path error:", err)
            }
        }
        for k, path := range watcher.Paths {
            watcher.Paths[k], err = filepath.Abs(path)
            if err != nil {
                log.Fatalln("file path error:", err)
            }
        }
    }

    if appConfig.Pid == "" {
        appConfig.Pid = "app.pid"
    }
    return appConfig
}
