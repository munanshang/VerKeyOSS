package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"verkeyoss/internal/model"
	"verkeyoss/internal/service"
)

// CheckHandler 校验API处理器

type CheckHandler struct {
	service *service.CheckService
}

// NewCheckHandler 创建校验API处理器
func NewCheckHandler(service *service.CheckService) *CheckHandler {
	return &CheckHandler{service: service}
}

// CheckLegality 校验AKey和VKey合法性接口
func (h *CheckHandler) CheckLegality(c *gin.Context) {
	// 绑定请求体
	var checkRequest model.CheckRequest

	if err := c.ShouldBindJSON(&checkRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层进行校验
	result, err := h.service.CheckLegality(checkRequest.AKey, checkRequest.VKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "校验失败"))
		return
	}

	// 根据结果返回相应的状态码
	statusCode := http.StatusOK
	if !result.Legal {
		statusCode = http.StatusNotFound
	}

	// 返回响应
	c.JSON(statusCode, SuccessResponse(map[string]interface{}{
		"legal":   result.Legal,
		"message": result.Message,
	}))
}

// CheckUpdate 检查是否有新版本接口
func (h *CheckHandler) CheckUpdate(c *gin.Context) {
	// 绑定请求体
	var checkRequest model.CheckRequest

	if err := c.ShouldBindJSON(&checkRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层检查更新
	result, err := h.service.CheckUpdate(checkRequest.AKey, checkRequest.VKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "检查更新失败"))
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, SuccessResponse(result))
}