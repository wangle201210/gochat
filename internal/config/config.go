package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	AI        AIConfig        `json:"ai"`
	Assistant AssistantConfig `json:"assistant"`
	UI        UIConfig        `json:"ui"`
}

// AIConfig AI 相关配置
type AIConfig struct {
	Provider string `json:"provider"` // 例如: "ark", "openai", "anthropic"
	Model    string `json:"model"`    // 模型名称
	APIKey   string `json:"api_key"`  // API Key
	BaseURL  string `json:"base_url"` // API Base URL
}

// AssistantConfig 助手模型配置（用于生成会话标题等辅助任务）
type AssistantConfig struct {
	Provider string `json:"provider"` // 例如: "ark", "openai", "anthropic"
	Model    string `json:"model"`    // 模型名称
	APIKey   string `json:"api_key"`  // API Key
	BaseURL  string `json:"base_url"` // API Base URL
}

// UIConfig UI 相关配置
type UIConfig struct {
	WindowWidth  int `json:"window_width"`
	WindowHeight int `json:"window_height"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		AI: AIConfig{
			Provider: "openai",
			Model:    "gpt-3.5-turbo",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "APIKey",
		},
		Assistant: AssistantConfig{
			Provider: "openai",
			Model:    "gpt-3.5-turbo",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "APIKey",
		},
		UI: UIConfig{
			WindowWidth:  800,
			WindowHeight: 600,
		},
	}
}

// Load 从文件加载配置
func Load(configPath string) (*Config, error) {
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}

// Save 保存配置到文件
func (c *Config) Save(configPath string) error {
	// 确保配置目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(homeDir, ".gochat", "config.json")
}
