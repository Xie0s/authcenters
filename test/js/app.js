document.addEventListener('DOMContentLoaded', async () => {
    // 初始化事件处理器
    initEventHandlers();
    
    // 检查并刷新Token状态
    console.log('页面加载完成，开始检查Token状态...');
    await checkAndRefreshToken();
    
    // 初始化Bootstrap的Tab事件
    const triggerTabList = document.querySelectorAll('#myTab a');
    triggerTabList.forEach(tabEl => {
        tabEl.addEventListener('click', e => {
            e.preventDefault();
            new bootstrap.Tab(tabEl).show();
        });
    });
});

// 初始化所有事件处理器
function initEventHandlers() {
    // 认证相关
    initAuthHandlers();
    
    // 用户相关
    initUserHandlers();
    
    // 角色相关
    initRoleHandlers();
    
    // 权限相关
    initPermissionHandlers();
}

// 认证相关事件处理器
function initAuthHandlers() {
    // 注册表单
    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('regUsername').value;
        const email = document.getElementById('regEmail').value;
        const password = document.getElementById('regPassword').value;
        
        const response = await authAPI.register(username, email, password);
        document.getElementById('registerResponse').textContent = formatJSON(response);
    });
    
    // 登录表单
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const identifier = document.getElementById('loginIdentifier').value;
        const password = document.getElementById('loginPassword').value;
        const type = document.getElementById('loginType').value;

        const response = await authAPI.login(identifier, password, type);
        document.getElementById('loginResponse').textContent = formatJSON(response);
    });
    
    // 刷新Token按钮
    document.getElementById('refreshToken').addEventListener('click', async () => {
        console.log('=== 刷新Token按钮被点击 ===');
        const response = await authAPI.refreshToken();
        console.log('=== 刷新Token响应 ===', response);
        document.getElementById('refreshResponse').textContent = formatJSON(response);
    });
    
    // 验证Token按钮
    document.getElementById('verifyToken').addEventListener('click', async () => {
        const response = await authAPI.verifyToken();
        document.getElementById('verifyResponse').textContent = formatJSON(response);
    });
}

// 登出
function logout() {
    authAPI.logout().then(response => {
        alert('已成功登出');
    }).catch(error => {
        console.error('登出失败:', error);
    });
}

// 调试函数
function showDebugInfo() {
    const debugInfo = document.getElementById('debugInfo');
    const debugContent = document.getElementById('debugContent');
    
    // 从localStorage读取
    const accessToken = localStorage.getItem('access_token');
    const refreshToken = localStorage.getItem('refresh_token');
    const userId = localStorage.getItem('user_id');
    
    // 从sessionStorage读取
    const sessionAccessToken = sessionStorage.getItem('access_token');
    const sessionRefreshToken = sessionStorage.getItem('refresh_token');
    const sessionUserId = sessionStorage.getItem('user_id');
    
    console.log('=== 调试信息 ===');
    console.log('localStorage Access Token:', accessToken ? '存在 (长度: ' + accessToken.length + ')' : '不存在');
    console.log('localStorage Refresh Token:', refreshToken ? '存在 (长度: ' + refreshToken.length + ')' : '不存在');
    console.log('localStorage User ID:', userId || '不存在');
    console.log('sessionStorage Access Token:', sessionAccessToken ? '存在 (长度: ' + sessionAccessToken.length + ')' : '不存在');
    console.log('sessionStorage Refresh Token:', sessionRefreshToken ? '存在 (长度: ' + sessionRefreshToken.length + ')' : '不存在');
    console.log('sessionStorage User ID:', sessionUserId || '不存在');
    
    debugContent.innerHTML = `
        <h5>localStorage 存储:</h5>
        <div><strong>Access Token:</strong> ${accessToken ? '存在 (长度: ' + accessToken.length + ')' : '不存在'}</div>
        <div><strong>Refresh Token:</strong> ${refreshToken ? '存在 (长度: ' + refreshToken.length + ')' : '不存在'}</div>
        <div><strong>User ID:</strong> ${userId || '不存在'}</div>
        <h5>sessionStorage 存储:</h5>
        <div><strong>Access Token:</strong> ${sessionAccessToken ? '存在 (长度: ' + sessionAccessToken.length + ')' : '不存在'}</div>
        <div><strong>Refresh Token:</strong> ${sessionRefreshToken ? '存在 (长度: ' + sessionRefreshToken.length + ')' : '不存在'}</div>
        <div><strong>User ID:</strong> ${sessionUserId || '不存在'}</div>
        <div><strong>当前时间:</strong> ${new Date().toLocaleString()}</div>
    `;
    
    debugInfo.style.display = 'block';
}

