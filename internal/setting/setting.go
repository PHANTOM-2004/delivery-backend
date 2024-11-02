package setting

import (
	"delivery-backend/pkg/utils"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-ini/ini"
)

type Server struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	SSLKeyPath   string
	SSLCertPath  string
}

type Test struct {
	CATest            bool
	HTTPPort          int
	LocalhostKeyPath  string
	LocalhostCertPath string
}

type Log struct {
	Level string
}

type App struct {
	Salt             string
	JWTSecretKey     string
	AdminToken       string
	AdminAKAge       int
	AdminRKAge       int
	MerchantAKAge    int
	MerchantRKAge    int
	LicensePhotoSize int
	LicenseStorePath string
	LicenseAllowExts []string
	LicensePageSize  int
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

type Redis struct {
	Host      string
	Password  string
	MaxIdle   int
	MaxActive int
	Secret    string
}

var cfg *ini.File

var (
	DatabaseSetting = &Database{}
	RedisSetting    = &Redis{}
	ServerSetting   = &Server{}
	TestSetting     = &Test{}
	AppSetting      = &App{}
	LogSetting      = &Log{}
)

const FallbackPreset = "localdebug"

var Preset = map[string]string{
	"localdebug": "conf/app.ini",
	"dockertest": "conf/app_test_docker.ini",
}

func parseConfigModeSetting() {
	n := len(os.Args)
	preset := FallbackPreset
	if n <= 1 {
		// 没有给出多余参数
		log.Warnf("No preset given, fallback:[%s]", FallbackPreset)
	} else {
		// trim the leading --
		preset = strings.TrimPrefix(os.Args[1], "--")
	}

	preset = strings.ToLower(preset)
	path, ok := Preset[preset]
	if !ok {
		log.Fatalf("Unknown preset given, supported: %v", Preset)
	}

	var err error
	cfg, err = ini.Load(path)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to parse [%s]", path)
	}
}

func Setup() {
	// 首先确定模式
	parseConfigModeSetting()

	parseLogSetting()
	parseServerSetting()
	parseDatabaseSetting()
	parseRedisSetting()
	parseTestSetting()
	parseAppSetting()
}

func logCurrentConf(s any, section string) {
	kv, err := utils.StructToStr(s)
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

func parseRedisSetting() {
	err := cfg.Section("redis").StrictMapTo(RedisSetting)
	if err != nil {
		log.Fatal(err)
	}

	logCurrentConf(RedisSetting, "Redis")
}
