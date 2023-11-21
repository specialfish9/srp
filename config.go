package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Mappings []Mapping
}

type Mapping struct {
	Key         string
	Destination string
}

func LoadConfig() *Config {

	var config = Config{}
	err := cleanenv.ReadConfig("config.yml", &config)

	if err != nil {
		log.Fatal(err)
	}

	return &config
}
