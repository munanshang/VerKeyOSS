package api

import (
	"net/http"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员功能处理器
// 负责处理所有需要管理员权限的操作

type AdminHandler struct {
	userService  *service.UserService
	appService   *service.AppService
	versionService *service.VersionService
}

// NewAdminHandler 创建管理员功能处理器实例
// 参数 userService 用户服务实例
// 参数 appService 应用服务实例
// 参数 versionService 版本服务实例
func NewAdminHandler(userService *service.UserService, appService *service.AppService, versionService *service.VersionService) *AdminHandler {
	return &AdminHandler{
		userService:  userService,
		appService:   appService,
		versionService: versionService,
	}
}

// GetAllUsers 获取所有用户列表
// 路由: GET /api/admin/users
// 需要管理员权限
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	// 获取查询参数
	page := 1
	size := 10
	if c.Query("page") != "" {
		c.ShouldBindQuery(&page)
	}
	if c.Query("size") != "" {
		c.ShouldBindQuery(&size)
	}

	// 调用服务层获取所有用户列表
	users, total, err := h.userService.GetAllUsers(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取用户列表失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"total": total,
		"list":  users,
		"page":  page,
		"size":  size,
	}))
}

// BanUser 封禁用户
// 路由: PUT /api/admin/users/:id/ban
// 需要管理员权限
func (h *AdminHandler) BanUser(c *gin.Context) {
	// 获取用户ID
	userId := c.Param("id")

	// 绑定请求体
	var banRequest struct {
		Banned     bool   `json:"banned" binding:"required"`
		BanReason  string `json:"ban_reason"`
	}

	if err := c.ShouldBindJSON(&banRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层封禁/解封用户
	err := h.userService.BanUser(userId, banRequest.Banned, banRequest.BanReason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse("操作成功"))
}

// GetAllApps 管理员获取所有应用列表
// 路由: GET /api/admin/apps
// 需要管理员权限
func (h *AdminHandler) GetAllApps(c *gin.Context) {
	// 获取查询参数
	page := 1
	size := 10
	if c.Query("page") != "" {
		c.ShouldBindQuery(&page)
	}
	if c.Query("size") != "" {
		c.ShouldBindQuery(&size)
	}

	// 调用服务层获取所有应用列表
	apps, total, err := h.appService.GetAllApps(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取应用列表失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"total": total,
		"list":  apps,
		"page":  page,
		"size":  size,
	}))
}

// BanApp 管理员封禁应用
// 路由: PUT /api/admin/apps/:akey/ban
// 需要管理员权限
func (h *AdminHandler) BanApp(c *gin.Context) {
	// 获取应用AKey
	akey := c.Param("akey")

	// 绑定请求体
	var banRequest struct {
		Banned     bool   `json:"banned" binding:"required"`
		BanReason  string `json:"ban_reason"`
	}

	if err := c.ShouldBindJSON(&banRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层封禁/解封应用
	err := h.appService.BanApp(akey, banRequest.Banned, banRequest.BanReason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse("操作成功"))
}