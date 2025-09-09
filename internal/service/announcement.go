package service

import (
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

// AnnouncementService 公告服务
// 提供获取系统公告的功能

type AnnouncementService struct {
	store store.AnnouncementStore
}

// NewAnnouncementService 创建公告服务实例
// 参数 store 为公告存储接口的实现
func NewAnnouncementService(store store.AnnouncementStore) *AnnouncementService {
	return &AnnouncementService{
		store: store,
	}
}

// GetActiveAnnouncements 获取激活的公告列表
// 返回所有状态为激活的公告，按发布日期降序排序
func (s *AnnouncementService) GetActiveAnnouncements() ([]*model.Announcement, error) {
	// 调用存储层获取激活的公告列表
	return s.store.GetActiveAnnouncements()
}
