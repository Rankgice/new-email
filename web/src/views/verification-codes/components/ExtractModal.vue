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
              提取验证码
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
        <div class="bg-white dark:bg-gray-800 px-6 py-4">
          <!-- 选择提取方式 -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              提取方式
            </label>
            <div class="flex space-x-4">
              <label class="flex items-center">
                <input
                  v-model="extractMode"
                  type="radio"
                  value="single"
                  class="form-radio"
                />
                <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">单个邮件</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="extractMode"
                  type="radio"
                  value="batch"
                  class="form-radio"
                />
                <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">批量提取</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="extractMode"
                  type="radio"
                  value="auto"
                  class="form-radio"
                />
                <span class="ml-2 text-sm text-gray-700 dark:text-gray-300">自动提取</span>
              </label>
            </div>
          </div>

          <!-- 单个邮件提取 -->
          <div v-if="extractMode === 'single'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                邮件ID
              </label>
              <input
                v-model="singleEmailId"
                type="number"
                placeholder="请输入邮件ID"
                class="input w-full"
              />
            </div>
            
            <div class="flex justify-end space-x-3">
              <button
                @click="$emit('close')"
                class="btn btn-ghost"
              >
                取消
              </button>
              <button
                @click="extractSingle"
                :disabled="!singleEmailId || loading"
                class="btn btn-primary"
              >
                <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
                提取验证码
              </button>
            </div>
          </div>

          <!-- 批量提取 -->
          <div v-if="extractMode === 'batch'" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                邮件筛选条件
              </label>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">邮箱</label>
                  <select v-model="batchParams.mailboxId" class="input w-full">
                    <option value="">全部邮箱</option>
                    <option v-for="mailbox in mailboxes" :key="mailbox.id" :value="mailbox.id">
                      {{ mailbox.email }}
                    </option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">发件人</label>
                  <input
                    v-model="batchParams.fromEmail"
                    type="text"
                    placeholder="发件人邮箱"
                    class="input w-full"
                  />
                </div>
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">开始日期</label>
                  <input
                    v-model="batchParams.startDate"
                    type="date"
                    class="input w-full"
                  />
                </div>
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">结束日期</label>
                  <input
                    v-model="batchParams.endDate"
                    type="date"
                    class="input w-full"
                  />
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                提取限制
              </label>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">最大邮件数</label>
                  <input
                    v-model="batchParams.maxEmails"
                    type="number"
                    min="1"
                    max="1000"
                    placeholder="最多处理邮件数量"
                    class="input w-full"
                  />
                </div>
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">最小置信度</label>
                  <input
                    v-model="batchParams.minConfidence"
                    type="number"
                    min="0"
                    max="100"
                    placeholder="最小置信度 (0-100)"
                    class="input w-full"
                  />
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-3">
              <button
                @click="$emit('close')"
                class="btn btn-ghost"
              >
                取消
              </button>
              <button
                @click="previewBatch"
                :disabled="loading"
                class="btn btn-secondary"
              >
                <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
                预览邮件
              </button>
              <button
                @click="extractBatch"
                :disabled="loading"
                class="btn btn-primary"
              >
                <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
                批量提取
              </button>
            </div>
          </div>

          <!-- 自动提取 -->
          <div v-if="extractMode === 'auto'" class="space-y-4">
            <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <i class="fas fa-info-circle text-blue-400"></i>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-blue-800 dark:text-blue-200">
                    自动提取说明
                  </h3>
                  <div class="mt-2 text-sm text-blue-700 dark:text-blue-300">
                    <p>自动提取将扫描最近的邮件，智能识别并提取验证码。</p>
                    <ul class="list-disc list-inside mt-2 space-y-1">
                      <li>扫描最近7天的邮件</li>
                      <li>只处理未提取过的邮件</li>
                      <li>自动过滤低置信度结果</li>
                      <li>支持多种验证码格式</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                自动提取设置
              </label>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">扫描天数</label>
                  <select v-model="autoParams.scanDays" class="input w-full">
                    <option value="1">最近1天</option>
                    <option value="3">最近3天</option>
                    <option value="7">最近7天</option>
                    <option value="30">最近30天</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm text-gray-600 dark:text-gray-400 mb-1">最小置信度</label>
                  <select v-model="autoParams.minConfidence" class="input w-full">
                    <option value="50">50% (宽松)</option>
                    <option value="70">70% (标准)</option>
                    <option value="90">90% (严格)</option>
                  </select>
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-3">
              <button
                @click="$emit('close')"
                class="btn btn-ghost"
              >
                取消
              </button>
              <button
                @click="extractAuto"
                :disabled="loading"
                class="btn btn-primary"
              >
                <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
                开始自动提取
              </button>
            </div>
          </div>

          <!-- 预览结果 -->
          <div v-if="previewEmails.length > 0" class="mt-6">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              预览邮件 ({{ previewEmails.length }} 封)
            </h4>
            <div class="max-h-60 overflow-y-auto border border-gray-200 dark:border-gray-700 rounded-lg">
              <table class="w-full text-sm">
                <thead class="bg-gray-50 dark:bg-gray-700">
                  <tr>
                    <th class="px-3 py-2 text-left">主题</th>
                    <th class="px-3 py-2 text-left">发件人</th>
                    <th class="px-3 py-2 text-left">时间</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
                  <tr v-for="email in previewEmails" :key="email.id">
                    <td class="px-3 py-2">{{ email.subject }}</td>
                    <td class="px-3 py-2">{{ email.fromEmail }}</td>
                    <td class="px-3 py-2">{{ formatDate(email.createdAt) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- 提取结果 -->
          <div v-if="extractResult" class="mt-6">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              提取结果
            </h4>
            <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <i class="fas fa-check-circle text-green-400"></i>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-green-800 dark:text-green-200">
                    提取完成
                  </h3>
                  <div class="mt-2 text-sm text-green-700 dark:text-green-300">
                    <p v-if="extractMode === 'single'">
                      从邮件中提取到 {{ extractResult.codes?.length || 0 }} 个验证码
                    </p>
                    <p v-else>
                      处理了 {{ extractResult.processedEmails }} 封邮件，
                      提取到 {{ extractResult.extractedCodes }} 个验证码
                    </p>
                    <div v-if="extractResult.errors?.length" class="mt-2">
                      <p class="font-medium">错误信息：</p>
                      <ul class="list-disc list-inside">
                        <li v-for="error in extractResult.errors" :key="error">{{ error }}</li>
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { verificationCodeApi } from '@/api/verification-code'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'
import type { 
  VerificationCodeExtractResponse, 
  VerificationCodeBatchExtractResponse,
  Email,
  Mailbox
} from '@/types'

// 事件定义
const emit = defineEmits<{
  close: []
  extracted: []
}>()

// 响应式数据
const extractMode = ref<'single' | 'batch' | 'auto'>('single')
const loading = ref(false)
const singleEmailId = ref<number>()
const previewEmails = ref<Email[]>([])
const extractResult = ref<VerificationCodeExtractResponse | VerificationCodeBatchExtractResponse | null>(null)
const mailboxes = ref<Mailbox[]>([])

// 批量提取参数
const batchParams = reactive({
  mailboxId: '',
  fromEmail: '',
  startDate: '',
  endDate: '',
  maxEmails: 100,
  minConfidence: 70
})

// 自动提取参数
const autoParams = reactive({
  scanDays: 7,
  minConfidence: 70
})

// 通知
const { showSuccess, showError } = useNotification()

// 生命周期
onMounted(() => {
  loadMailboxes()
  setDefaultDates()
})

// 方法
async function loadMailboxes() {
  try {
    const response = await mailboxApi.list({ pageSize: 100 })
    mailboxes.value = response.list
  } catch (error) {
    console.error('Load mailboxes error:', error)
  }
}

function setDefaultDates() {
  const today = new Date()
  const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000)
  
  batchParams.endDate = today.toISOString().split('T')[0]
  batchParams.startDate = weekAgo.toISOString().split('T')[0]
}

async function extractSingle() {
  if (!singleEmailId.value) {
    showError('请输入邮件ID')
    return
  }

  try {
    loading.value = true
    const result = await verificationCodeApi.extract({ emailId: singleEmailId.value })
    extractResult.value = result
    
    if (result.codes.length > 0) {
      showSuccess(`成功提取到 ${result.codes.length} 个验证码`)
      emit('extracted')
    } else {
      showError('未找到验证码')
    }
  } catch (error) {
    showError('提取失败')
    console.error('Extract single error:', error)
  } finally {
    loading.value = false
  }
}

async function previewBatch() {
  try {
    loading.value = true
    const params = {
      mailboxId: batchParams.mailboxId ? Number(batchParams.mailboxId) : undefined,
      fromEmail: batchParams.fromEmail || undefined,
      createdAtStart: batchParams.startDate || undefined,
      createdAtEnd: batchParams.endDate || undefined,
      pageSize: batchParams.maxEmails || 100
    }
    
    const response = await emailApi.list(params)
    previewEmails.value = response.list
    
    if (previewEmails.value.length === 0) {
      showError('没有找到符合条件的邮件')
    }
  } catch (error) {
    showError('预览失败')
    console.error('Preview batch error:', error)
  } finally {
    loading.value = false
  }
}

async function extractBatch() {
  if (previewEmails.value.length === 0) {
    showError('请先预览邮件')
    return
  }

  try {
    loading.value = true
    const emailIds = previewEmails.value.map(email => email.id)
    const result = await verificationCodeApi.batchExtract({ emailIds })
    extractResult.value = result
    
    showSuccess(`批量提取完成，处理了 ${result.processedEmails} 封邮件，提取到 ${result.extractedCodes} 个验证码`)
    emit('extracted')
  } catch (error) {
    showError('批量提取失败')
    console.error('Extract batch error:', error)
  } finally {
    loading.value = false
  }
}

async function extractAuto() {
  try {
    loading.value = true
    
    // 获取最近的邮件
    const endDate = new Date()
    const startDate = new Date(endDate.getTime() - autoParams.scanDays * 24 * 60 * 60 * 1000)
    
    const response = await emailApi.list({
      createdAtStart: startDate.toISOString(),
      createdAtEnd: endDate.toISOString(),
      pageSize: 500
    })
    
    if (response.list.length === 0) {
      showError('没有找到符合条件的邮件')
      return
    }
    
    // 批量提取
    const emailIds = response.list.map(email => email.id)
    const result = await verificationCodeApi.batchExtract({ emailIds })
    
    // 过滤低置信度结果
    const filteredResults = result.results.map(r => ({
      ...r,
      codes: r.codes.filter(c => c.confidence >= autoParams.minConfidence)
    })).filter(r => r.codes.length > 0)
    
    extractResult.value = {
      ...result,
      results: filteredResults,
      extractedCodes: filteredResults.reduce((sum, r) => sum + r.codes.length, 0)
    }
    
    showSuccess(`自动提取完成，处理了 ${result.processedEmails} 封邮件，提取到 ${extractResult.value.extractedCodes} 个验证码`)
    emit('extracted')
  } catch (error) {
    showError('自动提取失败')
    console.error('Extract auto error:', error)
  } finally {
    loading.value = false
  }
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>

<style scoped>
.form-radio {
  @apply h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 dark:bg-gray-700;
}

.input {
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:placeholder-gray-400;
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

.btn-ghost {
  @apply text-gray-700 bg-transparent hover:bg-gray-100 focus:ring-gray-500 dark:text-gray-300 dark:hover:bg-gray-800;
}
</style>
