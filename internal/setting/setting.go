package setting

import (
	"delivery-backend/pkg/common"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-ini/ini"
)

type Server struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Test struct {
	CATest            bool
	HTTPSPort         int
	LocalhostKeyPath  string
	LocalhostCertPath string
}

type Log struct {
	Level string
}

type App struct {
	Salt string
}

type Database struct {
	Type         string
	User         string
	Password     string
	Host         string
	Name         string
	TablePrefix  string
	MaxIdleConns int
	MaxOpenConns int
}

var cfg *ini.File

const conf_path = "conf/app.ini"

var (
	DatabaseSetting = &Database{}
	ServerSetting   = &Server{}
	TestSetting     = &Test{}
	AppSetting      = &App{}
	LogSetting      = &Log{}
)

func Setup() {
	var err error
	cfg, err = ini.Load(conf_path)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to parse [%s]", conf_path)
	}

	parseLogSetting()
	parseServerSetting()
	parseDatabaseSetting()
	parseTestSetting()
	parseAppSetting()
}

func logCurrentConf(s any, section string) {
	kv, err := common.StructToStr(s)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(section + "Setting:\n" + kv)
}

func parseAppSetting() {
	err := cfg.Section("app").StrictMapTo(AppSetting)
	if err != nil {
		log.Fatal(err)
	}
	logCurrentConf(AppSetting, "App")
}

// NOTE: loglevel会在setup时设置
func parseLogSetting() {
	err := cfg.Section("log").StrictMapTo(LogSetting)
	if err != nil {
		log.Fatal(err)
	}

	level, err := log.ParseLevel(LogSetting.Level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	logCurrentConf(LogSetting, "Log")
}

func parseTestSetting() {
	err := cfg.Section("test").StrictMapTo(TestSetting)
	if err != nil {
		log.Fatal(err)
	}

	logCurrentConf(TestSetting, "Test")
}

func parseDatabaseSetting() {
	err := cfg.Section("database").StrictMapTo(DatabaseSetting)
	if err != nil {
		log.Fatal(err)
	}

	logCurrentConf(DatabaseSetting, "Database")
}

func parseServerSetting() {
	err := cfg.Section("server").StrictMapTo(ServerSetting)
	if err != nil {
		log.Fatal(err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	logCurrentConf(ServerSetting, "Server")
}
