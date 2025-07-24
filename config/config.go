/*
 * @Author: Justin
 * @Date: 2025-07-21 03:14:23
 * @filename:
 * @version:
 * @Description:
 * @LastEditTime: 2025-07-21 07:04:51
 */
package config

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	clog "github.com/kpango/glg"
)

// 加载 server的配置文件
func LoadConfig(path string) *ServerConfig {
	var cfg ServerConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		clog.Errorf("Failed to load config: %v", err)
	}
	clog.Info("load config")
	return &cfg
}

// 解码获得 body中的配置文件
func DecodeSingleConfig(body []byte) *SingleConfig {
	var cfg SingleConfig
	if err := json.Unmarshal(body, &cfg); err != nil {
		clog.Errorf("Failed to decode config: %v", err)
	}
	return &cfg
}
