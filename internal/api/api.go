package api

import (
	"net/http"

	"verkeyoss/internal/errors"
	"verkeyoss/internal/logger"
	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// BaseHandler 基础处理器

type BaseHandler struct {
	// 基础字段和方法
}

// AuthMiddleware 管理员认证中间件
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取令牌
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
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
		valid, err := authService.VerifyToken(token)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "令牌无效或已过期",
			})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}

// AdminMiddleware 管理员权限验证中间件 (在单管理员系统中，等同于AuthMiddleware)
func AdminMiddleware(authService *service.AuthService) gin.HandlerFunc {
	// 直接使用AuthMiddleware，因为系统只有一个管理员
	return AuthMiddleware(authService)
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) gin.H {
	return gin.H{
		"code": 200,
		"data": data,
	}
}

// respondSuccess 返回成功响应
func respondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse(data))
}

// respondError 返回错误响应
func respondError(c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		// 如果是应用错误，使用定义的错误码和消息
		logger.Errorf("API错误: %v", appErr)
		c.JSON(appErr.Code, ErrorResponse(appErr.Code, appErr.Message))
	} else {
		// 其他错误作为内部服务器错误处理
		logger.Errorf("内部错误: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse(500, "服务器内部错误"))
	}
}
