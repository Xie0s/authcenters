// 使用 Node.js 执行 MongoDB 初始化脚本
const { MongoClient, ObjectId } = require('mongodb');

// 连接 URL
const url = 'mongodb://localhost:27017';
const dbName = 'auth_center';

async function initializeDB() {
  let client;
  try {
    console.log('正在连接到 MongoDB...');
    client = await MongoClient.connect(url);
    console.log('MongoDB 连接成功！');
    
    const db = client.db(dbName);
    console.log('开始初始化数据库...');
    
    // 2. 初始化权限数据
    console.log("正在初始化权限数据...");
    // 检查是否已存在权限数据
    const permCount = await db.collection('permissions').countDocuments();
    if (permCount > 0) {
      console.log(`权限数据已存在 (${permCount} 条记录)，跳过初始化`);
    } else {
      await db.collection('permissions').deleteMany({});  // 清空现有数据
      
      const permissionsList = [
        // 知识库内容权限
        { 
          name: "KNOWLEDGE_READ", 
          resource: "knowledge", 
          action: "READ", 
          description: "查看知识库文档", 
          category: "knowledge_content",
          created_at: new Date()
        },
        { 
          name: "KNOWLEDGE_CREATE", 
          resource: "knowledge", 
          action: "CREATE", 
          description: "创建知识库文档", 
          category: "knowledge_content",
          created_at: new Date()
        },
        { 
          name: "KNOWLEDGE_UPDATE", 
          resource: "knowledge", 
          action: "UPDATE", 
          description: "编辑知识库文档", 
          category: "knowledge_content",
          created_at: new Date()
        },
        { 
          name: "KNOWLEDGE_DELETE", 
          resource: "knowledge", 
          action: "DELETE", 
          description: "删除知识库文档", 
          category: "knowledge_content",
          created_at: new Date()
        },
        { 
          name: "KNOWLEDGE_PUBLISH", 
          resource: "knowledge", 
          action: "PUBLISH", 
          description: "发布知识库文档", 
          category: "knowledge_content",
          created_at: new Date()
        },
        { 
          name: "KNOWLEDGE_APPROVE", 
          resource: "knowledge", 
          action: "APPROVE", 
          description: "审核知识库内容", 
          category: "knowledge_content",
          created_at: new Date()
        },
        
        // 系统管理权限
        { 
          name: "USER_MANAGE", 
          resource: "user", 
          action: "MANAGE", 
          description: "用户管理", 
          category: "system_management",
          created_at: new Date()
        },
        { 
          name: "ROLE_MANAGE", 
          resource: "role", 
          action: "MANAGE", 
          description: "角色管理", 
          category: "system_management",
          created_at: new Date()
        },
        { 
          name: "CATEGORY_MANAGE", 
          resource: "category", 
          action: "MANAGE", 
          description: "分类管理（层级式）", 
          category: "system_management",
          created_at: new Date()
        },
        { 
          name: "TAG_CREATE", 
          resource: "tag", 
          action: "CREATE", 
          description: "创建标签（灵活标记）", 
          category: "content_organization",
          created_at: new Date()
        },
        { 
          name: "TAG_MANAGE", 
          resource: "tag", 
          action: "MANAGE", 
          description: "标签管理（编辑、删除）", 
          category: "content_organization",
          created_at: new Date()
        },
        { 
          name: "SYSTEM_CONFIG", 
          resource: "system", 
          action: "CONFIG", 
          description: "系统配置", 
          category: "system_management",
          created_at: new Date()
        },
        
        // 交互功能权限
        { 
          name: "COMMENT", 
          resource: "knowledge", 
          action: "COMMENT", 
          description: "评论文档", 
          category: "interaction",
          created_at: new Date()
        },
        { 
          name: "FAVORITE", 
          resource: "knowledge", 
          action: "FAVORITE", 
          description: "收藏文档", 
          category: "interaction",
          created_at: new Date()
        },
        { 
          name: "SEARCH", 
          resource: "knowledge", 
          action: "SEARCH", 
          description: "搜索文档", 
          category: "interaction",
          created_at: new Date()
        },
        { 
          name: "AI_ASSISTANT", 
          resource: "ai", 
          action: "USE", 
          description: "使用AI助手功能", 
          category: "interaction",
          created_at: new Date()
        }
      ];
      
      const permissionsResult = await db.collection('permissions').insertMany(permissionsList);
      console.log(`权限数据初始化完成，共插入 ${permissionsResult.insertedCount} 条权限记录`);
    }
    
    // 3. 获取权限并创建权限映射
    console.log("正在创建权限映射...");
    const permissions = await db.collection('permissions').find({}).toArray();
    const permMap = {};
    permissions.forEach((perm) => {
      permMap[perm.name] = perm;
    });
    
    // 4. 定义角色权限映射
    const rolePermissions = {
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
    
    // 5. 初始化角色数据
    console.log("正在初始化角色数据...");
    const roleCount = await db.collection('roles').countDocuments();
    if (roleCount > 0) {
      console.log(`角色数据已存在 (${roleCount} 条记录)，跳过初始化`);
    } else {
      await db.collection('roles').deleteMany({});  // 清空现有数据
      
      const roles = [
        {
          name: "Admin",
          display_name: "系统管理员",
          description: "拥有最高权限，可管理所有系统功能",
          level: 4,
          status: "active",
          permissions: rolePermissions["Admin"].map((name) => {
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
          permissions: rolePermissions["Editor"].map((name) => {
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
          permissions: rolePermissions["Author"].map((name) => {
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
          permissions: rolePermissions["User"].map((name) => {
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
      ];
      
      const rolesResult = await db.collection('roles').insertMany(roles);
      console.log(`角色数据初始化完成，共插入 ${rolesResult.insertedCount} 条角色记录`);
    }
    
    // 6. 初始化示例分类数据
    console.log("正在初始化分类数据...");
    const categoryCount = await db.collection('categories').countDocuments();
    if (categoryCount > 0) {
      console.log(`分类数据已存在 (${categoryCount} 条记录)，跳过初始化`);
    } else {
      await db.collection('categories').deleteMany({});  // 清空现有数据
      
      const categoriesResult = await db.collection('categories').insertMany([
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
        },
        {
          name: "管理制度",
          parent_id: null,
          path: "/管理制度",
          level: 0,
          sort_order: 3,
          description: "公司管理制度文档",
          status: "active",
          children: [],
          document_count: 0,
          created_at: new Date(),
          updated_at: new Date()
        }
      ]);
      
      console.log(`分类数据初始化完成，共插入 ${categoriesResult.insertedCount} 条分类记录`);
    }
    
    // 7. 初始化示例标签数据
    console.log("正在初始化标签数据...");
    const tagCount = await db.collection('tags').countDocuments();
    if (tagCount > 0) {
      console.log(`标签数据已存在 (${tagCount} 条记录)，跳过初始化`);
    } else {
      await db.collection('tags').deleteMany({});  // 清空现有数据
      
      const adminUserId = new ObjectId(); // 假设的管理员ID
      
      const tagsResult = await db.collection('tags').insertMany([
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
        },
        {
          name: "微服务",
          color: "#FF6B6B",
          description: "微服务架构相关",
          usage_count: 0,
          created_by: adminUserId,
          created_by_name: "系统管理员",
          related_tags: ["Docker", "Kubernetes", "API"],
          created_at: new Date(),
          last_used_at: new Date()
        },
        {
          name: "API设计",
          color: "#4ECDC4",
          description: "API设计和开发",
          usage_count: 0,
          created_by: adminUserId,
          created_by_name: "系统管理员",
          related_tags: ["RESTful", "GraphQL", "文档"],
          created_at: new Date(),
          last_used_at: new Date()
        }
      ]);
      
      console.log(`标签数据初始化完成，共插入 ${tagsResult.insertedCount} 条标签记录`);
    }
    
    // 8. 创建索引
    console.log("正在创建数据库索引...");
    
    // 用户索引
    await db.collection('users').createIndex({ "username": 1 }, { unique: true });
    await db.collection('users').createIndex({ "email": 1 }, { unique: true, sparse: true });
    await db.collection('users').createIndex({ "phone": 1 }, { unique: true, sparse: true });
    await db.collection('users').createIndex({ "roles.role_id": 1 });
    await db.collection('users').createIndex({ "status": 1 });
    
    // 角色索引
    await db.collection('roles').createIndex({ "name": 1 }, { unique: true });
    await db.collection('roles').createIndex({ "level": 1 });
    await db.collection('roles').createIndex({ "permissions.name": 1 });
    
    // 权限索引
    await db.collection('permissions').createIndex({ "name": 1 }, { unique: true });
    await db.collection('permissions').createIndex({ "resource": 1, "action": 1 });
    await db.collection('permissions').createIndex({ "category": 1 });
    
    // 分类索引
    await db.collection('categories').createIndex({ "parent_id": 1 });
    await db.collection('categories').createIndex({ "path": 1 });
    await db.collection('categories').createIndex({ "level": 1, "sort_order": 1 });
    await db.collection('categories').createIndex({ "status": 1 });
    
    // 标签索引
    await db.collection('tags').createIndex({ "name": 1 }, { unique: true });
    await db.collection('tags').createIndex({ "created_by": 1 });
    await db.collection('tags').createIndex({ "usage_count": -1 });
    await db.collection('tags').createIndex({ "last_used_at": -1 });
    
    // 会话索引
    await db.collection('sessions').createIndex({ "session_id": 1 }, { unique: true });
    await db.collection('sessions').createIndex({ "user_id": 1 });
    await db.collection('sessions').createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 }); // TTL索引
    await db.collection('sessions').createIndex({ "is_revoked": 1 });
    
    console.log("数据库索引创建完成");
    
    // 9. 显示初始化结果
    console.log("\n=== 数据库初始化完成 ===");
    console.log(`权限数量: ${await db.collection('permissions').countDocuments()}`);
    console.log(`角色数量: ${await db.collection('roles').countDocuments()}`);
    console.log(`分类数量: ${await db.collection('categories').countDocuments()}`);
    console.log(`标签数量: ${await db.collection('tags').countDocuments()}`);
    console.log(`用户数量: ${await db.collection('users').countDocuments()}`);
    
    console.log("\n角色权限分配:");
    const rolesList = await db.collection('roles').find({}, {projection: {name: 1, display_name: 1, permissions: 1}}).toArray();
    rolesList.forEach(role => {
      console.log(`- ${role.display_name} (${role.name}): ${role.permissions.length} 个权限`);
    });
    
    console.log("\n数据库初始化脚本执行完成！");
    console.log("您现在可以启动AuthCenter服务并开始使用完整的RBAC功能。");
    
  } catch (err) {
    console.error('执行脚本时出错:', err);
  } finally {
    if (client) {
      await client.close();
    }
  }
}

// 执行初始化
initializeDB();
