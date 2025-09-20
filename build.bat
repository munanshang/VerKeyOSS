@echo off
chcp 65001 >nul
echo 构建VerKeyOSS应用...
echo.

REM 获取版本号
set /p VERSION="请输入打包版本号 (例如: 1.0.0): "
if "%VERSION%"=="" (
    echo 版本号不能为空！
    pause
    exit /b 1
)

echo 打包版本: %VERSION%
echo.

REM 创建输出目录
set OUTPUT_DIR=bin\%VERSION%
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
)

REM 强制重新构建前端
echo 清理前端构建文件...
if exist "frontend\dist" (
    rmdir /s /q "frontend\dist"
)

echo 构建前端...
cd frontend
call npm run build
if %errorlevel% neq 0 (
    echo 前端构建失败！
    cd ..
    pause
    exit /b 1
)
cd ..
echo 前端构建完成
echo.

REM 构建Windows版本
echo 构建Windows版本...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\verkeyoss-windows-amd64.exe" .
if %errorlevel% neq 0 (
    echo Windows版本构建失败！
    pause
    exit /b 1
)

REM 构建Linux版本
echo 构建Linux版本...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\verkeyoss-linux-amd64" .
if %errorlevel% neq 0 (
    echo Linux版本构建失败！
    pause
    exit /b 1
)

REM 重置环境变量
set GOOS=
set GOARCH=

echo.
echo 构建完成！
echo 输出目录: %OUTPUT_DIR%
echo.
echo 生成的文件:
for %%f in ("%OUTPUT_DIR%\*") do (
    echo   %%~nxf - %%~zf 字节
)
echo.
pause