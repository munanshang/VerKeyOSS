package service

import (
	"errors"
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

var (
	ErrSoftwareNotFound = errors.New("软件不存在")
	ErrSoftwareExists   = errors.New("软件已存在")
)

// SoftwareService 软件服务

type SoftwareService struct {
	store store.SoftwareStore
}

// NewSoftwareService 创建软件服务实例
func NewSoftwareService(store store.SoftwareStore) *SoftwareService {
	return &SoftwareService{store: store}
}

// CreateSoftware 创建新软件
func (s *SoftwareService) CreateSoftware(name, description string) (*model.Software, error) {
	software := &model.Software{
		Name:        name,
		Description: description,
	}

	err := s.store.CreateSoftware(software)
	if err != nil {
		return nil, err
	}

	return software, nil
}

// GetSoftwareList 获取软件列表
func (s *SoftwareService) GetSoftwareList(page, size int) ([]*model.Software, int64, error) {
	// 分页参数校验
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	return s.store.GetSoftwareList(page, size)
}

// GetSoftwareInfo 获取软件信息
func (s *SoftwareService) GetSoftwareInfo(akey string) (*model.Software, error) {
	software, err := s.store.GetSoftwareByAKey(akey)
	if err != nil {
		return nil, ErrSoftwareNotFound
	}

	return software, nil
}

// UpdateSoftware 更新软件信息
func (s *SoftwareService) UpdateSoftware(akey, name, description string) error {
	// 获取软件信息
	software, err := s.store.GetSoftwareByAKey(akey)
	if err != nil {
		return ErrSoftwareNotFound
	}

	// 更新字段
	if name != "" {
		software.Name = name
	}
	if description != "" {
		software.Description = description
	}

	return s.store.UpdateSoftware(software)
}

// DeleteSoftware 删除软件
func (s *SoftwareService) DeleteSoftware(akey string) error {
	// 检查软件是否存在
	_, err := s.store.GetSoftwareByAKey(akey)
	if err != nil {
		return ErrSoftwareNotFound
	}

	return s.store.DeleteSoftware(akey)
}