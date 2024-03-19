package config

import (
	"sync"

	"github.com/nected/go-lib/generators"
)

type Config struct {
	AppConfig AppConfig    `json:"app" yaml:"app" `
	Server    ServerConfig `json:"server" yaml:"server"`
	Raft      RaftConfig   `json:"raft" yaml:"raft"`
	Test      string       `json:"test" yaml:"test" default:"test"`
}

type AppConfig struct {
	Name string `json:"name" yaml:"name" default:"sanchaalak"`
}

type ServerConfig struct {
	Host string `json:"host" yaml:"host" default:"127.0.0.1"`
	Port int    `json:"port" yaml:"port" default:"8080"`
}

type RaftConfig struct {
	StoragePath string `json:"storagePath" yaml:"storagePath" default:"/tmp/sanchaalak/raft"`
	NodeInfo    struct {
		ID      string `json:"id" yaml:"id" default:"node1"`
		Address string `json:"address" yaml:"address" default:"127.0.0.1:1234"`
	} `json:"node" yaml:"node"`
}

var (
	config *Config
	lock   = &sync.Mutex{}
)

func NewConfig() *Config {
	return &Config{}
}

func SetConfig(cfg *Config) {
	lock.Lock()
	defer lock.Unlock()
	config = cfg
}

func GetConfig() *Config {
	return config
}

func (c *Config) GenerateDefaults() {
	generators.GenerateDefaults(c)
}
