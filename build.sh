#!/bin/bash
set -e

# Navigate into source_code
cd source_code

# Download dependencies and build
go mod tidy
go build -o ../simulator main.go
