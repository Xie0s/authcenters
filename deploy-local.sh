#!/bin/bash

# AuthCenter 本地部署脚本
# 支持用户名和邮箱登录的企业级认证授权系统

set -e

echo "🚀 AuthCenter 本地部署脚本"
echo "================================"

# 检查必要的工具
check_requirements() {
    echo "📋 检查系统要求..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装。请先安装 Go 1.19 或更高版本"
        echo "   下载地址: https://golang.org/dl/"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo "✅ Go 版本: $GO_VERSION"
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker 未安装。请先安装 Docker"
        echo "   下载地址: https://www.docker.com/get-started"
        exit 1
    fi
    
    echo "✅ Docker 已安装"
    
    # 检查 Node.js (用于数据库初始化)
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js 未安装。请先安装 Node.js"
        echo "   下载地址: https://nodejs.org/"
        exit 1
    fi
    
    NODE_VERSION=$(node --version)
    echo "✅ Node.js 版本: $NODE_VERSION"
}

# 启动 MongoDB
start_mongodb() {
    echo "🗄️  启动 MongoDB..."
    
    # 检查是否已有 MongoDB 容器在运行
    if docker ps | grep -q "mongodb"; then
        echo "✅ MongoDB 容器已在运行"
    else
        # 停止并删除可能存在的旧容器
        docker stop mongodb 2>/dev/null || true
        docker rm mongodb 2>/dev/null || true
        
        # 启动新的 MongoDB 容器
        docker run -d \
            --name mongodb \
            -p 27017:27017 \
            -v mongodb_data:/data/db \
            mongo:latest
        
        echo "✅ MongoDB 容器已启动"
        
        # 等待 MongoDB 启动完成
        echo "⏳ 等待 MongoDB 启动完成..."
        sleep 10
    fi
}

# 初始化数据库
init_database() {
    echo "🔧 初始化数据库..."
    
    cd scripts
    
    # 安装依赖
    if [ ! -d "node_modules" ]; then
        echo "📦 安装 Node.js 依赖..."
        npm install
    fi
    
    # 运行初始化脚本
    echo "🏗️  执行数据库初始化..."
    node run_init_new.js
    
    cd ..
    echo "✅ 数据库初始化完成"
}

# 安装 Go 依赖
install_go_deps() {
    echo "📦 安装 Go 依赖..."
    go mod tidy
    echo "✅ Go 依赖安装完成"
}

# 构建项目
build_project() {
    echo "🔨 构建项目..."
    go build -o authcenter cmd/server/main.go
    echo "✅ 项目构建完成"
}

# 启动服务
start_server() {
    echo "🚀 启动 AuthCenter 服务..."
    echo ""
    echo "服务将在以下地址启动:"
    echo "  - API 服务: http://localhost:8080"
    echo "  - 测试页面: http://localhost:8080/test/"
    echo "  - 健康检查: http://localhost:8080/health"
    echo ""
    echo "按 Ctrl+C 停止服务"
    echo ""
    
    ./authcenter
}

# 显示使用说明
show_usage() {
    echo ""
    echo "🎉 AuthCenter 部署完成！"
    echo "========================"
    echo ""
    echo "🌟 新功能特性:"
    echo "  ✅ 注册时必须提供用户名"
    echo "  ✅ 支持用户名登录"
    echo "  ✅ 支持邮箱登录"
    echo "  ✅ 支持自动识别登录（用户名或邮箱）"
    echo ""
    echo "🔗 访问地址:"
    echo "  - API 文档: http://localhost:8080/health"
    echo "  - 测试页面: http://localhost:8080/test/"
    echo ""
    echo "📝 API 示例:"
    echo "  # 注册用户"
    echo "  curl -X POST http://localhost:8080/api/v1/auth/register \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}'"
    echo ""
    echo "  # 用户名登录"
    echo "  curl -X POST http://localhost:8080/api/v1/auth/login \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"testuser\",\"password\":\"password123\",\"type\":\"username\"}'"
    echo ""
    echo "  # 邮箱登录"
    echo "  curl -X POST http://localhost:8080/api/v1/auth/login \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"email\":\"test@example.com\",\"password\":\"password123\",\"type\":\"email\"}'"
    echo ""
    echo "  # 自动识别登录"
    echo "  curl -X POST http://localhost:8080/api/v1/auth/login \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"username\":\"testuser\",\"password\":\"password123\",\"type\":\"auto\"}'"
    echo ""
    echo "🛠️  管理命令:"
    echo "  - 停止服务: Ctrl+C"
    echo "  - 停止 MongoDB: docker stop mongodb"
    echo "  - 查看日志: docker logs mongodb"
    echo "  - 重新初始化数据库: cd scripts && node run_init_new.js"
    echo ""
}

# 主函数
main() {
    case "${1:-all}" in
        "check")
            check_requirements
            ;;
        "mongodb")
            start_mongodb
            ;;
        "init")
            init_database
            ;;
        "build")
            install_go_deps
            build_project
            ;;
        "start")
            start_server
            ;;
        "all")
            check_requirements
            start_mongodb
            init_database
            install_go_deps
            build_project
            show_usage
            start_server
            ;;
        *)
            echo "用法: $0 [check|mongodb|init|build|start|all]"
            echo ""
            echo "  check    - 检查系统要求"
            echo "  mongodb  - 启动 MongoDB"
            echo "  init     - 初始化数据库"
            echo "  build    - 构建项目"
            echo "  start    - 启动服务"
            echo "  all      - 执行完整部署流程（默认）"
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
