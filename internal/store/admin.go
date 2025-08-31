package store

import (
	"verkeyoss/internal/model"
)

// AdminStoreImpl 管理员存储实现

type AdminStoreImpl struct {
	*Store
}

// NewAdminStore 创建管理员存储实例
func (s *Store) NewAdminStore() *AdminStoreImpl {
	return &AdminStoreImpl{Store: s}
}

// GetAdminByUsername 根据用户名获取管理员信息
func (s *AdminStoreImpl) GetAdminByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := s.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateAdmin 更新管理员信息
func (s *AdminStoreImpl) UpdateAdmin(admin *model.Admin) error {
	return s.DB.Save(admin).Error
}
