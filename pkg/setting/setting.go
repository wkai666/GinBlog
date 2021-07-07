package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File

	RunMode string
	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string

	HTTPPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration

	PageSize int
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Failed to parse 'conf/app.ini': %v ", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadLog()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Failed get to section 'server' :%v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(9099)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("failed to get section 'app':  %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func LoadLog() {
	sec, err := Cfg.GetSection("log")
	if err != nil {
		log.Fatalf("Failed get to section 'log' %v:", err)
	}

	LogSavePath = sec.Key("LOG_SAVE_PATH").MustString("runtime/logs/")
	LogSaveName = sec.Key("LOG_SAVE_NAME").MustString("ginBlog_")
	LogFileExt  = sec.Key("LOG_FILE_EXT").MustString("log")
	TimeFormat  = sec.Key("TIME_FORMAT").MustString("20060102")
}