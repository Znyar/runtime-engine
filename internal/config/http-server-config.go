package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type HttpServerConfig struct {
	Env        string `yaml:"env" env-default:"development"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address        string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout        time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ContextTimeout time.Duration `yaml:"context_timeout" env-default:"10s"`
	CachedExecutor `yaml:"cached_executor"`
}

type CachedExecutor struct {
	Ttl         time.Duration `yaml:"ttl" env-default:"60s"`
	MaxParallel int           `yaml:"max_parallel" env-default:"5"`
}

func MustLoad() *HttpServerConfig {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg HttpServerConfig

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
