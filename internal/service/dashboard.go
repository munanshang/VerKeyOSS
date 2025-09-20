package service

import (
	"verkeyoss/internal/store"
)

// DashboardService 仪表盘服务
// 提供获取系统统计信息的功能
// 在单管理员模式下，只提供系统总应用数

type DashboardService struct {
	store store.DashboardStore
}

// NewDashboardService 创建仪表盘服务实例
// 参数 store 为仪表盘存储接口的实现
func NewDashboardService(store store.DashboardStore) *DashboardService {
	return &DashboardService{
		store: store,
	}
}

// GetDashboardData 获取仪表盘数据
// 返回系统总应用数、总版本数、最近应用和最近版本
func (s *DashboardService) GetDashboardData() (map[string]interface{}, error) {
	// 创建返回结果
	dashboardData := make(map[string]interface{})

	// 获取系统总应用数
	totalApps, err := s.store.GetTotalApps()
	if err != nil {
		return nil, err
	}
	dashboardData["total_apps"] = totalApps

	// 获取系统总版本数
	totalVersions, err := s.store.GetTotalVersions()
	if err != nil {
		return nil, err
	}
	dashboardData["total_versions"] = totalVersions

	// 获取最近的应用（最多5个）
	recentApps, err := s.store.GetRecentApps(5)
	if err != nil {
		return nil, err
	}
	dashboardData["recent_apps"] = recentApps

	// 获取最近的版本（最多5个）
	recentVersions, err := s.store.GetRecentVersions(5)
	if err != nil {
		return nil, err
	}
	dashboardData["recent_versions"] = recentVersions

	return dashboardData, nil
}
