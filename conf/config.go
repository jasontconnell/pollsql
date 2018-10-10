package conf

import (
	"github.com/jasontconnell/conf"
)

type Config struct {
	ConnectionString string `json:"connectionString"`
}

func LoadConfig(file string) Config {
	cfg := Config{}

	conf.LoadConfig(file, &cfg)

	return cfg
}
