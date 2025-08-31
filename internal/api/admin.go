package api

import (
	"net/http"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员API处理器

type AdminHandler struct {
	service *service.AdminService
}

// NewAdminHandler 创建管理员API处理器
func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// Login 管理员登录接口
func (h *AdminHandler) Login(c *gin.Context) {
	// 绑定请求体
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层处理登录
	token, expiresAt, err := h.service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "用户名或密码错误"))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt,
	}))
}

// ChangePassword 修改管理员密码接口
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	// 获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 绑定请求体
	var passwordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&passwordRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 调用服务层修改密码
	err := h.service.ChangePassword(username.(string), passwordRequest.OldPassword, passwordRequest.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
	})
}
