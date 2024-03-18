package config

import "sync"

type Config struct {
	AppConfig AppConfig    `json:"app" yaml:"app" `
	Server    ServerConfig `json:"server" yaml:"server"`
	Raft      RaftConfig   `json:"raft" yaml:"raft"`
}

type AppConfig struct {
	Name string `json:"name" yaml:"name" default:"sanchaalak"`
}

type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type RaftConfig struct {
	StoragePath string `json:"storagePath" yaml:"storagePath"`
	NodeInfo    struct {
		ID      string `json:"id" yaml:"id"`
		Address string `json:"address" yaml:"address"`
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
	c.AppConfig.Name = "sanchaalak"
	c.Server.Host = "localhost"
	c.Server.Port = 8080
	c.Raft.StoragePath = "/tmp/sanchaalak"
	c.Raft.NodeInfo.ID = "node-1"
	c.Raft.NodeInfo.Address = "localhost:8080"
}
