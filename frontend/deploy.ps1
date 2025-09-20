# VerKeyOSS Frontend 部署脚本

# 构建项目
Write-Host "开始构建项目..." -ForegroundColor Green
npm run build

if ($LASTEXITCODE -eq 0) {
    Write-Host "项目构建成功！" -ForegroundColor Green
    Write-Host "构建文件位于 dist/ 目录" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "部署说明：" -ForegroundColor Cyan
    Write-Host "1. 将 dist/ 目录中的所有文件复制到您的 Web 服务器" -ForegroundColor White
    Write-Host "2. 配置 Web 服务器支持 SPA 路由（History 模式）" -ForegroundColor White
    Write-Host "3. 确保后端 API 服务正在运行" -ForegroundColor White
    Write-Host "4. 检查 CORS 配置是否正确" -ForegroundColor White
} else {
    Write-Host "项目构建失败！请检查错误信息。" -ForegroundColor Red
}