package configs

import (
	redisConfig "calendar/internal/app/redis"
	"calendar/internal/app/storage"
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)

var Conf *Config
var once sync.Once

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Storage *storage.Config
	RedisConfig *redisConfig.Config
}

func NewConfig() *Config {

	once.Do(func() {
		Conf = &Config {
			BindAddr: ":8080",
			LogLevel: "info",
			Storage: storage.NewConfig(),
			RedisConfig: redisConfig.NewConfig(),
		}
		_, err := toml.DecodeFile("configs/config.toml", Conf)
		if err != nil {
			log.Fatal(err)
		}
	})

	return Conf
}
