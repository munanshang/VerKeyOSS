package api

import (
	"net/http"
	"strconv"

	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// AppHandler 应用API处理器

type AppHandler struct {
	appService *service.AppService
}

// NewAppHandler 创建应用API处理器
func NewAppHandler(appService *service.AppService) *AppHandler {
	return &AppHandler{appService: appService}
}

// CreateApp 创建应用接口
func (h *AppHandler) CreateApp(c *gin.Context) {
	// 绑定请求体
	var softwareRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsPaid      bool   `json:"is_paid"`
	}

	var bindErr error
	if bindErr = c.ShouldBindJSON(&softwareRequest); bindErr != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层创建应用
	app, err := h.appService.CreateApp(1, softwareRequest.Name, softwareRequest.Description, softwareRequest.IsPaid) // 使用固定的管理员ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "创建应用失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"akey":        app.AKey,
		"name":        app.Name,
		"description": app.Description,
		"is_paid":     app.IsPaid,
		"created_at":  app.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}))
}

// GetAppList 获取应用列表接口
func (h *AppHandler) GetAppList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用服务层获取应用列表
	apps, total, err := h.appService.GetAppList(page, size)
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
			"is_paid":       app.IsPaid,
			"version_count": app.VersionCount,
			"created_at":    app.CreatedAt.Format("2006-01-02T15:04:05Z"),
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

	// 绑定请求体
	var updateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPaid      bool   `json:"is_paid"`
	}

	var bindErr error
	if bindErr = c.ShouldBindJSON(&updateRequest); bindErr != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层更新应用
	err := h.appService.UpdateApp(akey, updateRequest.Name, updateRequest.Description, updateRequest.IsPaid)
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

	// 调用服务层删除应用
	err := h.appService.DeleteApp(akey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "删除应用失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"message": "删除成功",
	}))
}
