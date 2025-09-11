package service

import (
	"errors"
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

var (
	ErrVersionNotFound = errors.New("版本不存在")
	ErrAKeyNotFound    = errors.New("软件标识不存在")
)

// VersionService 版本服务

type VersionService struct {
	store store.VersionStore
}

// NewVersionService 创建版本服务实例
func NewVersionService(store store.VersionStore) *VersionService {
	return &VersionService{store: store}
}

// CreateVersion 创建新版本
func (s *VersionService) CreateVersion(akey, version, description string, isLatest bool, isForcedUpdate bool) (*model.Version, error) {
	newVersion := &model.Version{
		AKey:           akey,
		Version:        version,
		Description:    description,
		IsLatest:       isLatest,
		IsForcedUpdate: isForcedUpdate,
	}

	err := s.store.CreateVersion(newVersion)
	if err != nil {
		return nil, err
	}

	return newVersion, nil
}

// GetVersionList 获取版本列表
func (s *VersionService) GetVersionList(akey string, page, size int) ([]*model.Version, int64, error) {
	// 分页参数校验
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	return s.store.GetVersionListByAKey(akey, page, size)
}

// GetVersionInfo 获取版本信息
func (s *VersionService) GetVersionInfo(vkey string) (*model.Version, error) {
	version, err := s.store.GetVersionByVKey(vkey)
	if err != nil {
		return nil, ErrVersionNotFound
	}

	return version, nil
}

// UpdateVersion 更新版本信息
func (s *VersionService) UpdateVersion(vkey, version, description string, isLatest bool, isForcedUpdate bool) error {
	// 获取版本信息
	versionInfo, err := s.store.GetVersionByVKey(vkey)
	if err != nil {
		return ErrVersionNotFound
	}

	// 更新字段
	if version != "" {
		versionInfo.Version = version
	}
	if description != "" {
		versionInfo.Description = description
	}
	// 只有明确设置了isLatest才更新
	versionInfo.IsLatest = isLatest
	// 更新强制更新字段
	versionInfo.IsForcedUpdate = isForcedUpdate

	return s.store.UpdateVersion(versionInfo)
}

// DeleteVersion 删除版本
func (s *VersionService) DeleteVersion(vkey string) error {
	// 检查版本是否存在
	_, err := s.store.GetVersionByVKey(vkey)
	if err != nil {
		return ErrVersionNotFound
	}

	return s.store.DeleteVersion(vkey)
}
