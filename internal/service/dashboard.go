package service

import (
	"verkeyoss/internal/store"
)

// DashboardService 仪表盘服务
// 提供获取系统统计信息的功能
// 支持普通用户获取自己的应用统计，管理员获取系统全局统计

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
// 参数 userID 为当前用户ID
// 参数 isAdmin 表示当前用户是否为管理员
// 返回仪表盘数据，包括用户应用数、系统用户总数（仅管理员可见）、系统应用总数（仅管理员可见）
func (s *DashboardService) GetDashboardData(userID uint, isAdmin bool) (map[string]interface{}, error) {
	// 创建返回结果
	dashboardData := make(map[string]interface{})

	// 获取用户的应用数量
	userAppCount, err := s.store.GetUserAppCount(userID)
	if err != nil {
		return nil, err
	}
	dashboardData["user_app_count"] = userAppCount

	// 如果是管理员，获取系统级别的统计数据
	if isAdmin {
		// 获取系统总用户数
		totalUsers, err := s.store.GetTotalUsers()
		if err != nil {
			return nil, err
		}
		dashboardData["total_users"] = totalUsers

		// 获取系统总应用数
		totalApps, err := s.store.GetTotalApps()
		if err != nil {
			return nil, err
		}
		dashboardData["total_apps"] = totalApps
	}

	return dashboardData, nil
}
