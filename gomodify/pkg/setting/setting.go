package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type AppSetting struct {
	ENV         string
	CORS        string
	OS          string
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

type DatabaseSetting struct {
	LocalURL    string
	DockerURL   string
	Retries     int
	MaxActive   int
	MaxIdle     int
	MaxLifetime time.Duration
}

const (
	ENV_DEVELOPMENT = "development"
	ENV_STAGING     = "staging"
	ENV_PRODUCTION  = "production"

	OS_MAC   = "mac"
	OS_LINUX = "linux"
)

var (
	App      = &AppSetting{}
	Database = &DatabaseSetting{}
)

var (
	filepath = "conf/app.ini"
)

// Setup initialize the configuration instance
func Setup() {
	cfg, err := ini.Load(filepath)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo(cfg, "app", App)
	mapTo(cfg, "database", Database)

	App.ENV = getEnv(App.ENV)
	App.OS = getOS(App.OS)

	Database.MaxLifetime = Database.MaxLifetime * time.Second
}

// mapTo map section
func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func getEnv(env string) string {
	switch env {
	case ENV_PRODUCTION:
		return ENV_PRODUCTION
	case ENV_STAGING:
		return ENV_STAGING
	default:
		return ENV_DEVELOPMENT
	}
}

func getOS(os string) string {
	switch os {
	case OS_MAC:
		return OS_MAC
	default:
		return OS_LINUX
	}
}
