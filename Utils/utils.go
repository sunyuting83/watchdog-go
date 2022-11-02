package utils

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v2"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

type Config struct {
	Watch    []Task   `yaml:"Watch"`
	Kill     []string `yaml:"Kill"`
	ScanTime int      `yaml:"ScanTime"`
}

type Task struct {
	TaskName string `yaml:"TaskName"`
	RunPath  string `yaml:"RunPath"`
}

// RunCommand run command
func RunCommand(app string) (k string, err error) {
	App := strings.Join([]string{"name=", `"`, app, `"`}, "")
	cmd := exec.Command("wmic", "process", "where", App)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//fmt.Println(cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		return "", err
	}

	if len(bytesErr) != 0 {
		return "", errors.New("0")

	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return ConvertByte2String(bytes, "GB18030"), nil
}

// CompressStr
func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// RunUpdate command
func RunCMD(cmdExec string) {
	cmd := exec.Command(cmdExec)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

// CheckConfig check config
func CheckConfig(OS, CurrentPath string) (conf *Config, err error) {
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	return confYaml, nil
}

func KillTask(task string) {
	cmd := exec.Command("taskkill.exe", "/F", "/IM", task)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}

func HomeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
