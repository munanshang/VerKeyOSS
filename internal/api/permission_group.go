package api

import (
	"net/http"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// PermissionGroupHandler 权限组API处理器
type PermissionGroupHandler struct {
	permissionGroupService *service.PermissionGroupService
}

// NewPermissionGroupHandler 创建权限组API处理器实例
func NewPermissionGroupHandler(permissionGroupService *service.PermissionGroupService) *PermissionGroupHandler {
	return &PermissionGroupHandler{
		permissionGroupService: permissionGroupService,
	}
}

// GetAllPermissionGroups 获取所有权限组信息
// 路由: GET /api/permission-groups
// 需要登录认证
func (h *PermissionGroupHandler) GetAllPermissionGroups(c *gin.Context) {
	// 调用服务层获取所有权限组信息
	groups, err := h.permissionGroupService.GetAllPermissionGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取权限组信息失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(groups))
}