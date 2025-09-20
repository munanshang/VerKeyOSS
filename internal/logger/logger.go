package logger

import (
	"log"
	"os"
)

// Logger 日志记录器结构
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var (
	// AppLogger 全局日志记录器实例
	AppLogger *Logger
)

// Init 初始化日志记录器
func Init(debug bool) {
	// 创建或打开日志文件
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("无法创建日志文件: %v", err)
		// 如果无法创建文件，使用标准输出
		logFile = os.Stdout
	}

	AppLogger = &Logger{
		infoLogger:  log.New(logFile, "[INFO] ", log.LstdFlags|log.Lshortfile),
		errorLogger: log.New(logFile, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debugLogger: log.New(logFile, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
	}

	// 调试模式下同时输出到控制台
	if debug {
		AppLogger.infoLogger.SetOutput(os.Stdout)
		AppLogger.errorLogger.SetOutput(os.Stderr)
		AppLogger.debugLogger.SetOutput(os.Stdout)
	}
}

// Info 记录信息日志
func (l *Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

// Error 记录错误日志
func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

// Debug 记录调试日志
func (l *Logger) Debug(v ...interface{}) {
	l.debugLogger.Println(v...)
}

// Infof 格式化记录信息日志
func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Errorf 格式化记录错误日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debugf 格式化记录调试日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

// 便捷函数
func Info(v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Info(v...)
	}
}

func Error(v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Error(v...)
	}
}

func Debug(v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debug(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Infof(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Errorf(format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if AppLogger != nil {
		AppLogger.Debugf(format, v...)
	}
}
