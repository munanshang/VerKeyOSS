package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"verkeyoss/internal/config"
	"verkeyoss/internal/initializer"
	"verkeyoss/internal/logger"
	"verkeyoss/internal/router"
	"verkeyoss/internal/service"
	"verkeyoss/internal/store"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 版本信息，在构建时通过 -ldflags 注入
var version = "dev"

// 嵌入前端构建后的文件
//
//go:embed all:frontend/dist
var frontendFS embed.FS

func main() {
	// 显示版本信息
	log.Printf("VerKeyOSS 版本: %s", version)

	// 加载配置文件
	appConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("配置文件加载失败: %v", err)
		log.Println("请检查 config.yaml 文件是否存在并且格式正确")
		log.Println("可以复制 config.example.yaml 为 config.yaml 并修改其中的配置")
		log.Println("按任意键退出...")
		fmt.Scanln()
		return
	}

	// 初始化日志系统
	logger.Init(appConfig.Server.Debug)
	logger.Info("应用启动中...")

	// 数据库连接
	logger.Info("正在连接数据库...")
	log.Printf("正在连接数据库: %s:%d/%s", appConfig.DB.Host, appConfig.DB.Port, appConfig.DB.Name)
	db, err := initDB(appConfig)
	if err != nil {
		logger.Error("数据库连接失败:", err)
		log.Printf("数据库连接失败: %v", err)
		log.Println("")
		log.Println("可能的解决方案:")
		log.Println("1. 检查MySQL服务是否已启动")
		log.Println("2. 验证数据库连接配置 (config.yaml)")
		log.Println("3. 确认数据库用户名密码正确")
		log.Println("4. 确认目标数据库已创建")
		log.Println("")
		log.Printf("当前配置: %s:%d, 用户: %s, 数据库: %s",
			appConfig.DB.Host, appConfig.DB.Port, appConfig.DB.User, appConfig.DB.Name)
		log.Println("")
		log.Println("按任意键退出...")
		fmt.Scanln()
		return
	}
	logger.Info("数据库连接成功")

	// 执行数据库初始化
	logger.Info("正在初始化数据库...")
	initializer.Initialize(db)
	logger.Info("数据库初始化完成")

	// 初始化存储层
	store := store.NewStore(db)

	// 初始化服务层
	services := service.NewServices(store, appConfig.JWT.Secret, appConfig.JWT.ExpireHours)

	// 初始化路由
	r := router.SetupRouter(services, appConfig, version, StaticFileHandler(), FrontendHandler())

	// 启动服务器
	port := appConfig.Server.Port
	if port == 0 {
		port = 8913 // 默认端口
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		logger.Infof("服务器启动在 http://localhost:%d", port)
		log.Printf("服务器启动在 http://localhost:%d\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务器启动失败:", err)
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("接收到关闭信号，正在关闭服务器...")
	log.Println("正在关闭服务器...")

	// 等待服务器关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器强制关闭:", err)
		log.Fatalf("服务器强制关闭: %v", err)
	}

	logger.Info("服务器已安全关闭")
	log.Println("服务器已关闭")
}

// initDB 初始化数据库连接
func initDB(appConfig *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		appConfig.DB.User, appConfig.DB.Password, appConfig.DB.Host, appConfig.DB.Port, appConfig.DB.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 配置GORM参数
		NowFunc: time.Now,
	})

	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	// 设置最大开放连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置连接的最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// GetFrontendFS 获取前端文件系统
func GetFrontendFS() (fs.FS, error) {
	return fs.Sub(frontendFS, "frontend/dist")
}

// FrontendHandler 处理前端路由
func FrontendHandler() gin.HandlerFunc {
	frontendFS, err := GetFrontendFS()
	if err != nil {
		panic("无法获取前端文件系统: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(frontendFS))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查文件是否存在
		if _, err := frontendFS.Open(strings.TrimPrefix(path, "/")); err != nil {
			// 如果是API路径，返回404
			if strings.HasPrefix(path, "/api") {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			// 对于前端路由，返回index.html（SPA应用）
			c.Request.URL.Path = "/"
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// StaticFileHandler 处理静态文件
func StaticFileHandler() gin.HandlerFunc {
	frontendFS, err := GetFrontendFS()
	if err != nil {
		panic("无法获取前端文件系统: " + err.Error())
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查是否是静态资源文件
		ext := filepath.Ext(path)
		if ext == ".js" || ext == ".css" || ext == ".png" || ext == ".jpg" || ext == ".ico" || ext == ".svg" {
			// 设置缓存头
			c.Header("Cache-Control", "public, max-age=31536000")
		}

		// 尝试打开文件
		if _, err := frontendFS.Open(strings.TrimPrefix(path, "/")); err == nil {
			fileServer := http.FileServer(http.FS(frontendFS))
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// 文件不存在，继续到下一个handler
		c.Next()
	}
}
