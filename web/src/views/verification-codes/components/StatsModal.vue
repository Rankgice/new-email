<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <!-- 背景遮罩 -->
      <div
        class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        @click="$emit('close')"
      ></div>

      <!-- 模态框内容 -->
      <div class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
        <!-- 头部 -->
        <div class="bg-white dark:bg-gray-800 px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">
              验证码统计信息
            </h3>
            <button
              @click="$emit('close')"
              class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>

        <!-- 内容 -->
        <div class="bg-white dark:bg-gray-800 px-6 py-6">
          <!-- 总体统计 -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div class="text-center">
              <div class="text-3xl font-bold text-blue-600 dark:text-blue-400">
                {{ stats.totalCodes }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                总验证码数
              </div>
            </div>
            
            <div class="text-center">
              <div class="text-3xl font-bold text-green-600 dark:text-green-400">
                {{ stats.usedCodes }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                已使用
              </div>
              <div class="text-xs text-gray-500 dark:text-gray-500 mt-1">
                {{ getPercentage(stats.usedCodes, stats.totalCodes) }}%
              </div>
            </div>
            
            <div class="text-center">
              <div class="text-3xl font-bold text-yellow-600 dark:text-yellow-400">
                {{ stats.unusedCodes }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                未使用
              </div>
              <div class="text-xs text-gray-500 dark:text-gray-500 mt-1">
                {{ getPercentage(stats.unusedCodes, stats.totalCodes) }}%
              </div>
            </div>
            
            <div class="text-center">
              <div class="text-3xl font-bold text-purple-600 dark:text-purple-400">
                {{ stats.todayCodes }}
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">
                今日新增
              </div>
            </div>
          </div>

          <!-- 图表区域 -->
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <!-- 验证码类型统计 -->
            <div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-6">
              <h4 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
                验证码类型分布
              </h4>
              <div class="space-y-3">
                <div
                  v-for="typeStat in stats.typeStats"
                  :key="typeStat.type"
                  class="flex items-center justify-between"
                >
                  <div class="flex items-center">
                    <div
                      class="w-4 h-4 rounded mr-3"
                      :style="{ backgroundColor: getTypeColor(typeStat.type) }"
                    ></div>
                    <span class="text-sm text-gray-700 dark:text-gray-300">
                      {{ getTypeDisplayName(typeStat.type) }}
                    </span>
                  </div>
                  <div class="flex items-center space-x-2">
                    <span class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ typeStat.count }}
                    </span>
                    <span class="text-xs text-gray-500 dark:text-gray-400">
                      ({{ getPercentage(typeStat.count, stats.totalCodes) }}%)
                    </span>
                  </div>
                </div>
              </div>
              
              <!-- 类型分布条形图 -->
              <div class="mt-6 space-y-2">
                <div
                  v-for="typeStat in stats.typeStats"
                  :key="typeStat.type"
                  class="flex items-center"
                >
                  <div class="w-20 text-xs text-gray-600 dark:text-gray-400 truncate">
                    {{ getTypeDisplayName(typeStat.type) }}
                  </div>
                  <div class="flex-1 mx-3">
                    <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                      <div
                        class="h-2 rounded-full transition-all duration-300"
                        :style="{
                          width: `${getPercentage(typeStat.count, stats.totalCodes)}%`,
                          backgroundColor: getTypeColor(typeStat.type)
                        }"
                      ></div>
                    </div>
                  </div>
                  <div class="w-12 text-xs text-gray-600 dark:text-gray-400 text-right">
                    {{ typeStat.count }}
                  </div>
                </div>
              </div>
            </div>

            <!-- 验证码来源统计 -->
            <div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-6">
              <h4 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
                验证码来源分布
              </h4>
              <div class="space-y-3">
                <div
                  v-for="sourceStat in stats.sourceStats.slice(0, 10)"
                  :key="sourceStat.fromEmail"
                  class="flex items-center justify-between"
                >
                  <div class="flex items-center min-w-0 flex-1">
                    <div
                      class="w-4 h-4 rounded mr-3"
                      :style="{ backgroundColor: getSourceColor(sourceStat.fromEmail) }"
                    ></div>
                    <span class="text-sm text-gray-700 dark:text-gray-300 truncate">
                      {{ sourceStat.fromEmail }}
                    </span>
                  </div>
                  <div class="flex items-center space-x-2 ml-2">
                    <span class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ sourceStat.count }}
                    </span>
                    <span class="text-xs text-gray-500 dark:text-gray-400">
                      ({{ getPercentage(sourceStat.count, stats.totalCodes) }}%)
                    </span>
                  </div>
                </div>
              </div>
              
              <!-- 来源分布条形图 -->
              <div class="mt-6 space-y-2">
                <div
                  v-for="sourceStat in stats.sourceStats.slice(0, 8)"
                  :key="sourceStat.fromEmail"
                  class="flex items-center"
                >
                  <div class="w-24 text-xs text-gray-600 dark:text-gray-400 truncate">
                    {{ getShortEmail(sourceStat.fromEmail) }}
                  </div>
                  <div class="flex-1 mx-3">
                    <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                      <div
                        class="h-2 rounded-full transition-all duration-300"
                        :style="{
                          width: `${getPercentage(sourceStat.count, stats.totalCodes)}%`,
                          backgroundColor: getSourceColor(sourceStat.fromEmail)
                        }"
                      ></div>
                    </div>
                  </div>
                  <div class="w-12 text-xs text-gray-600 dark:text-gray-400 text-right">
                    {{ sourceStat.count }}
                  </div>
                </div>
              </div>
              
              <div v-if="stats.sourceStats.length > 10" class="mt-4 text-center">
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  还有 {{ stats.sourceStats.length - 10 }} 个来源...
                </span>
              </div>
            </div>
          </div>

          <!-- 使用率环形图 -->
          <div class="mt-8 bg-gray-50 dark:bg-gray-900 rounded-lg p-6">
            <h4 class="text-lg font-medium text-gray-900 dark:text-white mb-4 text-center">
              验证码使用率
            </h4>
            <div class="flex items-center justify-center">
              <div class="relative w-48 h-48">
                <!-- SVG 环形图 -->
                <svg class="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
                  <!-- 背景圆环 -->
                  <circle
                    cx="50"
                    cy="50"
                    r="40"
                    stroke="currentColor"
                    stroke-width="8"
                    fill="transparent"
                    class="text-gray-200 dark:text-gray-700"
                  />
                  <!-- 已使用部分 -->
                  <circle
                    cx="50"
                    cy="50"
                    r="40"
                    stroke="currentColor"
                    stroke-width="8"
                    fill="transparent"
                    class="text-green-500"
                    :stroke-dasharray="circumference"
                    :stroke-dashoffset="circumference - (usagePercentage / 100) * circumference"
                    stroke-linecap="round"
                  />
                </svg>
                <!-- 中心文字 -->
                <div class="absolute inset-0 flex items-center justify-center">
                  <div class="text-center">
                    <div class="text-2xl font-bold text-gray-900 dark:text-white">
                      {{ usagePercentage }}%
                    </div>
                    <div class="text-sm text-gray-600 dark:text-gray-400">
                      使用率
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 图例 -->
            <div class="flex justify-center space-x-8 mt-6">
              <div class="flex items-center">
                <div class="w-4 h-4 bg-green-500 rounded mr-2"></div>
                <span class="text-sm text-gray-700 dark:text-gray-300">
                  已使用 ({{ stats.usedCodes }})
                </span>
              </div>
              <div class="flex items-center">
                <div class="w-4 h-4 bg-gray-300 dark:bg-gray-600 rounded mr-2"></div>
                <span class="text-sm text-gray-700 dark:text-gray-300">
                  未使用 ({{ stats.unusedCodes }})
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 底部 -->
        <div class="bg-gray-50 dark:bg-gray-900 px-6 py-4">
          <div class="flex justify-end">
            <button
              @click="$emit('close')"
              class="btn btn-primary"
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { VerificationCodeStats } from '@/types'

