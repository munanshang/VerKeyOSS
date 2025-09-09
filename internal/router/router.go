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

	// API路由组
	apiGroup := r.Group("/api")

	// 认证接口
	authGroup := apiGroup.Group("/auth")
	{
		authHandler := api.NewAuthHandler(services.UserService)
		authGroup.POST("/login", authHandler.Login)
		// 修改密码需要认证
		authGroup.PUT("/password", api.AuthMiddleware(services.UserService), authHandler.ChangePassword)
		// 通过token获取用户信息接口
		authGroup.GET("/user-info", api.AuthMiddleware(services.UserService), authHandler.GetUserInfoByToken)
	}

	// 应用管理接口
	appGroup := apiGroup.Group("/app")
	appHandler := api.NewAppHandler(services.AppService, services.UserService)
	{
		appGroup.Use(api.AuthMiddleware(services.UserService))
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
		versionDetailGroup.Use(api.AuthMiddleware(services.UserService))
		versionDetailGroup.PUT("/:vkey", versionDetailHandler.UpdateVersion)
		versionDetailGroup.DELETE("/:vkey", versionDetailHandler.DeleteVersion)
	}

	// 校验接口
	checkGroup := apiGroup.Group("/check")
	checkHandler := api.NewCheckHandler(services.CheckService)
	{
		checkGroup.POST("/legality", checkHandler.CheckLegality)
		checkGroup.POST("/update", checkHandler.CheckUpdate)
	}

	// 仪表盘接口
	dashboardGroup := apiGroup.Group("/dashboard")
	dashboardHandler := api.NewDashboardHandler(services.DashboardService, services.AnnouncementService)
	{
		dashboardGroup.Use(api.AuthMiddleware(services.UserService))
		dashboardGroup.GET("", dashboardHandler.GetDashboardData)
		dashboardGroup.GET("/announcements", dashboardHandler.GetAnnouncements)
	}

	// 权限组接口
	permissionGroupGroup := apiGroup.Group("/permission-groups")
	permissionGroupHandler := api.NewPermissionGroupHandler(services.PermissionGroupService)
	{
		permissionGroupGroup.Use(api.AuthMiddleware(services.UserService))
		permissionGroupGroup.GET("", permissionGroupHandler.GetAllPermissionGroups) // 获取所有权限组信息
	}

	// 管理员接口
	adminGroup := apiGroup.Group("/admin")
	adminHandler := api.NewAdminHandler(services.UserService, services.AppService, services.VersionService)
	{
		adminGroup.Use(api.AdminMiddleware(services.UserService))
		
		// 用户管理接口
		userGroup := adminGroup.Group("/users")
		{
			userGroup.GET("", adminHandler.GetAllUsers) // 获取所有用户列表
			userGroup.PUT("/:id/ban", adminHandler.BanUser) // 封禁/解封用户
		}
		
		// 应用管理接口
		appGroup := adminGroup.Group("/apps")
		{
			appGroup.GET("", adminHandler.GetAllApps) // 获取所有应用列表
			appGroup.PUT("/:akey/ban", adminHandler.BanApp) // 封禁/解封应用
		}
	}

	return r
}
