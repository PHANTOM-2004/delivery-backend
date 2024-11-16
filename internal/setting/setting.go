package setting

import (
	"delivery-backend/pkg/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"github.com/go-ini/ini"
)

type Rabbitmq struct {
	DialURL string
}

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

// 目前支持一台服务器的情况
type Email struct {
	EmailOn        bool
	SMTPHost       string
	SMTPPort       int
	SenderEmail    string
	SenderPassword string
	CCEmail        string
	TemplatePath   string
}

type Log struct {
	Level string
}

type App struct {
	Salt                string
	JWTSecretKey        string
	AdminAliveMinute    int
	MerchantAliveMinute int
	AdminToken          string
	AdminAKAge          int
	AdminRKAge          int
	MerchantAKAge       int
	MerchantRKAge       int
	MaxImageSize        int
	LicenseStorePath    string
	LicenseAllowExts    []string
	LicensePageSize     int
	DishImageAllowExts  []string
	DishImageStorePath  string
}

func (a *App) GetLicenseStorePath(name string) string {
	path := a.LicenseStorePath + "/" + name
	return path
}

func (a *App) GenLicenseName() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Warn(err)
	}
	path := "merchant-license-" + id.String()
	return path
}

func (a *App) checkExt(allows []string, name string) (string, bool) {
	ext := filepath.Ext(name)
	for i := range allows {
		if ext == allows[i] {
			return ext, true
		}
	}
	return ext, false
}

func (a *App) CheckLicenseExt(name string) (string, bool) {
	return a.checkExt(a.LicenseAllowExts, name)
}

func (a *App) CheckDishImageExt(name string) (string, bool) {
	return a.checkExt(a.DishImageAllowExts, name)
}

func (a *App) GetDishImageStorePath(name string) string {
	path := a.DishImageStorePath + "/" + name
	return path
}

func (a *App) GenDishImageName() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Warn(err)
	}
	path := "dish-" + id.String()
	return path
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
	LogLevel     string
}

func (d *Database) GetLogLevel() logger.LogLevel {
	l := strings.ToLower(d.LogLevel)
	switch l {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
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
	RabitmqSetting  = &Rabbitmq{}
	ServerSetting   = &Server{}
	TestSetting     = &Test{}
	AppSetting      = &App{}
	LogSetting      = &Log{}
	EmailSetting    = &Email{}
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
  parseRabbitmqSetting()
	parseRedisSetting()
	parseTestSetting()
	parseAppSetting()
	parseEmailSetting()
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

func parseRabbitmqSetting() {
	err := cfg.Section("rabbitmq").StrictMapTo(RabitmqSetting)
	if err != nil {
		log.Fatal(err)
	}

	logCurrentConf(RabitmqSetting, "rabbitmq")
}

func parseEmailSetting() {
	err := cfg.Section("email").StrictMapTo(EmailSetting)
	if err != nil {
		log.Fatal(err)
	}
	logCurrentConf(EmailSetting, "Email")
	if EmailSetting.SenderEmail == "" {
		log.Fatal("邮件发送人地址不可为空")
	}

	_, err = os.Open(EmailSetting.TemplatePath)
	if err != nil {
		log.Fatal("parse email setting", err)
	}
}
