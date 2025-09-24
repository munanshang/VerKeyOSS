#!/bin/bash

echo "Build VerKeyOSS Application..."
echo

# Get version number
read -p "Please enter build version (e.g., 1.0.0): " VERSION
if [ -z "$VERSION" ]; then
    echo "Version number cannot be empty!"
    exit 1
fi

echo "Build version: $VERSION"
echo

# Create output directory
OUTPUT_DIR="bin/$VERSION"
if [ ! -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR"
fi

# Force rebuild frontend
echo "Cleaning frontend build files..."
if [ -d "frontend/dist" ]; then
    rm -rf "frontend/dist"
fi

echo "Building frontend..."
cd frontend
npm run build
if [ $? -ne 0 ]; then
    echo "Frontend build failed!"
    cd ..
    exit 1
fi
cd ..
echo "Frontend build completed"
echo

# Build Windows version
echo "Building Windows version..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$VERSION" -o "$OUTPUT_DIR/verkeyoss-$VERSION-windows-amd64.exe" .
if [ $? -ne 0 ]; then
    echo "Windows build failed!"
    exit 1
fi

# Build Linux version
echo "Building Linux version..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$VERSION" -o "$OUTPUT_DIR/verkeyoss-$VERSION-linux-amd64" .
if [ $? -ne 0 ]; then
    echo "Linux build failed!"
    exit 1
fi

echo
echo "Build completed!"
echo "Output directory: $OUTPUT_DIR"
echo
echo "Generated files:"
for file in "$OUTPUT_DIR"/*; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        filesize=$(du -h "$file" | cut -f1)
        echo "  $filename - $filesize"
    fi
done
echo