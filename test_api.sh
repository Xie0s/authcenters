#!/bin/bash

# API测试脚本
BASE_URL="http://localhost:8080/api/v1/auth"

echo "=== AuthCenter API 测试 ==="
echo "基础URL: $BASE_URL"
echo

# 测试1: 注册用户（包含用户名）
echo "--- 测试1: 注册用户 ---"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "email": "testuser123@example.com",
    "password": "password123"
  }')

echo "注册响应: $REGISTER_RESPONSE"
echo

# 测试2: 用户名登录
echo "--- 测试2: 用户名登录 ---"
USERNAME_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "password": "password123",
    "type": "username"
  }')

echo "用户名登录响应: $USERNAME_LOGIN_RESPONSE"
echo

# 测试3: 邮箱登录
echo "--- 测试3: 邮箱登录 ---"
EMAIL_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser123@example.com",
    "password": "password123",
    "type": "email"
  }')

echo "邮箱登录响应: $EMAIL_LOGIN_RESPONSE"
echo

# 测试4: 自动识别登录（用户名）
echo "--- 测试4: 自动识别登录（用户名） ---"
AUTO_USERNAME_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "password": "password123",
    "type": "auto"
  }')

echo "自动识别用户名登录响应: $AUTO_USERNAME_LOGIN_RESPONSE"
echo

# 测试5: 自动识别登录（邮箱）
echo "--- 测试5: 自动识别登录（邮箱） ---"
AUTO_EMAIL_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser123@example.com",
    "password": "password123",
    "type": "auto"
  }')

echo "自动识别邮箱登录响应: $AUTO_EMAIL_LOGIN_RESPONSE"
echo

# 测试6: 错误情况 - 空用户名注册
echo "--- 测试6: 错误情况 - 空用户名注册 ---"
EMPTY_USERNAME_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "empty@example.com",
    "password": "password123"
  }')

echo "空用户名注册响应: $EMPTY_USERNAME_RESPONSE"
echo

# 测试7: 错误情况 - 重复用户名
echo "--- 测试7: 错误情况 - 重复用户名 ---"
DUPLICATE_USERNAME_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser123",
    "email": "another@example.com",
    "password": "password123"
  }')

echo "重复用户名注册响应: $DUPLICATE_USERNAME_RESPONSE"
echo

echo "=== 测试完成 ==="
