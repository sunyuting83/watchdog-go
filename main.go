package main

import (
	"strings"
	"time"
	utils "watchdog/Utils"
)

func WatchDog(AppName, RunPath string) {
	task := func() {
		line := "\\"
		if !isProcessExist(AppName) {
			var Exec string
			runPath := "C:\\chiain"
			Exec = strings.Join([]string{runPath, "chiaStart.bat"}, line)
			go utils.RunCommand(Exec)
		}
	}
	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(30))
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
	p = utils.CompressStr(p)
	pList := strings.Split(p, "\r\n")
	for _, v := range pList {
		if len(v) > 0 {
			if strings.Contains(v, "=") {
				s = strings.Split(v, "=")[1]
				break
			}
		}
	}
	return
}

func isProcessExist(appName string) bool {
	command := `wmic process where name="` + appName + `" get commandline 2>nul | find "daemon" 1>nul 2>nul && echo 1 || echo 0`
	c := GetPublicWinCommandLine(command)
	if c == "1" {
		return c == "1"
	}
	return false
}

func main() {
	RunPath := "c:\\chia\\daemon\\"
	WatchDog("chia.exe", RunPath)
}
