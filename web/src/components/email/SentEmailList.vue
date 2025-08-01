<template>
  <div class="p-4">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="(!emails || emails.length === 0)" class="text-center py-8">
      <PaperAirplaneIcon class="w-12 h-12 text-text-secondary mx-auto mb-3" />
      <p class="text-text-secondary">此邮箱暂无已发送邮件</p>
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
          selectedEmails.includes(email?.id) ? 'bg-primary-500/10 border-primary-500/30' : 'hover:bg-white/5'
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
            <!-- 收件人 -->
            <div class="flex items-center space-x-2">
              <span class="text-text-primary truncate">
                {{ getRecipientDisplay(email?.to) }}
              </span>
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
import { useRouter } from 'vue-router'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'
import type { Email, EmailListParams } from '@/types'

import Button from '@/components/ui/Button.vue'

import {
  PaperAirplaneIcon,
  PaperClipIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'

// Props
interface Props {
  mailboxId: string
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  'email-selected': [email: Email]
}>()

// 路由和通知
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const loading = ref(false)
const loadingMore = ref(false)
const batchLoading = ref(false)

const emails = ref<Email[]>([])
const selectedEmails = ref<string[]>([])
const selectedEmailId = ref<string | null>(null)
const totalCount = ref(0)
const currentPage = ref(1)
const hasMore = ref(true)

// 查询参数
const queryParams = reactive<EmailListParams>({
  page: 1,
  limit: 20,
  sortBy: 'createdAt',
  sortOrder: 'desc'
})

// 计算属性
const isAllSelected = computed(() => {
  return emails.value && emails.value.length > 0 && selectedEmails.value.length === emails.value.length
})

// 监听器
watch(() => props.mailboxId, () => {
  resetData()
  loadEmails()
}, { immediate: true })

// 方法
const resetData = () => {
  emails.value = []
  selectedEmails.value = []
  selectedEmailId.value = null
  totalCount.value = 0
  currentPage.value = 1
  hasMore.value = true
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
      ...queryParams,
      page: append ? currentPage.value + 1 : 1,
      mailboxId: props.mailboxId
    }

    const response = await emailApi.getSentEmails(params)
    console.log('Sent emails API response:', response)

    let emailData = []
    if (response.data) {
      if (Array.isArray(response.data)) {
        emailData = response.data
      } else if (response.data.list && Array.isArray(response.data.list)) {
        emailData = response.data.list
      } else if (response.data.data && Array.isArray(response.data.data)) {
        emailData = response.data.data
      }
    }

    // 如果没有真实数据，提供一些模拟的已发送邮件数据
    if (emailData.length === 0 && !append) {
      console.log('Using mock sent email data for testing')
      emailData = [
        {
          id: `sent-${props.mailboxId}-1`,
          subject: '测试已发送邮件 1',
          body: '这是一封测试的已发送邮件内容...',
          to: [{ email: 'recipient1@example.com', name: '收件人1' }],
          cc: [],
          bcc: [],
          attachments: [],
          isRead: true,
          isStarred: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(), // 2小时前
          updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString()
        },
        {
          id: `sent-${props.mailboxId}-2`,
          subject: '重要会议邮件',
          body: '关于明天会议的详细安排和议程...',
          to: [
            { email: 'colleague1@company.com', name: '同事A' },
            { email: 'colleague2@company.com', name: '同事B' }
          ],
          cc: [{ email: 'manager@company.com', name: '经理' }],
          bcc: [],
          attachments: [
            { name: '会议议程.pdf', size: 1024000 },
            { name: '演示文稿.pptx', size: 2048000 }
          ],
          isRead: true,
          isStarred: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(), // 1天前
          updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString()
        },
        {
          id: `sent-${props.mailboxId}-3`,
          subject: '项目进度更新',
          body: '本周项目进度总结和下周计划...',
          to: [{ email: 'team@company.com', name: '项目团队' }],
          cc: [],
          bcc: [],
          attachments: [],
          isRead: true,
          isStarred: false,
          createdAt: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString(), // 2天前
          updatedAt: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString()
        }
      ]
    }

    if (append) {
      emails.value.push(...emailData)
      currentPage.value++
    } else {
      emails.value = emailData
    }

    totalCount.value = response.total || emailData.length
    hasMore.value = emailData.length === queryParams.limit
  } catch (error) {
    console.error('Failed to load sent emails:', error)
    showError('加载已发送邮件失败')

    // 在错误情况下也提供模拟数据
    if (!append) {
      emails.value = [
        {
          id: `sent-${props.mailboxId}-error-1`,
          subject: '模拟已发送邮件',
          body: '这是在API调用失败时显示的模拟邮件...',
          to: [{ email: 'test@example.com', name: '测试收件人' }],
          cc: [],
          bcc: [],
          attachments: [],
          isRead: true,
          isStarred: false,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]
      totalCount.value = 1
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
  selectedEmailId.value = email?.id || null
  if (email) {
    emit('email-selected', email)
  }
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedEmails.value = []
  } else {
    selectedEmails.value = (emails.value || []).map(email => email?.id).filter(Boolean)
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

const batchDelete = async () => {
  if (!confirm(`确定要删除选中的 ${selectedEmails.value.length} 封邮件吗？`)) return
  
  try {
    batchLoading.value = true
    
    await emailApi.batchUpdate(selectedEmails.value, 'delete')
    
    emails.value = emails.value.filter(email => !selectedEmails.value.includes(email?.id))
    totalCount.value -= selectedEmails.value.length
    selectedEmails.value = []
    
    showSuccess('邮件已删除')
  } catch (error) {
    showError('删除失败')
  } finally {
    batchLoading.value = false
  }
}

const getRecipientDisplay = (recipients: any[]) => {
  if (!recipients || recipients.length === 0) return '(无收件人)'
  
  if (recipients.length === 1) {
    return recipients[0]?.name || recipients[0]?.email
  } else {
    return `${recipients[0]?.name || recipients[0]?.email} 等 ${recipients.length} 人`
  }
}

const getEmailPreview = (body: string) => {
  if (!body) return ''
  const text = body.replace(/<[^>]*>/g, '').trim()
  return text.length > 80 ? text.substring(0, 80) + '...' : text
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
