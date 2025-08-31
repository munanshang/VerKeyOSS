package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	gorm.Model
	Username string `gorm:"size:50;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"` // 不对外暴露密码
}

// BeforeCreate 在创建前对密码进行加密
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	// 使用 bcrypt 对密码进行加密
	if a.Password != "" && !isBcryptHash(a.Password) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashedPassword)
	}
	return nil
}

// BeforeSave 在保存前对密码进行加密（适用于创建和更新操作）
func (a *Admin) BeforeSave(tx *gorm.DB) error {
	// 只在密码不为空且不是bcrypt哈希时才进行加密
	if a.Password != "" && !isBcryptHash(a.Password) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		a.Password = string(hashedPassword)
	}
	return nil
}

// HashPassword 单独的密码加密函数，用于需要预先加密密码的场景
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// isBcryptHash 判断字符串是否为 bcrypt 哈希
func isBcryptHash(password string) bool {
	// bcrypt哈希长度通常为60个字符，并且以 $2a$, $2b$ 或 $2y$ 开头
	if len(password) != 60 {
		return false
	}

	// 检查前缀
	prefix := password[:4]
	return prefix == "$2a$" || prefix == "$2b$" || prefix == "$2y$"
}

// Software 软件模型

type Software struct {
	gorm.Model
	AKey         string    `gorm:"size:100;not null;uniqueIndex" json:"akey"` // 软件唯一标识
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  string    `gorm:"size:500" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	VersionCount int64     `gorm:"-" json:"version_count"` // 版本数量，不映射到数据库字段
	// 关联版本（一对多）
	Versions []Version `gorm:"foreignKey:AKey;references:AKey;constraint:OnDelete:CASCADE" json:"versions,omitempty"`
}

// Version 版本模型

type Version struct {
	gorm.Model
	VKey        string    `gorm:"size:100;not null;uniqueIndex" json:"vkey"` // 版本唯一标识
	AKey        string    `gorm:"size:100;not null;index" json:"akey"`
	Version     string    `gorm:"size:50;not null" json:"version"` // 版本号
	Description string    `gorm:"size:500" json:"description"`
	IsLatest    bool      `gorm:"not null;default:false" json:"is_latest"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// CheckRequest 校验请求模型

type CheckRequest struct {
	AKey string `json:"akey" binding:"required"`
	VKey string `json:"vkey" binding:"required"`
}

// LegalityResponse 合法性校验响应模型

type LegalityResponse struct {
	Legal   bool   `json:"legal"`
	Message string `json:"message"`
}
