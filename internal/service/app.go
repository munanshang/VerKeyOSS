package service

import (
	"errors"
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

var (
	ErrAppNotFound = errors.New("应用不存在")
	ErrAppExists   = errors.New("应用已存在")
)

// AppService 应用服务
type AppService struct {
	store store.AppStore
}

// NewAppService 创建应用服务实例
func NewAppService(store store.AppStore) *AppService {
	return &AppService{store: store}
}

// CreateApp 创建新应用
func (s *AppService) CreateApp(userId uint, name, description string) (*model.App, error) {
	app := &model.App{
		UserID:      userId,
		Name:        name,
		Description: description,
	}

	err := s.store.CreateApp(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// GetAppList 获取应用列表
func (s *AppService) GetAppList(page, size int, userID uint, isAdmin bool) ([]*model.App, int64, error) {
	// 分页参数校验
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 如果是管理员，获取所有应用；否则只获取用户自己的应用
	if isAdmin {
		return s.store.GetAppList(page, size)
	} else {
		return s.store.GetAppListByUserID(userID, page, size)
	}
}

// GetAppByAKey 根据AKey获取应用信息
func (s *AppService) GetAppByAKey(akey string) (*model.App, error) {
	app, err := s.store.GetAppByAKey(akey)
	if err != nil {
		return nil, ErrAppNotFound
	}

	return app, nil
}

// UpdateApp 更新应用信息
func (s *AppService) UpdateApp(akey, name, description string) error {
	app, err := s.store.GetAppByAKey(akey)
	if err != nil {
		return ErrAppNotFound
	}

	// 更新应用信息
	app.Name = name
	app.Description = description

	err = s.store.UpdateApp(app)
	if err != nil {
		return err
	}

	return nil
}

// DeleteApp 删除应用
func (s *AppService) DeleteApp(akey string) error {
	// 检查应用是否存在
	_, err := s.store.GetAppByAKey(akey)
	if err != nil {
		return ErrAppNotFound
	}

	// 删除应用
	err = s.store.DeleteApp(akey)
	if err != nil {
		return err
	}

	return nil
}

// GetAllApps 管理员获取所有应用列表
// 仅管理员可调用
func (s *AppService) GetAllApps(page, size int) ([]*model.App, int64, error) {
	// 分页参数校验
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	return s.store.GetAllApps(page, size)
}

// BanApp 管理员封禁/解封应用
// 仅管理员可调用
func (s *AppService) BanApp(akey string, banned bool, banReason string) error {
	// 检查应用是否存在
	_, err := s.store.GetAppByAKey(akey)
	if err != nil {
		return ErrAppNotFound
	}

	// 封禁/解封应用
	return s.store.BanApp(akey, banned, banReason)
}
