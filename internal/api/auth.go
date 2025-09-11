package api

import (
	"net/http"
	"verkeyoss/internal/config"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证API处理器

type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler 创建认证API处理器
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login 管理员登录接口
func (h *AuthHandler) Login(c *gin.Context) {
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

	// 返回成功响应，包含令牌和过期时间
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt,
		"user_info": map[string]interface{}{
			"username": loginRequest.Username,
		},
	}))
}

// ChangePassword 修改管理员密码接口
func (h *AuthHandler) ChangePassword(c *gin.Context) {
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
	err := h.service.ChangePassword(passwordRequest.OldPassword, passwordRequest.NewPassword)
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

// GetUserInfoByToken 通过token获取管理员信息接口
func (h *AuthHandler) GetUserInfoByToken(c *gin.Context) {
	// 获取实际的管理员配置
	adminConfig, err := config.GetAdminConfig()
	if err != nil {
		// 发生错误时返回默认管理员信息
		c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
			"username": "verkeyoss",
		}))
		return
	}
	// 返回实际的管理员用户名
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"username": adminConfig.Username,
	}))
}
