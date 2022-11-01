package main

import (
	"strings"
	"time"
	utils "watchdog/Utils"
)

func WatchDog(AppName, RunPath string) {
	task := func() {
		line := "\\"
		// fmt.Println("start")
		if !isProcessExist(AppName) {
			// fmt.Println("no")
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
	s = utils.CompressStr(p)
	return
}

func isProcessExist(appName string) bool {
	command := `wmic process where name="` + appName + `"`
	c := GetPublicWinCommandLine(command)
	// fmt.Println(c)
	// fmt.Println("a")
	// fmt.Println(appName)
	// fmt.Println(strings.Contains(c, appName))
	return strings.Contains(c, appName)
}

func main() {
	RunPath := "c:\\chia\\daemon\\"
	WatchDog("start_harvester.exe", RunPath)
}
