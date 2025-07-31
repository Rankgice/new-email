<template>
  <div class="admin-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">管理员仪表板</h1>
          <p class="text-gray-600 dark:text-gray-400 mt-1">系统概览和关键指标</p>
        </div>
        <div class="flex items-center space-x-3">
          <button
            @click="refreshData"
            :disabled="loading"
            class="btn btn-secondary"
          >
            <i :class="loading ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'" class="mr-2"></i>
            刷新数据
          </button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- 用户统计 -->
      <div class="stats-card">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-blue-100 dark:bg-blue-900">
            <i class="fas fa-users text-blue-600 dark:text-blue-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总用户数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.userStats.totalUsers }}</p>
            <p class="text-xs text-green-600 dark:text-green-400">
              <i class="fas fa-arrow-up"></i>
              今日新增 {{ stats.userStats.newUsers }}
            </p>
          </div>
        </div>
      </div>

      <!-- 邮件统计 -->
      <div class="stats-card">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-green-100 dark:bg-green-900">
            <i class="fas fa-envelope text-green-600 dark:text-green-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总邮件数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.emailStats.totalEmails }}</p>
            <p class="text-xs text-green-600 dark:text-green-400">
              <i class="fas fa-arrow-up"></i>
              今日新增 {{ stats.emailStats.todayEmails }}
            </p>
          </div>
        </div>
      </div>

      <!-- 邮箱统计 -->
      <div class="stats-card">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-purple-100 dark:bg-purple-900">
            <i class="fas fa-inbox text-purple-600 dark:text-purple-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总邮箱数</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.mailboxStats.totalMailboxes }}</p>
            <p class="text-xs text-green-600 dark:text-green-400">
              <i class="fas fa-check"></i>
              活跃 {{ stats.mailboxStats.activeMailboxes }}
            </p>
          </div>
        </div>
      </div>

      <!-- 系统状态 -->
      <div class="stats-card">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-yellow-100 dark:bg-yellow-900">
            <i class="fas fa-server text-yellow-600 dark:text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">系统状态</p>
            <p class="text-lg font-bold text-green-600 dark:text-green-400">正常运行</p>
            <p class="text-xs text-gray-600 dark:text-gray-400">
              运行时间 {{ stats.systemStats.uptime }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 系统资源使用情况 -->
    <div class="glass-card p-6 mb-8">
      <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">系统资源使用情况</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- CPU使用率 -->
        <div>
          <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-2">
            <span>CPU使用率</span>
            <span>{{ stats.systemStats.cpuUsage.toFixed(1) }}%</span>
          </div>
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
            <div
              class="h-3 rounded-full transition-all duration-300"
              :class="getUsageColor(stats.systemStats.cpuUsage)"
              :style="{ width: `${stats.systemStats.cpuUsage}%` }"
            ></div>
          </div>
        </div>

        <!-- 内存使用率 -->
        <div>
          <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-2">
            <span>内存使用率</span>
            <span>{{ stats.systemStats.memUsage.toFixed(1) }}%</span>
          </div>
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
            <div
              class="h-3 rounded-full transition-all duration-300"
              :class="getUsageColor(stats.systemStats.memUsage)"
              :style="{ width: `${stats.systemStats.memUsage}%` }"
            ></div>
          </div>
        </div>

        <!-- 磁盘使用率 -->
        <div>
          <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-2">
            <span>磁盘使用率</span>
            <span>{{ stats.systemStats.diskUsage.toFixed(1) }}%</span>
          </div>
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-3">
            <div
              class="h-3 rounded-full transition-all duration-300"
              :class="getUsageColor(stats.systemStats.diskUsage)"
              :style="{ width: `${stats.systemStats.diskUsage}%` }"
            ></div>
          </div>
        </div>
      </div>

      <!-- 系统信息 -->
      <div class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
        <h4 class="text-sm font-medium text-gray-900 dark:text-white mb-3">系统信息</h4>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div>
            <span class="text-gray-600 dark:text-gray-400">版本:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ stats.systemStats.version }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">Go版本:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ stats.systemStats.goVersion }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">平台:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ stats.systemStats.platform }}</span>
          </div>
          <div>
            <span class="text-gray-600 dark:text-gray-400">启动时间:</span>
            <span class="ml-2 text-gray-900 dark:text-white">{{ formatDate(stats.systemStats.startTime) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminApi } from '@/api/admin'
import { useNotification } from '@/composables/useNotification'
import type { AdminSystemStats } from '@/types'

// 响应式数据
const loading = ref(false)

// 仪表板数据
const stats = ref<AdminSystemStats>({
  userStats: { totalUsers: 0, activeUsers: 0, newUsers: 0, onlineUsers: 0 },
  emailStats: { totalEmails: 0, todayEmails: 0, sentEmails: 0, receivedEmails: 0 },
  mailboxStats: { totalMailboxes: 0, activeMailboxes: 0, imapMailboxes: 0, pop3Mailboxes: 0 },
  systemStats: { version: '', startTime: '', uptime: '', goVersion: '', platform: '', cpuUsage: 0, memUsage: 0, diskUsage: 0 }
})

// 通知
const { showSuccess, showError } = useNotification()

// 生命周期
onMounted(() => {
  loadDashboardData()
})

// 方法
async function loadDashboardData() {
  try {
    loading.value = true
    const dashboardStats = await adminApi.getSystemStats()
    stats.value = dashboardStats
  } catch (error) {
    showError('加载仪表板数据失败')
    console.error('Load dashboard error:', error)
  } finally {
    loading.value = false
  }
}

async function refreshData() {
  await loadDashboardData()
  showSuccess('数据已刷新')
}

function getUsageColor(usage: number): string {
  if (usage < 50) return 'bg-green-500'
  if (usage < 80) return 'bg-yellow-500'
  return 'bg-red-500'
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.admin-dashboard {
  @apply p-6 space-y-6;
}

.page-header {
  @apply mb-6;
}

.stats-card {
  @apply bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm border border-gray-200/50 dark:border-gray-700/50 rounded-lg shadow-sm p-6;
}

.glass-card {
  @apply bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm border border-gray-200/50 dark:border-gray-700/50 rounded-lg shadow-sm;
}

.btn {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors;
}

.btn-primary {
  @apply text-white bg-blue-600 hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply text-gray-700 bg-gray-200 hover:bg-gray-300 focus:ring-gray-500 dark:text-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600;
}
</style>
