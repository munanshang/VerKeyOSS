package api

import (
	"net/http"
	"strconv"

	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// VersionHandler 版本API处理器

type VersionHandler struct {
	service *service.VersionService
}

// NewVersionHandler 创建版本API处理器
func NewVersionHandler(service *service.VersionService) *VersionHandler {
	return &VersionHandler{service: service}
}

// CreateVersion 创建新版本接口
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 绑定请求体
	var versionRequest struct {
		Version        string `json:"version" binding:"required"`
		Description    string `json:"description"`
		IsLatest       bool   `json:"is_latest"`
		IsForcedUpdate bool   `json:"is_forced_update"`
	}

	if err := c.ShouldBindJSON(&versionRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层创建版本
	version, err := h.service.CreateVersion(akey, versionRequest.Version, versionRequest.Description, versionRequest.IsLatest, versionRequest.IsForcedUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "创建版本失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"vkey":             version.VKey,
		"akey":             version.AKey,
		"version":          version.Version,
		"description":      version.Description,
		"is_latest":        version.IsLatest,
		"is_forced_update": version.IsForcedUpdate,
		"created_at":       version.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}))
}

// GetVersionList 获取版本列表接口
func (h *VersionHandler) GetVersionList(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用服务层获取版本列表
	versions, total, err := h.service.GetVersionList(akey, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取版本列表失败"))
		return
	}

	// 格式化返回数据，不包含akey字段，但包含版本描述
	var resultList []map[string]interface{}
	for _, version := range versions {
		resultList = append(resultList, map[string]interface{}{
			"vkey":             version.VKey,
			"version":          version.Version,
			"description":      version.Description,
			"is_latest":        version.IsLatest,
			"is_forced_update": version.IsForcedUpdate,
			"created_at":       version.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"total": total,
		"list":  resultList,
	}))
}

// UpdateVersion 更新版本信息接口
func (h *VersionHandler) UpdateVersion(c *gin.Context) {
	// 获取VKey
	vkey := c.Param("vkey")

	// 绑定请求体
	var updateRequest struct {
		Version        string `json:"version"`
		Description    string `json:"description"`
		IsLatest       bool   `json:"is_latest"`
		IsForcedUpdate bool   `json:"is_forced_update"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层更新版本
	err := h.service.UpdateVersion(vkey, updateRequest.Version, updateRequest.Description, updateRequest.IsLatest, updateRequest.IsForcedUpdate)
	if err != nil {
		if err.Error() == "版本不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse(404, "VKey不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse(500, "更新版本失败"))
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteVersion 删除版本接口
func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	// 获取VKey
	vkey := c.Param("vkey")

	// 调用服务层删除版本
	err := h.service.DeleteVersion(vkey)
	if err != nil {
		if err.Error() == "版本不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse(404, "VKey不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse(500, "删除版本失败"))
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
