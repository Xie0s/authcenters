#!/bin/bash

# MongoDB初始化脚本
# 用于创建初始数据和索引

echo "正在连接MongoDB并初始化数据..."

# 检查MongoDB连接
mongo --eval "db.runCommand('ping')" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "错误: 无法连接到MongoDB"
    exit 1
fi

echo "MongoDB连接成功，开始初始化..."

# 执行初始化脚本
mongo auth_center <<EOF

// 删除现有集合（如果存在）
print("清理现有数据...");
db.permissions.drop();
db.roles.drop();
db.categories.drop();
db.tags.drop();

// 插入权限数据
print("插入权限数据...");
db.permissions.insertMany([
  // 知识库内容权限
  { name: "KNOWLEDGE_READ", resource: "knowledge", action: "READ", description: "查看知识库文档", category: "knowledge_content", created_at: new Date() },
  { name: "KNOWLEDGE_CREATE", resource: "knowledge", action: "CREATE", description: "创建知识库文档", category: "knowledge_content", created_at: new Date() },
  { name: "KNOWLEDGE_UPDATE", resource: "knowledge", action: "UPDATE", description: "编辑知识库文档", category: "knowledge_content", created_at: new Date() },
  { name: "KNOWLEDGE_DELETE", resource: "knowledge", action: "DELETE", description: "删除知识库文档", category: "knowledge_content", created_at: new Date() },
  { name: "KNOWLEDGE_PUBLISH", resource: "knowledge", action: "PUBLISH", description: "发布知识库文档", category: "knowledge_content", created_at: new Date() },
  { name: "KNOWLEDGE_APPROVE", resource: "knowledge", action: "APPROVE", description: "审核知识库内容", category: "knowledge_content", created_at: new Date() },
  
  // 系统管理权限
  { name: "USER_MANAGE", resource: "user", action: "MANAGE", description: "用户管理", category: "system_management", created_at: new Date() },
  { name: "ROLE_MANAGE", resource: "role", action: "MANAGE", description: "角色管理", category: "system_management", created_at: new Date() },
  { name: "CATEGORY_MANAGE", resource: "category", action: "MANAGE", description: "分类管理（层级式）", category: "system_management", created_at: new Date() },
  { name: "TAG_CREATE", resource: "tag", action: "CREATE", description: "创建标签（灵活标记）", category: "content_organization", created_at: new Date() },
  { name: "TAG_MANAGE", resource: "tag", action: "MANAGE", description: "标签管理（编辑、删除）", category: "content_organization", created_at: new Date() },
  { name: "SYSTEM_CONFIG", resource: "system", action: "CONFIG", description: "系统配置", category: "system_management", created_at: new Date() },
  
  // 交互功能权限
  { name: "COMMENT", resource: "knowledge", action: "COMMENT", description: "评论文档", category: "interaction", created_at: new Date() },
  { name: "FAVORITE", resource: "knowledge", action: "FAVORITE", description: "收藏文档", category: "interaction", created_at: new Date() },
  { name: "SEARCH", resource: "knowledge", action: "SEARCH", description: "搜索文档", category: "interaction", created_at: new Date() },
  { name: "AI_ASSISTANT", resource: "ai", action: "USE", description: "使用AI助手功能", category: "interaction", created_at: new Date() }
]);

print("权限数据插入完成");

// 获取权限ID映射
print("创建权限映射...");
var permissions = db.permissions.find({}).toArray();
var permMap = {};
permissions.forEach(function(perm) {
    permMap[perm.name] = perm;
});

// 插入角色数据
print("插入角色数据...");

