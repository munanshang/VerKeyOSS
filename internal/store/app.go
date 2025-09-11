package store

import (
	"verkeyoss/internal/model"

	"github.com/google/uuid"
)

// AppStoreImpl 应用存储实现
type AppStoreImpl struct {
	*Store
}

// NewAppStore 创建应用存储实例
func (s *Store) NewAppStore() *AppStoreImpl {
	return &AppStoreImpl{Store: s}
}

// CreateApp 创建新应用
func (s *AppStoreImpl) CreateApp(app *model.App) error {
	// 生成唯一的AKey
	app.AKey = "app_" + uuid.New().String()
	return s.DB.Create(app).Error
}

// GetAppList 获取应用列表（分页）
func (s *AppStoreImpl) GetAppList(page, size int) ([]*model.App, int64, error) {
	var apps []*model.App
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	s.DB.Model(&model.App{}).Count(&total)

	// 查询列表
	err := s.DB.Limit(size).Offset(offset).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	// 为每个应用获取版本数量
	for i := range apps {
		var versionCount int64
		s.DB.Model(&model.Version{}).Where("a_key = ?", apps[i].AKey).Count(&versionCount)
		apps[i].VersionCount = versionCount
	}

	return apps, total, nil
}

// GetAppByAKey 根据AKey获取应用信息
func (s *AppStoreImpl) GetAppByAKey(akey string) (*model.App, error) {
	var app model.App
	err := s.DB.Where("a_key = ?", akey).First(&app).Error
	if err != nil {
		return nil, err
	}

	// 获取版本数量
	var versionCount int64
	s.DB.Model(&model.Version{}).Where("a_key = ?", akey).Count(&versionCount)
	app.VersionCount = versionCount

	return &app, nil
}

// GetAppListByUserID 根据用户ID获取应用列表（分页）
func (s *AppStoreImpl) GetAppListByUserID(userID uint, page, size int) ([]*model.App, int64, error) {
	var apps []*model.App
	var total int64

	// 计算偏移量
	offset := (page - 1) * size

	// 查询总数
	s.DB.Model(&model.App{}).Where("user_id = ?", userID).Count(&total)

	// 查询列表
	err := s.DB.Where("user_id = ?", userID).Limit(size).Offset(offset).Find(&apps).Error
	if err != nil {
		return nil, 0, err
	}

	// 为每个应用获取版本数量
	for i := range apps {
		var versionCount int64
		s.DB.Model(&model.Version{}).Where("a_key = ?", apps[i].AKey).Count(&versionCount)
		apps[i].VersionCount = versionCount
	}

	return apps, total, nil
}

// UpdateApp 更新应用信息
func (s *AppStoreImpl) UpdateApp(app *model.App) error {
	return s.DB.Model(&model.App{}).Where("a_key = ?", app.AKey).Updates(app).Error
}

// DeleteApp 删除应用（同时删除关联的版本）
func (s *AppStoreImpl) DeleteApp(akey string) error {
	// 开始事务
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 删除关联的版本
	if err := tx.Where("a_key = ?", akey).Delete(&model.Version{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除应用
	if err := tx.Where("a_key = ?", akey).Delete(&model.App{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()
	return nil
}
