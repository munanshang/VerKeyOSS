package store

import (
	"verkeyoss/internal/model"
)

// PermissionGroupStoreImpl 权限组存储实现

type PermissionGroupStoreImpl struct {
	*Store
}

// NewPermissionGroupStore 创建权限组存储实例
func (s *Store) NewPermissionGroupStore() *PermissionGroupStoreImpl {
	return &PermissionGroupStoreImpl{Store: s}
}

// GetAllPermissionGroups 获取所有权限组
func (s *PermissionGroupStoreImpl) GetAllPermissionGroups() ([]*model.PermissionGroup, error) {
	var groups []*model.PermissionGroup
	err := s.DB.Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// GetPermissionGroupByID 根据ID获取权限组
func (s *PermissionGroupStoreImpl) GetPermissionGroupByID(id uint) (*model.PermissionGroup, error) {
	var group model.PermissionGroup
	err := s.DB.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}