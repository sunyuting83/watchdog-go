package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	utils "watchdog/Utils"
)

func WatchDog(conf *utils.Config) {
	task := func() {
		line := "\\"
		if len(conf.Watch) > 0 {
			for i, Task := range conf.Watch {
				if !isProcessExist(Task.TaskName) {
					if len(conf.Kill) > 0 && i == 0 {
						for _, Task := range conf.Kill {
							utils.KillTask(Task)
							time.Sleep(time.Duration(3) * time.Second)
						}
					}
					time.Sleep(time.Duration(3) * time.Second)
					var Exec string = Task.TaskName
					if len(Task.RunPath) > 0 {
						Exec = strings.Join([]string{Task.RunPath, Task.TaskName}, line)
					}
					utils.RunCMD(Exec)
				}
			}
		}
	}
	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(conf.ScanTime))
	go func() {
		for range ticker.C {
			task()
		}
		ch <- 1
	}()
	<-ch
}

func GetPublicWinCommandLine(command string) (s string) {
	p, _ := utils.RunCommand(command)
	s = utils.CompressStr(p)
	return
}

func isProcessExist(appName string) bool {
	// command := `wmic process where name="` + appName + `"`
	c := GetPublicWinCommandLine(appName)
	return strings.Contains(c, appName)
}

func main() {
	CurrentPath, _ := utils.GetCurrentPath()
	OS := runtime.GOOS

	lnk := strings.Join([]string{CurrentPath, "watchdog.lnk"}, "\\")
	if !utils.Exists(lnk) {
		fmt.Println("快捷方式不存在")
		os.Exit(0)
	}
	home, _ := utils.HomeWindows()
	p := "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"
	startPath := strings.Join([]string{home, p, "watchdog.lnk"}, "\\")
	if !utils.Exists(startPath) {
		s, _ := os.ReadFile(lnk)
		os.WriteFile(startPath, s, 0644)
	}
	yaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	WatchDog(yaml)
}
