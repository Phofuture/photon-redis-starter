package redis

import "github.com/dennesshen/photon-core-starter/configuration"

func init() {
	configuration.Register(&config)
}

type ClientType string

const (
	ClientTypeStandalone ClientType = "standalone"
	ClientTypeCluster    ClientType = "cluster"
)

type Config struct {
	Redis struct {
		Type     string   `yaml:"type"`
		Hosts    []string `yaml:"hosts"`
		Password string   `yaml:"password"`
	} `yaml:"redis"`
}

var config Config
