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
		Port  int  `yaml:"port"`
		Debug bool `yaml:"debug"` // 是否为调试模式
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
// 当配置文件中缺少某些字段时，会自动使用默认值
func LoadConfig(configPath string) (*Config, error) {
	// 创建默认配置
	defaultConfig := createDefaultConfig()

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置并创建配置文件
		log.Println("配置文件不存在，正在创建默认配置...")
		if err := saveConfig(defaultConfig, configPath); err != nil {
			return nil, err
		}
		appConfig = defaultConfig
	} else {
		// 配置文件存在，读取配置
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

		// 合并配置，对缺失的字段使用默认值
		mergeConfigWithDefaults(&config, defaultConfig)

		appConfig = &config
	}

	// 设置管理员配置
	SetAdminConfigFromAppConfig(appConfig.Admin.Username, appConfig.Admin.Password)

	return appConfig, nil
}

// mergeConfigWithDefaults 合并配置文件中的配置和默认配置
// 对于配置文件中缺失的字段，使用默认配置中的对应值
func mergeConfigWithDefaults(config *Config, defaults *Config) {
	// 合并数据库配置
	if config.DB.Host == "" {
		config.DB.Host = defaults.DB.Host
	}
	if config.DB.Port == 0 {
		config.DB.Port = defaults.DB.Port
	}
	if config.DB.User == "" {
		config.DB.User = defaults.DB.User
	}
	if config.DB.Password == "" {
		config.DB.Password = defaults.DB.Password
	}
	if config.DB.Name == "" {
		config.DB.Name = defaults.DB.Name
	}

	// 合并服务器配置
	if config.Server.Port == 0 {
		config.Server.Port = defaults.Server.Port
	}
	// 注意：布尔值不需要特殊处理，因为它们在配置文件中缺失时会被正确解析为false

	// 合并JWT配置
	if config.JWT.Secret == "" {
		config.JWT.Secret = defaults.JWT.Secret
	}
	if config.JWT.ExpireHours == 0 {
		config.JWT.ExpireHours = defaults.JWT.ExpireHours
	}

	// 合并管理员配置
	if config.Admin.Username == "" {
		config.Admin.Username = defaults.Admin.Username
	}
	if config.Admin.Password == "" || config.Admin.Password == "verkeyoss" {
		config.Admin.Password = defaults.Admin.Password
	}
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
	config.Server.Port = 8913
	config.Server.Debug = false // 默认非调试模式
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
	log.Printf("请使用默认密码登录后立即修改密码")

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