function toggleDebug() {
    const debugInfo = document.getElementById('debugInfo');
    debugInfo.style.display = debugInfo.style.display === 'none' ? 'block' : 'none';
}

function clearAllTokens() {
    if (confirm('确定要清除所有Token吗？')) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_id');
        
        // 重新初始化
        initTokenState();
        checkToken();
        
        alert('所有Token已清除');
    }
}

// 用户相关事件处理器
function initUserHandlers() {
    // 获取用户列表
    document.getElementById('getUsersForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const page = document.getElementById('usersPage').value;
        const limit = document.getElementById('usersLimit').value;
        
        const response = await userAPI.getUsers(page, limit);
        document.getElementById('getUsersResponse').textContent = formatJSON(response);
    });
    
    // 获取单个用户
    document.getElementById('getUserForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = document.getElementById('userId').value;
        
        const response = await userAPI.getUser(userId);
        document.getElementById('getUserResponse').textContent = formatJSON(response);
    });
    
    // 更新用户
    document.getElementById('updateUserForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = document.getElementById('updateUserId').value;
        const username = document.getElementById('updateUsername').value;
        const email = document.getElementById('updateEmail').value;
        
        const userData = {};
        if (username) userData.username = username;
        if (email) userData.email = email;
        
        const response = await userAPI.updateUser(userId, userData);
        document.getElementById('updateUserResponse').textContent = formatJSON(response);
    });
    
    // 删除用户
    document.getElementById('deleteUserForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = document.getElementById('deleteUserId').value;
        
        if (confirm('确定要删除此用户吗？此操作不可撤销！')) {
            const response = await userAPI.deleteUser(userId);
            document.getElementById('deleteUserResponse').textContent = formatJSON(response);
        }
    });
    
    // 分配角色
    document.getElementById('assignRoleForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = document.getElementById('roleUserId').value;
        const roleId = document.getElementById('roleId').value;
        
        const response = await userAPI.assignRole(userId, roleId);
        document.getElementById('assignRoleResponse').textContent = formatJSON(response);
    });
    
    // 移除角色
    document.getElementById('removeRoleForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = document.getElementById('removeRoleUserId').value;
        const roleId = document.getElementById('removeRoleId').value;
        
        const response = await userAPI.removeRole(userId, roleId);
        document.getElementById('removeRoleResponse').textContent = formatJSON(response);
    });
    
    // 获取用户权限
    document.getElementById('getUserPermissionsForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = document.getElementById('permUserId').value;
        
        const response = await userAPI.getUserPermissions(userId);
        document.getElementById('getUserPermissionsResponse').textContent = formatJSON(response);
    });
}

