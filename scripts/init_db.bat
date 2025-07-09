@echo off
echo MongoDB初始化脚本 (Windows版本)

REM 检查MongoDB连接
echo 正在检查MongoDB连接...
mongo --eval "db.runCommand('ping')" >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 无法连接到MongoDB
    exit /b 1
)

echo MongoDB连接成功，开始初始化...

REM 执行初始化脚本
mongo auth_center < "%~dp0init_data.js"

echo MongoDB初始化完成！
pause
