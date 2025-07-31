import { api } from '@/utils/api'
import type {
  VerificationCode,
  VerificationCodeExtractRequest,
  VerificationCodeExtractResponse,
  VerificationCodeBatchExtractRequest,
  VerificationCodeBatchExtractResponse,
  VerificationCodeStats,
  VerificationCodeListParams,
  PageResponse
} from '@/types'

// 验证码API服务
export const verificationCodeApi = {
  // 获取验证码列表
  async list(params: VerificationCodeListParams): Promise<PageResponse<VerificationCode>> {
    const response = await api.get('/user/verification-codes', { params })
    return response.data
  },

  // 获取验证码详情
  async getById(id: number): Promise<VerificationCode> {
    const response = await api.get(`/user/verification-codes/${id}`)
    return response.data
  },

  // 从邮件中提取验证码
  async extract(data: VerificationCodeExtractRequest): Promise<VerificationCodeExtractResponse> {
    const response = await api.post('/user/verification-codes/extract', data)
    return response.data
  },

  // 批量提取验证码
  async batchExtract(data: VerificationCodeBatchExtractRequest): Promise<VerificationCodeBatchExtractResponse> {
    const response = await api.post('/user/verification-codes/batch-extract', data)
    return response.data
  },

  // 标记验证码为已使用/未使用
  async markUsed(id: number, used: boolean): Promise<void> {
    await api.put(`/user/verification-codes/${id}/used`, { used })
  },

  // 获取验证码统计信息
  async getStats(): Promise<VerificationCodeStats> {
    const response = await api.get('/user/verification-codes/stats')
    return response.data
  },

  // 获取最新验证码
  async getLatest(source: string): Promise<VerificationCode> {
    const response = await api.get('/user/verification-codes/latest', {
      params: { source }
    })
    return response.data
  },

  // 搜索验证码
  async search(params: {
    query?: string
    mailboxId?: number
    codeType?: string
    used?: boolean
    startDate?: string
    endDate?: string
    page?: number
    pageSize?: number
  }): Promise<PageResponse<VerificationCode>> {
    const response = await api.get('/user/verification-codes', { params })
    return response.data
  },

  // 删除验证码
  async delete(id: number): Promise<void> {
    await api.delete(`/user/verification-codes/${id}`)
  },

  // 批量删除验证码
  async batchDelete(ids: number[]): Promise<void> {
    await api.post('/user/verification-codes/batch-delete', { ids })
  },

  // 导出验证码
  async export(params: {
    format?: 'csv' | 'json'
    mailboxId?: number
    source?: string
    used?: boolean
    startDate?: string
    endDate?: string
  }): Promise<{
    fileName: string
    fileSize: number
    recordCount: number
    downloadUrl: string
  }> {
    const response = await api.get('/user/verification-codes/export', { params })
    return response.data
  },

  // 获取验证码类型统计
  async getTypeStats(): Promise<Array<{ type: string; count: number; description: string }>> {
    const response = await api.get('/user/verification-codes/type-stats')
    return response.data
  },

  // 获取验证码来源统计
  async getSourceStats(): Promise<Array<{ source: string; count: number; lastCode: string; lastTime: string }>> {
    const response = await api.get('/user/verification-codes/source-stats')
    return response.data
  },

  // 清理过期验证码
  async cleanExpired(): Promise<{ deletedCount: number }> {
    const response = await api.post('/user/verification-codes/clean-expired')
    return response.data
  },

  // 验证验证码格式
  async validateCode(code: string): Promise<{ valid: boolean; message?: string }> {
    const response = await api.post('/user/verification-codes/validate', { code })
    return response.data
  },

  // 获取支持的验证码类型
  async getSupportedTypes(): Promise<Array<{ type: string; description: string; pattern: string }>> {
    const response = await api.get('/user/verification-codes/supported-types')
    return response.data
  }
}
