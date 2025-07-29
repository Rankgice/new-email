import type { ApiResponse, User, Email, EmailListParams, PaginatedResponse } from '@/types'

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const API_TIMEOUT = 10000

// 请求拦截器
class ApiClient {
  private baseURL: string
  private timeout: number

  constructor(baseURL: string, timeout: number = API_TIMEOUT) {
    this.baseURL = baseURL
    this.timeout = timeout
  }

  // 获取认证头
  private getAuthHeaders(): Record<string, string> {
    const token = localStorage.getItem('auth_token')
    return token ? { Authorization: `Bearer ${token}` } : {}
  }

  // 通用请求方法
  private async request<T = any>(
    endpoint: string,
    options: RequestInit = {},
    customHeaders?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`

    const config: RequestInit = {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...this.getAuthHeaders(),
        ...customHeaders,
        ...options.headers,
      },
      signal: AbortSignal.timeout(this.timeout),
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }

      const data = await response.json()
      return data
    } catch (error: any) {
      console.error(`API Error [${endpoint}]:`, error)
      
      if (error.name === 'AbortError') {
        throw new Error('请求超时，请稍后重试')
      }
      
      if (error.message.includes('401')) {
        // 清除过期的认证信息
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
        window.location.href = '/auth/login'
        throw new Error('登录已过期，请重新登录')
      }
      
      throw error
    }
  }

  // GET 请求
  async get<T = any>(endpoint: string, params?: Record<string, any>, options?: { headers?: Record<string, string> }): Promise<ApiResponse<T>> {
    const url = params ? `${endpoint}?${new URLSearchParams(params)}` : endpoint
    return this.request<T>(url, { method: 'GET' }, options?.headers)
  }

  // POST 请求
  async post<T = any>(endpoint: string, data?: any, options?: { headers?: Record<string, string> }): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    }, options?.headers)
  }

  // PUT 请求
  async put<T = any>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  // DELETE 请求
  async delete<T = any>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }

  // 文件上传
  async upload<T = any>(endpoint: string, file: File, data?: Record<string, any>): Promise<ApiResponse<T>> {
    const formData = new FormData()
    formData.append('file', file)
    
    if (data) {
      Object.entries(data).forEach(([key, value]) => {
        formData.append(key, String(value))
      })
    }

    return this.request<T>(endpoint, {
      method: 'POST',
      body: formData,
      headers: {
        // 不设置 Content-Type，让浏览器自动设置
        ...this.getAuthHeaders(),
      },
    })
  }
}

// 创建 API 客户端实例
const apiClient = new ApiClient(API_BASE_URL)

// 认证相关 API
export const authApi = {
  // 用户登录
  login: (credentials: { email: string; password: string }) =>
    apiClient.post<{ user: User; token: string }>('/public/user/login', credentials),

  // 用户注册
  register: (userData: { username: string; email: string; password: string; nickname?: string }) =>
    apiClient.post<{ user: User }>('/public/user/register', userData),

  // 管理员登录
  adminLogin: (credentials: { username: string; password: string }) =>
    apiClient.post<{ user: User; token: string }>('/public/admin/login', credentials),

  // 登出
  logout: () => apiClient.post('/user/logout'),

  // 验证 token
  validateToken: () => apiClient.get<{ user: User }>('/user/profile'),

  // 发送验证码
  sendVerificationCode: (data: { email: string; type: string }) =>
    apiClient.post('/public/send-code', data),

  // 验证验证码
  verifyCode: (data: { email: string; code: string; type: string }) =>
    apiClient.post('/public/verify-code', data),

  // 忘记密码
  forgotPassword: (data: { email: string }) =>
    apiClient.post('/public/send-code', { ...data, type: 'reset_password' }),

  // 重置密码
  resetPassword: (data: { email: string; code: string; password: string }) =>
    apiClient.post('/public/reset-password', data),
}

// 邮件相关 API
export const emailApi = {
  // 获取邮件列表
  getEmails: (params: EmailListParams) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails', params),

  // 获取邮件详情
  getEmail: (id: string) =>
    apiClient.get<Email>(`/user/emails/${id}`),

  // 发送邮件
  sendEmail: (emailData: any) =>
    apiClient.post<Email>('/user/emails/send', emailData),

  // 标记已读/未读
  markAsRead: (id: string, isRead: boolean) =>
    apiClient.put(`/user/emails/${id}/read`, { is_read: isRead }),

  // 标记星标
  markAsStarred: (id: string, isStarred: boolean) =>
    apiClient.put(`/user/emails/${id}/star`, { is_starred: isStarred }),

  // 删除邮件
  deleteEmail: (id: string) =>
    apiClient.delete(`/user/emails/${id}`),

  // 批量操作
  batchUpdate: (emailIds: string[], action: string, data?: any) =>
    apiClient.post('/user/emails/batch', { email_ids: emailIds, action, ...data }),

  // 搜索邮件
  searchEmails: (query: string, filters?: any) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails', { search: query, ...filters }),

  // 上传附件
  uploadAttachment: (file: File) =>
    apiClient.upload<{ url: string; filename: string; size: number }>('/user/attachments', file),

  // 草稿管理
  getDrafts: () =>
    apiClient.get<PaginatedResponse<Email>>('/user/drafts'),

  saveDraft: (draftData: any) =>
    apiClient.post<Email>('/user/drafts', draftData),

  updateDraft: (id: string, draftData: any) =>
    apiClient.put<Email>(`/user/drafts/${id}`, draftData),

  deleteDraft: (id: string) =>
    apiClient.delete(`/user/drafts/${id}`),

  sendDraft: (id: string) =>
    apiClient.post(`/user/drafts/${id}/send`),
}

// 用户相关 API
export const userApi = {
  // 获取用户信息
  getProfile: () => apiClient.get<User>('/user/profile'),

  // 更新用户信息
  updateProfile: (userData: Partial<User>) =>
    apiClient.put<User>('/user/profile', userData),

  // 更改密码
  changePassword: (data: { old_password: string; new_password: string }) =>
    apiClient.post('/user/change-password', data),

  // 上传头像
  uploadAvatar: (file: File) =>
    apiClient.upload<{ url: string }>('/user/avatar', file),

  // 邮箱管理
  getMailboxes: () => apiClient.get('/user/mailboxes'),

  createMailbox: (mailboxData: any) =>
    apiClient.post('/user/mailboxes', mailboxData),

  updateMailbox: (id: string, mailboxData: any) =>
    apiClient.put(`/user/mailboxes/${id}`, mailboxData),

  deleteMailbox: (id: string) =>
    apiClient.delete(`/user/mailboxes/${id}`),

  syncMailbox: (id: string) =>
    apiClient.post(`/user/mailboxes/${id}/sync`),

  testMailbox: (id: string) =>
    apiClient.post(`/user/mailboxes/${id}/test`),

  // 邮件模板
  getTemplates: () => apiClient.get('/user/templates'),

  createTemplate: (templateData: any) =>
    apiClient.post('/user/templates', templateData),

  updateTemplate: (id: string, templateData: any) =>
    apiClient.put(`/user/templates/${id}`, templateData),

  deleteTemplate: (id: string) =>
    apiClient.delete(`/user/templates/${id}`),

  // 邮件签名
  getSignatures: () => apiClient.get('/user/signatures'),

  createSignature: (signatureData: any) =>
    apiClient.post('/user/signatures', signatureData),

  // 验证码记录
  getVerificationCodes: () => apiClient.get('/user/verification-codes'),

  getVerificationCode: (id: string) =>
    apiClient.get(`/user/verification-codes/${id}`),

  markCodeAsUsed: (id: string) =>
    apiClient.put(`/user/verification-codes/${id}/used`),

  // API密钥管理
  getApiKeys: () => apiClient.get('/user/api-keys'),

  createApiKey: (keyData: any) =>
    apiClient.post('/user/api-keys', keyData),

  deleteApiKey: (id: string) =>
    apiClient.delete(`/user/api-keys/${id}`),

  // 日志查询
  getOperationLogs: () => apiClient.get('/user/logs/operation'),

  getEmailLogs: () => apiClient.get('/user/logs/email'),
}

// 管理员相关 API
export const adminApi = {
  // 获取统计数据
  getStats: () => apiClient.get('/admin/stats'),

  // 用户管理
  getUsers: (params?: any) =>
    apiClient.get<PaginatedResponse<User>>('/admin/users', params),

  getUser: (id: string) =>
    apiClient.get<User>(`/admin/users/${id}`),

  createUser: (userData: any) =>
    apiClient.post<User>('/admin/users', userData),

  updateUser: (id: string, userData: Partial<User>) =>
    apiClient.put<User>(`/admin/users/${id}`, userData),

  deleteUser: (id: string) =>
    apiClient.delete(`/admin/users/${id}`),

  resetUserPassword: (id: string, newPassword: string) =>
    apiClient.post(`/admin/users/${id}/reset-password`, { new_password: newPassword }),

  // 邮箱管理
  getAllMailboxes: () => apiClient.get('/admin/mailboxes'),

  getMailbox: (id: string) =>
    apiClient.get(`/admin/mailboxes/${id}`),

  forceSyncMailbox: (id: string) =>
    apiClient.post(`/admin/mailboxes/${id}/force-sync`),

  updateMailboxStatus: (id: string, status: string) =>
    apiClient.put(`/admin/mailboxes/${id}/status`, { status }),

  // 邮件监控
  getAllEmails: (params?: any) =>
    apiClient.get<PaginatedResponse<Email>>('/admin/emails', params),

  getEmailStats: () => apiClient.get('/admin/emails/stats'),

  deleteEmailAsAdmin: (id: string) =>
    apiClient.delete(`/admin/emails/${id}`),

  // 系统管理
  getSystemConfig: () => apiClient.get('/admin/config'),

  updateSystemConfig: (config: any) =>
    apiClient.put('/admin/config', config),

  // 日志管理
  getSystemLogs: (params?: any) =>
    apiClient.get('/admin/logs/system', params),

  getOperationLogs: (params?: any) =>
    apiClient.get('/admin/logs/operation', params),

  getErrorLogs: (params?: any) =>
    apiClient.get('/admin/logs/error', params),
}

// API Key 相关接口 (使用 X-API-Key 认证)
export const apiKeyApi = {
  // 发送邮件
  sendEmail: (emailData: any, apiKey: string) =>
    apiClient.post('/v1/emails/send', emailData, {
      headers: { 'X-API-Key': apiKey }
    }),

  // 获取邮件列表
  getEmails: (apiKey: string, params?: any) =>
    apiClient.get('/v1/emails', params, {
      headers: { 'X-API-Key': apiKey }
    }),

  // 获取验证码
  getVerificationCodes: (apiKey: string, params?: any) =>
    apiClient.get('/v1/verification-codes', params, {
      headers: { 'X-API-Key': apiKey }
    }),
}

// 导出统一的 API 对象
export const api = {
  ...authApi,
  ...emailApi,
  ...userApi,
  ...adminApi,

  // API Key 接口
  apiKey: apiKeyApi,

  // 通用方法
  get: apiClient.get.bind(apiClient),
  post: apiClient.post.bind(apiClient),
  put: apiClient.put.bind(apiClient),
  delete: apiClient.delete.bind(apiClient),
  upload: apiClient.upload.bind(apiClient),
}
