<template>
  <div class="h-screen flex bg-background-primary">
    <!-- 侧边栏 -->
    <EmailSidebar />

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部工具栏 -->
      <EmailToolbar />

      <!-- 草稿箱列表 -->
      <div class="flex-1 overflow-hidden p-4">
        <div class="h-full flex flex-col">
          <!-- 页面标题 -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-3">
              <DocumentIcon class="w-6 h-6 text-text-primary" />
              <h1 class="text-2xl font-bold text-text-primary">草稿箱</h1>
              <span v-if="totalCount > 0" class="text-sm text-text-secondary">
                ({{ totalCount }})
              </span>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center space-x-3">
              <!-- 新建草稿 -->
              <Button
                variant="primary"
                @click="createNewDraft"
              >
                <PlusIcon class="w-4 h-4 mr-2" />
                新建草稿
              </Button>

              <!-- 刷新按钮 -->
              <Button
                variant="ghost"
                size="sm"
                @click="refreshDrafts"
                :loading="loading"
              >
                <ArrowPathIcon class="w-4 h-4 mr-2" />
                刷新
              </Button>

              <!-- 批量操作 -->
              <div v-if="selectedDrafts.length > 0" class="flex items-center space-x-2">
                <Button
                  variant="secondary"
                  size="sm"
                  @click="batchDelete"
                  :loading="batchLoading"
                >
                  <TrashIcon class="w-4 h-4 mr-2" />
                  删除 ({{ selectedDrafts.length }})
                </Button>
              </div>
            </div>
          </div>

          <!-- 草稿列表 -->
          <div class="flex-1 overflow-hidden">
            <GlassCard padding="none" class="h-full">
              <!-- 空状态 -->
              <div v-if="!loading && drafts.length === 0" class="h-full flex items-center justify-center">
                <div class="text-center">
                  <DocumentIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
                  <h3 class="text-lg font-medium text-text-primary mb-2">
                    暂无草稿
                  </h3>
                  <p class="text-text-secondary mb-4">
                    开始撰写邮件，系统会自动保存为草稿
                  </p>
                  <Button
                    variant="primary"
                    @click="createNewDraft"
                  >
                    <PlusIcon class="w-4 h-4 mr-2" />
                    新建草稿
                  </Button>
                </div>
              </div>

              <!-- 草稿列表 -->
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
                    <div class="col-span-3">收件人</div>
                    <div class="col-span-5">主题</div>
                    <div class="col-span-2">附件</div>
                    <div class="col-span-2">保存时间</div>
                  </div>
                </div>

                <!-- 草稿项目 -->
                <div class="flex-1 overflow-y-auto">
                  <div
                    v-for="draft in drafts"
                    :key="draft?.id || Math.random()"
                    :class="[
                      'flex items-center px-4 py-3 border-b border-glass-border hover:bg-white/5 cursor-pointer transition-colors',
                      draft?.id && selectedDrafts.includes(draft.id) ? 'bg-primary-500/10' : ''
                    ]"
                    @click="editDraft(draft)"
                  >
                    <!-- 选择框 -->
                    <input
                      type="checkbox"
                      :checked="draft?.id && selectedDrafts.includes(draft.id)"
                      @change="draft?.id && toggleDraftSelection(draft.id)"
                      @click.stop
                      class="rounded border-gray-600 bg-gray-700 text-primary-500 focus:ring-primary-500 mr-3"
                    />

                    <!-- 草稿信息 -->
                    <div class="flex-1 grid grid-cols-12 gap-4 items-center">
                      <!-- 收件人 -->
                      <div class="col-span-3">
                        <div class="text-text-primary truncate">
                          {{ getRecipientDisplay(draft?.to) }}
                        </div>
                        <div v-if="draft?.cc && draft.cc.length > 0" class="text-xs text-text-secondary truncate">
                          抄送: {{ draft.cc.map(c => c?.email).join(', ') }}
                        </div>
                      </div>

                      <!-- 主题 -->
                      <div class="col-span-5">
                        <div class="text-text-primary truncate">
                          {{ draft?.subject || '(无主题)' }}
                        </div>
                        <div class="text-xs text-text-secondary truncate">
                          {{ getEmailPreview(draft?.body) }}
                        </div>
                      </div>

                      <!-- 附件 -->
                      <div class="col-span-2">
                        <div v-if="draft?.attachments && draft.attachments.length > 0" class="flex items-center text-text-secondary">
                          <PaperClipIcon class="w-4 h-4 mr-1" />
                          <span class="text-xs">{{ draft.attachments.length }}</span>
                        </div>
                      </div>

                      <!-- 保存时间 -->
                      <div class="col-span-2 text-right">
                        <div class="text-xs text-text-secondary">
                          {{ formatDate(draft?.updatedAt || draft?.createdAt) }}
                        </div>
                      </div>
                    </div>

                    <!-- 操作按钮 -->
                    <div class="flex items-center space-x-2 ml-4">
                      <Button
                        variant="ghost"
                        size="sm"
                        @click.stop="draft && editDraft(draft)"
                      >
                        <PencilIcon class="w-4 h-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        @click.stop="draft?.id && deleteDraft(draft.id)"
                      >
                        <TrashIcon class="w-4 h-4" />
                      </Button>
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
  DocumentIcon,
  PlusIcon,
  ArrowPathIcon,
  TrashIcon,
  PaperClipIcon,
  PencilIcon
} from '@heroicons/vue/24/outline'

