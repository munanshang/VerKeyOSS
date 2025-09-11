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
