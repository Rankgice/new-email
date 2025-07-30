import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { User } from '@/types'
import { api } from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isLoading = ref(false)
  const isInitialized = ref(false)
  const redirectPath = ref<string>('/inbox')

  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const userInitials = computed(() => {
    if (!user.value) return ''
    const name = user.value.nickname || user.value.username
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
  })

  // 初始化认证状态
  const initAuth = async () => {
    const savedToken = localStorage.getItem('auth_token')
    const savedUser = localStorage.getItem('auth_user')

    if (savedToken && savedUser) {
      try {
        token.value = savedToken
        user.value = JSON.parse(savedUser)

        // 验证 token 有效性
        const isValid = await validateToken()
        if (!isValid) {
          clearAuth()
        }
      } catch (error) {
        console.error('Failed to parse saved user data:', error)
        clearAuth()
      }
    }

    isInitialized.value = true
  }

  // 验证 token 有效性
  const validateToken = async () => {
    if (!token.value) return false
    
    try {
      const response = await api.validateToken()
      if (!response.success) {
        clearAuth()
        return false
      }
      return true
    } catch (error) {
      console.error('Token validation failed:', error)
      clearAuth()
      return false
    }
  }

  // 登录
  const login = async (credentials: { username: string; password: string }) => {
    isLoading.value = true
    
    try {
      const response = await api.login(credentials)
      
      if (response.success && response.data) {
        token.value = response.data.token
        user.value = response.data.user
        
        // 保存到本地存储
        localStorage.setItem('auth_token', response.data.token)
        localStorage.setItem('auth_user', JSON.stringify(response.data.user))
        
        return { success: true }
      } else {
        return { 
          success: false, 
          message: response.message || '登录失败' 
        }
      }
    } catch (error: any) {
      console.error('Login error:', error)
      return { 
        success: false, 
        message: error.message || '网络错误，请稍后重试' 
      }
    } finally {
      isLoading.value = false
    }
  }

  // 注册
  const register = async (userData: {
    username: string
    email: string
    password: string
    nickname?: string
  }) => {
    isLoading.value = true
    
    try {
      const response = await api.register(userData)
      
      if (response.success) {
        return { success: true, message: '注册成功，请登录' }
      } else {
        return { 
          success: false, 
          message: response.message || '注册失败' 
        }
      }
    } catch (error: any) {
      console.error('Register error:', error)
      return { 
        success: false, 
        message: error.message || '网络错误，请稍后重试' 
      }
    } finally {
      isLoading.value = false
    }
  }

  // 登出
  const logout = async () => {
    try {
      if (token.value) {
        await api.logout()
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuth()
    }
  }

  // 清除认证信息
  const clearAuth = () => {
    user.value = null
    token.value = null
    isInitialized.value = true
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
  }

  // 更新用户信息
  const updateUser = (userData: Partial<User>) => {
    if (user.value) {
      user.value = { ...user.value, ...userData }
      localStorage.setItem('auth_user', JSON.stringify(user.value))
    }
  }

  // 设置重定向路径
  const setRedirectPath = (path: string) => {
    redirectPath.value = path
  }

  // 获取并清除重定向路径
  const getAndClearRedirectPath = () => {
    const path = redirectPath.value
    redirectPath.value = '/inbox'
    return path
  }

  // 忘记密码
  const forgotPassword = async (email: string) => {
    isLoading.value = true
    
    try {
      const response = await api.forgotPassword({ email })
      
      if (response.success) {
        return { success: true, message: '重置密码邮件已发送' }
      } else {
        return { 
          success: false, 
          message: response.message || '发送失败' 
        }
      }
    } catch (error: any) {
      console.error('Forgot password error:', error)
      return { 
        success: false, 
        message: error.message || '网络错误，请稍后重试' 
      }
    } finally {
      isLoading.value = false
    }
  }

  // 重置密码
  const resetPassword = async (data: {
    token: string
    password: string
    confirmPassword: string
  }) => {
    isLoading.value = true
    
    try {
      const response = await api.resetPassword(data)
      
      if (response.success) {
        return { success: true, message: '密码重置成功' }
      } else {
        return { 
          success: false, 
          message: response.message || '重置失败' 
        }
      }
    } catch (error: any) {
      console.error('Reset password error:', error)
      return { 
        success: false, 
        message: error.message || '网络错误，请稍后重试' 
      }
    } finally {
      isLoading.value = false
    }
  }

  return {
    // 状态
    user: readonly(user),
    token: readonly(token),
    isLoading: readonly(isLoading),
    isInitialized: readonly(isInitialized),
    redirectPath: readonly(redirectPath),

    // 计算属性
    isAuthenticated,
    isAdmin,
    userInitials,

    // 方法
    initAuth,
    validateToken,
    login,
    register,
    logout,
    clearAuth,
    updateUser,
    setRedirectPath,
    getAndClearRedirectPath,
    forgotPassword,
    resetPassword
  }
})
