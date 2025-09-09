package api

import (
	"net/http"

	"verkeyoss/internal/service"

	"github.com/gin-gonic/gin"
)

// BaseHandler 基础处理器

type BaseHandler struct {
	// 基础字段和方法
}

// AuthMiddleware 用户认证中间件
func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
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

		// 验证令牌并获取用户信息
		userId, isAdmin, err := userService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "令牌无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户ID和管理员标识存储到上下文中
		c.Set("user_id", userId)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}

// AdminMiddleware 管理员权限验证中间件
func AdminMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先进行用户认证
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

		// 验证令牌并获取用户信息
		userId, isAdmin, err := userService.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "令牌无效或已过期",
			})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "无管理员权限",
			})
			c.Abort()
			return
		}

		// 将用户ID和管理员标识存储到上下文中
		c.Set("user_id", userId)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
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
