@echo off
chcp 65001

:: 编译可执行文件
SET GOOS=windows
SET GOARCH=amd64
go build -o ./dist/Win64/ReleaseDelivr.exe main.go

echo ---------------------- Start ----------------------

:: 运行可执行文件 文件名务必更改
cd dist/Win64

ReleaseDelivr.exe

echo ----------------------- End -----------------------
