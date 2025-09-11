package router

import (
	"verkeyoss/internal/api"
	"verkeyoss/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置路由
// 参数 services 提供所有需要的服务实例
// 返回配置好的 gin.Engine 实例
func SetupRouter(services *service.Services) *gin.Engine {
	// 创建gin引擎实例
	r := gin.Default()

	// 配置跨域设置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 配置前端资源服务
	r.LoadHTMLFiles("web/dist/index.html")
	r.Static("/assets", "./web/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 调试页面路由 - 独立于API路由组
	r.GET("/debug", func(c *gin.Context) {
		c.File("frontend/debug.html")
	})

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
	}

	// 仪表盘接口
	dashboardGroup := apiGroup.Group("/dashboard")
	dashboardHandler := api.NewDashboardHandler(services.DashboardService, services.AnnouncementService)
	{
		dashboardGroup.Use(api.AuthMiddleware(services.AuthService))
		dashboardGroup.GET("/stats", dashboardHandler.GetDashboardData)
		dashboardGroup.GET("/announcements", dashboardHandler.GetAnnouncements)
	}

	return r
}
