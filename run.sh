#!/bin/bash
# 创建必要的目录
mkdir -p data uploads/resumes uploads/jds

# 检查是否有.env文件，如果没有则从示例创建
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please edit .env file and set your OpenAI API Key"
    exit 1
fi

# 运行应用
go run cmd/server/main.go
