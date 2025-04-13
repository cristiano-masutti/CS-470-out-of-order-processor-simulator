#!/bin/bash
set -e

echo "🔧 Starting build..."

GO_VERSION=1.23.1

# Clean any old Go installs
rm -rf /usr/local/go

# Download Go for ARM64
wget -q https://go.dev/dl/go${GO_VERSION}.linux-arm64.tar.gz
tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz
export PATH=/usr/local/go/bin:$PATH

# Confirm Go version
go version

echo "⚙️ Building ARM64 binary..."

# Go to source directory and build
cd source_code
go mod tidy
GOARCH=arm64 GOOS=linux go build -o ../simulator

echo "✅ Build complete. Binary available at ./simulator"
