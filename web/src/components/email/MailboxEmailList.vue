<template>
  <div class="p-4">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="emails.length === 0" class="text-center py-8">
      <EnvelopeIcon class="w-12 h-12 text-text-secondary mx-auto mb-3" />
      <p class="text-text-secondary">此邮箱暂无邮件</p>
    </div>

    <!-- 邮件列表 -->
    <div v-else class="space-y-2">
      <!-- 邮件操作栏 -->
      <div class="flex items-center justify-between mb-4">
        <div class="text-sm text-text-secondary">
          共 {{ totalCount }} 封邮件
        </div>
        <div class="flex items-center space-x-2">
          <!-- 批量操作 -->
          <div v-if="selectedEmails.length > 0" class="flex items-center space-x-2">
            <Button
              variant="secondary"
              size="sm"
              @click="batchMarkAsRead"
              :loading="batchLoading"
            >
              <CheckIcon class="w-4 h-4 mr-1" />
              标记已读 ({{ selectedEmails.length }})
            </Button>
            <Button
              variant="secondary"
              size="sm"
              @click="batchDelete"
              :loading="batchLoading"
            >
              <TrashIcon class="w-4 h-4 mr-1" />
              删除 ({{ selectedEmails.length }})
            </Button>
          </div>
          
          <!-- 全选 -->
          <input
            type="checkbox"
            :checked="isAllSelected"
            @change="toggleSelectAll"
            class="rounded border-gray-600 bg-gray-700 text-primary-500 focus:ring-primary-500"
          />
        </div>
      </div>

      <!-- 邮件项目 -->
      <div
        v-for="email in emails"
        :key="email?.id || Math.random()"
        :class="[
          'flex items-center p-3 rounded-lg border border-glass-border cursor-pointer transition-colors',
          selectedEmails.includes(email?.id) ? 'bg-primary-500/10 border-primary-500/30' : 'hover:bg-white/5',
          !email?.isRead ? 'font-semibold' : ''
        ]"
        @click="selectEmail(email)"
      >
        <!-- 选择框 -->
        <input
          type="checkbox"
          :checked="email?.id && selectedEmails.includes(email.id)"
          @change="email?.id && toggleEmailSelection(email.id)"
          @click.stop
          class="rounded border-gray-600 bg-gray-700 text-primary-500 focus:ring-primary-500 mr-3"
        />

        <!-- 邮件信息 -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between mb-1">
            <!-- 发件人 -->
            <div class="flex items-center space-x-2">
              <span class="text-text-primary truncate">
                {{ email?.from?.name || email?.from?.email || '未知发件人' }}
              </span>
              <!-- 星标 -->
              <button
                v-if="email?.isStarred"
                @click.stop="toggleStar(email.id, false)"
                class="text-yellow-400 hover:text-yellow-300"
              >
                <StarIcon class="w-4 h-4 fill-current" />
              </button>
              <button
                v-else
                @click.stop="toggleStar(email.id, true)"
                class="text-text-secondary hover:text-yellow-400"
              >
                <StarIcon class="w-4 h-4" />
              </button>
            </div>
            
            <!-- 时间 -->
            <div class="text-xs text-text-secondary">
              {{ formatDate(email?.createdAt) }}
            </div>
          </div>
          
          <!-- 主题 -->
          <div class="text-text-primary truncate mb-1">
            {{ email?.subject || '(无主题)' }}
          </div>
          
          <!-- 预览 -->
          <div class="text-sm text-text-secondary truncate">
            {{ getEmailPreview(email?.body) }}
          </div>
          
          <!-- 附件标识 -->
          <div v-if="email?.attachments && email.attachments.length > 0" class="flex items-center mt-1">
            <PaperClipIcon class="w-4 h-4 text-text-secondary mr-1" />
            <span class="text-xs text-text-secondary">{{ email.attachments.length }} 个附件</span>
          </div>
        </div>

        <!-- 未读标识 -->
        <div v-if="!email?.isRead" class="w-2 h-2 bg-primary-500 rounded-full ml-3"></div>
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore" class="text-center pt-4">
        <Button
          variant="ghost"
          @click="loadMore"
          :loading="loadingMore"
        >
          加载更多
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'
import type { Email, EmailListParams } from '@/types'

import Button from '@/components/ui/Button.vue'

import {
  EnvelopeIcon,
  CheckIcon,
  TrashIcon,
  StarIcon,
  PaperClipIcon
} from '@heroicons/vue/24/outline'

// Props
interface Props {
  mailboxId: number
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  emailSelected: [email: Email]
}>()

// 通知
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const loading = ref(false)
const loadingMore = ref(false)
const batchLoading = ref(false)

const emails = ref<Email[]>([])
const selectedEmails = ref<string[]>([])
const totalCount = ref(0)
const currentPage = ref(1)
const hasMore = ref(true)

// 查询参数
const queryParams = reactive<EmailListParams>({
  page: 1,
  pageSize: 10
})

// 计算属性
const isAllSelected = computed(() => {
  return emails.value.length > 0 && selectedEmails.value.length === emails.value.length
})

// 监听邮箱ID变化
watch(() => props.mailboxId, () => {
  if (props.mailboxId) {
    resetAndLoad()
  }
}, { immediate: true })

// 生命周期
onMounted(() => {
  if (props.mailboxId) {
    loadEmails()
  }
})

// 方法
const resetAndLoad = () => {
  emails.value = []
  selectedEmails.value = []
  currentPage.value = 1
  hasMore.value = true
  loadEmails()
}

