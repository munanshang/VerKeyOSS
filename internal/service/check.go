package service

import (
	"verkeyoss/internal/model"
	"verkeyoss/internal/store"
)

// CheckService 校验服务

type CheckService struct {
	versionStore store.VersionStore
}

// NewCheckService 创建校验服务实例
func NewCheckService(versionStore store.VersionStore) *CheckService {
	return &CheckService{versionStore: versionStore}
}

// CheckLegality 校验AKey和VKey的合法性
func (s *CheckService) CheckLegality(akey, vkey string) (*model.LegalityResponse, error) {
	// 校验AKey和VKey是否存在对应关系
	legal, err := s.versionStore.CheckLegality(akey, vkey)
	if err != nil {
		return &model.LegalityResponse{
			Legal:   false,
			Message: "校验失败",
		}, err
	}

	if legal {
		return &model.LegalityResponse{
			Legal:   true,
			Message: "AKey和VKey合法",
		}, nil
	}

	// 检查是AKey不存在还是VKey不存在
	// 先检查VKey是否存在
	_, err = s.versionStore.GetVersionByVKey(vkey)
	if err == nil {
		// VKey存在但AKey不匹配
		return &model.LegalityResponse{
			Legal:   false,
			Message: "AKey和VKey不匹配",
		}, nil
	}

	// 尝试检查AKey是否存在对应的版本
	versions, _, err := s.versionStore.GetVersionListByAKey(akey, 1, 1)
	if err != nil || len(versions) == 0 {
		return &model.LegalityResponse{
			Legal:   false,
			Message: "AKey不存在",
		}, nil
	}

	// 如果AKey存在但VKey不存在
	return &model.LegalityResponse{
		Legal:   false,
		Message: "VKey不存在",
	}, nil
}

// CheckUpdate 检查是否有新版本
func (s *CheckService) CheckUpdate(akey, vkey string) (map[string]interface{}, error) {
	// 首先检查AKey和VKey的合法性
	legal, err := s.versionStore.CheckLegality(akey, vkey)
	if err != nil || !legal {
		return map[string]interface{}{
			"has_update": false,
			"message":    "AKey或VKey无效",
		}, nil
	}

	// 检查当前版本是否是最新版本
	isLatest, latestVersion, err := s.versionStore.IsVersionLatest(akey, vkey)
	if err != nil {
		return map[string]interface{}{
			"has_update": false,
			"message":    "检查更新失败",
		}, err
	}

	if isLatest {
		return map[string]interface{}{
			"has_update": false,
			"message":    "当前已是最新版本",
		}, nil
	}

	// 存在新版本
	return map[string]interface{}{
		"has_update":     true,
		"latest_version": latestVersion.Version,
		"release_time":   latestVersion.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
