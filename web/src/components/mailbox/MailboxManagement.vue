<template>
  <div class="space-y-6">
    <!-- 页面标题和操作 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-text-primary">邮箱管理</h2>
        <p class="text-text-secondary mt-1">管理您的邮箱账户配置</p>
      </div>
      <Button
        variant="primary"
        @click="showAddModal = true"
      >
        <PlusIcon class="w-4 h-4 mr-2" />
        添加邮箱
      </Button>
    </div>

    <!-- 统计信息 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-glass-light border border-glass-border rounded-lg p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <EnvelopeIcon class="w-8 h-8 text-primary-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">总邮箱数</p>
            <p class="text-2xl font-bold text-text-primary">{{ stats.totalMailboxes }}</p>
          </div>
        </div>
      </div>
      <div class="bg-glass-light border border-glass-border rounded-lg p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <CheckCircleIcon class="w-8 h-8 text-green-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">活跃邮箱</p>
            <p class="text-2xl font-bold text-text-primary">{{ stats.activeMailboxes }}</p>
          </div>
        </div>
      </div>
      <div class="bg-glass-light border border-glass-border rounded-lg p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <BuildingOfficeIcon class="w-8 h-8 text-blue-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">自建邮箱</p>
            <p class="text-2xl font-bold text-text-primary">{{ stats.selfMailboxes }}</p>
          </div>
        </div>
      </div>
      <div class="bg-glass-light border border-glass-border rounded-lg p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <ServerIcon class="w-8 h-8 text-purple-400" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-text-secondary">第三方邮箱</p>
            <p class="text-2xl font-bold text-text-primary">{{ stats.thirdMailboxes }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 筛选和搜索 -->
    <div class="bg-glass-light border border-glass-border rounded-lg p-4">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0">
        <div class="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
          <!-- 搜索框 -->
          <div class="relative">
            <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-secondary" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索邮箱地址..."
              class="pl-10 pr-4 py-2 bg-glass-medium border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
            />
          </div>

          <!-- 类型筛选 -->
          <select
            v-model="typeFilter"
            class="px-3 py-2 bg-glass-medium border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
          >
            <option value="">所有类型</option>
            <option value="self">自建邮箱</option>
            <option value="third">第三方邮箱</option>
          </select>

          <!-- 状态筛选 -->
          <select
            v-model="statusFilter"
            class="px-3 py-2 bg-glass-medium border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
          >
            <option value="">所有状态</option>
            <option value="1">启用</option>
            <option value="2">禁用</option>
          </select>
        </div>

        <div class="flex items-center space-x-2">
          <Button
            variant="ghost"
            size="sm"
            @click="loadMailboxes"
            :loading="loading"
          >
            <ArrowPathIcon class="w-4 h-4 mr-2" />
            刷新
          </Button>
        </div>
      </div>
    </div>

    <!-- 邮箱列表 -->
    <div class="space-y-4">
      <div v-if="loading && (mailboxes?.length || 0) === 0" class="text-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto"></div>
        <p class="mt-2 text-text-secondary">加载中...</p>
      </div>

      <div v-else-if="filteredMailboxes.length === 0" class="text-center py-12">
        <EnvelopeIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
        <h3 class="text-lg font-medium text-text-primary mb-2">
          {{ searchQuery || typeFilter || statusFilter ? '没有找到匹配的邮箱' : '暂无邮箱' }}
        </h3>
        <p class="text-text-secondary mb-4">
          {{ searchQuery || typeFilter || statusFilter ? '请尝试调整筛选条件' : '开始添加您的第一个邮箱账户' }}
        </p>
        <Button v-if="!searchQuery && !typeFilter && !statusFilter" variant="primary" @click="showAddModal = true">
          添加邮箱
        </Button>
      </div>

      <MailboxCard
        v-else
        v-for="mailbox in filteredMailboxes"
        :key="mailbox.id"
        :mailbox="mailbox"
        @edit="editMailbox"
        @delete="deleteMailbox"
        @sync="syncMailbox"
        @test="testConnection"
        @toggle-status="toggleMailboxStatus"
      />
    </div>

    <!-- 分页 -->
    <div v-if="totalPages > 1" class="flex justify-center">
      <Pagination
        :current-page="currentPage"
        :total-pages="totalPages"
        @page-change="handlePageChange"
      />
    </div>

    <!-- 添加/编辑邮箱模态框 -->
    <MailboxModal
      v-if="showAddModal || showEditModal"
      :visible="showAddModal || showEditModal"
      :mailbox="editingMailbox"
      :providers="providers"
      @close="closeModal"
      @save="saveMailbox"
    />

    <!-- 删除确认模态框 -->
    <ConfirmModal
      v-if="showDeleteModal"
      :visible="showDeleteModal"
      title="删除邮箱"
      :message="`确定要删除邮箱 ${deletingMailbox?.email} 吗？此操作不可恢复。`"
      confirm-text="删除"
      confirm-variant="danger"
      @close="showDeleteModal = false"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { 
  Mailbox, 
  MailboxCreateRequest, 
  MailboxUpdateRequest, 
  MailboxProvider,
  MailboxStats,
  MailboxListRequest
} from '@/types'
import { mailboxApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'

// Icons
import {
  PlusIcon,
  EnvelopeIcon,
  CheckCircleIcon,
  BuildingOfficeIcon,
  ServerIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline'

// Components
import Button from '@/components/ui/Button.vue'
import MailboxCard from './MailboxCard.vue'
import MailboxModal from './MailboxModal.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import Pagination from '@/components/ui/Pagination.vue'

// Composables
const { showNotification } = useNotification()

// State
const loading = ref(false)
const mailboxes = ref<Mailbox[]>([])
const stats = ref<MailboxStats>({
  totalMailboxes: 0,
  activeMailboxes: 0,
  selfMailboxes: 0,
  thirdMailboxes: 0
})
const providers = ref<MailboxProvider[]>([])

// 模态框状态
const showAddModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingMailbox = ref<Mailbox | null>(null)
const deletingMailbox = ref<Mailbox | null>(null)

// 筛选和搜索
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')

// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 计算属性
const totalPages = computed(() => Math.ceil((total.value || 0) / (pageSize.value || 10)))

const filteredMailboxes = computed(() => {
  let filtered = mailboxes.value || []

  if (searchQuery.value) {
    filtered = filtered.filter(mailbox =>
      mailbox.email?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  if (typeFilter.value) {
    filtered = filtered.filter(mailbox => mailbox.type === typeFilter.value)
  }

  if (statusFilter.value) {
    filtered = filtered.filter(mailbox => mailbox.status?.toString() === statusFilter.value)
  }

  return filtered
})

// 方法
const loadMailboxes = async () => {
  loading.value = true
  try {
    const params: MailboxListRequest = {
      page: currentPage.value,
      pageSize: pageSize.value
    }

    const response = await mailboxApi.list(params)
    if (response.success && response.data) {
      mailboxes.value = response.data.list
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Failed to load mailboxes:', error)
    showNotification({
      type: 'error',
      title: '加载失败',
      message: '无法加载邮箱列表'
    })
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await mailboxApi.getStats()
    if (response.success && response.data) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

const loadProviders = async () => {
  try {
    const response = await mailboxApi.getProviders()
    if (response.success && response.data) {
      providers.value = response.data
    }
  } catch (error) {
    console.error('Failed to load providers:', error)
  }
}

const editMailbox = (mailbox: Mailbox) => {
  editingMailbox.value = mailbox
  showEditModal.value = true
}

const deleteMailbox = (mailbox: Mailbox) => {
  deletingMailbox.value = mailbox
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!deletingMailbox.value) return

  try {
    const response = await mailboxApi.delete(deletingMailbox.value.id)
    if (response.success) {
      showNotification({
        type: 'success',
        title: '删除成功',
        message: '邮箱已删除'
      })
      loadMailboxes()
      loadStats()
    }
  } catch (error) {
    console.error('Failed to delete mailbox:', error)
    showNotification({
      type: 'error',
      title: '删除失败',
      message: '无法删除邮箱'
    })
  } finally {
    showDeleteModal.value = false
    deletingMailbox.value = null
  }
}

const syncMailbox = async (mailbox: Mailbox) => {
  try {
    const response = await mailboxApi.sync({ id: mailbox.id })
    if (response.success && response.data) {
      showNotification({
        type: 'success',
        title: '同步成功',
        message: `同步了 ${response.data.syncCount} 封邮件`
      })
      loadMailboxes()
    }
  } catch (error) {
    console.error('Failed to sync mailbox:', error)
    showNotification({
      type: 'error',
      title: '同步失败',
      message: '无法同步邮箱'
    })
  }
}

const testConnection = async (mailbox: Mailbox) => {
  // 这里需要用户重新输入密码，暂时跳过
  showNotification({
    type: 'info',
    title: '测试连接',
    message: '请编辑邮箱配置后进行连接测试'
  })
}

const toggleMailboxStatus = async (mailbox: Mailbox) => {
  try {
    const newStatus = mailbox.status === 1 ? 2 : 1
    const response = await mailboxApi.update({
      id: mailbox.id,
      status: newStatus
    })
    
    if (response.success) {
      showNotification({
        type: 'success',
        title: '状态更新成功',
        message: `邮箱已${newStatus === 1 ? '启用' : '禁用'}`
      })
      loadMailboxes()
      loadStats()
    }
  } catch (error) {
    console.error('Failed to toggle mailbox status:', error)
    showNotification({
      type: 'error',
      title: '状态更新失败',
      message: '无法更新邮箱状态'
    })
  }
}

const saveMailbox = async (data: MailboxCreateRequest | MailboxUpdateRequest) => {
  try {
    let response
    if ('id' in data) {
      // 更新邮箱
      response = await mailboxApi.update(data as MailboxUpdateRequest)
    } else {
      // 创建邮箱
      response = await mailboxApi.create(data as MailboxCreateRequest)
    }

    if (response.success) {
      showNotification({
        type: 'success',
        title: 'id' in data ? '更新成功' : '创建成功',
        message: '邮箱已保存'
      })
      closeModal()
      loadMailboxes()
      loadStats()
    }
  } catch (error) {
    console.error('Failed to save mailbox:', error)
    showNotification({
      type: 'error',
      title: '保存失败',
      message: '无法保存邮箱'
    })
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingMailbox.value = null
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadMailboxes()
}

// 生命周期
onMounted(() => {
  loadMailboxes()
  loadStats()
  loadProviders()
})
</script>
