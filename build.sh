#!/bin/bash

echo "构建VerKeyOSS应用..."
echo

# 获取版本号
read -p "请输入打包版本号 (例如: 1.0.0): " VERSION
if [ -z "$VERSION" ]; then
    echo "版本号不能为空！"
    exit 1
fi

echo "打包版本: $VERSION"
echo

# 创建输出目录
OUTPUT_DIR="bin/$VERSION"
if [ ! -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR"
fi

# 强制重新构建前端
echo "清理前端构建文件..."
if [ -d "frontend/dist" ]; then
    rm -rf "frontend/dist"
fi

echo "构建前端..."
cd frontend
npm run build
if [ $? -ne 0 ]; then
    echo "前端构建失败！"
    cd ..
    exit 1
fi
cd ..
echo "前端构建完成"
echo

# 构建Windows版本
echo "构建Windows版本..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$VERSION" -o "$OUTPUT_DIR/verkeyoss-windows-amd64.exe" .
if [ $? -ne 0 ]; then
    echo "Windows版本构建失败！"
    exit 1
fi

# 构建Linux版本
echo "构建Linux版本..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$VERSION" -o "$OUTPUT_DIR/verkeyoss-linux-amd64" .
if [ $? -ne 0 ]; then
    echo "Linux版本构建失败！"
    exit 1
fi

echo
echo "构建完成！"
echo "输出目录: $OUTPUT_DIR"
echo
echo "生成的文件:"
for file in "$OUTPUT_DIR"/*; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        filesize=$(du -h "$file" | cut -f1)
        echo "  $filename - $filesize"
    fi
done
echo