package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// HashPassword 单独的密码加密函数，用于需要预先加密密码的场景
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// IsBcryptHash 判断字符串是否为 bcrypt 哈希
func IsBcryptHash(password string) bool {
	// bcrypt哈希长度通常为60个字符，并且以 $2a$, $2b$ 或 $2y$ 开头
	if len(password) != 60 {
		return false
	}

	// 检查前缀
	prefix := password[:4]
	return prefix == "$2a$" || prefix == "$2b$" || prefix == "$2y$"
}

// App 应用模型
type App struct {
	gorm.Model
	UserID       uint      `gorm:"not null;index" json:"user_id"`             // 创建者ID，固定为管理员ID
	AKey         string    `gorm:"size:100;not null;uniqueIndex" json:"akey"` // 应用唯一标识
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  string    `gorm:"size:500" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	VersionCount int64     `gorm:"-" json:"version_count"`                // 版本数量，不映射到数据库字段
	IsPaid       bool      `gorm:"not null;default:false" json:"is_paid"` // 是否收费
	// 关联版本（一对多）
	Versions []Version `gorm:"foreignKey:AKey;references:AKey;constraint:OnDelete:CASCADE" json:"versions,omitempty"`
}

// Version 版本模型
type Version struct {
	gorm.Model
	VKey           string    `gorm:"size:100;not null;uniqueIndex" json:"vkey"`                      // 版本唯一标识
	AKey           string    `gorm:"size:100;not null;index;foreignKey;references:AKey" json:"akey"` // 外键，关联App结构体
	Version        string    `gorm:"size:50;not null" json:"version"`                                // 版本号
	Description    string    `gorm:"size:500" json:"description"`
	IsLatest       bool      `gorm:"not null;default:false" json:"is_latest"`
	IsForcedUpdate bool      `gorm:"not null;default:false" json:"is_forced_update"` // 是否强制更新
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Announcement 公告模型
type Announcement struct {
	gorm.Model
	Title       string    `gorm:"size:100;not null" json:"title"`         // 公告标题
	Content     string    `gorm:"type:text;not null" json:"content"`      // 公告内容
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"` // 是否激活
	PublishDate time.Time `json:"publish_date"`                           // 发布日期
	URL         string    `gorm:"size:500" json:"url,omitempty"`          // 公告链接，可选
}

// CheckRequest 校验请求模型

type CheckRequest struct {
	AKey string `json:"akey" binding:"required"`
	VKey string `json:"vkey" binding:"required"`
}

// ValidationResponse 合法性校验响应模型

type ValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
	AppName string `json:"app_name,omitempty"`
	Version string `json:"version,omitempty"`
}
