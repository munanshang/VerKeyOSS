package service

import (
	"errors"
	"time"
	"verkeyoss/internal/config"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

// AuthService 管理员认证服务
// 负责处理管理员登录、密码修改和令牌验证等功能
// 基于本地文件存储管理员账号信息

type AuthService struct {
	jwtSecret  []byte
	expireTime time.Duration
}

// NewAuthService 创建认证服务实例
// 参数 jwtSecret JWT令牌的密钥
// 参数 expireHours 令牌有效期（小时）
func NewAuthService(jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		jwtSecret:  []byte(jwtSecret),
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

// Login 管理员登录
// 参数 username 用户名
// 参数 password 密码
// 返回 token JWT令牌（包含admin:true声明，用于验证管理员身份）
// 返回 expiresAt 过期时间
// 返回 error 错误信息
func (s *AuthService) Login(username, password string) (string, string, error) {
	// 获取管理员配置
	adminConfig, err := config.GetAdminConfig()
	if err != nil {
		return "", "", errors.New("获取管理员配置失败")
	}

	// 验证用户名
	if username != adminConfig.Username {
		return "", "", ErrInvalidCredentials
	}

	// 验证密码
	if pwErr := bcrypt.CompareHashAndPassword([]byte(adminConfig.Password), []byte(password)); pwErr != nil {
		return "", "", ErrInvalidCredentials
	}

	// 生成JWT令牌
	expirationTime := time.Now().Add(s.expireTime)

	// 创建声明
	claims := &jwt.MapClaims{
		"admin": true,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	// 创建token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", errors.New("生成令牌失败")
	}

	// 格式化过期时间
	expiresAt := expirationTime.Format(time.RFC3339)

	// 返回token和过期时间
	return tokenString, expiresAt, nil
}

// ChangePassword 修改管理员密码
// 参数 oldPassword 旧密码
// 参数 newPassword 新密码
// 返回 error 错误信息
func (s *AuthService) ChangePassword(oldPassword, newPassword string) error {
	// 获取管理员配置
	adminConfig, err := config.GetAdminConfig()
	if err != nil {
		return errors.New("获取管理员配置失败")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(adminConfig.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// 更新密码
	return config.UpdateAdminPassword(newPassword)
}

// VerifyToken 验证管理员令牌
// 参数 tokenString JWT令牌字符串
// 返回 bool 是否有效
// 返回 error 错误信息
// 验证内容包括令牌签名、过期时间和admin:true声明
func (s *AuthService) VerifyToken(tokenString string) (bool, error) {
	// 检查token是否为空
	if tokenString == "" {
		return false, errors.New("令牌为空")
	}

	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return false, errors.New("令牌无效或已过期")
	}

	// 验证token是否有效并检查是否为管理员
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 检查是否包含管理员标识
		adminClaim, hasAdminClaim := claims["admin"].(bool)
		if !hasAdminClaim || !adminClaim {
			return false, errors.New("无效的管理员令牌")
		}
		return true, nil
	} else {
		return false, errors.New("令牌无效或已过期")
	}
}
