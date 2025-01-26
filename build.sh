#!/bin/bash

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o ./dist/data-generator-windows.exe main.go

# Build for macOS (Apple M1/M2)
GOOS=darwin GOARCH=arm64 go build -o ./dist/data-generator-arm64 main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o ./dist/data-generator-linux main.go

echo "Builds completed for Windows, macOS (arm64), and Linux!"
