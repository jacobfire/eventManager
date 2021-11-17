package redisConfig

type Config struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
	Password string `toml:"password"`
	DB int `toml:"db"`
}

func NewConfig() *Config {
	return &Config{}
}