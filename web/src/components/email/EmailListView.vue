<template>
  <div class="h-full flex flex-col">
    <!-- 工具栏 -->
    <div class="flex items-center justify-between p-4 border-b border-glass-border">
      <div class="flex items-center space-x-4">
        <h2 class="text-lg font-semibold text-text-primary">{{ title }}</h2>
        <span class="text-sm text-text-secondary">
          {{ total > 0 ? `共 ${total} 封邮件` : '暂无邮件' }}
        </span>
      </div>
      
      <div class="flex items-center space-x-2">
        <!-- 刷新按钮 -->
        <Button
          variant="ghost"
          size="sm"
          :loading="loading"
          @click="refresh"
        >
          <ArrowPathIcon class="w-4 h-4 mr-2" />
          刷新
        </Button>
        
        <!-- 搜索框 -->
        <div class="relative">
          <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-secondary" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索邮件..."
            class="pl-10 pr-4 py-2 w-64 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
            @keyup.enter="handleSearch"
          />
        </div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="flex items-center justify-between p-4 bg-glass-light border-b border-glass-border">
      <div class="flex items-center space-x-4">
        <!-- 全选 -->
        <input
          v-model="selectAll"
          type="checkbox"
          class="w-4 h-4 text-primary-600 bg-glass-light border-glass-border rounded focus:ring-primary-500"
          @change="handleSelectAll"
        />
        
        <!-- 批量操作 -->
        <div v-if="selectedEmails.length > 0" class="flex items-center space-x-2">
          <span class="text-sm text-text-secondary">
            已选择 {{ selectedEmails.length }} 封邮件
          </span>
          <Button
            variant="ghost"
            size="sm"
            @click="handleBatchMarkRead"
          >
            标记已读
          </Button>
          <Button
            variant="ghost"
            size="sm"
            @click="handleBatchStar"
          >
            添加星标
          </Button>
          <Button
            variant="ghost"
            size="sm"
            @click="handleBatchDelete"
          >
            删除
          </Button>
        </div>
      </div>
      
      <div class="flex items-center space-x-2">
        <!-- 筛选选项 -->
        <select
          v-model="filterRead"
          class="px-3 py-1 bg-glass-medium border border-glass-border rounded text-sm text-text-primary"
        >
          <option value="">全部</option>
          <option value="false">未读</option>
          <option value="true">已读</option>
        </select>
        
        <select
          v-model="filterStarred"
          class="px-3 py-1 bg-glass-medium border border-glass-border rounded text-sm text-text-primary"
        >
          <option value="">全部</option>
          <option value="true">已标星</option>
          <option value="false">未标星</option>
        </select>
      </div>
    </div>

    <!-- 邮件列表 -->
    <div class="flex-1 overflow-auto">
      <div v-if="loading && emails.length === 0" class="flex items-center justify-center h-64">
        <div class="text-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto mb-4"></div>
          <p class="text-text-secondary">加载中...</p>
        </div>
      </div>
      
      <div v-else-if="emails.length === 0" class="flex items-center justify-center h-64">
        <div class="text-center">
          <component :is="emptyIcon" class="w-16 h-16 text-text-secondary mx-auto mb-4" />
          <h3 class="text-lg font-medium text-text-primary mb-2">{{ emptyTitle }}</h3>
          <p class="text-text-secondary">{{ emptyMessage }}</p>
        </div>
      </div>
      
      <div v-else class="p-4 space-y-2">
        <EmailItem
          v-for="email in emails"
          :key="email.id"
          :email="email"
          :selected="selectedEmails.includes(email.id)"
          @click="handleEmailClick"
          @select="handleEmailSelect"
          @star="handleEmailStar"
          @delete="handleEmailDelete"
        />
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="totalPages > 1" class="p-4 border-t border-glass-border">
      <Pagination
        :current-page="currentPage"
        :total-pages="totalPages"
        @page-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Email, EmailListParams } from '@/types'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'

// Icons
import {
  ArrowPathIcon,
  MagnifyingGlassIcon,
  InboxIcon,
  PaperAirplaneIcon
} from '@heroicons/vue/24/outline'

// Components
import Button from '@/components/ui/Button.vue'
import EmailItem from './EmailItem.vue'
import Pagination from '@/components/ui/Pagination.vue'

// Props
interface Props {
  title: string
  emailType: 'inbox' | 'sent'
  emptyTitle?: string
  emptyMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  emptyTitle: '暂无邮件',
  emptyMessage: '这里还没有邮件'
})

// Composables
const router = useRouter()
const { showNotification } = useNotification()

