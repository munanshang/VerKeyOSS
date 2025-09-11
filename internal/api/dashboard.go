package api

import (
	"net/http"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 仪表盘处理器
// 处理仪表盘相关的HTTP请求

type DashboardHandler struct {
	dashboardService    *service.DashboardService
	announcementService *service.AnnouncementService
}

// NewDashboardHandler 创建仪表盘处理器实例
// 参数 dashboardService 为仪表盘服务实例
// 参数 announcementService 为公告服务实例
func NewDashboardHandler(dashboardService *service.DashboardService, announcementService *service.AnnouncementService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService:    dashboardService,
		announcementService: announcementService,
	}
}

// GetDashboardData 获取仪表盘数据
// 路由: GET /api/dashboard/stats
// 需要认证
func (h *DashboardHandler) GetDashboardData(c *gin.Context) {
	// 调用服务层获取仪表盘数据
	dashboardData, err := h.dashboardService.GetDashboardData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取仪表盘数据失败", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": dashboardData,
	})
}

// GetAnnouncements 获取公告列表
// 路由: GET /api/dashboard/announcements
// 需要认证
func (h *DashboardHandler) GetAnnouncements(c *gin.Context) {
	// 调用服务层获取激活的公告列表
	announcements, err := h.announcementService.GetActiveAnnouncements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取公告列表失败", "error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": announcements,
	})
}
