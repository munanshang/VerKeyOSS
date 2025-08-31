package service

import (
	"errors"
	"time"

	"verkeyoss/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrAdminNotFound      = errors.New("管理员不存在")
)

// AdminService 管理员服务
type AdminService struct {
	store      store.AdminStore
	jwtSecret  []byte
	expireTime time.Duration
}

// NewAdminService 创建管理员服务实例
func NewAdminService(store store.AdminStore, jwtSecret string, expireHours int) *AdminService {
	return &AdminService{
		store:      store,
		jwtSecret:  []byte(jwtSecret),
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

// Login 管理员登录
func (s *AdminService) Login(username, password string) (string, string, error) {
	// 获取管理员信息
	admin, err := s.store.GetAdminByUsername(username)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	// 生成JWT令牌
	expirationTime := time.Now().Add(s.expireTime)

	// 创建声明
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
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

	return tokenString, expiresAt, nil
}

// ChangePassword 修改管理员密码
func (s *AdminService) ChangePassword(username, oldPassword, newPassword string) error {
	// 获取管理员信息
	admin, err := s.store.GetAdminByUsername(username)
	if err != nil {
		return ErrAdminNotFound
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	// 更新密码（注意：密码加密会在 BeforeSave 钩子中自动处理）
	admin.Password = newPassword
	return s.store.UpdateAdmin(admin)
}

// VerifyToken 验证管理员令牌
func (s *AdminService) VerifyToken(tokenString string) (string, error) {
	// 检查token是否为空
	if tokenString == "" {
		return "", errors.New("令牌为空")
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
		return "", errors.New("令牌无效或已过期")
	}

	// 验证token是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 从声明中获取用户名
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("令牌中没有用户名")
		}

		// 验证管理员是否存在
		_, err := s.store.GetAdminByUsername(username)
		if err != nil {
			return "", ErrAdminNotFound
		}

		return username, nil
	} else {
		return "", errors.New("令牌无效或已过期")
	}
}
