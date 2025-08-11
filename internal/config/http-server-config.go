package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type HttpServerConfig struct {
	Env            string `yaml:"env" env-default:"development"`
	HTTPServer     `yaml:"http_server"`
	CachedExecutor `yaml:"cached_executor"`
	Runner         `yaml:"runner"`
}

type HTTPServer struct {
	Address        string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout        time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ContextTimeout time.Duration `yaml:"context_timeout" env-default:"10s"`
}

type CachedExecutor struct {
	Ttl         time.Duration `yaml:"ttl" env-default:"60s"`
	MaxParallel int           `yaml:"max_parallel" env-default:"5"`
}

type Runner struct {
	CompileTimeout       int64 `yaml:"compile_timeout"`
	RunTimeout           int64 `yaml:"run_timeout"`
	RunCPUTimeout        int64 `yaml:"run_cpu_timeout"`
	CompileCPUTimeout    int64 `yaml:"compile_cpu_timeout"`
	CompileMemoryLimitKB int64 `yaml:"compile_memory_limit_KB"`
	RunMemoryLimitKB     int64 `yaml:"run_memory_limit_KB"`
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
