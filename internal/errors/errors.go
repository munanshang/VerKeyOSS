package errors

import (
	"errors"
	"fmt"
)

// 定义错误类型
type ErrorType int

const (
	// ErrTypeValidation 参数验证错误
	ErrTypeValidation ErrorType = iota
	// ErrTypeNotFound 资源不存在错误
	ErrTypeNotFound
	// ErrTypeUnauthorized 未授权错误
	ErrTypeUnauthorized
	// ErrTypeInternal 内部服务器错误
	ErrTypeInternal
	// ErrTypeConflict 资源冲突错误
	ErrTypeConflict
)

// AppError 应用错误结构
type AppError struct {
	Type    ErrorType
	Code    int
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持error unwrapping
func (e *AppError) Unwrap() error {
	return e.Err
}

// 预定义错误
var (
	ErrAppNotFound     = NewNotFoundError("应用不存在")
	ErrVersionNotFound = NewNotFoundError("版本不存在")
	ErrInvalidParams   = NewValidationError("参数错误")
	ErrUnauthorized    = NewUnauthorizedError("未授权访问")
	ErrTokenInvalid    = NewUnauthorizedError("令牌无效或已过期")
	ErrAppExists       = NewConflictError("应用已存在")
	ErrVersionExists   = NewConflictError("版本已存在")
)

// NewValidationError 创建参数验证错误
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:    ErrTypeValidation,
		Code:    400,
		Message: message,
	}
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    ErrTypeNotFound,
		Code:    404,
		Message: message,
	}
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrTypeUnauthorized,
		Code:    401,
		Message: message,
	}
}

// NewInternalError 创建内部服务器错误
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrTypeInternal,
		Code:    500,
		Message: message,
		Err:     err,
	}
}

// NewConflictError 创建资源冲突错误
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrTypeConflict,
		Code:    409,
		Message: message,
	}
}

// WrapError 包装现有错误
func WrapError(err error, message string) *AppError {
	return &AppError{
		Type:    ErrTypeInternal,
		Code:    500,
		Message: message,
		Err:     err,
	}
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
