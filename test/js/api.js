// API基础URL
const API_BASE_URL = 'http://localhost:8080';  // 后端服务地址

// 存储Token - 从localStorage初始化
let accessToken = null;
let refreshToken = null;
let userId = null;

// 判断浏览器是否支持localStorage
const isLocalStorageSupported = (() => {
    try {
        const testKey = '__test__';
        localStorage.setItem(testKey, testKey);
        const result = localStorage.getItem(testKey) === testKey;
        localStorage.removeItem(testKey);
        return result;
    } catch (e) {
        console.error('localStorage不被支持:', e);
        return false;
    }
})();

// 初始化Token状态
function initTokenState() {
    if (!isLocalStorageSupported) {
        console.error('浏览器不支持localStorage，无法初始化Token状态');
        return;
    }
    
    try {
        accessToken = localStorage.getItem('access_token') || null;
        refreshToken = localStorage.getItem('refresh_token') || null;
        userId = localStorage.getItem('user_id') || null;
        
        // 使用sessionStorage作为备份机制
        if (!accessToken) {
            accessToken = sessionStorage.getItem('access_token') || null;
        }
        if (!refreshToken) {
            refreshToken = sessionStorage.getItem('refresh_token') || null;
        }
        if (!userId) {
            userId = sessionStorage.getItem('user_id') || null;
        }
        
        console.log('Token状态初始化:', {
            hasAccessToken: !!accessToken,
            hasRefreshToken: !!refreshToken,
            userId: userId
        });
    } catch (e) {
        console.error('初始化Token状态出错:', e);
    }
}

// 页面加载时初始化
initTokenState();

// 检查Token并更新UI
function checkToken() {
    try {
        // 重新从localStorage读取token，确保最新状态
        if (isLocalStorageSupported) {
            accessToken = localStorage.getItem('access_token') || null;
            refreshToken = localStorage.getItem('refresh_token') || null;
            userId = localStorage.getItem('user_id') || null;
        }
        
        // 如果localStorage没有，则从sessionStorage尝试获取
        if (!accessToken) {
            accessToken = sessionStorage.getItem('access_token') || null;
        }
        if (!refreshToken) {
            refreshToken = sessionStorage.getItem('refresh_token') || null;
        }
        if (!userId) {
            userId = sessionStorage.getItem('user_id') || null;
        }
        
        console.log('checkToken: access_token:', accessToken ? '存在' : '不存在');
        console.log('checkToken: refresh_token:', refreshToken ? '存在' : '不存在');
        
        if (accessToken) {
            document.getElementById('tokenStatus').style.display = 'block';
            document.getElementById('tokenInfo').textContent = `用户ID: ${userId || 'Unknown'}`;
            return true;
        } else {
            document.getElementById('tokenStatus').style.display = 'none';
            return false;
        }
    } catch (e) {
        console.error('检查Token状态出错:', e);
        document.getElementById('tokenStatus').style.display = 'none';
        return false;
    }
}

// 检查Token有效性并自动刷新
async function checkAndRefreshToken() {
    console.log('检查Token状态...');

    try {
        // 重新从localStorage读取最新的token
        if (isLocalStorageSupported) {
            accessToken = localStorage.getItem('access_token') || null;
            refreshToken = localStorage.getItem('refresh_token') || null;
            userId = localStorage.getItem('user_id') || null;
        }
        
        // 如果localStorage没有，则从sessionStorage尝试获取
        if (!accessToken) {
            accessToken = sessionStorage.getItem('access_token') || null;
        }
        if (!refreshToken) {
            refreshToken = sessionStorage.getItem('refresh_token') || null;
        }
        if (!userId) {
            userId = sessionStorage.getItem('user_id') || null;
        }
    } catch (e) {
        console.error('读取Token状态出错:', e);
    }

    if (!accessToken) {
        console.log('没有access_token');
        checkToken();
        return false;
    }

    // 尝试验证当前Token
    try {
        const verifyResponse = await authAPI.verifyToken();
        if (verifyResponse.status === 200) {
            console.log('Token验证成功');
            checkToken();
            return true;
        }
    } catch (error) {
        console.log('Token验证失败:', error);
    }

    // 如果验证失败，尝试刷新Token
    const currentRefreshToken = localStorage.getItem('refresh_token') || sessionStorage.getItem('refresh_token');
    if (currentRefreshToken) {
        console.log('尝试刷新Token...');
        try {
            const refreshResponse = await authAPI.refreshToken();
            if (refreshResponse.status === 200) {
                console.log('Token刷新成功');
                checkToken();
                return true;
            } else {
                console.log('Token刷新失败:', refreshResponse);
            }
        } catch (error) {
            console.log('Token刷新异常:', error);
        }
    }

    // 如果都失败了，清除所有Token
    console.log('清除所有Token');
    accessToken = null;
    refreshToken = null;
    userId = null;

    try {
        if (isLocalStorageSupported) {
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            localStorage.removeItem('user_id');
        }
        
        // 同时清除sessionStorage
        sessionStorage.removeItem('access_token');
        sessionStorage.removeItem('refresh_token');
        sessionStorage.removeItem('user_id');
    } catch (e) {
        console.error('清除Token出错:', e);
    }

    checkToken();
    return false;
}

