package service

import (
	"verkeyoss/internal/store"
)

// Services 服务层结构体

type Services struct {
	UserService         *UserService
	AppService          *AppService
	VersionService      *VersionService
	CheckService        *CheckService
	DashboardService    *DashboardService
	AnnouncementService *AnnouncementService
	PermissionGroupService *PermissionGroupService
}

// NewServices 创建新的服务层实例
func NewServices(store *store.Store, jwtSecret string, expireHours int) *Services {
	userService := NewUserService(store.NewUserStore(), jwtSecret, expireHours)
	appService := NewAppService(store.NewAppStore())
	versionService := NewVersionService(store.NewVersionStore())
	checkService := NewCheckService(store.NewVersionStore())
	dashboardService := NewDashboardService(store.NewDashboardStore())
	announcementService := NewAnnouncementService(store.NewAnnouncementStore())
	permissionGroupService := NewPermissionGroupService(store.NewPermissionGroupStore())

	return &Services{
		UserService:         userService,
		AppService:          appService,
		VersionService:      versionService,
		CheckService:        checkService,
		DashboardService:    dashboardService,
		AnnouncementService: announcementService,
		PermissionGroupService: permissionGroupService,
	}
}
