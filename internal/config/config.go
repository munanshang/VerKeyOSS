package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret      string `yaml:"secret"`
		ExpireHours int    `yaml:"expire_hours"`
	} `yaml:"jwt"`
	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"` // 存储加密后的密码
	} `yaml:"admin"`
}

// 全局变量存储应用配置
var appConfig *Config

// LoadConfig 加载配置文件并初始化全局配置
func LoadConfig(configPath string) (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		log.Println("配置文件不存在，正在创建默认配置...")
		defaultConfig := createDefaultConfig()
		if err := saveConfig(defaultConfig, configPath); err != nil {
			return nil, err
		}
		appConfig = defaultConfig
	} else {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var config Config
		d := yaml.NewDecoder(file)
		if err := d.Decode(&config); err != nil {
			return nil, err
		}

		appConfig = &config
	}

	// 设置管理员配置
	SetAdminConfigFromAppConfig(appConfig.Admin.Username, appConfig.Admin.Password)

	return appConfig, nil
}

// GetAppConfig 获取应用配置
func GetAppConfig() (*Config, error) {
	if appConfig == nil {
		return nil, fmt.Errorf("应用配置未初始化")
	}
	return appConfig, nil
}

// 创建默认配置
func createDefaultConfig() *Config {
	// 生成随机JWT密钥
	jwtSecret, _ := generateRandomString(32)
	// 默认管理员用户名和密码
	defaultUsername := "verkeyoss"
	defaultPassword := "verkeyoss"

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)

	config := &Config{}
	config.DB.Host = "localhost"
	config.DB.Port = 3306
	config.DB.User = "verkeyoss"
	config.DB.Password = "verkeyoss"
	config.DB.Name = "verkeyoss"
	config.Server.Port = 8080
	config.JWT.Secret = jwtSecret
	config.JWT.ExpireHours = 24
	config.Admin.Username = defaultUsername
	config.Admin.Password = string(hashedPassword)

	return config
}

// 保存配置到文件
func saveConfig(config *Config, filePath string) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(filePath, content, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	log.Printf("默认配置已创建，请在首次登录后修改默认密码")
	log.Printf("用户名: %s", config.Admin.Username)
	log.Printf("密码: verkeyoss")

	return nil
}

// 生成随机字符串
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
