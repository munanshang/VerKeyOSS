package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"verkeyoss/internal/service"
)

// SoftwareHandler 软件API处理器

type SoftwareHandler struct {
	service *service.SoftwareService
}

// NewSoftwareHandler 创建软件API处理器
func NewSoftwareHandler(service *service.SoftwareService) *SoftwareHandler {
	return &SoftwareHandler{service: service}
}

// CreateSoftware 创建软件接口
func (h *SoftwareHandler) CreateSoftware(c *gin.Context) {
	// 绑定请求体
	var softwareRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&softwareRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层创建软件
	software, err := h.service.CreateSoftware(softwareRequest.Name, softwareRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "创建软件失败"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"akey":        software.AKey,
		"name":        software.Name,
		"description": software.Description,
		"created_at":  software.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}))
}

// GetSoftwareList 获取软件列表接口
func (h *SoftwareHandler) GetSoftwareList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用服务层获取软件列表
	softwares, total, err := h.service.GetSoftwareList(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取软件列表失败"))
		return
	}

	// 格式化返回数据
	var resultList []map[string]interface{}
	for _, software := range softwares {
		resultList = append(resultList, map[string]interface{}{
			"akey":          software.AKey,
			"name":          software.Name,
			"created_at":    software.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"version_count": software.VersionCount,
		})
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"total": total,
		"list":  resultList,
	}))
}

// GetSoftwareInfo 获取软件信息接口
func (h *SoftwareHandler) GetSoftwareInfo(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 调用服务层获取软件信息
	software, err := h.service.GetSoftwareInfo(akey)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(404, "AKey不存在"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"akey":          software.AKey,
		"name":          software.Name,
		"description":   software.Description,
		"created_at":    software.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"version_count": software.VersionCount,
	}))
}

// UpdateSoftware 更新软件信息接口
func (h *SoftwareHandler) UpdateSoftware(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 绑定请求体
	var updateRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层更新软件
	err := h.service.UpdateSoftware(akey, updateRequest.Name, updateRequest.Description)
	if err != nil {
		if err.Error() == "软件不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse(404, "AKey不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse(500, "更新软件失败"))
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "更新成功",
	})
}

// DeleteSoftware 删除软件接口
func (h *SoftwareHandler) DeleteSoftware(c *gin.Context) {
	// 获取AKey
	akey := c.Param("akey")

	// 调用服务层删除软件
	err := h.service.DeleteSoftware(akey)
	if err != nil {
		if err.Error() == "软件不存在" {
			c.JSON(http.StatusNotFound, ErrorResponse(404, "AKey不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse(500, "删除软件失败"))
		}
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "删除成功",
	})
}