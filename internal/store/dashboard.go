package store

import (
	"verkeyoss/internal/model"
)

// DashboardStoreImpl 仪表盘存储实现
type DashboardStoreImpl struct {
	*Store
}

// NewDashboardStore 创建仪表盘存储实例
func (s *Store) NewDashboardStore() *DashboardStoreImpl {
	return &DashboardStoreImpl{Store: s}
}

// GetTotalApps 获取总应用数
// 返回系统中创建的所有应用总数
func (s *DashboardStoreImpl) GetTotalApps() (int64, error) {
	var total int64
	// 计算App表中的记录总数
	err := s.DB.Model(&model.App{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetTotalVersions 获取总版本数
// 返回系统中创建的所有版本总数
func (s *DashboardStoreImpl) GetTotalVersions() (int64, error) {
	var total int64
	// 计算Version表中的记录总数
	err := s.DB.Model(&model.Version{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetRecentApps 获取最近创建的应用（最多5个）
func (s *DashboardStoreImpl) GetRecentApps(limit int) ([]*model.App, error) {
	var apps []*model.App
	err := s.DB.Model(&model.App{}).
		Order("created_at DESC").
		Limit(limit).
		Find(&apps).Error
	if err != nil {
		return nil, err
	}

	// 为每个应用计算版本数量
	for _, app := range apps {
		var versionCount int64
		err := s.DB.Model(&model.Version{}).Where("a_key = ?", app.AKey).Count(&versionCount).Error
		if err != nil {
			return nil, err
		}
		app.VersionCount = versionCount
	}

	return apps, nil
}

// GetRecentVersions 获取最近创建的版本（最多5个）
func (s *DashboardStoreImpl) GetRecentVersions(limit int) ([]*model.Version, error) {
	var versions []*model.Version
	err := s.DB.Model(&model.Version{}).
		Order("created_at DESC").
		Limit(limit).
		Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}
