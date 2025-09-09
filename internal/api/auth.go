package api

import (
	"net/http"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证API处理器

type AuthHandler struct {
	service *service.UserService
}

// NewAuthHandler 创建认证API处理器
func NewAuthHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login 登录接口
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
	token, expiresAt, userId, username, err := h.service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "用户名或密码错误"))
		return
	}

	// 获取完整用户信息以获取权限组数据
	user, err := h.service.GetUserByID(userId)
	if err != nil {
		// 如果获取用户信息失败，返回错误响应，而不是默认信息
		// 避免普通用户获取到管理员才能看到的信息
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "获取用户权限信息失败"))
		return
	}

	// 返回成功响应，包含令牌、过期时间和用户信息
	// 只返回group_id，前端通过单独接口获取权限组详细信息
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"token":      token,
		"expires_at": expiresAt,
		"user_info": map[string]interface{}{
			"user_id":      userId,
			"username":     username,
			"group_id":     user.GroupID,
		},
	}))
}

// ChangePassword 修改密码接口
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// 获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 绑定请求体
	var passwordRequest struct {
		OldPassword      string `json:"old_password" binding:"required"`
		NewPassword      string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&passwordRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(400, "参数错误"))
		return
	}

	// 通过用户ID获取用户信息
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "用户ID格式错误"))
		return
	}

	user, err := h.service.GetUserByID(userIdUint)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "用户不存在"))
		return
	}

	// 暂时跳过二级密码验证（保留字段，作为预留功能）
	// 后续版本可根据需求重新启用此验证

	// 调用服务层修改密码
	err = h.service.ChangePassword(user.Username, passwordRequest.OldPassword, passwordRequest.NewPassword)
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

// GetUserInfoByToken 通过token获取用户信息接口
func (h *AuthHandler) GetUserInfoByToken(c *gin.Context) {
	// 获取用户ID
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "未授权访问"))
		return
	}

	// 通过用户ID获取用户信息
	userIdUint, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "用户ID格式错误"))
		return
	}

	user, err := h.service.GetUserByID(userIdUint)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse(401, "用户不存在"))
		return
	}

	// 返回用户信息
	// 只返回group_id，前端通过单独接口获取权限组详细信息
	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"user_id":      user.ID,
		"username":     user.Username,
		"group_id":     user.GroupID,
	}))
}