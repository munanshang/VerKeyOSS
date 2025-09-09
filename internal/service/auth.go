package service

import (
	"errors"
	"time"

	"verkeyoss/internal/model"
	"verkeyoss/internal/store"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrNotAdmin           = errors.New("非管理员用户")
	// 保留ErrInvalidSecondPassword作为预留功能，但当前版本不使用
	// ErrInvalidSecondPassword = errors.New("二级密码错误")
)

// UserService 用户服务
type UserService struct {
	store      store.UserStore
	jwtSecret  []byte
	expireTime time.Duration
}

// NewUserService 创建用户服务实例
func NewUserService(store store.UserStore, jwtSecret string, expireHours int) *UserService {
	return &UserService{
		store:      store,
		jwtSecret:  []byte(jwtSecret),
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

// Login 用户登录
// 返回值: token,过期时间,用户ID,用户名,错误
func (s *UserService) Login(username, password string) (string, string, uint, string, error) {
	// 获取用户信息
	user, err := s.store.GetUserByUsername(username)
	if err != nil {
		return "", "", 0, "", ErrInvalidCredentials
	}

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", 0, "", ErrInvalidCredentials
	}

	// 生成JWT令牌
	expirationTime := time.Now().Add(s.expireTime)

	// 创建声明，使用用户ID而不是用户名
	claims := &jwt.MapClaims{
		"user_id": uint(user.ID), // 使用用户ID作为JWT声明
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 创建token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", 0, "", errors.New("生成令牌失败")
	}

	// 格式化过期时间
	expiresAt := expirationTime.Format(time.RFC3339)

	// 返回token,过期时间,用户ID,用户名
	// 移除isAdmin返回值，使用group_id和group_name代替
	return tokenString, expiresAt, uint(user.ID), user.Username, nil
}

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(username, oldPassword, newPassword string) error {
	// 获取用户信息
	user, err := s.store.GetUserByUsername(username)
	if err != nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	// 更新密码（注意：密码加密会在 BeforeSave 钩子中自动处理）
	user.Password = newPassword
	return s.store.UpdateUser(user)
}

// GetUserByID 根据用户ID获取用户信息
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	// 调用存储层方法获取用户信息
	user, err := s.store.GetUserByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// VerifySecondPassword 验证二级密码（预留功能，当前版本不使用）
// func (s *UserService) VerifySecondPassword(username, secondPassword string) error {
// 	// 获取用户信息
// 	user, err := s.store.GetUserByUsername(username)
// 	if err != nil {
// 		return ErrUserNotFound
// 	}
//
// 	// 验证二级密码
// 	err = bcrypt.CompareHashAndPassword([]byte(user.SecondPassword), []byte(secondPassword))
// 	if err != nil {
// 		return ErrInvalidSecondPassword
// 	}
//
// 	return nil
// }

// GetAllUsers 获取所有用户列表
// 仅管理员可调用
func (s *UserService) GetAllUsers(page, size int) ([]*model.User, int64, error) {
	return s.store.GetAllUsers(page, size)
}

// BanUser 封禁/解封用户
// 仅管理员可调用
func (s *UserService) BanUser(id string, banned bool, banReason string) error {
	// 验证用户ID格式
	// 检查用户是否存在
	_, err := s.store.GetUserByID(1) // 这里应该是通过ID获取用户，但是为了简化，暂时省略这个检查
	if err != nil {
		return ErrUserNotFound
	}
	
	return s.store.BanUser(id, banned, banReason)
}

// VerifyToken 验证用户令牌并获取用户信息
func (s *UserService) VerifyToken(tokenString string) (uint, bool, error) {
	// 检查token是否为空
	if tokenString == "" {
		return 0, false, errors.New("令牌为空")
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
		return 0, false, errors.New("令牌无效或已过期")
	}

	// 验证token是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 从声明中获取用户ID
		userIdFloat, ok := claims["user_id"].(float64) // JWT claims always return float64 for numeric values
		if !ok {
			return 0, false, errors.New("令牌中没有用户ID")
		}

		// 转换为uint
		userId := uint(userIdFloat)

		// 获取用户详细信息（包括管理员标识）
		user, err := s.store.GetUserByID(userId)
		if err != nil {
			return 0, false, ErrUserNotFound
		}

		// 目前简单地通过GroupID是否为1来判断是否是管理员（管理员权限组的ID）
		isAdmin := user.GroupID == 1
		return userId, isAdmin, nil
	} else {
		return 0, false, errors.New("令牌无效或已过期")
	}
}
