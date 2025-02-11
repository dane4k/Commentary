package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	Port int `yaml:"port" env-required:"true"`
}

type Database struct {
	Host      string `yaml:"host" env-required:"true"`
	Port      int    `yaml:"port" env-required:"true"`
	User      string `yaml:"user" env-required:"true"`
	Password  string `yaml:"password" env-required:"true"`
	Name      string `yaml:"name" env-required:"true"`
	StoreInDB bool   `yaml:"store_in_db"`
}

type Logger struct {
	FileName string `yaml:"filename" env-required:"true"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Logger   Logger   `yaml:"logger"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.Fatal(err.Error())
		return nil, err
	}

	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	logrus.Info("Loaded app configuration")

	//return &Config{
	//	Server: Server{
	//		Port: 8080,
	//	},
	//	Database: Database{
	//		Host:      "localhost",
	//		Port:      5432,
	//		User:      "postgres",
	//		Password:  "admin",
	//		Name:      "commentary",
	//		StoreInDB: false,
	//	},
	//	Logger: Logger{
	//		FileName: "Commentary.log",
	//	},
	//}, nil

	return &config, nil
}
