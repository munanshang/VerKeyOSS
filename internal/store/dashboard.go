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

// GetTotalUsers 获取总用户数
// 返回系统中注册的所有用户总数
func (s *DashboardStoreImpl) GetTotalUsers() (int64, error) {
	var total int64
	// 计算User表中的记录总数
	err := s.DB.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
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

// GetUserAppCount 获取指定用户的应用数
// 参数 userID 为用户ID
// 返回该用户创建的应用数量
func (s *DashboardStoreImpl) GetUserAppCount(userID uint) (int64, error) {
	var count int64
	// 计算特定用户ID的应用总数
	err := s.DB.Model(&model.App{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
