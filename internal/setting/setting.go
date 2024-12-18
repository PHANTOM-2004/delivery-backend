package setting

import (
	"delivery-backend/pkg/utils"
	"fmt"
	"io"
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
	Level    string
	SavePath string
}

type Wechat struct {
	AppID                string
	AppSecret            string
	TokenRefreshInterval int
	SessionAge           int
	ImageExt             []string
	ImageStorePath       string
	code2SessionURL      string
	accesstokenURL       string
}

func (w *Wechat) GetCode2SessionURL(js_code string) string {
	res := fmt.Sprintf(w.code2SessionURL,
		w.AppID, w.AppSecret, js_code)
	return res
}

func (w *Wechat) CheckImageExt(name string) (string, bool) {
	return checkExt(w.ImageExt, name)
}

func (w *Wechat) GetImageName(prefix string) string {
	res := prefix + "-" + uuid.NewString()
	return res
}

func (w *Wechat) GetImagePath(name string) string {
	path := w.ImageStorePath + "/" + name
	return path
}

func (w *Wechat) GetAccessTokenURL() string {
	res := fmt.Sprintf(w.accesstokenURL, w.AppID, w.AppSecret)
	return res
}

type App struct {
	Salt                 string
	JWTSecretKey         string
	AdminAliveMinute     int
	MerchantAliveMinute  int
	AdminToken           string
	AdminAKAge           int
	AdminRKAge           int
	MerchantAKAge        int
	MerchantRKAge        int
	MaxImageSize         int
	ApplicationStorePath string
	ApplicationAllowExts []string
	ApplicationPageSize  int
	DishImageAllowExts   []string
	DishImageStorePath   string
}

func (a *App) GetLicenseStorePath(name string) string {
	path := a.ApplicationStorePath + "/" + name
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

func checkExt(allows []string, name string) (string, bool) {
	ext := filepath.Ext(name)
	for i := range allows {
		if ext == allows[i] {
			return ext, true
		}
	}
	return ext, false
}

func (a *App) CheckLicenseExt(name string) (string, bool) {
	return checkExt(a.ApplicationAllowExts, name)
}

func (a *App) CheckDishImageExt(name string) (string, bool) {
	return checkExt(a.DishImageAllowExts, name)
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
	WechatSetting   = &Wechat{}
)

const FallbackPreset = "localdebug"

var Preset = map[string]string{
	"localdebug":   "conf/app.ini",
	"dockertest":   "conf/app_test_docker.ini",
	"dockerdeploy": "conf/app_deploy_docker.ini",
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
	log.Info("using preset: ", path)
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
	parseWechatSetting()
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

	// NOTE:对于log设置，我们要求写入屏幕以及文件
	//
	var logfile_path string
	logfilename := time.Now().Format("log-2006-01-02-15-04-05") + ".log"
	if LogSetting.SavePath == "" {
		logfile_path = logfilename
	} else {
		// 创建路径
		if _, err := os.Stat(LogSetting.SavePath); os.IsNotExist(err) {
			log.Fatal("Log path不存在", LogSetting.SavePath)
		}
		logfile_path = LogSetting.SavePath + "/" + logfilename
		log.Info("LOGGINT TO:", logfile_path)
	}

	logFile, err := os.OpenFile(logfile_path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error("fail to creating log file")
		log.Fatal(err)
	}

	// 设置控制台日志
	writer := io.MultiWriter(os.Stdout, logFile)
	// 首先设置timestamp格式
	log.SetOutput(writer)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      false,
		DisableColors:    true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  "2006-01-02 15:04:05",
	})
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

func parseWechatSetting() {
	err := cfg.Section("wechat").StrictMapTo(WechatSetting)
	if err != nil {
		log.Fatal(err)
	}
	if WechatSetting.AppID == "" {
		log.Fatal("wechat appid must be filled")
	}
	if WechatSetting.AppSecret == "" {
		log.Fatal("wechat appsecret must be filled")
	}

	url_fmt := "https://api.weixin.qq.com/sns/jscode2session?" +
		"appid=%s" +
		"&secret=%s" +
		"&js_code=%s&grant_type=authorization_code"
	WechatSetting.code2SessionURL = url_fmt

	WechatSetting.accesstokenURL = "https://api.weixin.qq.com/cgi-bin/token?" +
		"appid=%s" +
		"&secret=%s" +
		"&grant_type=client_credential"

	log.Info(WechatSetting.code2SessionURL)
	logCurrentConf(WechatSetting, "Wechat")
}
