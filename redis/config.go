package redis

import "github.com/dennesshen/photon-core-starter/configuration"

func init() {
	configuration.Register(&config)
}

type Config struct {
	Redis struct {
		Hosts    []string `yaml:"hosts"`
		Password string   `yaml:"password"`
	} `yaml:"redis"`
}

var config Config