// 通用API请求函数
async function apiRequest(endpoint, method = 'GET', data = null, withAuth = true) {
    console.log('=== apiRequest 开始 ===');
    console.log('endpoint:', endpoint);
    console.log('method:', method);
    console.log('data:', data);
    console.log('withAuth:', withAuth);
    
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = {
        'Content-Type': 'application/json',
    };
    
    // 添加授权头
    if (withAuth && accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
    }
    
    const options = {
        method,
        headers,
        credentials: 'include',
        mode: 'cors', // 启用CORS
    };
    
    if (data && (method === 'POST' || method === 'PUT')) {
        options.body = JSON.stringify(data);
    }
    
    console.log('请求URL:', url);
    console.log('请求选项:', options);
    
    try {
        console.log('开始发送请求...');
        const response = await fetch(url, options);
        console.log('收到响应，状态:', response.status);
        const responseData = await response.json();
        console.log('响应数据:', responseData);
        
        // 处理Token过期
        if (response.status === 401 && endpoint !== '/api/v1/auth/refresh' && endpoint !== '/api/v1/auth/verify') {
            // 重新获取最新的refresh_token
            const currentRefreshToken = localStorage.getItem('refresh_token');
            if (currentRefreshToken) {
                console.log('检测到401错误，尝试刷新Token...');
                const refreshed = await refreshTokenFunc();
                if (refreshed) {
                    console.log('Token刷新成功，重试原请求...');
                    // 重试请求，需要重新设置Authorization头
                    const retryHeaders = { ...headers };
                    if (withAuth && accessToken) {
                        retryHeaders['Authorization'] = `Bearer ${accessToken}`;
                    }
                    return apiRequest(endpoint, method, data, withAuth);
                } else {
                    console.log('Token刷新失败，返回401错误');
                }
            }
        }
        
        return {
            status: response.status,
            data: responseData
        };
    } catch (error) {
        console.error('API请求错误:', error);
        return {
            status: 500,
            data: {
                code: 500,
                message: `请求错误: ${error.message}`,
                error: error.message
            }
        };
    }
}