// 角色相关事件处理器
function initRoleHandlers() {
    // 获取角色列表
    document.getElementById('getRolesForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const page = document.getElementById('rolesPage').value;
        const limit = document.getElementById('rolesLimit').value;
        
        const response = await roleAPI.getRoles(page, limit);
        document.getElementById('getRolesResponse').textContent = formatJSON(response);
    });
    
    // 创建角色
    document.getElementById('createRoleForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('createRoleName').value;
        const description = document.getElementById('createRoleDesc').value;
        
        const response = await roleAPI.createRole(name, description);
        document.getElementById('createRoleResponse').textContent = formatJSON(response);
    });
    
    // 获取角色详情
    document.getElementById('getRoleForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const roleId = document.getElementById('getRoleId').value;
        
        const response = await roleAPI.getRole(roleId);
        document.getElementById('getRoleResponse').textContent = formatJSON(response);
    });
    
    // 更新角色
    document.getElementById('updateRoleForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const roleId = document.getElementById('updateRoleId').value;
        const name = document.getElementById('updateRoleName').value;
        const description = document.getElementById('updateRoleDesc').value;
        
        const roleData = {};
        if (name) roleData.name = name;
        if (description) roleData.description = description;
        
        const response = await roleAPI.updateRole(roleId, roleData);
        document.getElementById('updateRoleResponse').textContent = formatJSON(response);
    });
    
    // 删除角色
    document.getElementById('deleteRoleForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const roleId = document.getElementById('deleteRoleId').value;
        
        if (confirm('确定要删除此角色吗？此操作不可撤销！')) {
            const response = await roleAPI.deleteRole(roleId);
            document.getElementById('deleteRoleResponse').textContent = formatJSON(response);
        }
    });
    
    // 分配权限
    document.getElementById('assignPermissionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const roleId = document.getElementById('permRoleId').value;
        const permissionId = document.getElementById('permissionId').value;
        
        const response = await roleAPI.assignPermission(roleId, permissionId);
        document.getElementById('assignPermissionResponse').textContent = formatJSON(response);
    });
    
    // 移除权限
    document.getElementById('removePermissionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const roleId = document.getElementById('removePermRoleId').value;
        const permissionId = document.getElementById('removePermissionId').value;
        
        const response = await roleAPI.removePermission(roleId, permissionId);
        document.getElementById('removePermissionResponse').textContent = formatJSON(response);
    });
}

// 权限相关事件处理器
function initPermissionHandlers() {
    // 获取权限列表
    document.getElementById('getPermissions').addEventListener('click', async () => {
        const response = await permissionAPI.getPermissions();
        document.getElementById('getPermissionsResponse').textContent = formatJSON(response);
    });
    
    // 创建权限
    document.getElementById('createPermissionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('createPermName').value;
        const description = document.getElementById('createPermDesc').value;
        const resource = document.getElementById('createPermResource').value;
        const action = document.getElementById('createPermAction').value;
        
        const response = await permissionAPI.createPermission(name, description, resource, action);
        document.getElementById('createPermissionResponse').textContent = formatJSON(response);
    });
    
    // 获取权限详情
    document.getElementById('getPermissionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const permId = document.getElementById('getPermId').value;
        
        const response = await permissionAPI.getPermission(permId);
        document.getElementById('getPermissionResponse').textContent = formatJSON(response);
    });
    
    // 更新权限
    document.getElementById('updatePermissionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const permId = document.getElementById('updatePermId').value;
        const name = document.getElementById('updatePermName').value;
        const description = document.getElementById('updatePermDesc').value;
        const resource = document.getElementById('updatePermResource').value;
        const action = document.getElementById('updatePermAction').value;
        
        const permData = {};
        if (name) permData.name = name;
        if (description) permData.description = description;
        if (resource) permData.resource = resource;
        if (action) permData.action = action;
        
        const response = await permissionAPI.updatePermission(permId, permData);
        document.getElementById('updatePermissionResponse').textContent = formatJSON(response);
    });
    
    // 删除权限
    document.getElementById('deletePermissionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const permId = document.getElementById('deletePermId').value;
        
        if (confirm('确定要删除此权限吗？此操作不可撤销！')) {
            const response = await permissionAPI.deletePermission(permId);
            document.getElementById('deletePermissionResponse').textContent = formatJSON(response);
        }
    });
}