// State
const loading = ref(false)
const emails = ref<Email[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')
const filterRead = ref('')
const filterStarred = ref('')
const selectedEmails = ref<string[]>([])
const selectAll = ref(false)

// Computed
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const emptyIcon = computed(() => {
  return props.emailType === 'sent' ? PaperAirplaneIcon : InboxIcon
})

const queryParams = computed((): EmailListParams => {
  const params: EmailListParams = {
    page: currentPage.value,
    pageSize: pageSize.value
  }
  
  if (searchQuery.value) {
    params.subject = searchQuery.value
  }
  
  if (filterRead.value !== '') {
    params.isRead = filterRead.value === 'true'
  }
  
  if (filterStarred.value !== '') {
    params.isStarred = filterStarred.value === 'true'
  }
  
  return params
})

// Methods
const loadEmails = async () => {
  loading.value = true
  try {
    const apiMethod = props.emailType === 'sent' ? emailApi.getSentEmails : emailApi.getInboxEmails
    const response = await apiMethod(queryParams.value)
    
    if (response.success && response.data) {
      emails.value = response.data.list
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Failed to load emails:', error)
    showNotification({
      type: 'error',
      title: '加载失败',
      message: '无法加载邮件列表'
    })
  } finally {
    loading.value = false
  }
}

const refresh = () => {
  currentPage.value = 1
  loadEmails()
}

const handleSearch = () => {
  currentPage.value = 1
  loadEmails()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadEmails()
}

const handleEmailClick = (email: Email) => {
  router.push(`/email/${email.id}`)
}

const handleEmailSelect = (email: Email, selected: boolean) => {
  if (selected) {
    selectedEmails.value.push(email.id)
  } else {
    const index = selectedEmails.value.indexOf(email.id)
    if (index > -1) {
      selectedEmails.value.splice(index, 1)
    }
  }
  
  // 更新全选状态
  selectAll.value = selectedEmails.value.length === emails.value.length
}

const handleSelectAll = () => {
  if (selectAll.value) {
    selectedEmails.value = emails.value.map(email => email.id)
  } else {
    selectedEmails.value = []
  }
}

const handleEmailStar = async (email: Email) => {
  try {
    await emailApi.markAsStarred(email.id, !email.isStarred)
    email.isStarred = !email.isStarred
    showNotification({
      type: 'success',
      title: email.isStarred ? '已添加星标' : '已取消星标'
    })
  } catch (error) {
    console.error('Failed to toggle star:', error)
    showNotification({
      type: 'error',
      title: '操作失败',
      message: '无法更新星标状态'
    })
  }
}

const handleEmailDelete = async (email: Email) => {
  try {
    await emailApi.deleteEmail(email.id)
    emails.value = emails.value.filter(e => e.id !== email.id)
    total.value--
    showNotification({
      type: 'success',
      title: '邮件已删除'
    })
  } catch (error) {
    console.error('Failed to delete email:', error)
    showNotification({
      type: 'error',
      title: '删除失败',
      message: '无法删除邮件'
    })
  }
}

const handleBatchMarkRead = async () => {
  try {
    await emailApi.batchUpdate(selectedEmails.value, 'markRead', { isRead: true })
    emails.value.forEach(email => {
      if (selectedEmails.value.includes(email.id)) {
        email.isRead = true
      }
    })
    selectedEmails.value = []
    selectAll.value = false
    showNotification({
      type: 'success',
      title: '已标记为已读'
    })
  } catch (error) {
    console.error('Failed to batch mark read:', error)
    showNotification({
      type: 'error',
      title: '操作失败',
      message: '无法批量标记已读'
    })
  }
}

const handleBatchStar = async () => {
  try {
    await emailApi.batchUpdate(selectedEmails.value, 'markStar', { isStarred: true })
    emails.value.forEach(email => {
      if (selectedEmails.value.includes(email.id)) {
        email.isStarred = true
      }
    })
    selectedEmails.value = []
    selectAll.value = false
    showNotification({
      type: 'success',
      title: '已添加星标'
    })
  } catch (error) {
    console.error('Failed to batch star:', error)
    showNotification({
      type: 'error',
      title: '操作失败',
      message: '无法批量添加星标'
    })
  }
}

const handleBatchDelete = async () => {
  try {
    await emailApi.batchUpdate(selectedEmails.value, 'delete')
    emails.value = emails.value.filter(email => !selectedEmails.value.includes(email.id))
    total.value -= selectedEmails.value.length
    selectedEmails.value = []
    selectAll.value = false
    showNotification({
      type: 'success',
      title: '邮件已删除'
    })
  } catch (error) {
    console.error('Failed to batch delete:', error)
    showNotification({
      type: 'error',
      title: '删除失败',
      message: '无法批量删除邮件'
    })
  }
}

// Watchers
watch([filterRead, filterStarred], () => {
  currentPage.value = 1
  loadEmails()
})

// Lifecycle
onMounted(() => {
  loadEmails()
})
</script>
