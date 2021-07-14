package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type Server struct {
	HttpPort int
	ReadTimeOut time.Duration
	WriteTimeOut time.Duration
}

type App struct {
	RunMode string

	PrefixUrl string
	JwtSecret string
	PageSize int
	RuntimeRootPath string
	ExportSavePath string

	QrCodeSavePath string
	FontSavePath string

	ImageSavePath string
	ImageMaxSize int
	ImageAllowExts []string
}

type Log struct {
	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string
}

type Database struct {
	Type string
	Host string
	User string
	Password string
	DBName string
	TablePrefix string
}

type Redis struct {
	Host string
	Password string
	MaxIdle int
	MaxActive int
	IdleTimeout time.Duration
}

var (
	ServerSetting   = &Server{}
	AppSetting      = &App{}
	LogSetting      = &Log{}
	DatabaseSetting = &Database{}
	RedisSetting    = &Redis{}
)

func Setup() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Failed to parse 'app.ini': %v", err)
	}

	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeOut = ServerSetting.ReadTimeOut * time.Second
	ServerSetting.WriteTimeOut = ServerSetting.WriteTimeOut * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo Database err: %v", err)
	}

	err = Cfg.Section("log").MapTo(LogSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo Log err: %v", err)
	}

	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo Redis err: %v", err)
	}

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}
