package service

import (
	"verkeyoss/internal/store"
)

// Services 服务层结构体

type Services struct {
	AdminService    *AdminService
	SoftwareService *SoftwareService
	VersionService  *VersionService
	CheckService    *CheckService
}

// NewServices 创建新的服务层实例
func NewServices(store *store.Store, jwtSecret string, expireHours int) *Services {
	adminService := NewAdminService(store.NewAdminStore(), jwtSecret, expireHours)
	softwareService := NewSoftwareService(store.NewSoftwareStore())
	versionService := NewVersionService(store.NewVersionStore())
	checkService := NewCheckService(store.NewVersionStore())

	return &Services{
		AdminService:    adminService,
		SoftwareService: softwareService,
		VersionService:  versionService,
		CheckService:    checkService,
	}
}