// 认证API
const authAPI = {
    // 注册
    register: (username, email, password) => {
        return apiRequest('/api/v1/auth/register', 'POST', {
            username,
            email,
            password
        }, false);
    },
    
    // 登录
    login: async (email, password) => {
        const response = await apiRequest('/api/v1/auth/login', 'POST', {
            email,
            password,
            type: 'email'  // 添加登录类型
        }, false);
        
        console.log('登录响应:', response);
        
        // 处理后端响应格式 {code: 200, message: "success", data: {...}}
        if (response.status === 200) {
            let tokenData;
            if (response.data && response.data.code === 200 && response.data.data) {
                tokenData = response.data.data;
            } else if (response.data && response.data.access_token) {
                tokenData = response.data;
            }
            
            if (tokenData && tokenData.access_token) {
                accessToken = tokenData.access_token;
                refreshToken = tokenData.refresh_token;
                userId = tokenData.user_id;
                
                console.log('=== 登录Token保存前 ===');
                console.log('accessToken:', accessToken ? `存在(长度:${accessToken.length})` : '不存在');
                console.log('refreshToken:', refreshToken ? `存在(长度:${refreshToken.length})` : '不存在');
                console.log('userId:', userId);
                
                // 保存token到localStorage和sessionStorage双重备份
                try {
                    // 测试localStorage写入能力
                    if (isLocalStorageSupported) {
                        localStorage.setItem('test_write', 'test_value');
                        console.log('localStorage写入测试:', localStorage.getItem('test_write'));
                        
                        localStorage.setItem('access_token', accessToken);
                        localStorage.setItem('refresh_token', refreshToken);
                        localStorage.setItem('user_id', userId);
                    }
                    
                    // 同时保存到sessionStorage作为备份
                    sessionStorage.setItem('access_token', accessToken);
                    sessionStorage.setItem('refresh_token', refreshToken);
                    sessionStorage.setItem('user_id', userId);
                    
                } catch (e) {
                    console.error('Token保存失败:', e);
                    alert('登录成功但无法保存状态，请检查浏览器设置');
                }
                
                console.log('=== 登录Token保存后验证 ===');
                console.log('localStorage access_token:', localStorage.getItem('access_token') ? '存在' : '不存在');
                console.log('localStorage refresh_token:', localStorage.getItem('refresh_token') ? '存在' : '不存在');
                console.log('localStorage user_id:', localStorage.getItem('user_id'));
                
                // 再次验证具体内容
                console.log('=== 详细验证 ===');
                const storedAccessToken = localStorage.getItem('access_token');
                const storedRefreshToken = localStorage.getItem('refresh_token');
                console.log('stored access_token 长度:', storedAccessToken ? storedAccessToken.length : 0);
                console.log('stored refresh_token 长度:', storedRefreshToken ? storedRefreshToken.length : 0);
                console.log('refresh_token 前50字符:', storedRefreshToken ? storedRefreshToken.substring(0, 50) : 'null');
                
                console.log('登录成功，Token已保存');
                checkToken();
            } else {
                console.error('登录响应格式错误:', response.data);
            }
        }
        
        return response;
    },
    
    // 刷新Token
    refreshToken: async () => {
        console.log('=== refreshToken 函数开始执行 ===');
        
        // 从localStorage和sessionStorage中读取refresh_token，确保最新状态
        let currentRefreshToken = null;
        
        try {
            // 优先从localStorage读取
            if (isLocalStorageSupported) {
                currentRefreshToken = localStorage.getItem('refresh_token');
            }
            
            // 如果localStorage没有，则尝试从sessionStorage获取
            if (!currentRefreshToken) {
                currentRefreshToken = sessionStorage.getItem('refresh_token');
                console.log('从sessionStorage获取refresh_token');
            }
            
            console.log('读取到的refresh_token:', currentRefreshToken ? `存在(长度:${currentRefreshToken.length})` : '不存在');
        } catch (e) {
            console.error('获取refresh_token出错:', e);
        }
        
        if (!currentRefreshToken) {
            console.log('没有refresh_token，返回401错误');
            return {
                status: 401,
                data: {
                    code: 401,
                    message: '未登录，无法刷新Token'
                }
            };
        }

        console.log('正在刷新Token，当前refresh_token:', currentRefreshToken ? '存在' : '不存在');
        console.log('准备调用 apiRequest...');

        console.log('发送刷新Token请求，参数:', { refresh_token: currentRefreshToken });
        
        let response;
        // 尝试直接使用fetch而不是apiRequest，确保不会有其他影响
        try {
            const fetchResponse = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ refresh_token: currentRefreshToken }),
                credentials: 'include'
            });
            
            const data = await fetchResponse.json();
            console.log('刷新Token原始响应:', data);
            
            // 构造与apiRequest相同格式的响应
            response = {
                status: fetchResponse.status,
                data: data
            };
        } catch (error) {
            console.error('刷新Token请求出错:', error);
            return {
                status: 500,
                data: {
                    message: '刷新Token请求异常'
                }
            };
        }

        console.log('Token刷新响应:', response);

        // 检查响应格式，兼容后端的 {code, message, data} 格式
        if (response.status === 200) {
            // 后端返回的格式是 {code: 200, message: "success", data: {...}}
            let tokenData;
            if (response.data && response.data.code === 200 && response.data.data) {
                tokenData = response.data.data;
            } else if (response.data && response.data.access_token) {
                tokenData = response.data;
            }
            
            if (tokenData && tokenData.access_token) {
                accessToken = tokenData.access_token;
                refreshToken = tokenData.refresh_token;
                userId = tokenData.user_id; // 确保user_id也被更新

                // 保存到localStorage和sessionStorage双重备份
                try {
                    if (isLocalStorageSupported) {
                        localStorage.setItem('access_token', accessToken);
                        localStorage.setItem('refresh_token', refreshToken);
                        localStorage.setItem('user_id', userId);
                    }
                    
                    // 同时保存到sessionStorage
                    sessionStorage.setItem('access_token', accessToken);
                    sessionStorage.setItem('refresh_token', refreshToken);
                    sessionStorage.setItem('user_id', userId);
                } catch (e) {
                    console.error('刷新Token后保存失败:', e);
                }

                console.log('Token刷新成功');
                checkToken();
            } else {
                console.error('Token刷新响应格式错误:', response.data);
            }
        } else {
            console.error('Token刷新失败:', response);
            // 如果刷新失败，清除所有Token
            if (response.status === 401) {
                accessToken = null;
                refreshToken = null;
                userId = null;

                localStorage.removeItem('access_token');
                localStorage.removeItem('refresh_token');
                localStorage.removeItem('user_id');

                checkToken();
            }
        }

        return response;
    },
    
    // 验证Token
    verifyToken: () => {
        if (!accessToken) {
            return Promise.resolve({
                status: 401,
                data: {
                    code: 401,
                    message: '未登录，无法验证Token'
                }
            });
        }

        console.log('正在验证Token，当前access_token:', accessToken ? '存在' : '不存在');

        // 发送Token验证请求，使用Authorization头部而不是请求体
        return apiRequest('/api/v1/auth/verify', 'POST', {
            token: accessToken
        }, true);
    },
    
    // 登出
    logout: async () => {
        const response = await apiRequest('/api/v1/auth/logout', 'POST', null, true);
        
        accessToken = null;
        refreshToken = null;
        userId = null;
        
        try {
            // 清除localStorage
            if (isLocalStorageSupported) {
                localStorage.removeItem('access_token');
                localStorage.removeItem('refresh_token');
                localStorage.removeItem('user_id');
            }
            
            // 清除sessionStorage
            sessionStorage.removeItem('access_token');
            sessionStorage.removeItem('refresh_token');
            sessionStorage.removeItem('user_id');
            
            console.log('所有Token已清除');
        } catch (e) {
            console.error('清除Token失败:', e);
        }
        
        checkToken();
        
        return response;
    }
};

