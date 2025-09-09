package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"verkeyoss/internal/initializer"
	"verkeyoss/internal/router"
	"verkeyoss/internal/service"

	"verkeyoss/internal/store"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// //go:embed web
// var vueDist embed.FS

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

	// 执行数据库初始化
	initializer.Initialize(db)

	// 初始化存储层
	store := store.NewStore(db)

	// 初始化服务层
	services := service.NewServices(store, config.JWT.Secret, config.JWT.ExpireHours)

	// 初始化路由
	r := router.SetupRouter(services)

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
