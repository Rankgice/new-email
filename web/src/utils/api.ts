import type {
  ApiResponse,
  User,
  Email,
  EmailListParams,
  PaginatedResponse,
  PageResponse,
  Mailbox,
  MailboxCreateRequest,
  MailboxUpdateRequest,
  MailboxListRequest,
  MailboxSyncRequest,
  MailboxSyncResponse,
  MailboxStats,
  MailboxConnectionTestResponse,
  EmailSendRequest,
  EmailSendResponse,
  EmailAddress
} from '@/types'

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

      const backendResponse = await response.json()

      // 转换后端响应格式为前端期望的格式
      const apiResponse: ApiResponse<T> = {
        success: backendResponse.code === 0,
        data: backendResponse.data,
        message: backendResponse.msg,
        code: backendResponse.code
      }

      return apiResponse
    } catch (error: unknown) {
      console.error(`API Error [${endpoint}]:`, error)

      if (error instanceof DOMException && error.name === 'AbortError') {
        throw new Error('请求超时，请稍后重试')
      }

      if (error instanceof Error && error.message.includes('401')) {
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
  async get<T = any>(endpoint: string, paramsOrOptions?: Record<string, any> | { params?: Record<string, any>, headers?: Record<string, string> }, options?: { headers?: Record<string, string> }): Promise<ApiResponse<T>> {
    let url = endpoint
    let headers = options?.headers

    // 处理参数：支持两种调用方式
    // 1. get(url, params, options)
    // 2. get(url, { params, headers })
    let queryParams: Record<string, any> | undefined

    if (paramsOrOptions) {
      if (paramsOrOptions.params !== undefined) {
        // 第二种调用方式：{ params: {...}, headers: {...} }
        queryParams = paramsOrOptions.params
        headers = paramsOrOptions.headers || headers
      } else {
        // 第一种调用方式：直接传递params
        queryParams = paramsOrOptions
      }
    }

    if (queryParams) {
      // 过滤掉undefined和null值，并确保所有值都是字符串
      const cleanParams: Record<string, string> = {}
      Object.entries(queryParams).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          cleanParams[key] = String(value)
        }
      })

      if (Object.keys(cleanParams).length > 0) {
        url = `${endpoint}?${new URLSearchParams(cleanParams)}`
      }
    }

    return this.request<T>(url, { method: 'GET' }, headers)
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

// 导出 API 客户端供其他模块使用
export { apiClient }

const toEmailAddress = (email: string, name?: string): EmailAddress => ({
  email,
  ...(name ? { name } : {})
})

const normalizePageResponse = <T>(payload?: PaginatedResponse<T> | PageResponse<T> | T[]): PageResponse<T> => {
  if (Array.isArray(payload)) {
    return {
      list: payload,
      items: payload,
      data: payload,
      total: payload.length,
      page: 1,
      pageSize: payload.length || 1,
      limit: payload.length || 1,
    }
  }

  const list = payload?.list ?? payload?.items ?? payload?.data ?? []
  const pageSize = (payload?.pageSize ?? payload?.limit ?? list.length) || 1

  return {
    list,
    items: list,
    data: list,
    total: payload?.total ?? list.length,
    page: payload?.page ?? 1,
    pageSize,
    limit: pageSize,
  }
}

const normalizeEmail = (email: Partial<Email>): Email => {
  const fromEmail = email.fromEmail ?? email.from?.email ?? ''
  const fromName = email.fromName ?? email.from?.name
  const toEmails = email.toEmails ?? email.to?.map(item => item.email) ?? []
  const ccEmails = email.ccEmails ?? email.cc?.map(item => item.email) ?? []
  const bccEmails = email.bccEmails ?? email.bcc?.map(item => item.email) ?? []
  const content = email.content ?? email.body ?? ''

  return {
    id: email.id ?? '',
    userId: email.userId ?? 0,
    mailboxId: email.mailboxId ?? 0,
    subject: email.subject ?? '',
    fromEmail,
    fromName,
    toEmails,
    ccEmails,
    bccEmails,
    content,
    body: content,
    from: toEmailAddress(fromEmail, fromName),
    to: toEmails.map(item => toEmailAddress(item)),
    cc: ccEmails.map(item => toEmailAddress(item)),
    bcc: bccEmails.map(item => toEmailAddress(item)),
    contentType: email.contentType ?? 'text',
    attachments: email.attachments ?? [],
    status: email.status ?? 0,
    type: email.type ?? 'inbox',
    isRead: email.isRead ?? false,
    isStarred: email.isStarred ?? false,
    isImportant: email.isImportant,
    labels: email.labels ?? [],
    deletedAt: email.deletedAt,
    sentAt: email.sentAt,
    receivedAt: email.receivedAt,
    createdAt: email.createdAt ?? '',
    updatedAt: email.updatedAt ?? '',
  }
}

