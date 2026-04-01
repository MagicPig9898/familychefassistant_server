package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AppConfig 全局配置
type AppConfig struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
	WX     WXConfig     `yaml:"wx"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	Salt string `yaml:"salt"`
}

type WXConfig struct {
	AppID           string `yaml:"app_id"`
	AppSecret       string `yaml:"app_secret"`
	Code2SessionURL string `yaml:"code2session_url"`
}

var Cfg *AppConfig

// MustLoad 从指定路径加载 YAML 配置，任何错误直接终止程序
func MustLoad(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	Cfg = &cfg
	return nil
}
