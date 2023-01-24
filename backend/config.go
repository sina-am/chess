package main

import "github.com/Netflix/go-env"

type Config struct {
	SrvAddr   string `env:"SERVER_ADDRESS,default=:6969"`
	DBAddr    string `env:"DATABASE_ADDRESS,default=mongodb://localhost"`
	DBName    string `env:"DATABASE_NAME,default=chess"`
	SecretKey string `env:"SECRET_KEY,default=notsave"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	_, err := env.UnmarshalFromEnviron(config)
	return config, err
}