// Props
const props = defineProps<{
  stats: VerificationCodeStats
}>()

// 事件定义
const emit = defineEmits<{
  close: []
}>()

// 计算属性
const usagePercentage = computed(() => {
  if (props.stats.totalCodes === 0) return 0
  return Math.round((props.stats.usedCodes / props.stats.totalCodes) * 100)
})

const circumference = computed(() => 2 * Math.PI * 40) // r=40

// 方法
function getPercentage(value: number, total: number): number {
  if (total === 0) return 0
  return Math.round((value / total) * 100)
}

function getTypeDisplayName(type: string): string {
  const typeMap: Record<string, string> = {
    'numeric_6': '6位数字',
    'numeric_4': '4位数字',
    'numeric_8': '8位数字',
    'alphanumeric_6': '6位字母数字',
    'alphanumeric_4': '4位字母数字',
    'numeric_general': '数字验证码',
    'sms_format': '短信验证码',
    'login_code': '登录验证码',
    'register_code': '注册验证码',
    'reset_password_code': '重置密码验证码',
    'otp': '动态密码',
    'security_code': '安全码'
  }
  return typeMap[type] || type
}

function getTypeColor(type: string): string {
  const colors = [
    '#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6',
    '#06B6D4', '#F97316', '#84CC16', '#EC4899', '#6366F1'
  ]
  const index = type.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return colors[index % colors.length]
}

function getSourceColor(email: string): string {
  const colors = [
    '#3B82F6', '#EF4444', '#10B981', '#F59E0B', '#8B5CF6',
    '#06B6D4', '#F97316', '#84CC16', '#EC4899', '#6366F1'
  ]
  const index = email.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return colors[index % colors.length]
}

function getShortEmail(email: string): string {
  if (email.length <= 20) return email
  const [local, domain] = email.split('@')
  if (local.length > 10) {
    return `${local.substring(0, 8)}...@${domain}`
  }
  return email
}
</script>

<style scoped>
.btn {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors;
}

.btn-primary {
  @apply text-white bg-blue-600 hover:bg-blue-700 focus:ring-blue-500;
}
</style>
