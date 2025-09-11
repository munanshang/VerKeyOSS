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

// Validate 校验AKey和VKey的合法性
func (s *CheckService) Validate(akey, vkey string) (*model.ValidationResponse, error) {
	// 校验AKey和VKey是否存在对应关系
	legal, err := s.versionStore.Validate(akey, vkey)
	if err != nil {
		return &model.ValidationResponse{
			Valid:   false,
			Message: "校验失败",
		}, err
	}

	if legal {
		return &model.ValidationResponse{
			Valid:   true,
			Message: "校验成功",
		}, nil
	}

	return &model.ValidationResponse{
		Valid:   false,
		Message: "校验失败",
	}, nil
}

// CheckUpdate 检查是否有新版本
func (s *CheckService) CheckUpdate(akey, vkey string) (map[string]interface{}, error) {
	// 首先检查AKey和VKey的合法性
	legal, err := s.versionStore.Validate(akey, vkey)
	if err != nil || !legal {
		return map[string]interface{}{
			"has_update": false,
			"message":    "校验失败",
		}, nil
	}

	// 检查当前版本是否是最新版本
	isLatest, latestVersion, err := s.versionStore.IsVersionLatest(akey, vkey)
	if err != nil {
		return map[string]interface{}{
			"has_update": false,
			"message":    "校验失败",
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
