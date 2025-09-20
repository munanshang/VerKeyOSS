package api

import (
	"strconv"

	"verkeyoss/internal/errors"
	"verkeyoss/internal/logger"
	"verkeyoss/internal/service"
	"verkeyoss/internal/validator"

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
	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		IsPaid      bool   `json:"is_paid"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorf("创建应用请求参数错误: %v", err)
		respondError(c, errors.NewValidationError("请求参数错误"))
		return
	}

	// 验证输入参数
	if err := validator.ValidateAppName(request.Name); err != nil {
		logger.Errorf("应用名称验证失败: %v", err)
		respondError(c, err)
		return
	}

	if err := validator.ValidateDescription(request.Description); err != nil {
		logger.Errorf("应用描述验证失败: %v", err)
		respondError(c, err)
		return
	}

	// 调用服务层创建应用
	app, err := h.appService.CreateApp(1, request.Name, request.Description, request.IsPaid)
	if err != nil {
		logger.Errorf("创建应用失败: %v", err)
		respondError(c, errors.WrapError(err, "创建应用失败"))
		return
	}

	logger.Infof("成功创建应用: %s (AKey: %s)", app.Name, app.AKey)

	// 返回成功响应
	respondSuccess(c, map[string]interface{}{
		"akey":        app.AKey,
		"name":        app.Name,
		"description": app.Description,
		"is_paid":     app.IsPaid,
		"created_at":  app.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetAppList 获取应用列表接口
func (h *AppHandler) GetAppList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 验证分页参数
	validPage, validSize, err := validator.ValidatePagination(page, size)
	if err != nil {
		logger.Errorf("分页参数验证失败: %v", err)
		respondError(c, err)
		return
	}

	// 调用服务层获取应用列表
	apps, total, err := h.appService.GetAppList(validPage, validSize)
	if err != nil {
		logger.Errorf("获取应用列表失败: %v", err)
		respondError(c, errors.WrapError(err, "获取应用列表失败"))
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

	logger.Infof("成功获取应用列表: 共%d条记录", total)

	// 返回成功响应
	respondSuccess(c, map[string]interface{}{
		"list":  appList,
		"total": total,
		"page":  validPage,
		"size":  validSize,
	})
}

// UpdateApp 更新应用信息接口
func (h *AppHandler) UpdateApp(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")
	if err := validator.ValidateAKey(akey); err != nil {
		logger.Errorf("AKey验证失败: %v", err)
		respondError(c, err)
		return
	}

	// 绑定请求体
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPaid      bool   `json:"is_paid"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorf("更新应用请求参数错误: %v", err)
		respondError(c, errors.NewValidationError("请求参数错误"))
		return
	}

	// 验证输入参数
	if err := validator.ValidateAppName(request.Name); err != nil {
		logger.Errorf("应用名称验证失败: %v", err)
		respondError(c, err)
		return
	}

	if err := validator.ValidateDescription(request.Description); err != nil {
		logger.Errorf("应用描述验证失败: %v", err)
		respondError(c, err)
		return
	}

	// 调用服务层更新应用
	err := h.appService.UpdateApp(akey, request.Name, request.Description, request.IsPaid)
	if err != nil {
		logger.Errorf("更新应用失败 (AKey: %s): %v", akey, err)
		respondError(c, errors.WrapError(err, "更新应用失败"))
		return
	}

	logger.Infof("成功更新应用 (AKey: %s)", akey)

	// 返回成功响应
	respondSuccess(c, map[string]interface{}{
		"message": "更新成功",
	})
}

// DeleteApp 删除应用接口
func (h *AppHandler) DeleteApp(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")
	if err := validator.ValidateAKey(akey); err != nil {
		logger.Errorf("AKey验证失败: %v", err)
		respondError(c, err)
		return
	}

	// 调用服务层删除应用
	err := h.appService.DeleteApp(akey)
	if err != nil {
		logger.Errorf("删除应用失败 (AKey: %s): %v", akey, err)
		respondError(c, errors.WrapError(err, "删除应用失败"))
		return
	}

	logger.Infof("成功删除应用 (AKey: %s)", akey)

	// 返回成功响应
	respondSuccess(c, map[string]interface{}{
		"message": "删除成功",
	})
}
