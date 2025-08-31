package store

import (
	"github.com/google/uuid"
	"verkeyoss/internal/model"
)

// SoftwareStoreImpl 软件存储实现

type SoftwareStoreImpl struct {
	*Store
}

// NewSoftwareStore 创建软件存储实例
func (s *Store) NewSoftwareStore() *SoftwareStoreImpl {
	return &SoftwareStoreImpl{Store: s}
}

// CreateSoftware 创建新软件
func (s *SoftwareStoreImpl) CreateSoftware(software *model.Software) error {
	// 生成唯一的AKey
	software.AKey = "soft_" + uuid.New().String()
	return s.DB.Create(software).Error
}

// GetSoftwareList 获取软件列表（分页）
func (s *SoftwareStoreImpl) GetSoftwareList(page, size int) ([]*model.Software, int64, error) {
	var softwares []*model.Software
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	s.DB.Model(&model.Software{}).Count(&total)

	// 查询列表
	err := s.DB.Limit(size).Offset(offset).Find(&softwares).Error
	if err != nil {
		return nil, 0, err
	}

	// 为每个软件获取版本数量
	for i := range softwares {
		var versionCount int64
		s.DB.Model(&model.Version{}).Where("a_key = ?", softwares[i].AKey).Count(&versionCount)
		softwares[i].VersionCount = versionCount
	}

	return softwares, total, nil
}

// GetSoftwareByAKey 根据AKey获取软件信息
func (s *SoftwareStoreImpl) GetSoftwareByAKey(akey string) (*model.Software, error) {
	var software model.Software
	err := s.DB.Where("a_key = ?", akey).First(&software).Error
	if err != nil {
		return nil, err
	}

	// 获取版本数量
	var versionCount int64
	s.DB.Model(&model.Version{}).Where("a_key = ?", akey).Count(&versionCount)
	software.VersionCount = versionCount

	return &software, nil
}

// UpdateSoftware 更新软件信息
func (s *SoftwareStoreImpl) UpdateSoftware(software *model.Software) error {
	return s.DB.Model(&model.Software{}).Where("a_key = ?", software.AKey).Updates(software).Error
}

// DeleteSoftware 删除软件（同时删除关联的版本）
func (s *SoftwareStoreImpl) DeleteSoftware(akey string) error {
	// 开始事务
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除关联的版本
	if err := tx.Where("a_key = ?", akey).Delete(&model.Version{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除软件
	if err := tx.Where("a_key = ?", akey).Delete(&model.Software{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}