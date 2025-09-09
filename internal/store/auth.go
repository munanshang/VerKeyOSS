package store

import (
	"verkeyoss/internal/model"
)

// UserStoreImpl 用户存储实现

type UserStoreImpl struct {
	*Store
}

// NewUserStore 创建用户存储实例
func (s *Store) NewUserStore() *UserStoreImpl {
	return &UserStoreImpl{Store: s}
}

// GetUserByUsername 根据用户名获取用户信息
func (s *UserStoreImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户ID获取用户信息
// 预加载权限组信息，以便获取权限组名称
func (s *UserStoreImpl) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.DB.Where("id = ?", id).Preload("PermissionGroup").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserStoreImpl) UpdateUser(user *model.User) error {
	return s.DB.Save(user).Error
}

// GetAllUsers 获取所有用户列表（分页）
func (s *UserStoreImpl) GetAllUsers(page, size int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	s.DB.Model(&model.User{}).Count(&total)

	// 查询列表，预加载权限组信息
	err := s.DB.Limit(size).Offset(offset).Preload("PermissionGroup").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// BanUser 封禁/解封用户
func (s *UserStoreImpl) BanUser(id string, banned bool, banReason string) error {
	// 更新用户的封禁状态和封禁原因
	return s.DB.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_banned":   banned,
		"ban_reason":  banReason,
	}).Error
}
