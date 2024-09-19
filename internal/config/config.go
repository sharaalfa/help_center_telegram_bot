package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env  string `yaml:"environment" env-default:"local"`
	Conf Conf   `yaml:"conf"`
}

type Conf struct {
	TelegramToken      string `yaml:"TELEGRAM_TOKEN"`
	PostgresUrl        string `yaml:"POSTGRES_URL"`
	RedisUrl           string `yaml:"REDIS_URL"`
	MongoDbUrl         string `yaml:"MONGODB_URL"`
	SupportAdminChatID int64  `yaml:"SUPPORT_ADMIN_CHAT_ID"`
	ITAdminChatID      int64  `yaml:"IT_ADMIN_CHAT_ID"`
	BillingAdminChatID int64  `yaml:"BILLING_ADMIN_CHAT_ID"`
}

// LoadConfig loads the configuration from the file specified by the -config flag or the CONFIG_PATH environment variable.
// The configuration is loaded into a Config struct.
// If the configuration file does not exist, or if there is an error reading the configuration file, the function returns an error.
func LoadConfig() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	return LoadConfigPath(configPath)
}

func LoadConfigPath(configPath string) (*Config, error) {
	cfg := &Config{}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}
	return cfg, nil
}

func fetchConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}

// fetchConfigPath returns the path to the config file.
// The path is either the value of the -config flag, or the value of the CONFIG_PATH environment variable.
// If neither is set, an empty string is returned.