// 用户API
const userAPI = {
    // 获取用户列表
    getUsers: (page = 1, limit = 10) => {
        return apiRequest(`/api/v1/users?page=${page}&limit=${limit}`, 'GET');
    },
    
    // 获取单个用户
    getUser: (userId) => {
        return apiRequest(`/api/v1/users/${userId}`, 'GET');
    },
    
    // 更新用户
    updateUser: (userId, userData) => {
        return apiRequest(`/api/v1/users/${userId}`, 'PUT', userData);
    },
    
    // 删除用户
    deleteUser: (userId) => {
        return apiRequest(`/api/v1/users/${userId}`, 'DELETE');
    },
    
    // 分配角色
    assignRole: (userId, roleId) => {
        return apiRequest(`/api/v1/users/${userId}/roles`, 'POST', { role_id: roleId });
    },
    
    // 移除角色
    removeRole: (userId, roleId) => {
        return apiRequest(`/api/v1/users/${userId}/roles/${roleId}`, 'DELETE');
    },
    
    // 获取用户权限
    getUserPermissions: (userId) => {
        return apiRequest(`/api/v1/users/${userId}/permissions`, 'GET');
    }
};

// 角色API
const roleAPI = {
    // 获取角色列表
    getRoles: (page = 1, limit = 10) => {
        return apiRequest(`/api/v1/roles?page=${page}&limit=${limit}`, 'GET');
    },
    
    // 创建角色
    createRole: (name, description) => {
        return apiRequest('/api/v1/roles', 'POST', { name, description });
    },
    
    // 获取角色详情
    getRole: (roleId) => {
        return apiRequest(`/api/v1/roles/${roleId}`, 'GET');
    },
    
    // 更新角色
    updateRole: (roleId, roleData) => {
        return apiRequest(`/api/v1/roles/${roleId}`, 'PUT', roleData);
    },
    
    // 删除角色
    deleteRole: (roleId) => {
        return apiRequest(`/api/v1/roles/${roleId}`, 'DELETE');
    },
    
    // 分配权限
    assignPermission: (roleId, permissionId) => {
        return apiRequest(`/api/v1/roles/${roleId}/permissions`, 'POST', { permission_id: permissionId });
    },
    
    // 移除权限
    removePermission: (roleId, permissionId) => {
        return apiRequest(`/api/v1/roles/${roleId}/permissions/${permissionId}`, 'DELETE');
    }
};

// 权限API
const permissionAPI = {
    // 获取权限列表
    getPermissions: () => {
        return apiRequest('/api/v1/permissions', 'GET');
    },
    
    // 创建权限
    createPermission: (name, description, resource, action) => {
        return apiRequest('/api/v1/permissions', 'POST', { 
            name, 
            description, 
            resource, 
            action 
        });
    },
    
    // 获取权限详情
    getPermission: (permissionId) => {
        return apiRequest(`/api/v1/permissions/${permissionId}`, 'GET');
    },
    
    // 更新权限
    updatePermission: (permissionId, permData) => {
        return apiRequest(`/api/v1/permissions/${permissionId}`, 'PUT', permData);
    },
    
    // 删除权限
    deletePermission: (permissionId) => {
        return apiRequest(`/api/v1/permissions/${permissionId}`, 'DELETE');
    }
};