// 路由和通知
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const loading = ref(false)
const loadingMore = ref(false)
const batchLoading = ref(false)

const drafts = ref<Email[]>([])
const selectedDrafts = ref<string[]>([])
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
  return drafts.value.length > 0 && selectedDrafts.value.length === drafts.value.length
})

// 生命周期
onMounted(() => {
  loadDrafts()
})

// 方法
const loadDrafts = async (append = false) => {
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

    const response = await emailApi.getDrafts(params)
    console.log('Drafts API response:', response)

    if (response.success && response.data) {
      const draftList = response.data.list || response.data || []

      if (append) {
        drafts.value.push(...draftList)
        currentPage.value++
      } else {
        drafts.value = draftList
      }

      totalCount.value = response.data.total || draftList.length
      hasMore.value = draftList.length === queryParams.pageSize
    } else {
      if (!append) {
        drafts.value = []
      }
      totalCount.value = 0
      hasMore.value = false
    }
  } catch (error) {
    console.error('Failed to load drafts:', error)
    showError('加载草稿失败')

    // 不使用模拟数据，保持空状态
    if (!append) {
      drafts.value = []
      totalCount.value = 0
      hasMore.value = false
    }
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const refreshDrafts = () => {
  selectedDrafts.value = []
  loadDrafts()
}

const loadMore = () => {
  loadDrafts(true)
}

const createNewDraft = () => {
  router.push('/compose')
}

const editDraft = (draft: Email) => {
  router.push(`/compose?draft=${draft.id}`)
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedDrafts.value = []
  } else {
    selectedDrafts.value = drafts.value.map(draft => draft.id)
  }
}

const toggleDraftSelection = (draftId: string) => {
  const index = selectedDrafts.value.indexOf(draftId)
  if (index > -1) {
    selectedDrafts.value.splice(index, 1)
  } else {
    selectedDrafts.value.push(draftId)
  }
}

const deleteDraft = async (draftId: string) => {
  if (!confirm('确定要删除这个草稿吗？')) return

  try {
    await emailApi.deleteEmail(draftId)
    drafts.value = drafts.value.filter(draft => draft.id !== draftId)
    totalCount.value--
    selectedDrafts.value = selectedDrafts.value.filter(id => id !== draftId)
    showSuccess('草稿已删除')
  } catch (error) {
    showError('删除失败')
  }
}

const batchDelete = async () => {
  if (!confirm(`确定要删除选中的 ${selectedDrafts.value.length} 个草稿吗？`)) return

  try {
    batchLoading.value = true

    await emailApi.batchUpdate(selectedDrafts.value, 'delete')

    drafts.value = drafts.value.filter(draft => !selectedDrafts.value.includes(draft.id))
    totalCount.value -= selectedDrafts.value.length
    selectedDrafts.value = []

    showSuccess('草稿已删除')
  } catch (error) {
    showError('批量删除失败')
  } finally {
    batchLoading.value = false
  }
}

const getRecipientDisplay = (recipients: any[]) => {
  if (!recipients || recipients.length === 0) return '(无收件人)'

  if (recipients.length === 1) {
    return recipients[0].name || recipients[0].email
  } else {
    return `${recipients[0].name || recipients[0].email} 等 ${recipients.length} 人`
  }
}

const getEmailPreview = (body: string) => {
  if (!body) return ''
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
