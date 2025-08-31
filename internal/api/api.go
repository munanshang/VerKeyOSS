package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"verkeyoss/internal/service"
)

// BaseHandler 基础处理器

type BaseHandler struct {
	// 基础字段和方法
}

// AuthMiddleware 管理员认证中间件
func AuthMiddleware(adminService *service.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取令牌
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"message": "未授权访问",
		})
		c.Abort()
		return
	}

	// 简单处理Bearer令牌格式
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// 验证令牌
	username, err := adminService.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"message": "令牌无效或已过期",
		})
		c.Abort()
		return
	}

	// 将用户名存储到上下文中
	c.Set("username", username)
	c.Next()
	}
}

// 响应结构

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 返回成功响应
func SuccessResponse(data interface{}) Response {
	return Response{
		Code: 200,
		Data: data,
	}
}

// ErrorResponse 返回错误响应
func ErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}