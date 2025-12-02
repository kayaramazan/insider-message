package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Redis   *RedisConfig    `mapstructure:"redis"`
	Db      *DatabaseConfig `mapstructure:"db"`
	Timer   *TimerConfig    `mapstructure:"timer"`
	Webhook *WebhookConfig  `mapstructure:"webhook"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type TimerConfig struct {
	Interval        int `mapstructure:"interval"`
	MessagePerCycle int `mapstructure:"message_per_cycle"`
}

type WebhookConfig struct {
	Url string `mapstructure:"url"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {

	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.password", "")

	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", 5432)
	v.SetDefault("db.user", "postgres")
	v.SetDefault("db.password", "Aa123456")
	v.SetDefault("db.dbname", "postgres")
	v.SetDefault("db.sslmode", "disable")

	v.SetDefault("timer.interval", 2)
	v.SetDefault("timer.message_per_cycle", 2)

	v.SetDefault("webhook.url", "")
}
