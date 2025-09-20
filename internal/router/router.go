package router

import (
	"log"
	"net/http"

	"verkeyoss/internal/api"
	"verkeyoss/internal/config"
	"verkeyoss/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置路由
// 参数 services 提供所有需要的服务实例
// 参数 appConfig 提供应用配置信息
// 参数 version 应用版本号
// 参数 staticHandler 和 frontendHandler 用于处理前端文件
// 返回配置好的 gin.Engine 实例
func SetupRouter(services *service.Services, appConfig *config.Config, version string, staticHandler, frontendHandler gin.HandlerFunc) *gin.Engine {
	// 根据配置设置gin模式
	if appConfig.Server.Debug {
		gin.SetMode(gin.DebugMode)
		log.Println("当前为调试模式")
		log.Println("⚠️  注意：调试模式下已允许所有域名访问，正式环境请关闭 debug 模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Println("当前为正式模式")
	}

	// 创建gin引擎实例
	r := gin.New()

	// 添加中间件
	r.Use(gin.Recovery()) // 恢复中间件，处理panic
	if appConfig.Server.Debug {
		r.Use(gin.Logger()) // 调试模式下添加日志中间件
		// 调试模式下启用宽松的CORS，允许任何域名访问
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  true, // 允许所有来源
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
			AllowHeaders:     []string{"*"}, // 允许所有请求头
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: false, // 当AllowAllOrigins为true时，必须设置为false
		}))
		log.Println("CORS已启用：允许所有域名访问")
	}

	// API路由组
	apiGroup := r.Group("/api")

	// 认证接口
	authGroup := apiGroup.Group("/auth")
	{
		authHandler := api.NewAuthHandler(services.AuthService)
		authGroup.POST("/login", authHandler.Login)
		// 修改密码需要认证
		authGroup.PUT("/password", api.AuthMiddleware(services.AuthService), authHandler.ChangePassword)
		// 通过token获取用户信息接口
		authGroup.GET("/user-info", api.AuthMiddleware(services.AuthService), authHandler.GetUserInfoByToken)
	}

	// 应用管理接口
	appGroup := apiGroup.Group("/app")
	appHandler := api.NewAppHandler(services.AppService)
	{
		appGroup.Use(api.AuthMiddleware(services.AuthService))
		appGroup.POST("", appHandler.CreateApp)
		appGroup.GET("", appHandler.GetAppList)
		appGroup.PUT("/:akey", appHandler.UpdateApp)
		appGroup.DELETE("/:akey", appHandler.DeleteApp)

		// 版本管理接口
		versionGroup := appGroup.Group("/:akey/versions")
		versionHandler := api.NewVersionHandler(services.VersionService)
		{
			versionGroup.POST("", versionHandler.CreateVersion)
			versionGroup.GET("", versionHandler.GetVersionList)
		}
	}

	// 版本详情接口
	versionDetailGroup := apiGroup.Group("/versions")
	versionDetailHandler := api.NewVersionHandler(services.VersionService)
	{
		versionDetailGroup.Use(api.AuthMiddleware(services.AuthService))
		versionDetailGroup.PUT("/:vkey", versionDetailHandler.UpdateVersion)
		versionDetailGroup.DELETE("/:vkey", versionDetailHandler.DeleteVersion)
	}

	// 校验接口
	checkGroup := apiGroup.Group("/check")
	checkHandler := api.NewCheckHandler(services.CheckService)
	{
		checkGroup.POST("/validate", checkHandler.Validate)
		checkGroup.POST("/update", checkHandler.CheckUpdate)
		// 健康检查接口（不需要认证）
		checkGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "OK",
				"data": gin.H{
					"status":  "healthy",
					"service": "VerKeyOSS",
					"version": version,
				},
			})
		})
	}

	// 仪表盘接口
	dashboardGroup := apiGroup.Group("/dashboard")
	dashboardHandler := api.NewDashboardHandler(services.DashboardService, services.AnnouncementService)
	{
		dashboardGroup.Use(api.AuthMiddleware(services.AuthService))
		dashboardGroup.GET("/stats", dashboardHandler.GetDashboardData)
		dashboardGroup.GET("/announcements", dashboardHandler.GetAnnouncements)
	}

	// 静态文件处理（前端资源）
	r.Use(staticHandler)

	// 前端路由处理（SPA应用，必须放在最后）
	r.NoRoute(frontendHandler)

	return r
}
