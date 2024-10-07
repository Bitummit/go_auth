package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)


type Config struct {
	Env string `yaml:"env" env-default:"dev"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"0.0.0.0:8000"`
	Timeout time.Duration `yaml:"timeout" env-default:"3s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"40s"`
}


func InitConfig() *Config{
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file!")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("Empty path")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalln("Can not find config file")
	} 
	
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalln("Error in reading config file!")
	}
	
	return &cfg
	
}