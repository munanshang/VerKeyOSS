@echo off
chcp 65001 >nul 2>&1
setlocal

echo Build VerKeyOSS Application...
echo.

REM Get version number
echo Please enter build version (e.g., 1.0.0):
set /p VERSION=

if "%VERSION%"=="" (
    echo Version number cannot be empty!
    pause
    exit /b 1
)

echo Build version: %VERSION%
echo.

REM Create output directory
set OUTPUT_DIR=bin\%VERSION%
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%"
)

REM Force rebuild frontend
echo Cleaning frontend build files...
if exist "frontend\dist" (
    rmdir /s /q "frontend\dist"
)

echo Building frontend...
cd frontend
call npm run build
if %errorlevel% neq 0 (
    echo Frontend build failed!
    cd ..
    pause
    exit /b 1
)
cd ..
echo Frontend build completed
echo.

REM Build Windows version
echo Building Windows version...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\verkeyoss-windows-amd64.exe" .
if %errorlevel% neq 0 (
    echo Windows build failed!
    pause
    exit /b 1
)

REM Build Linux version
echo Building Linux version...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%OUTPUT_DIR%\verkeyoss-linux-amd64" .
if %errorlevel% neq 0 (
    echo Linux build failed!
    pause
    exit /b 1
)

REM Reset environment variables
set GOOS=
set GOARCH=

echo.
echo Build completed!
echo Output directory: %OUTPUT_DIR%
echo.
echo Generated files:
for %%f in ("%OUTPUT_DIR%\*") do (
    echo   %%~nxf - %%~zf bytes
)
echo.
pause