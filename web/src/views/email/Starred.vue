<template>
  <div class="h-screen flex bg-background-primary">
    <!-- 侧边栏 -->
    <EmailSidebar />

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部工具栏 -->
      <EmailToolbar />

      <!-- 星标邮件页面 -->
      <div class="flex-1 overflow-hidden p-4">
        <div class="h-full flex flex-col">
          <!-- 页面标题 -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-3">
              <StarIcon class="w-6 h-6 text-yellow-400" />
              <h1 class="text-2xl font-bold text-text-primary">已加星标</h1>
              <span v-if="totalCount > 0" class="text-sm text-text-secondary">
                ({{ totalCount }})
              </span>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center space-x-3">
              <!-- 刷新按钮 -->
              <Button
                variant="ghost"
                size="sm"
                @click="refreshEmails"
                :loading="loading"
              >
                <ArrowPathIcon class="w-4 h-4 mr-2" />
                刷新
              </Button>

              <!-- 批量操作 -->
              <div v-if="selectedEmails.length > 0" class="flex items-center space-x-2">
                <Button
                  variant="secondary"
                  size="sm"
                  @click="batchRemoveStar"
                  :loading="batchLoading"
                >
                  <StarIcon class="w-4 h-4 mr-2" />
                  取消星标 ({{ selectedEmails.length }})
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  @click="batchDelete"
                  :loading="batchLoading"
                >
                  <TrashIcon class="w-4 h-4 mr-2" />
                  删除 ({{ selectedEmails.length }})
                </Button>
              </div>
            </div>
          </div>

          <!-- 邮件列表 -->
          <div class="flex-1 overflow-hidden">
            <GlassCard padding="none" class="h-full">
              <!-- 空状态 -->
              <div v-if="!loading && (!emails || emails.length === 0)" class="h-full flex items-center justify-center">
                <div class="text-center">
                  <StarIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
                  <h3 class="text-lg font-medium text-text-primary mb-2">
                    暂无星标邮件
                  </h3>
                  <p class="text-text-secondary mb-4">
                    为重要邮件添加星标，方便快速查找
                  </p>
                  <Button
                    variant="primary"
                    @click="$router.push('/inbox')"
                  >
                    前往收件箱
                  </Button>
                </div>
              </div>

              <!-- 邮件列表 -->
              <div v-else class="h-full flex flex-col">
                <!-- 列表头部 -->
                <div class="flex items-center px-4 py-3 border-b border-glass-border bg-white/5">
                  <input
                    type="checkbox"
                    :checked="isAllSelected"
                    @change="toggleSelectAll"
                    class="rounded border-gray-600 bg-gray-700 text-primary-500 focus:ring-primary-500 mr-3"
                  />
                  <div class="flex-1 grid grid-cols-12 gap-4 text-sm font-medium text-text-secondary">
                    <div class="col-span-1">星标</div>
                    <div class="col-span-3">发件人</div>
                    <div class="col-span-5">主题</div>
                    <div class="col-span-2">附件</div>
                    <div class="col-span-1">时间</div>
                  </div>
                </div>

                <!-- 邮件项目 -->
                <div class="flex-1 overflow-y-auto">
                  <div
                    v-for="email in (emails || [])"
                    :key="email?.id || Math.random()"
                    :class="[
                      'flex items-center px-4 py-3 border-b border-glass-border hover:bg-white/5 cursor-pointer transition-colors',
                      email?.id && selectedEmails.includes(email.id) ? 'bg-primary-500/10' : '',
                      !email?.isRead ? 'font-semibold' : ''
                    ]"
                    @click="email && openEmail(email)"
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
                    <div class="flex-1 grid grid-cols-12 gap-4 items-center">
                      <!-- 星标 -->
                      <div class="col-span-1">
                        <button
                          @click.stop="email?.id && toggleStar(email.id, false)"
                          class="text-yellow-400 hover:text-yellow-300 transition-colors"
                        >
                          <StarIcon class="w-5 h-5 fill-current" />
                        </button>
                      </div>

                      <!-- 发件人 -->
                      <div class="col-span-3">
                        <div class="text-text-primary truncate">
                          {{ email?.from?.name || email?.from?.email || '未知发件人' }}
                        </div>
                        <div class="text-xs text-text-secondary truncate">
                          {{ email?.from?.email || '' }}
                        </div>
                      </div>

                      <!-- 主题 -->
                      <div class="col-span-5">
                        <div class="text-text-primary truncate">
                          {{ email?.subject || '(无主题)' }}
                        </div>
                        <div class="text-xs text-text-secondary truncate">
                          {{ getEmailPreview(email?.body) }}
                        </div>
                      </div>

                      <!-- 附件 -->
                      <div class="col-span-2">
                        <div v-if="email?.attachments && email.attachments.length > 0" class="flex items-center text-text-secondary">
                          <PaperClipIcon class="w-4 h-4 mr-1" />
                          <span class="text-xs">{{ email.attachments.length }}</span>
                        </div>
                      </div>

                      <!-- 时间 -->
                      <div class="col-span-1 text-right">
                        <div class="text-xs text-text-secondary">
                          {{ formatDate(email?.createdAt) }}
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 加载更多 -->
                  <div v-if="hasMore" class="p-4 text-center">
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
            </GlassCard>
          </div>
        </div>
      </div>
    </div>

    <!-- TODO: 添加邮件详情模态框 -->
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'
import type { Email, EmailListParams } from '@/types'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import EmailSidebar from '@/components/email/EmailSidebar.vue'
import EmailToolbar from '@/components/email/EmailToolbar.vue'

