package agridataintegration

import (
	"encoding/json"
	"os"
)

// Config 集成系统配置
type Config struct {
	// 农业平台配置
	AgriPlatform struct {
		BaseURL  string `json:"baseUrl"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"agriPlatform"`

	// Magistrala 平台配置
	Magistrala struct {
		BaseURL   string `json:"baseUrl"`
		UserToken string `json:"userToken"`
		DomainID  string `json:"domainId"`
		ChannelID string `json:"channelId"`
	} `json:"magistrala"`

	// 集成配置
	Integration struct {
		SyncInterval     int    `json:"syncInterval"`     // 数据同步间隔（秒）
		MappingFile      string `json:"mappingFile"`      // 映射表文件路径
		BackgroundImage  string `json:"backgroundImage"`  // 背景图片名称
		DefaultPartition string `json:"defaultPartition"` // 默认分区名称
	} `json:"integration"`

	// 服务器配置
	Server struct {
		Port string `json:"port"` // 服务端口
		Host string `json:"host"` // 服务主机
	} `json:"server"`
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		AgriPlatform: struct {
			BaseURL  string `json:"baseUrl"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			BaseURL:  "http://www.0531yun.com",
			Username: "",
			Password: "",
		},
		Magistrala: struct {
			BaseURL   string `json:"baseUrl"`
			UserToken string `json:"userToken"`
			DomainID  string `json:"domainId"`
			ChannelID string `json:"channelId"`
		}{
			BaseURL:   "http://localhost:9002",
			UserToken: "",
			DomainID:  "",
			ChannelID: "",
		},
		Integration: struct {
			SyncInterval     int    `json:"syncInterval"`
			MappingFile      string `json:"mappingFile"`
			BackgroundImage  string `json:"backgroundImage"`
			DefaultPartition string `json:"defaultPartition"`
		}{
			SyncInterval:     30, // 30秒同步一次
			MappingFile:      "sensor_mapping.json",
			BackgroundImage:  "farm_layout.jpg",
			DefaultPartition: "field_1",
		},
		Server: struct {
			Port string `json:"port"`
			Host string `json:"host"`
		}{
			Port: "8888",
			Host: "0.0.0.0",
		},
	}
}
