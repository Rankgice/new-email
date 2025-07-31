import { apiClient } from '@/utils/api'
import type {
  UserSettings,
  UserProfileUpdateRequest,
  SecuritySettingsUpdateRequest,
  NotificationSettingsUpdateRequest,
  EmailFilter,
  ApiKey
} from '@/types'

// 用户设置API
export const userSettingsApi = {
  // 获取用户设置
  async getSettings(): Promise<UserSettings> {
    const response = await apiClient.get('/user/settings')
    return response.data
  },

  // 更新用户资料
  async updateProfile(data: UserProfileUpdateRequest): Promise<void> {
    await apiClient.put('/user/profile', data)
  },

  // 更新安全设置
  async updateSecurity(data: SecuritySettingsUpdateRequest): Promise<void> {
    await apiClient.put('/user/security', data)
  },

  // 更新通知设置
  async updateNotifications(data: NotificationSettingsUpdateRequest): Promise<void> {
    await apiClient.put('/user/notifications', data)
  },

  // 更新邮件签名
  async updateEmailSignature(signature: string): Promise<void> {
    await apiClient.put('/user/email-signature', { signature })
  },

  // 更新主题设置
  async updateTheme(theme: string): Promise<void> {
    await apiClient.put('/user/theme', { theme })
  },

  // 邮件过滤器相关
  async getEmailFilters(): Promise<EmailFilter[]> {
    const response = await apiClient.get('/user/email-filters')
    return response.data
  },

  async createEmailFilter(filter: Omit<EmailFilter, 'id' | 'createdAt' | 'updatedAt'>): Promise<EmailFilter> {
    const response = await apiClient.post('/user/email-filters', filter)
    return response.data
  },

  async updateEmailFilter(id: string, filter: Partial<EmailFilter>): Promise<EmailFilter> {
    const response = await apiClient.put(`/user/email-filters/${id}`, filter)
    return response.data
  },

  async deleteEmailFilter(id: string): Promise<void> {
    await apiClient.delete(`/user/email-filters/${id}`)
  },

  async toggleEmailFilter(id: string, enabled: boolean): Promise<void> {
    await apiClient.put(`/user/email-filters/${id}/toggle`, { enabled })
  },

  // API密钥管理
  async getApiKeys(): Promise<ApiKey[]> {
    const response = await apiClient.get('/user/api-keys')
    return response.data
  },

  async createApiKey(data: { name: string; permissions: string[]; expiresAt?: string }): Promise<ApiKey> {
    const response = await apiClient.post('/user/api-keys', data)
    return response.data
  },

  async updateApiKey(id: string, data: { name?: string; permissions?: string[]; enabled?: boolean }): Promise<ApiKey> {
    const response = await apiClient.put(`/user/api-keys/${id}`, data)
    return response.data
  },

  async deleteApiKey(id: string): Promise<void> {
    await apiClient.delete(`/user/api-keys/${id}`)
  },

  async regenerateApiKey(id: string): Promise<ApiKey> {
    const response = await apiClient.post(`/user/api-keys/${id}/regenerate`)
    return response.data
  },

  // 头像上传
  async uploadAvatar(file: File): Promise<{ url: string }> {
    const formData = new FormData()
    formData.append('avatar', file)
    
    const response = await apiClient.post('/user/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  // 验证当前密码
  async verifyPassword(password: string): Promise<boolean> {
    try {
      const response = await apiClient.post('/user/verify-password', { password })
      return response.data.valid
    } catch {
      return false
    }
  },

  // 启用/禁用两步验证
  async setupTwoFactor(): Promise<{ qrCode: string; secret: string }> {
    const response = await apiClient.post('/user/2fa/setup')
    return response.data
  },

  async enableTwoFactor(token: string): Promise<{ recoveryCodes: string[] }> {
    const response = await apiClient.post('/user/2fa/enable', { token })
    return response.data
  },

  async disableTwoFactor(password: string): Promise<void> {
    await apiClient.post('/user/2fa/disable', { password })
  },

  // 获取登录历史
  async getLoginHistory(): Promise<Array<{
    id: string
    ip: string
    userAgent: string
    location?: string
    loginAt: string
    success: boolean
  }>> {
    const response = await apiClient.get('/user/login-history')
    return response.data
  },

  // 注销其他设备
  async logoutOtherDevices(password: string): Promise<void> {
    await apiClient.post('/user/logout-other-devices', { password })
  }
}
