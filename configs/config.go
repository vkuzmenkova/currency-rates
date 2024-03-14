package configs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Service ServiceConfig `mapstructure:"server"`
	DB      DBConfig      `mapstructure:"db"`
	Redis   RedisConfig   `mapstructure:"redis"`
}

type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	Host      string
	Port      string
	Namespace string
	JobRetry  uint
}

func NewConfig(path, name string) (Config, error) {
	deployment := os.Getenv("DEPLOYMENT")
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))

	var config Config

	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("viper.ReadInConfig: %w", err)
	}

	err = v.UnmarshalKey(deployment, &config)
	if err != nil {
		return Config{}, fmt.Errorf("viper.Unmarshal: %w", err)
	}

	return config, nil
}
