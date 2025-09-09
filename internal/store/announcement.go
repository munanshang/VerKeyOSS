package store

import (
	"verkeyoss/internal/model"
)

// AnnouncementStoreImpl 公告存储实现
type AnnouncementStoreImpl struct {
	*Store
}

// NewAnnouncementStore 创建公告存储实例
func (s *Store) NewAnnouncementStore() *AnnouncementStoreImpl {
	return &AnnouncementStoreImpl{Store: s}
}

// GetActiveAnnouncements 获取激活的公告列表
// 返回所有状态为激活的公告，按发布日期降序排序
func (s *AnnouncementStoreImpl) GetActiveAnnouncements() ([]*model.Announcement, error) {
	var announcements []*model.Announcement
	// 查询所有激活状态的公告，并按发布日期降序排序
	err := s.DB.Where("is_active = ?", true).Order("publish_date DESC").Find(&announcements).Error
	if err != nil {
		return nil, err
	}
	return announcements, nil
}