// 辅助函数 - 刷新Token
async function refreshTokenFunc() {
    // 重新从localStorage和sessionStorage读取refresh_token，确保最新状态
    let currentRefreshToken = null;
    
    try {
        if (isLocalStorageSupported) {
            currentRefreshToken = localStorage.getItem('refresh_token');
        }
        
        // 如果localStorage没有，则从sessionStorage尝试获取
        if (!currentRefreshToken) {
            currentRefreshToken = sessionStorage.getItem('refresh_token');
        }
    } catch (e) {
        console.error('获取refresh_token出错:', e);
    }
    
    if (!currentRefreshToken) {
        console.log('没有refresh_token，无法刷新');
        return false;
    }

    console.log('refreshTokenFunc: 开始刷新Token...');

    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ refresh_token: currentRefreshToken }),
            credentials: 'include'
        });

        const data = await response.json();
        console.log('refreshTokenFunc: 刷新响应:', data);

        // 处理后端响应格式 {code, message, data}
        if (response.ok && data.code === 200 && data.data) {
            const tokenData = data.data;
            if (tokenData.access_token) {
                accessToken = tokenData.access_token;
                refreshToken = tokenData.refresh_token;
                userId = tokenData.user_id; // 确保user_id也被更新

                // 保存到localStorage和sessionStorage双重备份
                try {
                    if (isLocalStorageSupported) {
                        localStorage.setItem('access_token', accessToken);
                        localStorage.setItem('refresh_token', refreshToken);
                        localStorage.setItem('user_id', userId);
                    }
                    
                    // 同时保存到sessionStorage
                    sessionStorage.setItem('access_token', accessToken);
                    sessionStorage.setItem('refresh_token', refreshToken);
                    sessionStorage.setItem('user_id', userId);
                } catch (e) {
                    console.error('刷新Token后保存失败:', e);
                }

                console.log('refreshTokenFunc: Token刷新成功');
                checkToken();
                return true;
            }
        } else {
            console.log('refreshTokenFunc: 刷新失败，清除所有Token');
            // 刷新失败，清除token
            accessToken = null;
            refreshToken = null;
            userId = null;

            try {
                // 从localStorage和sessionStorage都清除
                if (isLocalStorageSupported) {
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    localStorage.removeItem('user_id');
                }
                
                sessionStorage.removeItem('access_token');
                sessionStorage.removeItem('refresh_token');
                sessionStorage.removeItem('user_id');
            } catch (e) {
                console.error('清除Token出错:', e);
            }

            checkToken();
            return false;
        }
    } catch (error) {
        console.error('refreshTokenFunc: 刷新Token异常:', error);
        return false;
    }
}

// 格式化JSON显示
function formatJSON(json) {
    return JSON.stringify(json, null, 2);
}

// 初始化时检查Token
document.addEventListener('DOMContentLoaded', async () => {
    console.log('页面加载完成，开始检查Token状态...');
    console.log('=== DOMContentLoaded: 当前Token存储状态 ===');
    
    try {
        // 检查localStorage状态
        if (isLocalStorageSupported) {
            console.log('localStorage支持状态: 支持');
            console.log('localStorage access_token:', localStorage.getItem('access_token') ? '存在' : '不存在');
            console.log('localStorage refresh_token:', localStorage.getItem('refresh_token') ? '存在' : '不存在');
            console.log('localStorage user_id:', localStorage.getItem('user_id'));
        } else {
            console.log('localStorage支持状态: 不支持');
        }
        
        // 检查sessionStorage状态
        console.log('sessionStorage access_token:', sessionStorage.getItem('access_token') ? '存在' : '不存在');
        console.log('sessionStorage refresh_token:', sessionStorage.getItem('refresh_token') ? '存在' : '不存在');
        console.log('sessionStorage user_id:', sessionStorage.getItem('user_id'));
    } catch (e) {
        console.error('检查Token存储状态出错:', e);
    }
    
    // 只有当有access_token时才进行检查和刷新
    const hasAccessToken = !!(localStorage.getItem('access_token') || sessionStorage.getItem('access_token'));
    if (hasAccessToken) {
        await checkAndRefreshToken();
    } else {
        console.log('没有access_token，跳过Token检查');
        checkToken(); // 只更新UI状态
    }
});
