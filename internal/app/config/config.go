package config

import (
	"context"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Config Структура конфигурации;
// Содержит все конфигурационные данные о сервисе;
// автоподгружается при изменении исходного файла
type Config struct {
	VKToken  string
	AdminID  int
	LogLevel string
	Bot      BotConfig
	Postgres PostgresConfig
}

type BotConfig struct {
	GroupID string
	ChatID  int
}

type PostgresConfig struct {
	// from config file
	DialTimeout int
	ReadTimeout int
	// from env
	Host     string
	Port     int
	User     string
	Password string
}

// NewConfig Создаёт новый объект конфигурации, загружая данные из файла конфигурации
func NewConfig(ctx context.Context) (*Config, error) {
	var err error

	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("не вышло здесь")
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	cfg.setLogLevel(cfg.LogLevel)

	cfg.VKToken = os.Getenv("VK_TOKEN")

	viper.OnConfigChange(cfg.onConfigChange)

	return cfg, nil
}

// Запускает обновление данных в объекте конфигурации при изменении исходного файла с данными
func (c *Config) onConfigChange(_ fsnotify.Event) {
	err := viper.Unmarshal(c)
	if err != nil {
		return
	}

	c.setLogLevel(c.LogLevel)
}

// SetLogLevel setup log level for app
func (c *Config) setLogLevel(logLevel string) {
	foundLogLevel, ok := LogLevelMap[logLevel]
	if !ok {
		jww.ERROR.Printf("incorrect log level %s\n", logLevel)
		return
	}

	jww.SetLogThreshold(foundLogLevel)
	jww.SetStdoutThreshold(foundLogLevel)
}

// LogLevelMap Содержит разрешённые уровни логирования;
// Чем выше уровень, тем больше выводится логов (снизу вверх)
var LogLevelMap = map[string]jww.Threshold{
	"DEBUG": jww.LevelDebug,
	"INFO":  jww.LevelInfo,
	"WARN":  jww.LevelWarn,
	"ERROR": jww.LevelError,
	"FATAL": jww.LevelFatal,
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

func (c *Config) Validate() error {
	return nil
}