const loadEmails = async (append = false) => {
  if (!props.mailboxId) return
  
  try {
    if (!append) {
      loading.value = true
      currentPage.value = 1
    } else {
      loadingMore.value = true
    }

    const params = {
      page: append ? currentPage.value + 1 : 1,
      pageSize: queryParams.pageSize,
      mailboxId: props.mailboxId,
      direction: 'received' // 收件箱只显示接收的邮件
    }

    console.log('MailboxEmailList API params:', params)
    const response = await emailApi.getInboxEmails(params)
    console.log('MailboxEmailList API response:', response)

    if (response.success && response.data) {
      const emailList = response.data.list || response.data || []

      if (append) {
        emails.value.push(...emailList)
        currentPage.value++
      } else {
        emails.value = emailList
      }

      totalCount.value = response.data.total || emailList.length
      hasMore.value = emailList.length === queryParams.pageSize
    } else {
      if (!append) {
        emails.value = []
      }
      totalCount.value = 0
      hasMore.value = false
    }
  } catch (error) {
    console.error('Failed to load emails:', error)
    showError('加载邮件失败')

    // 提供模拟数据用于测试
    if (!append) {
      emails.value = [
        {
          id: `email-${props.mailboxId}-1`,
          mailboxId: props.mailboxId,
          subject: `测试邮件 1 - 邮箱 ${props.mailboxId}`,
          fromEmail: 'sender1@example.com',
          toEmails: 'recipient@example.com',
          body: '这是一封测试邮件的内容...',
          isRead: false,
          isStarred: true,
          attachments: [],
          from: {
            email: 'sender1@example.com',
            name: '发件人1'
          },
          to: [{
            email: 'recipient@example.com',
            name: '收件人'
          }],
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        },
        {
          id: `email-${props.mailboxId}-2`,
          mailboxId: props.mailboxId,
          subject: `测试邮件 2 - 邮箱 ${props.mailboxId}`,
          fromEmail: 'sender2@example.com',
          toEmails: 'recipient@example.com',
          body: '这是另一封测试邮件的内容，包含更多文字用于测试预览功能...',
          isRead: true,
          isStarred: false,
          attachments: [{ name: 'document.pdf', size: 1024 }],
          from: {
            email: 'sender2@example.com',
            name: '发件人2'
          },
          to: [{
            email: 'recipient@example.com',
            name: '收件人'
          }],
          createdAt: new Date(Date.now() - 86400000).toISOString(), // 1天前
          updatedAt: new Date(Date.now() - 86400000).toISOString()
        }
      ]
      totalCount.value = emails.value.length
      hasMore.value = false
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = () => {
  loadEmails(true)
}

const selectEmail = (email: Email) => {
  emit('emailSelected', email)
  
  // 标记为已读
  if (!email.isRead) {
    markAsRead(email.id, true)
  }
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedEmails.value = []
  } else {
    selectedEmails.value = emails.value.map(email => email?.id).filter(Boolean)
  }
}

const toggleEmailSelection = (emailId: string) => {
  const index = selectedEmails.value.indexOf(emailId)
  if (index > -1) {
    selectedEmails.value.splice(index, 1)
  } else {
    selectedEmails.value.push(emailId)
  }
}

const toggleStar = async (emailId: string, starred: boolean) => {
  try {
    await emailApi.markAsStarred(emailId, starred)
    
    // 更新本地状态
    const email = emails.value.find(e => e.id === emailId)
    if (email) {
      email.isStarred = starred
    }
    
    showSuccess(starred ? '已添加星标' : '已取消星标')
  } catch (error) {
    showError('操作失败')
  }
}

const markAsRead = async (emailId: string, isRead: boolean) => {
  try {
    await emailApi.markAsRead(emailId, isRead)
    
    // 更新本地状态
    const email = emails.value.find(e => e.id === emailId)
    if (email) {
      email.isRead = isRead
    }
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

const batchMarkAsRead = async () => {
  try {
    batchLoading.value = true
    
    await emailApi.batchUpdate(selectedEmails.value, 'read')
    
    // 更新本地状态
    emails.value.forEach(email => {
      if (selectedEmails.value.includes(email.id)) {
        email.isRead = true
      }
    })
    
    selectedEmails.value = []
    showSuccess('已标记为已读')
  } catch (error) {
    showError('批量操作失败')
  } finally {
    batchLoading.value = false
  }
}

const batchDelete = async () => {
  if (!confirm(`确定要删除选中的 ${selectedEmails.value.length} 封邮件吗？`)) return
  
  try {
    batchLoading.value = true
    
    await emailApi.batchUpdate(selectedEmails.value, 'delete')
    
    // 从列表中移除
    emails.value = emails.value.filter(email => !selectedEmails.value.includes(email.id))
    totalCount.value -= selectedEmails.value.length
    selectedEmails.value = []
    
    showSuccess('邮件已删除')
  } catch (error) {
    showError('删除失败')
  } finally {
    batchLoading.value = false
  }
}

const getEmailPreview = (body: string) => {
  if (!body) return ''
  const text = body.replace(/<[^>]*>/g, '').trim()
  return text.length > 100 ? text.substring(0, 100) + '...' : text
}

const formatDate = (dateString: string) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else if (diffDays === 1) {
    return '昨天'
  } else if (diffDays < 7) {
    return `${diffDays}天前`
  } else {
    return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
  }
}
</script>
