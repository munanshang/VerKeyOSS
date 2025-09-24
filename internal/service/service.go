package service

import (
	"verkeyoss/internal/store"
)

// Services 服务层结构体

type Services struct {
	AuthService         *AuthService
	AppService          *AppService
	VersionService      *VersionService
	CheckService        *CheckService
	DashboardService    *DashboardService
	AnnouncementService *AnnouncementService
}

// NewServices 创建新的服务层实例
func NewServices(store *store.Store, jwtSecret string, expireHours int) *Services {
	// 创建认证服务（替代用户服务）
	authService := NewAuthService(jwtSecret, expireHours)
	appService := NewAppService(store.NewAppStore())
	versionService := NewVersionService(store.NewVersionStore())
	checkService := NewCheckService(store.NewVersionStore(), store.NewAppStore())
	dashboardService := NewDashboardService(store.NewDashboardStore())
	announcementService := NewAnnouncementService(store.NewAnnouncementStore())

	return &Services{
		AuthService:         authService,
		AppService:          appService,
		VersionService:      versionService,
		CheckService:        checkService,
		DashboardService:    dashboardService,
		AnnouncementService: announcementService,
	}
}
