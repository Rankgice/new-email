import { ref, readonly } from 'vue'
import type { Notification } from '@/types'

// 全局通知状态
const notifications = ref<Notification[]>([])

export const useNotification = () => {
  // 显示通知
  const showNotification = (notification: Omit<Notification, 'id' | 'createdAt'>) => {
    const id = `notification-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    const createdAt = new Date().toISOString()
    
    const newNotification: Notification = {
      id,
      createdAt,
      duration: 5000, // 默认5秒
      ...notification
    }
    
    notifications.value.push(newNotification)
    
    // 自动移除通知
    if (newNotification.duration && newNotification.duration > 0) {
      setTimeout(() => {
        removeNotification(id)
      }, newNotification.duration)
    }
    
    return id
  }
  
  // 移除通知
  const removeNotification = (id: string) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }
  
  // 清除所有通知
  const clearNotifications = () => {
    notifications.value = []
  }
  
  // 成功通知
  const success = (title: string, message?: string, duration?: number) => {
    return showNotification({
      type: 'success',
      title,
      message,
      duration
    })
  }
  
  // 错误通知
  const error = (title: string, message?: string, duration?: number) => {
    return showNotification({
      type: 'error',
      title,
      message,
      duration: duration || 8000 // 错误通知显示更久
    })
  }
  
  // 警告通知
  const warning = (title: string, message?: string, duration?: number) => {
    return showNotification({
      type: 'warning',
      title,
      message,
      duration
    })
  }
  
  // 信息通知
  const info = (title: string, message?: string, duration?: number) => {
    return showNotification({
      type: 'info',
      title,
      message,
      duration
    })
  }
  
  return {
    notifications: readonly(notifications),
    showNotification,
    removeNotification,
    clearNotifications,
    success,
    error,
    warning,
    info
  }
}