import {
  StarIcon,
  ArrowPathIcon,
  TrashIcon,
  PaperClipIcon
} from '@heroicons/vue/24/outline'

// 路由和通知
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const loading = ref(false)
const loadingMore = ref(false)
const batchLoading = ref(false)

const emails = ref<Email[]>([])
const selectedEmail = ref<Email | null>(null)
const selectedEmails = ref<string[]>([])
const totalCount = ref(0)
const currentPage = ref(1)
const hasMore = ref(true)

// 查询参数
const queryParams = reactive<EmailListParams>({
  page: 1,
  pageSize: 20
})

// 计算属性
const isAllSelected = computed(() => {
  return emails.value && emails.value.length > 0 && selectedEmails.value.length === emails.value.length
})

// 生命周期
onMounted(() => {
  loadEmails()
})

// 方法
const loadEmails = async (append = false) => {
  try {
    if (!append) {
      loading.value = true
      currentPage.value = 1
    } else {
      loadingMore.value = true
    }

    const params = {
      page: append ? currentPage.value + 1 : 1,
      pageSize: queryParams.pageSize
    }

    const response = await emailApi.getStarredEmails(params)
    console.log('Starred emails API response:', response)

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
    console.error('Failed to load starred emails:', error)
    showError('加载星标邮件失败')

    // 不使用模拟数据，保持空状态
    if (!append) {
      emails.value = []
      totalCount.value = 0
      hasMore.value = false
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const refreshEmails = () => {
  selectedEmails.value = []
  loadEmails()
}

const loadMore = () => {
  loadEmails(true)
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

const openEmail = (email: Email) => {
  // TODO: 实现邮件详情查看
  console.log('查看邮件:', email)
  
  // 标记为已读
  if (!email.isRead) {
    markAsRead(email.id, true)
  }
}

const toggleStar = async (emailId: string, starred: boolean) => {
  try {
    await emailApi.markAsStarred(emailId, starred)
    
    if (!starred) {
      // 从星标列表中移除
      emails.value = emails.value.filter(email => email.id !== emailId)
      totalCount.value--
      selectedEmails.value = selectedEmails.value.filter(id => id !== emailId)
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

const batchRemoveStar = async () => {
  try {
    batchLoading.value = true
    
    await emailApi.batchUpdate(selectedEmails.value, 'unstar')
    
    // 从列表中移除
    emails.value = emails.value.filter(email => !selectedEmails.value.includes(email.id))
    totalCount.value -= selectedEmails.value.length
    selectedEmails.value = []
    
    showSuccess('已取消星标')
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

const handleStarChanged = (emailId: string, starred: boolean) => {
  if (!starred) {
    // 从星标列表中移除
    emails.value = emails.value.filter(email => email.id !== emailId)
    totalCount.value--
  }
}

const handleEmailDeleted = (emailId: string) => {
  emails.value = emails.value.filter(email => email.id !== emailId)
  totalCount.value--
}

const getEmailPreview = (body: string) => {
  if (!body) return ''
  // 移除HTML标签并截取前50个字符
  const text = body.replace(/<[^>]*>/g, '').trim()
  return text.length > 50 ? text.substring(0, 50) + '...' : text
}

const formatDate = (dateString: string) => {
  if (!dateString) return ''

  const date = new Date(dateString)
  if (isNaN(date.getTime())) return '' // 检查日期是否有效

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
