package config

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// AdminConfig 管理员配置结构
type AdminConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"` // 存储加密后的密码
}

// 全局变量存储管理员配置
var adminConfig *AdminConfig

// SetAdminConfigFromAppConfig 从应用配置设置管理员配置
func SetAdminConfigFromAppConfig(username, password string) {
	adminConfig = &AdminConfig{
		Username: username,
		Password: password,
	}
}

// GetAdminConfig 获取管理员配置
func GetAdminConfig() (*AdminConfig, error) {
	if adminConfig == nil {
		return nil, fmt.Errorf("管理员配置未初始化")
	}
	return adminConfig, nil
}

// UpdateAdminPassword 更新管理员密码
func UpdateAdminPassword(newPassword string) error {
	if adminConfig == nil {
		return fmt.Errorf("管理员配置未初始化")
	}

	// 加密新密码
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("加密新密码失败: %w", err)
	}

	// 更新密码
	adminConfig.Password = string(newHashedPassword)

	return nil
}
