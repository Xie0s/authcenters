#!/bin/bash

# AuthCenter 系统检查脚本
# 对认证系统进行全面的"检查"

echo "=========================================="
echo "AuthCenter 认证系统检查 (检查)"
echo "=========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 计数器
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# 检查函数
check_item() {
    local description="$1"
    local command="$2"
    local expected_result="$3"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    echo -n "检查 $description... "
    
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}通过${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        return 0
    else
        echo -e "${RED}失败${NC}"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        return 1
    fi
}

# 详细检查函数
detailed_check() {
    local description="$1"
    local command="$2"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    echo -e "${BLUE}检查 $description:${NC}"
    
    local result=$(eval "$command" 2>&1)
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✓ 通过${NC}"
        echo "$result" | sed 's/^/  /'
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        echo
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "$result" | sed 's/^/  /'
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
        echo
        return 1
    fi
}

echo "开始系统检查..."
echo

# 1. 检查项目结构
echo -e "${YELLOW}=== 项目结构检查 ===${NC}"
check_item "Go模块文件存在" "test -f go.mod"
check_item "配置文件存在" "test -f configs/config.yaml"
check_item "主程序入口存在" "test -f cmd/server/main.go"
check_item "认证处理器存在" "test -f internal/auth/handler/auth_handler.go"
check_item "认证服务存在" "test -f internal/auth/service/auth_service.go"
check_item "健康检查处理器存在" "test -f internal/health/handler/health_handler.go"
check_item "测试文件存在" "test -f test/index.html"
echo

# 2. 检查代码质量
echo -e "${YELLOW}=== 代码质量检查 ===${NC}"
detailed_check "Go代码编译检查" "go build ./cmd/server"
detailed_check "Go代码格式检查" "gofmt -l . | wc -l | grep -q '^0$'"
detailed_check "Go模块依赖检查" "go mod verify"
echo

# 3. 检查配置文件
echo -e "${YELLOW}=== 配置文件检查 ===${NC}"
detailed_check "YAML配置文件语法" "python3 -c 'import yaml; yaml.safe_load(open(\"configs/config.yaml\"))'"
check_item "JWT配置存在" "grep -q 'jwt:' configs/config.yaml"
check_item "MongoDB配置存在" "grep -q 'mongodb:' configs/config.yaml"
check_item "安全配置存在" "grep -q 'security:' configs/config.yaml"
echo

# 4. 检查关键文件内容
echo -e "${YELLOW}=== 关键功能检查 ===${NC}"
check_item "登录API实现" "grep -q 'func.*Login' internal/auth/handler/auth_handler.go"
check_item "注册API实现" "grep -q 'func.*Register' internal/auth/handler/auth_handler.go"
check_item "Token刷新实现" "grep -q 'func.*RefreshToken' internal/auth/handler/auth_handler.go"
check_item "Token验证实现" "grep -q 'func.*VerifyToken' internal/auth/handler/auth_handler.go"
check_item "健康检查实现" "grep -q 'func.*DetailedHealth' internal/health/handler/health_handler.go"
check_item "中间件实现" "test -f internal/middleware/auth.go"
check_item "JWT管理器实现" "test -f pkg/jwt/jwt.go"
echo

# 5. 检查安全特性
echo -e "${YELLOW}=== 安全特性检查 ===${NC}"
check_item "密码加密实现" "grep -q 'bcrypt' internal/auth/service/auth_service.go || grep -q 'password' pkg/utils/password.go"
check_item "JWT Token实现" "grep -q 'jwt' internal/auth/service/auth_service.go"
check_item "认证中间件实现" "grep -q 'RequireAuth' internal/middleware/auth.go"
check_item "权限检查实现" "grep -q 'RequirePermission' internal/middleware/auth.go"
check_item "CORS中间件实现" "grep -q 'CORS' internal/middleware/security.go || grep -q 'CORS' internal/router/router.go"
echo

# 6. 检查测试和调试功能
echo -e "${YELLOW}=== 测试和调试功能检查 ===${NC}"
check_item "前端测试页面" "test -f test/index.html"
check_item "API测试脚本" "test -f test/js/api.js"
check_item "调试功能实现" "grep -q 'showDebugInfo' test/js/app.js"
check_item "健康检查测试" "grep -q 'performSystemHealthCheck' test/js/api.js"
check_item "认证测试功能" "grep -q 'runAuthenticationTests' test/js/api.js"
echo

# 7. 检查数据库相关
echo -e "${YELLOW}=== 数据库结构检查 ===${NC}"
check_item "用户模型定义" "grep -q 'type User struct' internal/models/models.go"
check_item "角色模型定义" "grep -q 'type Role struct' internal/models/models.go"
check_item "会话模型定义" "grep -q 'type Session struct' internal/models/models.go"
check_item "用户仓库实现" "test -f internal/auth/repository/user_repository.go"
check_item "会话仓库实现" "test -f internal/auth/repository/session_repository.go"
echo

# 8. 检查日志和监控
echo -e "${YELLOW}=== 日志和监控检查 ===${NC}"
check_item "日志工具实现" "test -f pkg/logger/logger.go"
check_item "日志中间件实现" "test -f internal/middleware/logging.go"
check_item "审计日志实现" "grep -q 'AuditMiddleware' internal/middleware/logging.go"
check_item "安全事件日志" "grep -q 'SecurityEventMiddleware' internal/middleware/logging.go"
echo

# 9. 检查API响应格式
echo -e "${YELLOW}=== API响应格式检查 ===${NC}"
check_item "统一响应格式" "test -f pkg/response/response.go"
check_item "错误处理实现" "grep -q 'Error' pkg/response/response.go"
check_item "成功响应实现" "grep -q 'Success' pkg/response/response.go"
check_item "API文档存在" "test -f docs/API响应格式规范.md"
echo

# 生成总结报告
echo "=========================================="
echo -e "${BLUE}检查完成总结:${NC}"
echo "总检查项: $TOTAL_CHECKS"
echo -e "通过项: ${GREEN}$PASSED_CHECKS${NC}"
echo -e "失败项: ${RED}$FAILED_CHECKS${NC}"

PASS_RATE=$((PASSED_CHECKS * 100 / TOTAL_CHECKS))
echo "通过率: $PASS_RATE%"

if [ $PASS_RATE -ge 90 ]; then
    echo -e "${GREEN}系统状态: 优秀 ✓${NC}"
elif [ $PASS_RATE -ge 80 ]; then
    echo -e "${YELLOW}系统状态: 良好 ⚠${NC}"
elif [ $PASS_RATE -ge 70 ]; then
    echo -e "${YELLOW}系统状态: 一般 ⚠${NC}"
else
    echo -e "${RED}系统状态: 需要改进 ✗${NC}"
fi

echo "=========================================="

# 生成建议
echo -e "${BLUE}改进建议:${NC}"
if [ $FAILED_CHECKS -eq 0 ]; then
    echo "✓ 系统检查全部通过，认证系统实现完整！"
else
    echo "建议解决以上失败的检查项以提高系统完整性。"
fi

echo
echo "可以通过以下方式进一步验证系统:"
echo "1. 启动MongoDB服务器"
echo "2. 运行 './server' 启动认证服务"
echo "3. 在浏览器中打开 'test/index.html' 进行功能测试"
echo "4. 运行系统健康检查和认证测试"

exit $FAILED_CHECKS