package config

import "github.com/kelseyhightower/envconfig"

const envPrefix = "QA"

type Config struct {
	Host   string `split_words:"true" default:"localhost:8082"`
	DbHost string `split_words:"true" default:"localhost"`
	DbPort int    `split_words:"true" default:"5432"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	return cfg, err
}
