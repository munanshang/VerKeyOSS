package store

import (
	"verkeyoss/internal/model"

	"github.com/google/uuid"
)

// VersionStoreImpl 版本存储实现

type VersionStoreImpl struct {
	*Store
}

// NewVersionStore 创建版本存储实例
func (s *Store) NewVersionStore() *VersionStoreImpl {
	return &VersionStoreImpl{Store: s}
}

// CreateVersion 创建新版本
func (s *VersionStoreImpl) CreateVersion(version *model.Version) error {
	// 生成唯一的VKey
	version.VKey = "ver_" + uuid.New().String()

	// 开始事务
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 如果标记为最新版本，先将其他版本设为非最新
	if version.IsLatest {
		if err := tx.Model(&model.Version{}).Where("a_key = ?", version.AKey).Update("is_latest", false).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 检查是否应该自动设置为最新版本（如果是该软件的第一个版本）
		var existingVersionCount int64
		tx.Model(&model.Version{}).Where("a_key = ?", version.AKey).Count(&existingVersionCount)
		if existingVersionCount == 0 {
			version.IsLatest = true
		}
	}

	// 创建新版本
	if err := tx.Create(version).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}

// GetVersionListByAKey 根据AKey获取版本列表
func (s *VersionStoreImpl) GetVersionListByAKey(akey string, page, size int) ([]*model.Version, int64, error) {
	var versions []*model.Version
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	s.DB.Model(&model.Version{}).Where("a_key = ?", akey).Count(&total)

	// 查询列表（按创建时间倒序）
	err := s.DB.Where("a_key = ?", akey).Order("created_at DESC").Limit(size).Offset(offset).Find(&versions).Error
	if err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

// GetVersionByVKey 根据VKey获取版本信息
func (s *VersionStoreImpl) GetVersionByVKey(vkey string) (*model.Version, error) {
	var version model.Version
	err := s.DB.Where("v_key = ?", vkey).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// UpdateVersion 更新版本信息
func (s *VersionStoreImpl) UpdateVersion(version *model.Version) error {
	// 开始事务
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 如果要标记为最新版本，先将其他版本设为非最新
	if version.IsLatest {
		if err := tx.Model(&model.Version{}).Where("a_key = ? AND v_key != ?", version.AKey, version.VKey).Update("is_latest", false).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新版本信息
	if err := tx.Model(&model.Version{}).Where("v_key = ?", version.VKey).Updates(version).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}

// DeleteVersion 删除版本
func (s *VersionStoreImpl) DeleteVersion(vkey string) error {
	// 获取版本信息，用于后续检查
	version, err := s.GetVersionByVKey(vkey)
	if err != nil {
		return err
	}

	// 开始事务
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除版本
	if err := tx.Where("v_key = ?", vkey).Delete(&model.Version{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果删除的是最新版本，从剩余版本中选择最新的一个作为最新版本
	if version.IsLatest {
		var latestVersion model.Version
		if err := tx.Where("a_key = ?", version.AKey).Order("created_at DESC").First(&latestVersion).Error; err == nil {
			latestVersion.IsLatest = true
			if err := tx.Save(&latestVersion).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 提交事务
	tx.Commit()
	return nil
}

// GetLatestVersionByAKey 获取指定软件的最新版本
func (s *VersionStoreImpl) GetLatestVersionByAKey(akey string) (*model.Version, error) {
	var version model.Version
	err := s.DB.Where("a_key = ? AND is_latest = ?", akey, true).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// Validate 校验AKey和VKey的合法性
func (s *VersionStoreImpl) Validate(akey, vkey string) (bool, error) {
	var count int64
	err := s.DB.Model(&model.Version{}).Where("a_key = ? AND v_key = ?", akey, vkey).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsVersionLatest 检查指定版本是否是最新版本
func (s *VersionStoreImpl) IsVersionLatest(akey, vkey string) (bool, *model.Version, error) {
	// 获取当前版本信息
	currentVersion, err := s.GetVersionByVKey(vkey)
	if err != nil {
		return false, nil, err
	}

	// 检查AKey是否匹配
	if currentVersion.AKey != akey {
		return false, nil, nil
	}

	// 检查是否是最新版本
	if currentVersion.IsLatest {
		return true, nil, nil
	}

	// 获取最新版本
	latestVersion, err := s.GetLatestVersionByAKey(akey)
	if err != nil {
		return false, nil, err
	}

	return false, latestVersion, nil
}
