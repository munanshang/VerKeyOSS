package service

import (
	"verkeyoss/internal/errors"
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

// 预定义错误，使用统一的错误处理
var (
	ErrAppNotFound = errors.ErrAppNotFound
	ErrAppExists   = errors.ErrAppExists
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
func (s *AppService) CreateApp(userId uint, name, description string, isPaid bool) (*model.App, error) {
	app := &model.App{
		UserID:      userId,
		Name:        name,
		Description: description,
		IsPaid:      isPaid,
	}

	err := s.store.CreateApp(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// GetAppList 获取应用列表
func (s *AppService) GetAppList(page, size int) ([]*model.App, int64, error) {
	// 分页参数校验
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 在单管理员模式下，直接获取所有应用
	return s.store.GetAppList(page, size)
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
func (s *AppService) UpdateApp(akey, name, description string, isPaid bool) error {
	app, err := s.store.GetAppByAKey(akey)
	if err != nil {
		return ErrAppNotFound
	}

	// 更新应用信息
	app.Name = name
	app.Description = description
	app.IsPaid = isPaid

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
