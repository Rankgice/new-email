<template>
  <div class="verification-codes-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">验证码管理</h1>
          <p class="text-gray-600 dark:text-gray-400 mt-1">管理和查看从邮件中提取的验证码</p>
        </div>
        <div class="flex items-center space-x-3">
          <button
            @click="showExtractModal = true"
            class="btn btn-primary"
          >
            <i class="fas fa-search mr-2"></i>
            提取验证码
          </button>
          <button
            @click="showStatsModal = true"
            class="btn btn-secondary"
          >
            <i class="fas fa-chart-bar mr-2"></i>
            统计信息
          </button>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
      <div class="glass-card p-6">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-blue-100 dark:bg-blue-900">
            <i class="fas fa-key text-blue-600 dark:text-blue-400"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">总验证码</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.totalCodes }}</p>
          </div>
        </div>
      </div>
      
      <div class="glass-card p-6">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-green-100 dark:bg-green-900">
            <i class="fas fa-check text-green-600 dark:text-green-400"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">已使用</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.usedCodes }}</p>
          </div>
        </div>
      </div>
      
      <div class="glass-card p-6">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-yellow-100 dark:bg-yellow-900">
            <i class="fas fa-clock text-yellow-600 dark:text-yellow-400"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">未使用</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.unusedCodes }}</p>
          </div>
        </div>
      </div>
      
      <div class="glass-card p-6">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-purple-100 dark:bg-purple-900">
            <i class="fas fa-plus text-purple-600 dark:text-purple-400"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600 dark:text-gray-400">今日新增</p>
            <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.todayCodes }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="glass-card p-6 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            验证码
          </label>
          <input
            v-model="searchParams.code"
            type="text"
            placeholder="搜索验证码..."
            class="input"
            @input="debouncedSearch"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            来源
          </label>
          <input
            v-model="searchParams.source"
            type="text"
            placeholder="发件人邮箱..."
            class="input"
            @input="debouncedSearch"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            状态
          </label>
          <select
            v-model="searchParams.isUsed"
            class="input"
            @change="loadCodes"
          >
            <option value="">全部</option>
            <option :value="false">未使用</option>
            <option :value="true">已使用</option>
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            过期状态
          </label>
          <select
            v-model="searchParams.isExpired"
            class="input"
            @change="loadCodes"
          >
            <option value="">全部</option>
            <option :value="false">未过期</option>
            <option :value="true">已过期</option>
          </select>
        </div>
      </div>
      
      <div class="flex items-center justify-between mt-4">
        <div class="flex items-center space-x-2">
          <button
            @click="resetSearch"
            class="btn btn-ghost btn-sm"
          >
            <i class="fas fa-undo mr-1"></i>
            重置
          </button>
          <button
            @click="exportCodes"
            class="btn btn-ghost btn-sm"
          >
            <i class="fas fa-download mr-1"></i>
            导出
          </button>
        </div>
        
        <div class="text-sm text-gray-600 dark:text-gray-400">
          共 {{ pagination.total }} 条记录
        </div>
      </div>
    </div>

    <!-- 验证码列表 -->
    <div class="glass-card">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                <input
                  type="checkbox"
                  :checked="selectedCodes.length === codes.length && codes.length > 0"
                  @change="toggleSelectAll"
                  class="rounded border-gray-300"
                />
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                验证码
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                来源
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                类型
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                置信度
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                状态
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                创建时间
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                操作
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
            <tr
              v-for="code in codes"
              :key="code.id"
              class="hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <input
                  type="checkbox"
                  :checked="selectedCodes.includes(code.id)"
                  @change="toggleSelectCode(code.id)"
                  class="rounded border-gray-300"
                />
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <span class="font-mono text-lg font-bold text-blue-600 dark:text-blue-400">
                    {{ code.code }}
                  </span>
                  <button
                    @click="copyCode(code.code)"
                    class="ml-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                    title="复制验证码"
                  >
                    <i class="fas fa-copy"></i>
                  </button>
                </div>
                <div v-if="code.context" class="text-sm text-gray-500 dark:text-gray-400 mt-1">
                  {{ code.context.substring(0, 50) }}...
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900 dark:text-white">{{ code.source }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                  {{ getTypeDisplayName(code.type) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-16 bg-gray-200 dark:bg-gray-700 rounded-full h-2 mr-2">
                    <div
                      class="h-2 rounded-full"
                      :class="getConfidenceColor(code.confidence)"
                      :style="{ width: `${code.confidence}%` }"
                    ></div>
                  </div>
                  <span class="text-sm text-gray-600 dark:text-gray-400">{{ code.confidence }}%</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex flex-col space-y-1">
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
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                {{ formatDate(code.createdAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex items-center space-x-2">
                  <button
                    @click="toggleUsedStatus(code)"
                    class="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                    :title="code.isUsed ? '标记为未使用' : '标记为已使用'"
                  >
                    <i :class="code.isUsed ? 'fas fa-undo' : 'fas fa-check'"></i>
                  </button>
                  <button
                    @click="viewCodeDetail(code)"
                    class="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300"
                    title="查看详情"
                  >
                    <i class="fas fa-eye"></i>
                  </button>
                  <button
                    @click="deleteCode(code.id)"
                    class="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                    title="删除"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 分页 -->
      <div class="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700 dark:text-gray-300">
            显示第 {{ (pagination.page - 1) * pagination.pageSize + 1 }} 到 
            {{ Math.min(pagination.page * pagination.pageSize, pagination.total) }} 条，
            共 {{ pagination.total }} 条记录
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="changePage(pagination.page - 1)"
              :disabled="pagination.page <= 1"
              class="btn btn-ghost btn-sm"
            >
              上一页
            </button>
            <span class="text-sm text-gray-600 dark:text-gray-400">
              {{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) }}
            </span>
            <button
              @click="changePage(pagination.page + 1)"
              :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)"
              class="btn btn-ghost btn-sm"
            >
              下一页
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 批量操作栏 -->
    <div
      v-if="selectedCodes.length > 0"
      class="fixed bottom-6 left-1/2 transform -translate-x-1/2 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 px-6 py-3"
    >
      <div class="flex items-center space-x-4">
        <span class="text-sm text-gray-600 dark:text-gray-400">
          已选择 {{ selectedCodes.length }} 项
        </span>
        <button
          @click="batchMarkUsed(true)"
          class="btn btn-sm btn-primary"
        >
          批量标记已使用
        </button>
        <button
          @click="batchMarkUsed(false)"
          class="btn btn-sm btn-secondary"
        >
          批量标记未使用
        </button>
        <button
          @click="batchDelete"
          class="btn btn-sm btn-danger"
        >
          批量删除
        </button>
        <button
          @click="selectedCodes = []"
          class="btn btn-sm btn-ghost"
        >
          取消选择
        </button>
      </div>
    </div>

    <!-- 提取验证码模态框 -->
    <ExtractModal
      v-if="showExtractModal"
      @close="showExtractModal = false"
      @extracted="onCodesExtracted"
    />

    <!-- 统计信息模态框 -->
    <StatsModal
      v-if="showStatsModal"
      :stats="stats"
      @close="showStatsModal = false"
    />

    <!-- 验证码详情模态框 -->
    <CodeDetailModal
      v-if="showDetailModal"
      :code="selectedCode"
      @close="showDetailModal = false"
      @updated="onCodeUpdated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { verificationCodeApi } from '@/api/verification-code'
import { useNotification } from '@/composables/useNotification'
import { debounce } from 'lodash-es'
import type { VerificationCode, VerificationCodeStats, VerificationCodeListParams } from '@/types'

// 组件引入
import ExtractModal from './components/ExtractModal.vue'
import StatsModal from './components/StatsModal.vue'
import CodeDetailModal from './components/CodeDetailModal.vue'

// 响应式数据
const codes = ref<VerificationCode[]>([])
const selectedCodes = ref<number[]>([])
const selectedCode = ref<VerificationCode | null>(null)
const loading = ref(false)
const showExtractModal = ref(false)
const showStatsModal = ref(false)
const showDetailModal = ref(false)

// 统计信息
const stats = ref<VerificationCodeStats>({
  totalCodes: 0,
  usedCodes: 0,
  unusedCodes: 0,
  todayCodes: 0,
  typeStats: [],
  sourceStats: []
})

// 搜索参数
const searchParams = reactive<VerificationCodeListParams>({
  page: 1,
  pageSize: 20,
  code: '',
  source: '',
  isUsed: undefined,
  isExpired: undefined
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 通知
const { showSuccess, showError, showWarning } = useNotification()

// 计算属性
const hasSelectedCodes = computed(() => selectedCodes.value.length > 0)

// 防抖搜索
const debouncedSearch = debounce(() => {
  searchParams.page = 1
  pagination.page = 1
  loadCodes()
}, 500)

// 生命周期
onMounted(() => {
  loadCodes()
  loadStats()
})

// 方法
async function loadCodes() {
  try {
    loading.value = true
    const params = {
      ...searchParams,
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    
    const response = await verificationCodeApi.list(params)
    codes.value = response.list
    pagination.total = response.total
  } catch (error) {
    showError('加载验证码列表失败')
    console.error('Load codes error:', error)
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  try {
    stats.value = await verificationCodeApi.getStats()
  } catch (error) {
    console.error('Load stats error:', error)
  }
}

function resetSearch() {
  Object.assign(searchParams, {
    page: 1,
    pageSize: 20,
    code: '',
    source: '',
    isUsed: undefined,
    isExpired: undefined
  })
  pagination.page = 1
  loadCodes()
}

function changePage(page: number) {
  if (page < 1 || page > Math.ceil(pagination.total / pagination.pageSize)) {
    return
  }
  pagination.page = page
  searchParams.page = page
  loadCodes()
}

function toggleSelectAll() {
  if (selectedCodes.value.length === codes.value.length) {
    selectedCodes.value = []
  } else {
    selectedCodes.value = codes.value.map(code => code.id)
  }
}

function toggleSelectCode(id: number) {
  const index = selectedCodes.value.indexOf(id)
  if (index > -1) {
    selectedCodes.value.splice(index, 1)
  } else {
    selectedCodes.value.push(id)
  }
}

async function toggleUsedStatus(code: VerificationCode) {
  try {
    await verificationCodeApi.markUsed(code.id, !code.isUsed)
    code.isUsed = !code.isUsed
    if (code.isUsed) {
      code.usedAt = new Date().toISOString()
    } else {
      code.usedAt = undefined
    }
    showSuccess(`已${code.isUsed ? '标记为已使用' : '标记为未使用'}`)
    loadStats()
  } catch (error) {
    showError('操作失败')
    console.error('Toggle used status error:', error)
  }
}

async function deleteCode(id: number) {
  if (!confirm('确定要删除这个验证码吗？')) {
    return
  }
  
  try {
    await verificationCodeApi.delete(id)
    codes.value = codes.value.filter(code => code.id !== id)
    pagination.total--
    showSuccess('删除成功')
    loadStats()
  } catch (error) {
    showError('删除失败')
    console.error('Delete code error:', error)
  }
}

async function batchMarkUsed(used: boolean) {
  if (selectedCodes.value.length === 0) {
    showWarning('请先选择要操作的验证码')
    return
  }
  
  try {
    for (const id of selectedCodes.value) {
      await verificationCodeApi.markUsed(id, used)
    }
    
    // 更新本地数据
    codes.value.forEach(code => {
      if (selectedCodes.value.includes(code.id)) {
        code.isUsed = used
        if (used) {
          code.usedAt = new Date().toISOString()
        } else {
          code.usedAt = undefined
        }
      }
    })
    
    selectedCodes.value = []
    showSuccess(`批量${used ? '标记已使用' : '标记未使用'}成功`)
    loadStats()
  } catch (error) {
    showError('批量操作失败')
    console.error('Batch mark used error:', error)
  }
}

async function batchDelete() {
  if (selectedCodes.value.length === 0) {
    showWarning('请先选择要删除的验证码')
    return
  }
  
  if (!confirm(`确定要删除选中的 ${selectedCodes.value.length} 个验证码吗？`)) {
    return
  }
  
  try {
    await verificationCodeApi.batchDelete(selectedCodes.value)
    codes.value = codes.value.filter(code => !selectedCodes.value.includes(code.id))
    pagination.total -= selectedCodes.value.length
    selectedCodes.value = []
    showSuccess('批量删除成功')
    loadStats()
  } catch (error) {
    showError('批量删除失败')
    console.error('Batch delete error:', error)
  }
}

function viewCodeDetail(code: VerificationCode) {
  selectedCode.value = code
  showDetailModal.value = true
}

async function copyCode(code: string) {
  try {
    await navigator.clipboard.writeText(code)
    showSuccess('验证码已复制到剪贴板')
  } catch (error) {
    showError('复制失败')
    console.error('Copy error:', error)
  }
}

async function exportCodes() {
  try {
    const params = {
      format: 'csv' as const,
      ...searchParams
    }
    const result = await verificationCodeApi.export(params)
    
    // 创建下载链接
    const link = document.createElement('a')
    link.href = result.downloadUrl
    link.download = result.fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    showSuccess(`导出成功，共 ${result.recordCount} 条记录`)
  } catch (error) {
    showError('导出失败')
    console.error('Export error:', error)
  }
}

function getTypeDisplayName(type?: string): string {
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
  return typeMap[type || ''] || type || '未知'
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

function onCodesExtracted() {
  showExtractModal.value = false
  loadCodes()
  loadStats()
}

function onCodeUpdated() {
  showDetailModal.value = false
  loadCodes()
  loadStats()
}
</script>

<style scoped>
.verification-codes-page {
  @apply p-6 space-y-6;
}

.page-header {
  @apply mb-6;
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

.btn-ghost {
  @apply text-gray-700 bg-transparent hover:bg-gray-100 focus:ring-gray-500 dark:text-gray-300 dark:hover:bg-gray-800;
}

.btn-danger {
  @apply text-white bg-red-600 hover:bg-red-700 focus:ring-red-500;
}

.btn-sm {
  @apply px-3 py-1.5 text-xs;
}

.input {
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:placeholder-gray-400;
}
</style>
