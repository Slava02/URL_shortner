package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"` //  yaml - значение в файле, env - меременные среды
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http-server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := flag.String("conf", "config/local.yaml", "Путь до файла конфигурации. Чтобы использовать переменную окружения - необходимо передать ее название при запуске программы: -conf=$CONFIG_PATH")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", *configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
