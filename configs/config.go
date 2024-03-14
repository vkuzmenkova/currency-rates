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
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Host      string
	Port      string
	NameSpace string
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

	fmt.Println("config", config)

	return config, nil
}