// 定义角色权限映射
var rolePermissions = {
  "Admin": [
    "KNOWLEDGE_READ", "KNOWLEDGE_CREATE", "KNOWLEDGE_UPDATE", "KNOWLEDGE_DELETE", 
    "KNOWLEDGE_PUBLISH", "KNOWLEDGE_APPROVE", "USER_MANAGE", "ROLE_MANAGE", 
    "CATEGORY_MANAGE", "TAG_CREATE", "TAG_MANAGE", "SYSTEM_CONFIG", 
    "COMMENT", "FAVORITE", "SEARCH", "AI_ASSISTANT"
  ],
  "Editor": [
    "KNOWLEDGE_READ", "KNOWLEDGE_CREATE", "KNOWLEDGE_UPDATE", "KNOWLEDGE_DELETE", 
    "KNOWLEDGE_PUBLISH", "CATEGORY_MANAGE", "TAG_CREATE", 
    "COMMENT", "FAVORITE", "SEARCH", "AI_ASSISTANT"
  ],
  "Author": [
    "KNOWLEDGE_READ", "KNOWLEDGE_CREATE", "KNOWLEDGE_UPDATE", "TAG_CREATE",
    "COMMENT", "FAVORITE", "SEARCH", "AI_ASSISTANT"
  ],
  "User": [
    "KNOWLEDGE_READ", "COMMENT", "FAVORITE", "SEARCH", "AI_ASSISTANT"
  ]
};

// 插入角色
db.roles.insertMany([
  {
    name: "Admin",
    display_name: "系统管理员",
    description: "拥有最高权限，可管理所有系统功能",
    level: 4,
    status: "active",
    permissions: rolePermissions["Admin"].map(function(name) {
      return {
        permission_id: permMap[name]._id,
        name: name,
        resource: permMap[name].resource,
        action: permMap[name].action
      };
    }),
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Editor",
    display_name: "内容管理员",
    description: "负责知识库内容的全面管理",
    level: 3,
    status: "active",
    permissions: rolePermissions["Editor"].map(function(name) {
      return {
        permission_id: permMap[name]._id,
        name: name,
        resource: permMap[name].resource,
        action: permMap[name].action
      };
    }),
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "Author",
    display_name: "内容创作者",
    description: "专注于知识库内容的创作和编辑",
    level: 2,
    status: "active",
    permissions: rolePermissions["Author"].map(function(name) {
      return {
        permission_id: permMap[name]._id,
        name: name,
        resource: permMap[name].resource,
        action: permMap[name].action
      };
    }),
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "User",
    display_name: "普通用户",
    description: "知识库的日常使用者",
    level: 1,
    status: "active",
    permissions: rolePermissions["User"].map(function(name) {
      return {
        permission_id: permMap[name]._id,
        name: name,
        resource: permMap[name].resource,
        action: permMap[name].action
      };
    }),
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print("角色数据插入完成");

// 插入示例分类数据
print("插入分类数据...");
db.categories.insertMany([
  {
    name: "技术文档",
    parent_id: null,
    path: "/技术文档",
    level: 0,
    sort_order: 1,
    description: "技术相关的知识文档",
    status: "active",
    children: [],
    document_count: 0,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "产品文档",
    parent_id: null,
    path: "/产品文档",
    level: 0,
    sort_order: 2,
    description: "产品相关的文档",
    status: "active",
    children: [],
    document_count: 0,
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print("分类数据插入完成");

// 插入示例标签数据
print("插入标签数据...");
var adminUserId = ObjectId(); // 假设的管理员ID

db.tags.insertMany([
  {
    name: "Go语言",
    color: "#00ADD8",
    description: "Go编程语言相关内容",
    usage_count: 0,
    created_by: adminUserId,
    created_by_name: "系统管理员",
    related_tags: ["后端开发", "微服务", "高性能"],
    created_at: new Date(),
    last_used_at: new Date()
  },
  {
    name: "前端开发",
    color: "#61DAFB",
    description: "前端开发技术",
    usage_count: 0,
    created_by: adminUserId,
    created_by_name: "系统管理员",
    related_tags: ["JavaScript", "React", "Vue"],
    created_at: new Date(),
    last_used_at: new Date()
  },
  {
    name: "数据库",
    color: "#336791",
    description: "数据库相关技术",
    usage_count: 0,
    created_by: adminUserId,
    created_by_name: "系统管理员",
    related_tags: ["MongoDB", "MySQL", "NoSQL"],
    created_at: new Date(),
    last_used_at: new Date()
  }
]);

print("标签数据插入完成");

print("数据库初始化完成！");

EOF

echo "MongoDB初始化完成！"
