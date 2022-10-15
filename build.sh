#!/bin/bash

export GOPROXY="https://goproxy.io,direct"
export GOSUMDB="sum.golang.google.cn"

go mod init
go mod tidy

# Windows x64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -trimpath -ldflags "-s -w" -o redirector_windows_amd64.exe
# Linux x64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -trimpath -ldflags "-s -w" -o redirector_linux_amd64

# 上方参数中 -ldflags "-s -w" 用于清除调试符号 增加逆向难度, -trimpath 用来去除编译程序中与路径相关想信息 避免逆向分析人员从程序中获得你的用户名

# 根据需求 还可以编译 linux / windows / macOS / freebsd 的 x86 MIPS Arm/Arm64 等架构的可执行程序