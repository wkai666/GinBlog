package logging

import (
	"fmt"
	"ginApp/pkg/setting"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = setting.LogSavePath
	LogSaveName = setting.LogSaveName
	LogFileExt = setting.LogFileExt
	TimeFormat = setting.TimeFormat
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission: %v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile:  %v", err)
	}
	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	logDir := dir + "/" + getLogFilePath()
	_, err := os.Stat(logDir)

	if os.IsNotExist(err) {
		// 目录不存在
		err = os.Mkdir(dir + "/" + getLogFilePath(), os.ModePerm)
		if err != nil {
			// 创建目录出错
			panic(err)
		}
	}

}