package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string
	// check in env first
	configPath = os.Getenv("CONFIG_PATH")
	//if not in env then check in flags
	if configPath == "" {
		flags := flag.String("config", "", "Path to the configuration file")
		flag.Parse()

		// fmt.Println(flags)
		configPath = *flags
		// if not found on flag then show error message
		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}
	// now check config file present or not in  directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s not found", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config: %s", err.Error())
	}
	return &cfg
}
