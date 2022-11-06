package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// LoggingParams struct
type LoggingParams struct {
	IsDebug *bool `yaml:"is_debug" env-default:"true"`
}

// ListenParams struct
type ListenParams struct {
	BindIp          string        `yaml:"bind_ip" env-default:"127.0.0.1"`
	Port            string        `yaml:"port" env-default:"8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"5s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
}

// StorageParams struct
type StorageParams struct {
	Backend  string        `yaml:"backend" env-default:"mongodb"`
	Lifetime time.Duration `yaml:"lifetime" env-default:"10s"`
	MongoDB  MongodbParams `yaml:"mongodb"`
}

type MongodbParams struct {
	Host        string `yaml:"host" env-default:"127.0.0.1"`
	Port        int    `yaml:"port" env-default:"27017"`
	User        string `yaml:"user" env-default:""`
	Pass        string `yaml:"pass" env-default:""`
	MaxPoolSize int    `yaml:"pool_size" env-default:"10"`
	Database    string `yaml:"database" env-default:"configs"`
}

// Config struct
type Config struct {
	Logging LoggingParams `yaml:"logging"`
	Listen  ListenParams  `yaml:"listen"`
	Storage StorageParams `yaml:"storage"`
}

var instance *Config
var once sync.Once

func GetConfig(args ...string) (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{}
		path := "config.yml"
		if len(args) > 0 && args[0] != "" {
			path = args[0]
		}
		err = cleanenv.ReadConfig(path, instance)
		if err != nil {
			desc, _ := cleanenv.GetDescription(instance, nil)
			fmt.Println(desc)
		}
	})
	return instance, err
}
