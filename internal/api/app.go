package api

import (
	"net/http"
	"strconv"

	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// AppHandler 应用API处理器

type AppHandler struct {
	appService  *service.AppService
	userService *service.UserService
}

// NewAppHandler 创建应用API处理器
func NewAppHandler(appService *service.AppService, userService *service.UserService) *AppHandler {
	return &AppHandler{appService: appService, userService: userService}
}

// CreateApp 创建应用接口
func (h *AppHandler) CreateApp(c *gin.Context) {
	// 获取上下文中的用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 转换用户ID类型
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "用户ID格式错误"))
		return
	}

	// 绑定请求体
	var softwareRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	var bindErr error
	if bindErr = c.ShouldBindJSON(&softwareRequest); bindErr != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层创建应用
	app, err := h.appService.CreateApp(userIdUint, softwareRequest.Name, softwareRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "创建应用失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"akey":        app.AKey,
		"name":        app.Name,
		"description": app.Description,
		"user_id":     app.UserID,
		"created_at":  app.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}))
}

// GetAppList 获取应用列表接口
func (h *AppHandler) GetAppList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 获取上下文中的用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 转换用户ID类型
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "用户ID格式错误"))
		return
	}

	// 获取上下文中的管理员标识
	isAdmin, ok := c.Get("is_admin")
	if !ok {
		isAdmin = false
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		isAdminBool = false
	}

	// 调用服务层获取应用列表
	apps, total, err := h.appService.GetAppList(page, size, userIdUint, isAdminBool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取应用列表失败"))
		return
	}



	// 构造响应数据
	var appList []map[string]interface{}
	for _, app := range apps {
		appInfo := map[string]interface{}{
			"akey":          app.AKey,
			"name":          app.Name,
			"description":   app.Description,
			"version_count": app.VersionCount,
			"created_at":    app.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// 如果是管理员用户，添加用户ID和用户名
		if isAdminBool {
			appInfo["user_id"] = app.UserID
			// 获取用户名
			user, err := h.userService.GetUserByID(app.UserID)
			if err == nil {
				appInfo["username"] = user.Username
			}
		}

		appList = append(appList, appInfo)
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"list":  appList,
		"total": total,
		"page":  page,
		"size":  size,
	}))
}

// UpdateApp 更新应用信息接口
func (h *AppHandler) UpdateApp(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 获取上下文中的用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 转换用户ID类型
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "用户ID格式错误"))
		return
	}

	// 获取上下文中的管理员标识
	isAdmin, ok := c.Get("is_admin")
	if !ok {
		isAdmin = false
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		isAdminBool = false
	}

	// 获取应用信息
	app, err := h.appService.GetAppByAKey(akey)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(404, "应用不存在"))
		return
	}

	// 检查权限：如果不是管理员，必须是应用的创建者
	if !isAdminBool && app.UserID != userIdUint {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "无权限操作此应用"))
		return
	}

	// 绑定请求体
	var updateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var bindErr error
	if bindErr = c.ShouldBindJSON(&updateRequest); bindErr != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层更新应用
	err = h.appService.UpdateApp(akey, updateRequest.Name, updateRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "更新应用失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"message": "更新成功",
	}))
}

// DeleteApp 删除应用接口
func (h *AppHandler) DeleteApp(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 获取上下文中的用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 转换用户ID类型
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "用户ID格式错误"))
		return
	}

	// 获取上下文中的管理员标识
	isAdmin, ok := c.Get("is_admin")
	if !ok {
		isAdmin = false
	}

	isAdminBool, ok := isAdmin.(bool)
	if !ok {
		isAdminBool = false
	}

	// 获取应用信息
	app, err := h.appService.GetAppByAKey(akey)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(404, "应用不存在"))
		return
	}

	// 检查权限：如果不是管理员，必须是应用的创建者
	if !isAdminBool && app.UserID != userIdUint {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "无权限操作此应用"))
		return
	}

	// 调用服务层删除应用
	err = h.appService.DeleteApp(akey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "删除应用失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"message": "删除成功",
	}))
}
