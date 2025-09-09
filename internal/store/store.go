package store

import (
	"verkeyoss/internal/model"

	"gorm.io/gorm"
)

// Store 存储层结构体

type Store struct {
	DB *gorm.DB
}

// NewStore 创建新的存储层实例
func NewStore(db *gorm.DB) *Store {
	return &Store{DB: db}
}

// UserStore 用户存储接口

type UserStore interface {
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(user *model.User) error
	GetAllUsers(page, size int) ([]*model.User, int64, error)
	BanUser(id string, banned bool, banReason string) error
}

// AppStore 应用存储接口
type AppStore interface {
	CreateApp(app *model.App) error
	GetAppList(page, size int) ([]*model.App, int64, error)
	GetAppListByUserID(userID uint, page, size int) ([]*model.App, int64, error)
	GetAppByAKey(akey string) (*model.App, error)
	UpdateApp(app *model.App) error
	DeleteApp(akey string) error
	GetAllApps(page, size int) ([]*model.App, int64, error)
	BanApp(akey string, banned bool, banReason string) error
}

// VersionStore 版本存储接口

type VersionStore interface {
	CreateVersion(version *model.Version) error
	GetVersionListByAKey(akey string, page, size int) ([]*model.Version, int64, error)
	GetVersionByVKey(vkey string) (*model.Version, error)
	UpdateVersion(version *model.Version) error
	DeleteVersion(vkey string) error
	GetLatestVersionByAKey(akey string) (*model.Version, error)
	CheckLegality(akey, vkey string) (bool, error)
	IsVersionLatest(akey, vkey string) (bool, *model.Version, error)
}

// DashboardStore 仪表盘存储接口
type DashboardStore interface {
	// 获取总用户数
	GetTotalUsers() (int64, error)
	// 获取总应用数
	GetTotalApps() (int64, error)
	// 获取指定用户的应用数
	GetUserAppCount(userID uint) (int64, error)
}

// AnnouncementStore 公告存储接口
type AnnouncementStore interface {
	// 获取激活的公告列表
	GetActiveAnnouncements() ([]*model.Announcement, error)
}

// PermissionGroupStore 权限组存储接口
type PermissionGroupStore interface {
	// 获取所有权限组
	GetAllPermissionGroups() ([]*model.PermissionGroup, error)
	// 根据ID获取权限组
	GetPermissionGroupByID(id uint) (*model.PermissionGroup, error)
}
