import { api } from '@/utils/api'
import type {
  AdminLoginRequest,
  AdminLoginResponse,
  AdminDashboard,
  AdminSystemStats,
  AdminSystemSettings,
  User,
  UserCreateRequest,
  UserUpdateRequest,
  UserPasswordRequest,
  PageResponse,
  DailyStats,
  UserActivityStats,
  SystemHealth
} from '@/types'

// 管理员API服务
export const adminApi = {
  // 管理员登录
  async login(data: AdminLoginRequest): Promise<AdminLoginResponse> {
    const response = await api.post('/public/admin/login', data)
    return response.data
  },

  // 获取仪表板数据
  async getDashboard(): Promise<AdminDashboard> {
    const response = await api.get('/admin/dashboard')
    return response.data
  },

  // 获取系统统计信息
  async getSystemStats(): Promise<AdminSystemStats> {
    const response = await api.get('/admin/system/stats')
    return response.data
  },

  // 获取系统设置
  async getSystemSettings(): Promise<AdminSystemSettings> {
    const response = await api.get('/admin/system/settings')
    return response.data
  },

  // 更新系统设置
  async updateSystemSettings(data: AdminSystemSettings): Promise<AdminSystemSettings> {
    const response = await api.put('/admin/system/settings', data)
    return response.data
  },

  // 获取用户列表
  async getUsers(params: {
    page?: number
    pageSize?: number
    username?: string
    email?: string
    status?: number
    role?: string
    createdAtStart?: string
    createdAtEnd?: string
  }): Promise<PageResponse<User>> {
    const response = await api.get('/admin/users', { params })
    return response.data
  },

  // 获取用户详情
  async getUserById(id: number): Promise<{
    user: User
    mailboxes: any[]
    recentEmails: any[]
    statistics: any
  }> {
    const response = await api.get(`/admin/users/${id}`)
    return response.data
  },

  // 创建用户
  async createUser(data: UserCreateRequest): Promise<User> {
    const response = await api.post('/admin/users', data)
    return response.data
  },

  // 更新用户
  async updateUser(id: number, data: UserUpdateRequest): Promise<User> {
    const response = await api.put(`/admin/users/${id}`, data)
    return response.data
  },

  // 删除用户
  async deleteUser(id: number): Promise<void> {
    await api.delete(`/admin/users/${id}`)
  },

  // 重置用户密码
  async resetUserPassword(id: number, data: UserPasswordRequest): Promise<void> {
    await api.put(`/admin/users/${id}/password`, data)
  },

  // 批量操作用户
  async batchOperationUsers(data: {
    userIds: number[]
    action: 'enable' | 'disable' | 'delete'
  }): Promise<void> {
    await api.post('/admin/users/batch', data)
  },

  // 导入用户
  async importUsers(file: File): Promise<{
    total: number
    success: number
    failed: number
    errors: string[]
  }> {
    const formData = new FormData()
    formData.append('file', file)
    const response = await api.post('/admin/users/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  // 导出用户
  async exportUsers(params: {
    format?: 'csv' | 'excel'
    username?: string
    email?: string
    status?: number
    role?: string
    createdAtStart?: string
    createdAtEnd?: string
  }): Promise<{
    total: number
    format: string
    filename: string
    message: string
  }> {
    const response = await api.get('/admin/users/export', { params })
    return response.data
  },

  // 获取操作日志
  async getOperationLogs(params: {
    page?: number
    pageSize?: number
    userId?: number
    action?: string
    resource?: string
    method?: string
    status?: number
    createdAtStart?: string
    createdAtEnd?: string
  }): Promise<PageResponse<any>> {
    const response = await api.get('/admin/logs', { params })
    return response.data
  },

  // 获取用户增长统计
  async getUserGrowthStats(days: number = 30): Promise<DailyStats[]> {
    const response = await api.get('/admin/stats/user-growth', {
      params: { days }
    })
    return response.data
  },

  // 获取邮件增长统计
  async getEmailGrowthStats(days: number = 30): Promise<DailyStats[]> {
    const response = await api.get('/admin/stats/email-growth', {
      params: { days }
    })
    return response.data
  },

  // 获取活跃用户排行
  async getTopUsers(limit: number = 10): Promise<UserActivityStats[]> {
    const response = await api.get('/admin/stats/top-users', {
      params: { limit }
    })
    return response.data
  },

  // 获取系统健康状态
  async getSystemHealth(): Promise<SystemHealth> {
    const response = await api.get('/admin/system/health')
    return response.data
  },

  // 清理系统缓存
  async clearCache(): Promise<void> {
    await api.post('/admin/system/clear-cache')
  },

  // 重启系统服务
  async restartService(service: string): Promise<void> {
    await api.post('/admin/system/restart', { service })
  },

  // 备份数据库
  async backupDatabase(): Promise<{
    filename: string
    size: number
    createdAt: string
  }> {
    const response = await api.post('/admin/system/backup')
    return response.data
  },

  // 获取系统配置
  async getSystemConfig(): Promise<Record<string, any>> {
    const response = await api.get('/admin/system/config')
    return response.data
  },

  // 更新系统配置
  async updateSystemConfig(config: Record<string, any>): Promise<void> {
    await api.put('/admin/system/config', config)
  },

  // 获取邮箱统计
  async getMailboxStats(): Promise<{
    totalMailboxes: number
    activeMailboxes: number
    byProvider: Array<{ provider: string; count: number }>
    byStatus: Array<{ status: string; count: number }>
  }> {
    const response = await api.get('/admin/stats/mailboxes')
    return response.data
  },

  // 获取邮件统计
  async getEmailStats(): Promise<{
    totalEmails: number
    todayEmails: number
    byDirection: Array<{ direction: string; count: number }>
    byStatus: Array<{ status: string; count: number }>
  }> {
    const response = await api.get('/admin/stats/emails')
    return response.data
  },

  // 获取存储使用情况
  async getStorageUsage(): Promise<{
    totalSize: number
    usedSize: number
    freeSize: number
    byUser: Array<{ userId: number; username: string; size: number }>
  }> {
    const response = await api.get('/admin/stats/storage')
    return response.data
  },

  // 获取系统错误日志
  async getErrorLogs(params: {
    page?: number
    pageSize?: number
    level?: string
    startTime?: string
    endTime?: string
  }): Promise<PageResponse<any>> {
    const response = await api.get('/admin/logs/errors', { params })
    return response.data
  },

  // 获取性能监控数据
  async getPerformanceMetrics(timeRange: string = '1h'): Promise<{
    cpu: Array<{ time: string; value: number }>
    memory: Array<{ time: string; value: number }>
    disk: Array<{ time: string; value: number }>
    network: Array<{ time: string; in: number; out: number }>
  }> {
    const response = await api.get('/admin/metrics/performance', {
      params: { timeRange }
    })
    return response.data
  },

  // 发送系统通知
  async sendSystemNotification(data: {
    title: string
    message: string
    type: 'info' | 'warning' | 'error' | 'success'
    targetUsers?: number[]
    broadcast?: boolean
  }): Promise<void> {
    await api.post('/admin/notifications', data)
  },

  // 获取系统通知列表
  async getSystemNotifications(params: {
    page?: number
    pageSize?: number
    type?: string
    status?: string
  }): Promise<PageResponse<any>> {
    const response = await api.get('/admin/notifications', { params })
    return response.data
  }
}
