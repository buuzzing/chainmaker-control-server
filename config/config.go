package config

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	clog "github.com/kpango/glg"
)

// LoadConfig 加载 server 的配置文件
func LoadConfig(path string) *ServerConfig {
	var cfg ServerConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		clog.Errorf("Failed to load config: %v", err)
		return nil
	}
	return &cfg
}

// DecodeSingleConfig 解码获得 body 中的配置文件
func DecodeSingleConfig(body []byte) *SingleConfig {
	var cfg SingleConfig
	if err := json.Unmarshal(body, &cfg); err != nil {
		clog.Errorf("Failed to decode config: %v", err)
	}
	return &cfg
}
