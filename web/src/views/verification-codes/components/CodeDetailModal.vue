<template>
  <div class="fixed inset-0 z-50 overflow-y-auto">
    <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <!-- 背景遮罩 -->
      <div
        class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        @click="$emit('close')"
      ></div>

      <!-- 模态框内容 -->
      <div class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
        <!-- 头部 -->
        <div class="bg-white dark:bg-gray-800 px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">
              验证码详情
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
          <div class="space-y-6">
            <!-- 验证码信息 -->
            <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-lg font-bold text-blue-900 dark:text-blue-100 font-mono">
                    {{ code.code }}
                  </h4>
                  <p class="text-sm text-blue-700 dark:text-blue-300 mt-1">
                    {{ getTypeDisplayName(code.type) }}
                  </p>
                </div>
                <div class="flex items-center space-x-2">
                  <button
                    @click="copyCode"
                    class="btn btn-sm btn-secondary"
                    title="复制验证码"
                  >
                    <i class="fas fa-copy mr-1"></i>
                    复制
                  </button>
                  <button
                    @click="toggleUsedStatus"
                    class="btn btn-sm"
                    :class="code.isUsed ? 'btn-warning' : 'btn-success'"
                  >
                    <i :class="code.isUsed ? 'fas fa-undo mr-1' : 'fas fa-check mr-1'"></i>
                    {{ code.isUsed ? '标记未使用' : '标记已使用' }}
                  </button>
                </div>
              </div>
            </div>

            <!-- 基本信息 -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  来源邮箱
                </label>
                <div class="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-gray-700 rounded px-3 py-2">
                  {{ code.source }}
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  置信度
                </label>
                <div class="flex items-center">
                  <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2 mr-3">
                    <div
                      class="h-2 rounded-full transition-all duration-300"
                      :class="getConfidenceColor(code.confidence)"
                      :style="{ width: `${code.confidence}%` }"
                    ></div>
                  </div>
                  <span class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ code.confidence }}%
                  </span>
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  状态
                </label>
                <div class="flex space-x-2">
                  <span
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="code.isUsed ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'"
                  >
                    {{ code.isUsed ? '已使用' : '未使用' }}
                  </span>
                  <span
                    v-if="code.isExpired"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
                  >
                    已过期
                  </span>
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  匹配模式
                </label>
                <div class="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-gray-700 rounded px-3 py-2">
                  {{ code.pattern || '未知' }}
                </div>
              </div>
            </div>

            <!-- 时间信息 -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  创建时间
                </label>
                <div class="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-gray-700 rounded px-3 py-2">
                  {{ formatDate(code.createdAt) }}
                </div>
              </div>
              
              <div v-if="code.usedAt">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  使用时间
                </label>
                <div class="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-gray-700 rounded px-3 py-2">
                  {{ formatDate(code.usedAt) }}
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  过期时间
                </label>
                <div class="text-sm text-gray-900 dark:text-white bg-gray-50 dark:bg-gray-700 rounded px-3 py-2">
                  {{ formatDate(code.expiresAt) }}
                </div>
              </div>
            </div>

            <!-- 上下文信息 -->
            <div v-if="code.context">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                上下文信息
              </label>
              <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <p class="text-sm text-gray-900 dark:text-white leading-relaxed">
                  {{ code.context }}
                </p>
              </div>
            </div>

            <!-- 描述信息 -->
            <div v-if="code.description">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                描述信息
              </label>
              <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <p class="text-sm text-gray-900 dark:text-white">
                  {{ code.description }}
                </p>
              </div>
            </div>

            <!-- 相关邮件信息 -->
            <div v-if="emailInfo" class="border-t border-gray-200 dark:border-gray-700 pt-6">
              <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                相关邮件信息
              </h4>
              <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4 space-y-3">
                <div>
                  <span class="text-xs text-gray-500 dark:text-gray-400">主题：</span>
                  <span class="text-sm text-gray-900 dark:text-white ml-2">
                    {{ emailInfo.subject }}
                  </span>
                </div>
                <div>
                  <span class="text-xs text-gray-500 dark:text-gray-400">发件人：</span>
                  <span class="text-sm text-gray-900 dark:text-white ml-2">
                    {{ emailInfo.fromEmail }}
                  </span>
                </div>
                <div>
                  <span class="text-xs text-gray-500 dark:text-gray-400">时间：</span>
                  <span class="text-sm text-gray-900 dark:text-white ml-2">
                    {{ formatDate(emailInfo.createdAt) }}
                  </span>
                </div>
                <div class="flex justify-end">
                  <button
                    @click="viewEmail"
                    class="btn btn-sm btn-secondary"
                  >
                    <i class="fas fa-envelope mr-1"></i>
                    查看邮件
                  </button>
                </div>
              </div>
            </div>

            <!-- 操作历史 -->
            <div v-if="operationHistory.length > 0" class="border-t border-gray-200 dark:border-gray-700 pt-6">
              <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                操作历史
              </h4>
              <div class="space-y-2">
                <div
                  v-for="operation in operationHistory"
                  :key="operation.id"
                  class="flex items-center justify-between text-sm bg-gray-50 dark:bg-gray-700 rounded px-3 py-2"
                >
                  <div>
                    <span class="text-gray-900 dark:text-white">{{ operation.action }}</span>
                    <span class="text-gray-500 dark:text-gray-400 ml-2">
                      {{ formatDate(operation.createdAt) }}
                    </span>
                  </div>
                  <span class="text-xs text-gray-500 dark:text-gray-400">
                    {{ operation.ip }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 底部 -->
        <div class="bg-gray-50 dark:bg-gray-900 px-6 py-4">
          <div class="flex justify-between">
            <button
              @click="deleteCode"
              class="btn btn-danger"
            >
              <i class="fas fa-trash mr-2"></i>
              删除验证码
            </button>
            <button
              @click="$emit('close')"
              class="btn btn-secondary"
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
import { ref, onMounted } from 'vue'
import { verificationCodeApi } from '@/api/verification-code'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'
import { useRouter } from 'vue-router'
import type { VerificationCode, Email, OperationLog } from '@/types'

// Props
const props = defineProps<{
  code: VerificationCode
}>()

// 事件定义
const emit = defineEmits<{
  close: []
  updated: []
}>()

// 响应式数据
const emailInfo = ref<Email | null>(null)
const operationHistory = ref<OperationLog[]>([])
const loading = ref(false)

// 组合式API
const { showSuccess, showError } = useNotification()
const router = useRouter()

// 生命周期
onMounted(() => {
  loadEmailInfo()
  loadOperationHistory()
})

// 方法
async function loadEmailInfo() {
  try {
    emailInfo.value = await emailApi.getById(props.code.emailId)
  } catch (error) {
    console.error('Load email info error:', error)
  }
}

async function loadOperationHistory() {
  try {
    // 这里需要实现获取操作历史的API
    // operationHistory.value = await operationLogApi.getByResource('verification_code', props.code.id)
  } catch (error) {
    console.error('Load operation history error:', error)
  }
}

async function copyCode() {
  try {
    await navigator.clipboard.writeText(props.code.code)
    showSuccess('验证码已复制到剪贴板')
  } catch (error) {
    showError('复制失败')
    console.error('Copy error:', error)
  }
}

async function toggleUsedStatus() {
  try {
    loading.value = true
    await verificationCodeApi.markUsed(props.code.id, !props.code.isUsed)
    
    // 更新本地状态
    props.code.isUsed = !props.code.isUsed
    if (props.code.isUsed) {
      props.code.usedAt = new Date().toISOString()
    } else {
      props.code.usedAt = undefined
    }
    
    showSuccess(`已${props.code.isUsed ? '标记为已使用' : '标记为未使用'}`)
    emit('updated')
  } catch (error) {
    showError('操作失败')
    console.error('Toggle used status error:', error)
  } finally {
    loading.value = false
  }
}

async function deleteCode() {
  if (!confirm('确定要删除这个验证码吗？')) {
    return
  }
  
  try {
    loading.value = true
    await verificationCodeApi.delete(props.code.id)
    showSuccess('删除成功')
    emit('updated')
    emit('close')
  } catch (error) {
    showError('删除失败')
    console.error('Delete code error:', error)
  } finally {
    loading.value = false
  }
}

function viewEmail() {
  if (emailInfo.value) {
    router.push(`/emails/${emailInfo.value.id}`)
    emit('close')
  }
}

function getTypeDisplayName(type?: string): string {
  const typeMap: Record<string, string> = {
    'numeric_6': '6位数字验证码',
    'numeric_4': '4位数字验证码',
    'numeric_8': '8位数字验证码',
    'alphanumeric_6': '6位字母数字验证码',
    'alphanumeric_4': '4位字母数字验证码',
    'numeric_general': '数字验证码',
    'sms_format': '短信验证码',
    'login_code': '登录验证码',
    'register_code': '注册验证码',
    'reset_password_code': '重置密码验证码',
    'otp': '动态密码',
    'security_code': '安全码'
  }
  return typeMap[type || ''] || type || '未知类型'
}

function getConfidenceColor(confidence?: number): string {
  if (!confidence) return 'bg-gray-400'
  if (confidence >= 90) return 'bg-green-500'
  if (confidence >= 70) return 'bg-yellow-500'
  if (confidence >= 50) return 'bg-orange-500'
  return 'bg-red-500'
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.btn {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors;
}

.btn-primary {
  @apply text-white bg-blue-600 hover:bg-blue-700 focus:ring-blue-500;
}

.btn-secondary {
  @apply text-gray-700 bg-gray-200 hover:bg-gray-300 focus:ring-gray-500 dark:text-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600;
}

.btn-success {
  @apply text-white bg-green-600 hover:bg-green-700 focus:ring-green-500;
}

.btn-warning {
  @apply text-white bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500;
}

.btn-danger {
  @apply text-white bg-red-600 hover:bg-red-700 focus:ring-red-500;
}

.btn-sm {
  @apply px-3 py-1.5 text-xs;
}
</style>
