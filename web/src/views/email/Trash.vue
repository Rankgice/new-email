<template>
  <div class="h-screen flex bg-background-primary">
    <!-- 侧边栏 -->
    <EmailSidebar />

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部工具栏 -->
      <EmailToolbar />

      <!-- 垃圾箱列表 -->
      <div class="flex-1 overflow-hidden p-4">
        <div class="h-full flex flex-col">
          <!-- 页面标题 -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-3">
              <TrashIcon class="w-6 h-6 text-text-primary" />
              <h1 class="text-2xl font-bold text-text-primary">垃圾箱</h1>
              <span v-if="totalCount > 0" class="text-sm text-text-secondary">
                ({{ totalCount }})
              </span>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center space-x-3">
              <!-- 清空垃圾箱 -->
              <Button
                v-if="emails && emails.length > 0"
                variant="secondary"
                @click="emptyTrash"
                :loading="emptyLoading"
              >
                <TrashIcon class="w-4 h-4 mr-2" />
                清空垃圾箱
              </Button>

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
                  variant="primary"
                  size="sm"
                  @click="batchRestore"
                  :loading="batchLoading"
                >
                  <ArrowUturnLeftIcon class="w-4 h-4 mr-2" />
                  恢复 ({{ selectedEmails.length }})
                </Button>
                <Button
                  variant="secondary"
                  size="sm"
                  @click="batchPermanentDelete"
                  :loading="batchLoading"
                >
                  <XMarkIcon class="w-4 h-4 mr-2" />
                  永久删除 ({{ selectedEmails.length }})
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
                  <TrashIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
                  <h3 class="text-lg font-medium text-text-primary mb-2">
                    垃圾箱为空
                  </h3>
                  <p class="text-text-secondary mb-4">
                    已删除的邮件会在这里保存30天
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
                    <div class="col-span-3">发件人</div>
                    <div class="col-span-5">主题</div>
                    <div class="col-span-2">删除时间</div>
                    <div class="col-span-2">操作</div>
                  </div>
                </div>

                <!-- 邮件项目 -->
                <div class="flex-1 overflow-y-auto">
                  <div
                    v-for="email in (emails || [])"
                    :key="email?.id || Math.random()"
                    :class="[
                      'flex items-center px-4 py-3 border-b border-glass-border hover:bg-white/5 transition-colors',
                      email?.id && selectedEmails.includes(email.id) ? 'bg-primary-500/10' : ''
                    ]"
                  >
                    <!-- 选择框 -->
                    <input
                      type="checkbox"
                      :checked="email?.id && selectedEmails.includes(email.id)"
                      @change="email?.id && toggleEmailSelection(email.id)"
                      class="rounded border-gray-600 bg-gray-700 text-primary-500 focus:ring-primary-500 mr-3"
                    />

                    <!-- 邮件信息 -->
                    <div class="flex-1 grid grid-cols-12 gap-4 items-center">
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

                      <!-- 删除时间 -->
                      <div class="col-span-2">
                        <div class="text-xs text-text-secondary">
                          {{ formatDate(email?.deletedAt || email?.updatedAt) }}
                        </div>
                        <div class="text-xs text-orange-400">
                          {{ getDaysUntilPermanentDelete(email?.deletedAt) }}
                        </div>
                      </div>

                      <!-- 操作按钮 -->
                      <div class="col-span-2 flex items-center space-x-2">
                        <Button
                          variant="ghost"
                          size="sm"
                          @click="email?.id && restoreEmail(email.id)"
                          title="恢复邮件"
                        >
                          <ArrowUturnLeftIcon class="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          @click="email?.id && permanentDeleteEmail(email.id)"
                          title="永久删除"
                        >
                          <XMarkIcon class="w-4 h-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          @click="email && viewEmail(email)"
                          title="查看邮件"
                        >
                          <EyeIcon class="w-4 h-4" />
                        </Button>
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

import EmailSidebar from '@/components/email/EmailSidebar.vue'
import EmailToolbar from '@/components/email/EmailToolbar.vue'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'

import {
  TrashIcon,
  ArrowPathIcon,
  ArrowUturnLeftIcon,
  XMarkIcon,
  EyeIcon
} from '@heroicons/vue/24/outline'

// 路由和通知
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const loading = ref(false)
const loadingMore = ref(false)
const batchLoading = ref(false)
const emptyLoading = ref(false)
const emails = ref<Email[]>([])
const selectedEmail = ref<Email | null>(null)
const selectedEmails = ref<string[]>([])
const totalCount = ref(0)
const currentPage = ref(1)
const hasMore = ref(true)

