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

// AdminStore 管理员存储接口

type AdminStore interface {
	GetAdminByUsername(username string) (*model.Admin, error)
	UpdateAdmin(admin *model.Admin) error
}

// SoftwareStore 软件存储接口

type SoftwareStore interface {
	CreateSoftware(software *model.Software) error
	GetSoftwareList(page, size int) ([]*model.Software, int64, error)
	GetSoftwareByAKey(akey string) (*model.Software, error)
	UpdateSoftware(software *model.Software) error
	DeleteSoftware(akey string) error
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