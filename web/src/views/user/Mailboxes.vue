<template>
  <div class="min-h-screen bg-background-primary">
    <!-- 页面头部 -->
    <div class="bg-glass-light backdrop-blur-md border-b border-glass-border">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-2xl font-bold text-text-primary">邮箱管理</h1>
            <div class="ml-4 text-sm text-text-secondary">
              管理您的多个邮箱账户
            </div>
          </div>
          <div class="flex items-center space-x-4">
            <!-- 统计信息 -->
            <div class="hidden md:flex items-center space-x-6 text-sm">
              <div class="text-center">
                <div class="font-semibold text-text-primary">{{ stats.totalMailboxes }}</div>
                <div class="text-text-secondary">总邮箱</div>
              </div>
              <div class="text-center">
                <div class="font-semibold text-primary-400">{{ stats.activeMailboxes }}</div>
                <div class="text-text-secondary">活跃</div>
              </div>
            </div>
            <!-- 创建邮箱按钮 -->
            <Button
              variant="primary"
              @click="showAddModal = true"
              class="flex items-center"
            >
              <PlusIcon class="w-4 h-4 mr-2" />
              创建邮箱
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 筛选和搜索 -->
      <div class="mb-6">
        <GlassCard :level="2" padding="md" class="p-4">
          <div class="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0">
            <div class="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
              <!-- 搜索框 -->
              <div class="relative">
                <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-secondary" />
                <input
                  v-model="searchQuery"
                  type="text"
                  placeholder="搜索邮箱地址..."
                  class="pl-10 pr-4 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
                />
              </div>
              
              <!-- 类型筛选 -->
              <select
                v-model="filterType"
                class="px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
              >
                <option value="">所有类型</option>
                <option value="self">自建邮箱</option>
                <option value="third">第三方邮箱</option>
              </select>

              <!-- 状态筛选 -->
              <select
                v-model="filterStatus"
                class="px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
              >
                <option value="">所有状态</option>
                <option value="1">启用</option>
                <option value="0">禁用</option>
              </select>
            </div>

            <!-- 刷新按钮 -->
            <Button
              variant="ghost"
              @click="loadMailboxes"
              :loading="loading"
              class="flex items-center"
            >
              <ArrowPathIcon class="w-4 h-4 mr-2" />
              刷新
            </Button>
          </div>
        </GlassCard>
      </div>

      <!-- 邮箱列表 -->
      <div class="space-y-4">
        <div v-if="loading && mailboxes.length === 0" class="text-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto"></div>
          <p class="mt-2 text-text-secondary">加载中...</p>
        </div>

        <div v-else-if="mailboxes.length === 0" class="text-center py-12">
          <EnvelopeIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
          <h3 class="text-lg font-medium text-text-primary mb-2">暂无邮箱</h3>
          <p class="text-text-secondary mb-4">开始添加您的第一个邮箱账户</p>
          <Button variant="primary" @click="showAddModal = true">
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
        />
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="mt-8 flex justify-center">
        <Pagination
          :current-page="currentPage"
          :total-pages="totalPages"
          @page-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 添加/编辑邮箱模态框 -->
    <MailboxModal
      v-if="showAddModal || showEditModal"
      :visible="showAddModal || showEditModal"
      :mailbox="editingMailbox"
      @close="closeModal"
      @save="saveMailbox"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useNotification } from '@/composables/useNotification'
import { mailboxApi } from '@/utils/api'
import type { Mailbox, MailboxStats, MailboxCreateRequest, MailboxUpdateRequest } from '@/types'

// Icons
import {
  PlusIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon,
  EnvelopeIcon
} from '@heroicons/vue/24/outline'

// Components
import Button from '@/components/ui/Button.vue'
import GlassCard from '@/components/ui/GlassCard.vue'
import MailboxCard from '@/components/mailbox/MailboxCard.vue'
import MailboxModal from '@/components/mailbox/MailboxModal.vue'
import Pagination from '@/components/ui/Pagination.vue'

// 响应式数据
const loading = ref(false)
const mailboxes = ref<Mailbox[]>([])
const stats = ref<MailboxStats>({
  totalMailboxes: 0,
  activeMailboxes: 0,
  selfMailboxes: 0,
  thirdMailboxes: 0
})


// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)

// 筛选和搜索
const searchQuery = ref('')
const filterType = ref('')
const filterStatus = ref('')

// 模态框
const showAddModal = ref(false)
const showEditModal = ref(false)
const editingMailbox = ref<Mailbox | null>(null)

// 通知
const { showNotification } = useNotification()

// 计算属性
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value))

const filteredMailboxes = computed(() => {
  let filtered = mailboxes.value

  if (searchQuery.value) {
    filtered = filtered.filter(mailbox =>
      mailbox.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  if (filterType.value) {
    filtered = filtered.filter(mailbox => mailbox.type === filterType.value)
  }

  if (filterStatus.value !== '') {
    filtered = filtered.filter(mailbox => mailbox.status === parseInt(filterStatus.value))
  }

  return filtered
})

// 方法
const loadMailboxes = async () => {
  try {
    loading.value = true
    const response = await mailboxApi.list({
      page: currentPage.value,
      pageSize: pageSize.value
    })

    if (response.success && response.data) {
      mailboxes.value = response.data.items || []
      totalItems.value = response.data.total || 0
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



const editMailbox = (mailbox: Mailbox) => {
  editingMailbox.value = mailbox
  showEditModal.value = true
}

const deleteMailbox = async (mailbox: Mailbox) => {
  if (!confirm(`确定要删除邮箱 ${mailbox.email} 吗？`)) {
    return
  }

  try {
    const response = await mailboxApi.delete(mailbox.id)
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
  }
}

const syncMailbox = async (mailbox: Mailbox) => {
  try {
    const response = await mailboxApi.sync({
      id: mailbox.id,
      syncDays: 7
    })
    
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
  try {
    const response = await mailboxApi.testConnection({
      email: mailbox.email,
      password: '', // 需要用户重新输入密码
      imapHost: mailbox.imapHost,
      imapPort: mailbox.imapPort,
      imapSsl: mailbox.imapSsl,
      smtpHost: mailbox.smtpHost,
      smtpPort: mailbox.smtpPort,
      smtpSsl: mailbox.smtpSsl
    })
    
    if (response.success && response.data) {
      const result = response.data
      if (result.imapSuccess && result.smtpSuccess) {
        showNotification({
          type: 'success',
          title: '连接测试成功',
          message: result.message
        })
      } else {
        showNotification({
          type: 'warning',
          title: '连接测试部分失败',
          message: `IMAP: ${result.imapSuccess ? '成功' : result.imapError}, SMTP: ${result.smtpSuccess ? '成功' : result.smtpError}`
        })
      }
    }
  } catch (error) {
    console.error('Failed to test connection:', error)
    showNotification({
      type: 'error',
      title: '连接测试失败',
      message: '无法测试邮箱连接'
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

// 监听搜索和筛选变化
watch([searchQuery, filterType, filterStatus], () => {
  currentPage.value = 1
})

// 组件挂载时加载数据
onMounted(() => {
  loadMailboxes()
  loadStats()
})
</script>