// 查询参数
const queryParams = reactive<EmailListParams>({
  page: 1,
  limit: 20,
  sortBy: 'deletedAt',
  sortOrder: 'desc'
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
      ...queryParams,
      page: append ? currentPage.value + 1 : 1
    }

    const response = await emailApi.getTrashEmails(params)

    if (append) {
      emails.value.push(...response.data)
      currentPage.value++
    } else {
      emails.value = response.data
    }

    totalCount.value = response.total
    hasMore.value = response.data.length === queryParams.limit
  } catch (error) {
    showError('加载垃圾箱邮件失败')
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

const viewEmail = (email: Email) => {
  // TODO: 实现邮件详情查看
  console.log('查看邮件:', email)
}

const restoreEmail = async (emailId: string) => {
  try {
    await emailApi.restoreEmail(emailId)

    // 从垃圾箱列表中移除
    emails.value = emails.value.filter(email => email.id !== emailId)
    totalCount.value--
    selectedEmails.value = selectedEmails.value.filter(id => id !== emailId)

    showSuccess('邮件已恢复')
  } catch (error) {
    showError('恢复邮件失败')
  }
}

const permanentDeleteEmail = async (emailId: string) => {
  if (!confirm('确定要永久删除这封邮件吗？此操作不可撤销。')) return

  try {
    await emailApi.permanentDeleteEmail(emailId)

    // 从列表中移除
    emails.value = emails.value.filter(email => email.id !== emailId)
    totalCount.value--
    selectedEmails.value = selectedEmails.value.filter(id => id !== emailId)

    showSuccess('邮件已永久删除')
  } catch (error) {
    showError('删除失败')
  }
}

const batchRestore = async () => {
  try {
    batchLoading.value = true

    await emailApi.batchUpdate(selectedEmails.value, 'restore')

    // 从垃圾箱列表中移除
    emails.value = emails.value.filter(email => !selectedEmails.value.includes(email.id))
    totalCount.value -= selectedEmails.value.length
    selectedEmails.value = []

    showSuccess('邮件已恢复')
  } catch (error) {
    showError('批量恢复失败')
  } finally {
    batchLoading.value = false
  }
}

const batchPermanentDelete = async () => {
  if (!confirm(`确定要永久删除选中的 ${selectedEmails.value.length} 封邮件吗？此操作不可撤销。`)) return

  try {
    batchLoading.value = true

    await emailApi.batchUpdate(selectedEmails.value, 'permanent_delete')

    // 从列表中移除
    emails.value = emails.value.filter(email => !selectedEmails.value.includes(email.id))
    totalCount.value -= selectedEmails.value.length
    selectedEmails.value = []

    showSuccess('邮件已永久删除')
  } catch (error) {
    showError('批量删除失败')
  } finally {
    batchLoading.value = false
  }
}

const emptyTrash = async () => {
  if (!confirm('确定要清空垃圾箱吗？所有邮件将被永久删除，此操作不可撤销。')) return

  try {
    emptyLoading.value = true

    await emailApi.emptyTrash()

    emails.value = []
    totalCount.value = 0
    selectedEmails.value = []

    showSuccess('垃圾箱已清空')
  } catch (error) {
    showError('清空垃圾箱失败')
  } finally {
    emptyLoading.value = false
  }
}

const handleEmailDeleted = (emailId: string) => {
  emails.value = emails.value.filter(email => email.id !== emailId)
  totalCount.value--
}

const getEmailPreview = (body: string) => {
  if (!body) return ''
  const text = body.replace(/<[^>]*>/g, '').trim()
  return text.length > 50 ? text.substring(0, 50) + '...' : text
}

const formatDate = (dateString: string) => {
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

const getDaysUntilPermanentDelete = (deletedAt: string) => {
  if (!deletedAt) return ''

  const deleteDate = new Date(deletedAt)
  const permanentDeleteDate = new Date(deleteDate.getTime() + 30 * 24 * 60 * 60 * 1000) // 30天后
  const now = new Date()
  const diffTime = permanentDeleteDate.getTime() - now.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays <= 0) {
    return '即将永久删除'
  } else if (diffDays === 1) {
    return '1天后永久删除'
  } else {
    return `${diffDays}天后永久删除`
  }
}
</script>
