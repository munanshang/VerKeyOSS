package service

import (
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

// PermissionGroupService 权限组服务
type PermissionGroupService struct {
	store store.PermissionGroupStore
}

// NewPermissionGroupService 创建权限组服务实例
func NewPermissionGroupService(store store.PermissionGroupStore) *PermissionGroupService {
	return &PermissionGroupService{store: store}
}

// GetAllPermissionGroups 获取所有权限组信息
// 包括权限组ID、权限组名称、权限配置等
func (s *PermissionGroupService) GetAllPermissionGroups() ([]*model.PermissionGroup, error) {
	return s.store.GetAllPermissionGroups()
}

// GetPermissionGroupByID 根据ID获取权限组信息
func (s *PermissionGroupService) GetPermissionGroupByID(id uint) (*model.PermissionGroup, error) {
	return s.store.GetPermissionGroupByID(id)
}