const normalizeEmailPageResponse = (response: ApiResponse<PaginatedResponse<Email>>): ApiResponse<PageResponse<Email>> => {
  const page = normalizePageResponse(response.data)
  const list = page.list.map(normalizeEmail)

  return {
    ...response,
    data: {
      ...page,
      list,
      items: list,
      data: list,
    },
  }
}

const normalizeMailbox = (mailbox: Partial<Mailbox>): Mailbox => ({
  id: mailbox.id ?? 0,
  userId: mailbox.userId ?? 0,
  domainId: mailbox.domainId ?? 0,
  email: mailbox.email ?? '',
  name: mailbox.name,
  provider: mailbox.provider,
  type: mailbox.type ?? (mailbox.domainId ? 'self' : 'third'),
  imapHost: mailbox.imapHost,
  imapPort: mailbox.imapPort,
  imapSsl: mailbox.imapSsl,
  smtpHost: mailbox.smtpHost,
  smtpPort: mailbox.smtpPort,
  smtpSsl: mailbox.smtpSsl,
  autoReceive: mailbox.autoReceive ?? false,
  status: mailbox.status ?? 0,
  unreadCount: mailbox.unreadCount ?? 0,
  lastSyncAt: mailbox.lastSyncAt,
  createdAt: mailbox.createdAt ?? '',
  updatedAt: mailbox.updatedAt ?? '',
})

const normalizeMailboxPageResponse = (response: ApiResponse<PaginatedResponse<Mailbox>>): ApiResponse<PageResponse<Mailbox>> => {
  const page = normalizePageResponse(response.data)
  const list = page.list.map(normalizeMailbox)

  return {
    ...response,
    data: {
      ...page,
      list,
      items: list,
      data: list,
    },
  }
}

// 认证相关 API
export const authApi = {
  // 用户登录
  login: (credentials: { username: string; password: string }) =>
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
    apiClient.get<PaginatedResponse<Email>>('/user/emails', { params }).then(normalizeEmailPageResponse),

  list(params: EmailListParams) {
    return this.getEmails(params)
  },

  // 获取已发送邮件列表
  getSentEmails: (params: Omit<EmailListParams, 'direction'>) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails', {
      params: { ...params, direction: 'sent' }
    }).then(normalizeEmailPageResponse),

  // 获取收件箱邮件列表
  getInboxEmails: (params: Omit<EmailListParams, 'direction'>) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails', {
      params: { ...params, direction: 'received' }
    }).then(normalizeEmailPageResponse),

  // 获取星标邮件列表
  getStarredEmails: (params: EmailListParams) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails', {
      params: { ...params, isStarred: true }
    }).then(normalizeEmailPageResponse),

  // 获取垃圾箱邮件列表
  getTrashEmails: (params: EmailListParams) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails/trash', { params }).then(normalizeEmailPageResponse),

  // 获取邮件详情
  getEmail: (id: Email['id']) =>
    apiClient.get<Email>(`/user/emails/${id}`).then(response => ({
      ...response,
      data: response.data ? normalizeEmail(response.data) : response.data,
    })),

  getById(id: Email['id']) {
    return this.getEmail(id)
  },

  // 发送邮件
  sendEmail: (data: EmailSendRequest) =>
    apiClient.post<EmailSendResponse>('/user/emails/send', data),

  // 标记已读/未读
  markAsRead: (id: Email['id'], isRead: boolean) =>
    apiClient.put(`/user/emails/${id}/read`, { is_read: isRead }),

  // 标记星标
  markAsStarred: (id: Email['id'], isStarred: boolean) =>
    apiClient.put(`/user/emails/${id}/star`, { is_starred: isStarred }),

  // 删除邮件（移到垃圾箱）
  deleteEmail: (id: Email['id']) =>
    apiClient.delete(`/user/emails/${id}`),

  // 恢复邮件（从垃圾箱恢复）
  restoreEmail: (id: Email['id']) =>
    apiClient.put(`/user/emails/${id}/restore`),

  // 永久删除邮件
  permanentDeleteEmail: (id: Email['id']) =>
    apiClient.delete(`/user/emails/${id}/permanent`),

  // 清空垃圾箱
  emptyTrash: () =>
    apiClient.delete('/user/emails/trash/empty'),

  // 批量操作
  batchUpdate: (emailIds: Array<Email['id']>, action: string, data?: Record<string, unknown>) =>
    apiClient.post('/user/emails/batch', { email_ids: emailIds, action, ...data }),

  // 搜索邮件
  searchEmails: (query: string, filters?: any) =>
    apiClient.get<PaginatedResponse<Email>>('/user/emails', { search: query, ...filters }),

  // 上传附件
  uploadAttachment: (file: File) =>
    apiClient.upload<{ url: string; filename: string; size: number }>('/user/attachments', file),

  // 草稿管理
  getDrafts: (params?: EmailListParams) =>
    apiClient.get<PaginatedResponse<Email>>('/user/drafts', params ? { params } : undefined).then(normalizeEmailPageResponse),

  saveDraft: (draftData: any) =>
    apiClient.post<Email>('/user/drafts', draftData).then(response => ({
      ...response,
      data: response.data ? normalizeEmail(response.data) : response.data,
    })),

  updateDraft: (id: Email['id'], draftData: Record<string, unknown>) =>
    apiClient.put<Email>(`/user/drafts/${id}`, draftData).then(response => ({
      ...response,
      data: response.data ? normalizeEmail(response.data) : response.data,
    })),

  deleteDraft: (id: Email['id']) =>
    apiClient.delete(`/user/drafts/${id}`),

  sendDraft: (id: Email['id']) =>
    apiClient.post(`/user/drafts/${id}/send`),

  // 获取用户的活跃邮箱列表（用于发件人选择）
  getActiveMailboxes: () =>
    apiClient.get<Mailbox[]>('/user/mailboxes', {
      params: { status: 1, pageSize: 100 }
    }).then(response => ({
      ...response,
      data: normalizePageResponse(response.data).list.map(normalizeMailbox)
    }))
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

