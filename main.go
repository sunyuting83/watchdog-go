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
	// fmt.Println("aaa")
	task := func() {
		line := "\\"
		var run bool = false
		if len(conf.Watch) > 0 {
			for _, Task := range conf.Watch {
				if !isProcessExist(Task.TaskName) {
					if len(conf.Kill) > 0 {
						for _, Task := range conf.Kill {
							go utils.KillTask(Task)
							time.Sleep(200 * time.Millisecond)
						}
					}
					time.Sleep(time.Duration(3) * time.Second)
					run = true
					break
				}
			}
		}
		if run {
			if len(conf.Watch) > 0 {
				for i, Task := range conf.Watch {
					run = true
					var Exec string = Task.TaskName
					if len(Task.RunPath) > 0 {
						Exec = strings.Join([]string{Task.RunPath, Task.TaskName}, line)
					}
					utils.RunCMD(Exec)
					time.Sleep(time.Duration(1) * time.Second)
					if i == len(conf.Watch) - 1 {
						run = false
					}
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

func GetPublicWinCommandLine(cmd string) (s string, err error) {
	p, err := utils.RunCommand(cmd)
	s = utils.CompressStr(p)
	return
}

func isProcessExist(appName string) bool {
	// command := `wmic process where name="` + appName + `"`
	App := strings.Join([]string{"name=", `"`, appName, `"`}, "")
	cmd := strings.Join([]string{" process where", App, "get ProcessId"}, " ")
	// fmt.Println(cmd)
	c, err := GetPublicWinCommandLine(cmd)
	if err != nil {
		// fmt.Println("a" + err.Error())
		return true
	}
	return strings.Contains(c, "ProcessId")
}

func checkFirst(appName string) bool {
	App := strings.Join([]string{"name=", `"`, appName, `"`}, "")
	cmd := strings.Join([]string{" process where", App, "get ProcessId"}, " ")
	// fmt.Println(cmd)
	c, err := utils.RunCommand(cmd)
	if err != nil {
		// fmt.Println("a" + err.Error())
		return false
	}
	// fmt.Println(c)
	line := "\r\n"
	if !strings.Contains(c, line) {
		line = "\n"
	}
	var PID []string
	if strings.Contains(c, "ProcessId") {
		pidList := strings.Split(c, line)
		for _, item := range pidList {
			CompressStr := utils.CompressStr(item)
			if len(CompressStr) > 0 {
				if !strings.Contains(item, "ProcessId") {
					PID = append(PID, CompressStr)
				}
			}
		}
	}
	// fmt.Println(PID)
	// fmt.Println(len(PID))
	// fmt.Println(len(PID) <= 1)
	return len(PID) <= 1
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
	check := checkFirst("watchdog.exe")
	// WatchDog(yaml)
	if check {
		WatchDog(yaml)
	} else {
		os.Exit(0)
	}
}
