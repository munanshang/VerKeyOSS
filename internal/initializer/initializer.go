package initializer

import (
	"log"
	"time"

	"verkeyoss/internal/model"

	"gorm.io/gorm"
)

// Initialize 执行程序初始化，包括表创建和初始数据插入
// 如果初始化过程中发生错误，程序将退出
func Initialize(db *gorm.DB) {
	log.Println("开始执行程序初始化...")

	// 暂时禁用外键约束，以避免表创建顺序导致的问题
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	// 创建所有表
	createTableIfNotExists(db, &model.App{}, "应用")
	createTableIfNotExists(db, &model.Version{}, "版本")
	createTableIfNotExists(db, &model.Announcement{}, "公告")

	// 重新启用外键约束
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	// 初始化默认应用和版本
	initDefaultSoftwareAndVersion(db)
	// 初始化测试公告
	initTestAnnouncement(db)

	log.Println("程序初始化完成！")
}

// createTableIfNotExists 创建表并返回是否是新创建的表
// 参数 db 是数据库连接
// 参数 model 是要创建表的模型指针
// 参数 tableName 是表的中文名称
// 返回值表示表是否是新创建的
func createTableIfNotExists(db *gorm.DB, model interface{}, tableName string) bool {
	// 检查表是否存在
	if db.Migrator().HasTable(model) {
		return false
	}

	// 创建表
	if err := db.AutoMigrate(model); err != nil {
		log.Fatalf("创建%s表失败: %v", tableName, err)
	}
	log.Printf("%s表创建成功\n", tableName)
	return true
}

// 初始化默认应用和版本
func initDefaultSoftwareAndVersion(db *gorm.DB) {
	// 声明错误变量
	var err error
	// 检查默认应用是否已存在
	var app model.App
	err = db.Where("a_key = ?", "test").First(&app).Error

	if err == nil {
		// 已存在默认应用，跳过初始化
		return
	}

	// 创建默认应用
	defaultApp := model.App{
		UserID:      1, // 管理员ID固定为1
		AKey:        "test",
		Name:        "测试应用",
		Description: "这是一个用于测试的默认应用",
		IsPaid:      false,
	}

	// 保存应用到数据库
	if createErr := db.Create(&defaultApp).Error; createErr != nil {
		log.Fatalf("创建默认应用失败: %v", err)
	}
	log.Println("默认应用创建成功!")

	// 检查默认版本是否已存在
	var version model.Version
	err = db.Where("v_key = ?", "test").First(&version).Error

	if err == nil {
		// 已存在默认版本，跳过初始化
	} else {
		// 创建默认版本
		defaultVersion := model.Version{
			VKey:        "test",
			AKey:        "test",
			Version:     "1.0.0",
			Description: "这是测试应用的初始版本",
			IsLatest:    true,
		}

		// 保存版本到数据库
		if err := db.Create(&defaultVersion).Error; err != nil {
			log.Fatalf("创建默认版本失败: %v", err)
		}

		log.Println("默认版本创建成功!")
	}
}

// 初始化测试公告
func initTestAnnouncement(db *gorm.DB) {
	// 检查是否已存在公告
	var count int64
	db.Model(&model.Announcement{}).Count(&count)
	if count > 0 {
		// 公告已存在，跳过初始化
		return
	}

	// 创建测试公告
	announcement := model.Announcement{
		Title:       "欢迎使用 VerKeyOSS",
		Content:     "VerKeyOSS 是一款开源的应用版本管理系统，用于管理应用标识（AKey）和版本标识（VKey），支持合法性校验和新版本检测，适用于各类应用的版本管控场景。",
		IsActive:    true,
		PublishDate: time.Now(),
		URL:         "https://github.com/munanshang/verkeyoss", // 示例URL，指向项目GitHub地址
	}

	if err := db.Create(&announcement).Error; err != nil {
		log.Fatalf("创建测试公告失败: %v", err)
	}

	log.Println("测试公告创建成功!")
}
