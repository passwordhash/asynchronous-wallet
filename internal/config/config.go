package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App  AppConfig      `yaml:"app"`
	HTTP HttpConfig     `yaml:"http"`
	PG   PostgresConfig `yaml:"postgres"`
}

type AppConfig struct {
	Env string `env:"ENV" yaml:"env" env-required:"true"`
}

type HttpConfig struct {
	Port         int           `env:"PORT" yaml:"port" env-required:"true"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" yaml:"write_timeout" env-default:"10"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" yaml:"read_timeout" env-default:"10"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" yaml:"port" env-required:"true"`
	Username string `env:"POSTGRES_USER" yaml:"user" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" yaml:"password" env-required:"true"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d", p.Username, p.Password, p.Host, p.Port)
}

// MustLoad loads the configuration from a file specified by the `config` flag or
// the `CONFIG_PATH` environment variable. If the configuration file is not found
// or cannot be read, it panics with an error message.
func MustLoad() *Config {
	cfgPath := fetchConfigPath()
	if cfgPath == "" {
		panic("config path is not set")
	}

	_, err := os.Stat(cfgPath)
	if err != nil && os.IsPermission(err) {
		panic("no permission to config: " + cfgPath)
	}
	if err != nil && os.IsNotExist(err) {
		panic("there is no config file: " + cfgPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	return cfg
}

// fetchConfigPath retrieves the configuration file path from command line flags
// or the `CONFIG_PATH` environment variable. It returns the path as a string.
// If neither is provided, it returns an empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
