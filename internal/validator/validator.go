package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"verkeyoss/internal/errors"
)

// 正则表达式模式
var (
	// 应用名称：支持中英文、数字、下划线、短横线、空格、句点、圆括号
	appNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\p{Han}\s\.\(\)]{2,50}$`)
	// 版本号：语义化版本号格式
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+(-[a-zA-Z0-9\-]+)?(\+[a-zA-Z0-9\-]+)?$`)
	// 用户名：字母、数字、下划线
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

// ValidateAppName 验证应用名称
func ValidateAppName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.NewValidationError("应用名称不能为空")
	}

	if utf8.RuneCountInString(name) < 2 {
		return errors.NewValidationError("应用名称至少需要2个字符")
	}

	if utf8.RuneCountInString(name) > 50 {
		return errors.NewValidationError("应用名称不能超过50个字符")
	}

	if !appNamePattern.MatchString(name) {
		return errors.NewValidationError("应用名称包含无效字符")
	}

	return nil
}

// ValidateDescription 验证描述信息
func ValidateDescription(desc string) error {
	if utf8.RuneCountInString(desc) > 500 {
		return errors.NewValidationError("描述信息不能超过500个字符")
	}
	return nil
}

// ValidateVersion 验证版本号
func ValidateVersion(version string) error {
	version = strings.TrimSpace(version)
	if version == "" {
		return errors.NewValidationError("版本号不能为空")
	}

	if !versionPattern.MatchString(version) {
		return errors.NewValidationError("版本号格式无效，请使用语义化版本号格式（如：1.0.0）")
	}

	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.NewValidationError("用户名不能为空")
	}

	if !usernamePattern.MatchString(username) {
		return errors.NewValidationError("用户名只能包含字母、数字和下划线，长度为3-20个字符")
	}

	return nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.NewValidationError("密码长度不能少于6个字符")
	}

	if len(password) > 50 {
		return errors.NewValidationError("密码长度不能超过50个字符")
	}

	// 检查是否包含至少一个字母和一个数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	if !hasLetter || !hasNumber {
		return errors.NewValidationError("密码必须包含至少一个字母和一个数字")
	}

	return nil
}

// ValidatePagination 验证分页参数
func ValidatePagination(page, size int) (int, int, error) {
	if page <= 0 {
		page = 1
	}

	if size <= 0 {
		size = 10
	}

	if size > 100 {
		return 0, 0, errors.NewValidationError("每页数量不能超过100条")
	}

	return page, size, nil
}

// ValidateAKey 验证应用标识
func ValidateAKey(akey string) error {
	if akey == "" {
		return errors.NewValidationError("应用标识不能为空")
	}

	// 支持测试数据（最少 3 个字符）和正式数据（最多 100 个字符）
	if len(akey) < 3 || len(akey) > 100 {
		return errors.NewValidationError("应用标识长度无效（必须为 3-100 个字符）")
	}

	return nil
}

// ValidateVKey 验证版本标识
func ValidateVKey(vkey string) error {
	if vkey == "" {
		return errors.NewValidationError("版本标识不能为空")
	}

	// 支持测试数据（最少 3 个字符）和正式数据（最多 100 个字符）
	if len(vkey) < 3 || len(vkey) > 100 {
		return errors.NewValidationError("版本标识长度无效（必须为 3-100 个字符）")
	}

	return nil
}
