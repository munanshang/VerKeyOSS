package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"verkeyoss/internal/api"
	"verkeyoss/internal/service"

	"verkeyoss/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 配置结构体
type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret      string `yaml:"secret"`
		ExpireHours int    `yaml:"expire_hours"`
	} `yaml:"jwt"`
}

func main() {
	// 加载配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 数据库连接
	db, err := initDB(config)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化存储层
	store := store.NewStore(db)

	// 初始化服务层
	services := service.NewServices(store, config.JWT.Secret, config.JWT.ExpireHours)

	// 初始化路由
	r := setupRouter(services)

	// 启动服务器
	port := config.Server.Port
	if port == 0 {
		port = 8080 // 默认端口
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		log.Printf("服务器启动在 http://localhost:%d\n", port)
		log.Printf("API接口地址: http://localhost:%d/api\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 等待服务器关闭
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Fatalf("服务器强制关闭: %v", err)
	// }

	log.Println("服务器已关闭")
}

// 加载配置文件
func loadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// 初始化数据库连接
func initDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// 设置路由
func setupRouter(services *service.Services) *gin.Engine {
	r := gin.Default()
	// 允许跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API路由组
	apiGroup := r.Group("/api")

	// 管理员接口
	adminGroup := apiGroup.Group("/admin")
	{
		adminHandler := api.NewAdminHandler(services.AdminService)
		adminGroup.POST("/login", adminHandler.Login)
		// 修改密码需要认证
		adminGroup.PUT("/password", api.AuthMiddleware(services.AdminService), adminHandler.ChangePassword)
	}

	// 软件管理接口
	softwareGroup := apiGroup.Group("/software")
	softwareHandler := api.NewSoftwareHandler(services.SoftwareService)
	{
		softwareGroup.Use(api.AuthMiddleware(services.AdminService))
		softwareGroup.POST("", softwareHandler.CreateSoftware)
		softwareGroup.GET("", softwareHandler.GetSoftwareList)
		softwareGroup.GET("/:akey", softwareHandler.GetSoftwareInfo)
		softwareGroup.PUT("/:akey", softwareHandler.UpdateSoftware)
		softwareGroup.DELETE("/:akey", softwareHandler.DeleteSoftware)

		// 版本管理接口
		versionGroup := softwareGroup.Group("/:akey/versions")
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
		versionDetailGroup.Use(api.AuthMiddleware(services.AdminService))
		versionDetailGroup.GET("/:vkey", versionDetailHandler.GetVersionInfo)
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

	return r
}
