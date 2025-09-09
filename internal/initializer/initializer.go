package initializer

import (
	"log"
	"time"

	"verkeyoss/internal/model"

	"gorm.io/gorm"
)

// Initialize 执行数据库初始化，包括表创建和初始数据插入
// 如果初始化过程中发生错误，程序将退出
func Initialize(db *gorm.DB) {
	log.Println("开始执行数据库初始化...")

	// 暂时禁用外键约束，以避免表创建顺序导致的问题
	db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	// 创建所有表
	createTableIfNotExists(db, &model.PermissionGroup{}, "权限组")
	createTableIfNotExists(db, &model.User{}, "用户")
	createTableIfNotExists(db, &model.App{}, "应用")
	createTableIfNotExists(db, &model.Version{}, "版本")
	createTableIfNotExists(db, &model.Announcement{}, "公告")

	// 重新启用外键约束
	db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	// 检查是否首次初始化（通过检查权限组表是否包含数据）
	var adminGroup model.PermissionGroup
	err := db.Where("group_name = ?", "管理员").First(&adminGroup).Error

	if err != nil {
		// 首次初始化，创建所有默认数据
		initPermissionGroups(db)
		initAdminAccount(db)
		initDefaultSoftwareAndVersion(db)
		initTestAnnouncement(db)
	} else {
		// 不是首次初始化，只检查并创建可能缺失的关键数据，但不覆盖已存在的数据
		// 权限组检查已在上面完成
		// 管理员账户检查
		initAdminAccount(db)
		// 默认应用和版本检查
		initDefaultSoftwareAndVersion(db)
		// 测试公告检查
		initTestAnnouncement(db)
	}

	log.Println("数据库初始化完成！")
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

// 初始化权限组
func initPermissionGroups(db *gorm.DB) {
	// 检查是否已存在管理员权限组
	var adminGroup model.PermissionGroup
	err := db.Where("group_name = ?", "管理员").First(&adminGroup).Error

	if err == nil {
		// 已存在权限组
		log.Println("权限组已存在，跳过初始化")
		return
	}

	// 创建管理员权限组（ID=1）
	adminGroup = model.PermissionGroup{
		GroupName:   "管理员",
		Permission:  `{"all_access": true}`,
		Description: "拥有系统所有权限的超级管理员组",
	}

	if err := db.Create(&adminGroup).Error; err != nil {
		log.Fatalf("创建管理员权限组失败: %v", err)
	}

	// 创建普通用户权限组（ID=2）
	userGroup := model.PermissionGroup{
		GroupName:   "普通用户",
		Permission:  `{"all_access": false, "create_app": true, "manage_own_app": true}`,
		Description: "可以创建和管理自己应用的普通用户组",
	}

	if err := db.Create(&userGroup).Error; err != nil {
		log.Fatalf("创建普通用户权限组失败: %v", err)
	}

	log.Println("权限组初始化成功!")
}

// 初始化管理员账户
func initAdminAccount(db *gorm.DB) {
	// 检查是否已存在用户表和管理员用户
	var user model.User
	err := db.Where("username = ?", "admin").First(&user).Error

	if err == nil {
		// 已存在管理员用户
		log.Println("管理员用户已存在，跳过初始化")
		return
	}

	// 获取管理员权限组ID（应该是1）
	var adminGroup model.PermissionGroup
	err = db.Where("group_name = ?", "管理员").First(&adminGroup).Error
	if err != nil {
		log.Fatalf("获取管理员权限组失败: %v", err)
	}

	// 创建默认管理员用户
	// 注意：密码设置为明文，模型的 BeforeCreate 钩子会自动进行加密
	defaultUser := model.User{
		Username: "admin",
		Password: "admin123", // 默认登录密码
		Email:    "admin@example.com",
		IsBanned:   false,
		GroupID:  adminGroup.ID, // 设置为管理员权限组
		// 保留SecondPassword字段作为预留功能
		SecondPassword: "", // 空值，当前版本不使用
	}

	// 保存到数据库
	if err := db.Create(&defaultUser).Error; err != nil {
		log.Fatalf("创建默认管理员用户失败: %v", err)
	}

	log.Println("默认管理员用户创建成功!")
	log.Println("用户名: admin")
	log.Println("登录密码: admin123")
	log.Println("邮箱: admin@example.com")
	log.Println("首次登录后请立即修改密码")
}

// 初始化默认应用和版本
func initDefaultSoftwareAndVersion(db *gorm.DB) {
	// 声明错误变量
	var err error
	// 检查默认应用是否已存在
	var app model.App
	err = db.Where("a_key = ?", "test").First(&app).Error

	if err == nil {
		// 已存在默认应用
		log.Println("默认应用已存在，跳过初始化")
		return
	}

	// 获取管理员用户（通过用户名查找）
	var adminUser model.User
	err = db.Where("username = ?", "admin").First(&adminUser).Error
	if err != nil {
		log.Fatalf("获取管理员用户失败: %v", err)
	}

	// 创建默认应用，关联到管理员用户
	defaultApp := model.App{
		UserID:      adminUser.ID,
		AKey:        "test",
		Name:        "测试应用",
		Description: "这是一个用于测试的默认应用",
		IsBanned:      false,
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
		// 已存在默认版本
		log.Println("默认版本已存在，跳过初始化")
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
		log.Println("公告已存在，跳过初始化")
		return
	}

	// 创建测试公告
	announcement := model.Announcement{
		Title:       "欢迎使用 VerKeyOSS",
		Content:     "这是一条测试公告，VerKeyOSS 是一款开源的应用版本管理系统，用于管理应用标识（AKey）和版本标识（VKey），支持合法性校验和新版本检测，适用于各类应用的版本管控场景。",
		IsActive:    true,
		PublishDate: time.Now(),
		URL:         "https://github.com/munanshang/verkeyoss", // 示例URL，指向项目GitHub地址
	}

	if err := db.Create(&announcement).Error; err != nil {
		log.Fatalf("创建测试公告失败: %v", err)
	}

	log.Println("测试公告创建成功!")
}
