@echo off
REM Build for Windows
set GOOS=windows
set GOARCH=amd64
go build -o ./dist/data-generator-windows.exe main.go

REM Build for macOS (Apple M1/M2)
set GOOS=darwin
set GOARCH=arm64
go build -o ./dist/data-generator-arm64 main.go

REM Build for Linux
set GOOS=linux
set GOARCH=amd64
go build -o ./dist/data-generator-linux main.go

echo Builds completed for Windows, macOS, and Linux!