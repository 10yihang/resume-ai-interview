@echo off
REM 创建必要的目录
mkdir data 2>nul
mkdir uploads\resumes 2>nul
mkdir uploads\jds 2>nul

REM 检查是否有.env文件，如果没有则从示例创建
if not exist .env (
    echo Creating .env file from .env.example...
    copy .env.example .env
    echo Please edit .env file and set your OpenAI API Key
    pause
    exit /b
)

REM 运行应用
go run cmd/server/main.go
