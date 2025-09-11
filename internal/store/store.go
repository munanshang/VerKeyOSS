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

// AppStore 应用存储接口
type AppStore interface {
	CreateApp(app *model.App) error
	GetAppList(page, size int) ([]*model.App, int64, error)
	GetAppListByUserID(userID uint, page, size int) ([]*model.App, int64, error)
	GetAppByAKey(akey string) (*model.App, error)
	UpdateApp(app *model.App) error
	DeleteApp(akey string) error
}

// VersionStore 版本存储接口
type VersionStore interface {
	CreateVersion(version *model.Version) error
	GetVersionListByAKey(akey string, page, size int) ([]*model.Version, int64, error)
	GetVersionByVKey(vkey string) (*model.Version, error)
	UpdateVersion(version *model.Version) error
	DeleteVersion(vkey string) error
	GetLatestVersionByAKey(akey string) (*model.Version, error)
	Validate(akey, vkey string) (bool, error)
	IsVersionLatest(akey, vkey string) (bool, *model.Version, error)
}

// DashboardStore 仪表盘存储接口
type DashboardStore interface {
	// 获取总应用数
	GetTotalApps() (int64, error)
	// 获取总版本数
	GetTotalVersions() (int64, error)
}

// AnnouncementStore 公告存储接口
type AnnouncementStore interface {
	// 获取激活的公告列表
	GetActiveAnnouncements() ([]*model.Announcement, error)
}
