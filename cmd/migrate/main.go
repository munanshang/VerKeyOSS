package main

import (
	"fmt"
	"log"

	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"verkeyoss/internal/model"
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
}

func main() {
	// 加载配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		// 尝试从 config.example.yaml 加载
		config, err = loadConfig("config.example.yaml")
		if err != nil {
			log.Fatalf("加载配置文件失败: %v\n请确保存在 config.yaml 或 config.example.yaml 文件", err)
		}
	}

	// 数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("数据库连接失败: %v\n请检查数据库配置是否正确", err)
	}

	// 执行迁移
	log.Println("开始执行数据库迁移...")

	// 暂时禁用外键约束，以避免表创建顺序导致的问题
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	// 创建管理员表
	if err := db.AutoMigrate(&model.Admin{}); err != nil {
		log.Fatalf("创建管理员表失败: %v", err)
	}
	log.Println("管理员表创建成功")

	// 创建软件表
	if err := db.AutoMigrate(&model.Software{}); err != nil {
		log.Fatalf("创建软件表失败: %v", err)
	}
	log.Println("软件表创建成功")

	// 创建版本表
	if err := db.AutoMigrate(&model.Version{}); err != nil {
		log.Fatalf("创建版本表失败: %v", err)
	}
	log.Println("版本表创建成功")

	// 重新启用外键约束
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	// 初始化管理员账户
	initAdminAccount(db)

	// 初始化默认软件和版本
	initDefaultSoftwareAndVersion(db)

	log.Println("数据库迁移完成！")
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

// 初始化管理员账户
func initAdminAccount(db *gorm.DB) {
	// 检查是否已存在管理员账户
	var admin model.Admin
	err := db.First(&admin).Error

	if err == nil {
		// 已存在管理员账户
		log.Println("管理员账户已存在，跳过初始化")
		return
	}

	// 创建默认管理员账户
	// 注意：密码设置为明文，模型的 BeforeCreate 钩子会自动进行加密
	defaultAdmin := model.Admin{
		Username: "verkey",
		Password: "verkey", // 默认密码，用户首次登录后应修改
	}

	// 保存到数据库
	if err := db.Create(&defaultAdmin).Error; err != nil {
		log.Fatalf("创建默认管理员账户失败: %v", err)
	}

	log.Println("默认管理员账户创建成功!")
	log.Println("用户名: verkey")
	log.Println("密码: verkey")
	log.Println("首次登录后请立即修改密码")
}

// 初始化默认软件和版本
func initDefaultSoftwareAndVersion(db *gorm.DB) {
	// 检查默认软件是否已存在
	var software model.Software
	err := db.Where("a_key = ?", "test").First(&software).Error

	if err == nil {
		// 已存在默认软件
		log.Println("默认软件已存在，跳过初始化")
		return
	}

	// 创建默认软件
	defaultSoftware := model.Software{
		AKey:        "test",
		Name:        "测试软件",
		Description: "这是一个用于测试的默认软件",
	}

	// 保存软件到数据库
	if err := db.Create(&defaultSoftware).Error; err != nil {
		log.Fatalf("创建默认软件失败: %v", err)
	}

	// 创建默认版本
	defaultVersion := model.Version{
		VKey:        "test",
		AKey:        "test",
		Version:     "1.0.0",
		Description: "这是测试软件的初始版本",
		IsLatest:    true,
	}

	// 保存版本到数据库
	if err := db.Create(&defaultVersion).Error; err != nil {
		log.Fatalf("创建默认版本失败: %v", err)
	}

	log.Println("默认软件和版本创建成功!")
	log.Println("软件AKey: test")
	log.Println("版本VKey: test")
}