// 邮箱管理 API
export const mailboxApi = {
  // 获取邮箱列表
  list: (params?: MailboxListRequest) =>
    apiClient.get<PaginatedResponse<Mailbox>>('/user/mailboxes', { params }).then(normalizeMailboxPageResponse),

  // 获取邮箱列表（简化版本）
  getMailboxes: () =>
    apiClient.get<PaginatedResponse<Mailbox>>('/user/mailboxes').then(normalizeMailboxPageResponse),

  // 创建邮箱
  create: (data: MailboxCreateRequest) =>
    apiClient.post<Mailbox>('/user/mailboxes', data).then(response => ({
      ...response,
      data: response.data ? normalizeMailbox(response.data) : response.data,
    })),

  // 更新邮箱
  update: (data: MailboxUpdateRequest) =>
    apiClient.put<Mailbox>(`/user/mailboxes/${data.id}`, data).then(response => ({
      ...response,
      data: response.data ? normalizeMailbox(response.data) : response.data,
    })),

  // 删除邮箱
  delete: (id: number) =>
    apiClient.delete(`/user/mailboxes/${id}`),

  // 获取邮箱详情
  getById: (id: number) =>
    apiClient.get<Mailbox>(`/user/mailboxes/${id}`).then(response => ({
      ...response,
      data: response.data ? normalizeMailbox(response.data) : response.data,
    })),

  // 同步邮箱
  sync: (data: MailboxSyncRequest) =>
    apiClient.post<MailboxSyncResponse>(`/user/mailboxes/${data.id}/sync`),

  testConnection: (id: number) =>
    apiClient.post<MailboxConnectionTestResponse>(`/user/mailboxes/${id}/test`),

  // 获取邮箱统计
  getStats: () =>
    apiClient.get<Record<string, number>>('/user/mailboxes/stats').then(response => ({
      ...response,
      data: response.data ? {
        totalMailboxes: response.data.totalMailboxes ?? response.data.total ?? 0,
        activeMailboxes: response.data.activeMailboxes ?? response.data.active ?? 0,
        selfMailboxes: response.data.selfMailboxes ?? 0,
        thirdMailboxes: response.data.thirdMailboxes ?? 0,
      } satisfies MailboxStats : response.data,
    })),
}



// 导出统一的 API 对象
export const api = {
  ...authApi,
  ...emailApi,
  ...userApi,
  ...adminApi,

  // API Key 接口
  apiKey: apiKeyApi,

  // 邮箱管理接口
  mailbox: mailboxApi,

  // 通用方法
  get: apiClient.get.bind(apiClient),
  post: apiClient.post.bind(apiClient),
  put: apiClient.put.bind(apiClient),
  delete: apiClient.delete.bind(apiClient),
  upload: apiClient.upload.bind(apiClient),
}
