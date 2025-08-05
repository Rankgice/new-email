<template>
  <div class="w-96 bg-background-secondary border-l border-glass-border flex flex-col">
    <!-- 预览头部 -->
    <div class="p-4 border-b border-glass-border flex items-center justify-between">
      <h2 class="text-lg font-medium text-text-primary">
        邮件详情
      </h2>
      <div v-if="selectedEmail" class="flex items-center space-x-2">
        <!-- 星标按钮 -->
        <Button
          variant="ghost"
          size="sm"
          @click="toggleStar"
          :class="{ 'text-yellow-500': selectedEmail.isStarred }"
        >
          <StarIcon
            :class="selectedEmail.isStarred ? 'fill-current' : ''"
            class="w-4 h-4"
          />
        </Button>
        <!-- 删除按钮 -->
        <Button
          variant="ghost"
          size="sm"
          @click="handleDelete"
          class="text-red-500 hover:text-red-600"
        >
          <TrashIcon class="w-4 h-4" />
        </Button>
        <!-- 更多操作 -->
        <Button
          variant="ghost"
          size="sm"
          @click="showMoreActions = !showMoreActions"
        >
          <EllipsisVerticalIcon class="w-4 h-4" />
        </Button>
      </div>
    </div>

    <!-- 预览内容 -->
    <div class="flex-1 overflow-auto">
      <div
        v-if="selectedEmail"
        class="p-4"
      >
        <!-- 邮件头部信息 -->
        <div class="space-y-3 mb-6">
          <h3 class="text-lg font-medium text-text-primary">
            {{ selectedEmail.subject }}
          </h3>
          
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full flex items-center justify-center">
              <span class="text-white font-medium text-xs">
                {{ getInitials(selectedEmail.fromEmail) }}
              </span>
            </div>
            <div class="flex-1">
              <p class="text-sm font-medium text-text-primary">
                {{ getDisplayName(selectedEmail.fromEmail) }}
              </p>
              <p class="text-xs text-text-secondary">
                {{ formatDate(selectedEmail.createdAt) }}
              </p>
            </div>
          </div>

          <!-- 收件人信息 -->
          <div class="text-sm text-text-secondary space-y-1">
            <p>收件人: {{ selectedEmail.toEmails }}</p>
            <p v-if="selectedEmail.ccEmail">
              抄送: {{ selectedEmail.ccEmail }}
            </p>
            <p v-if="selectedEmail.bccEmail">
              密送: {{ selectedEmail.bccEmail }}
            </p>
          </div>
        </div>

        <!-- 邮件内容 -->
        <div class="prose prose-sm max-w-none text-text-primary">
          <div
            v-if="selectedEmail.contentType === 'html'"
            v-html="selectedEmail.content"
            class="email-content"
          ></div>
          <div
            v-else
            class="whitespace-pre-wrap"
          >{{ selectedEmail.content }}</div>
        </div>

        <!-- 附件 -->
        <div
          v-if="hasAttachments"
          class="mt-6 pt-4 border-t border-glass-border"
        >
          <h4 class="text-sm font-medium text-text-primary mb-3">
            附件 ({{ attachmentCount }})
          </h4>
          <div class="space-y-2">
            <div
              v-for="attachment in attachments"
              :key="attachment.id || attachment.filename"
              class="flex items-center space-x-3 p-2 glass-card rounded-lg hover:bg-glass-light transition-colors"
            >
              <PaperClipIcon class="w-4 h-4 text-text-secondary" />
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-text-primary truncate">
                  {{ attachment.filename }}
                </p>
                <p class="text-xs text-text-secondary">
                  {{ formatFileSize(attachment.size) }}
                </p>
              </div>
              <Button
                variant="ghost"
                size="sm"
                @click="downloadAttachment(attachment)"
                :loading="loading"
              >
                <ArrowDownTrayIcon class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>

        <!-- 更多操作菜单 -->
        <div
          v-if="showMoreActions"
          class="mt-4 p-3 glass-card rounded-lg space-y-2"
        >
          <Button
            variant="ghost"
            size="sm"
            class="w-full justify-start"
            @click="handlePrint"
          >
            <PrinterIcon class="w-4 h-4 mr-2" />
            打印邮件
          </Button>
          <Button
            variant="ghost"
            size="sm"
            class="w-full justify-start"
            @click="handleMarkAsUnread"
          >
            <EnvelopeIcon class="w-4 h-4 mr-2" />
            标记为未读
          </Button>
          <Button
            variant="ghost"
            size="sm"
            class="w-full justify-start"
            @click="handleMarkAsImportant"
          >
            <FlagIcon class="w-4 h-4 mr-2" />
            标记为重要
          </Button>
          <Button
            variant="ghost"
            size="sm"
            class="w-full justify-start"
            @click="handleCopyToClipboard"
          >
            <DocumentDuplicateIcon class="w-4 h-4 mr-2" />
            复制内容
          </Button>
        </div>
      </div>

      <!-- 空状态 -->
      <div
        v-else
        class="flex flex-col items-center justify-center h-full text-center p-8"
      >
        <EnvelopeIcon class="w-16 h-16 text-text-secondary mb-4" />
        <h3 class="text-lg font-medium text-text-primary mb-2">
          选择邮件查看
        </h3>
        <p class="text-text-secondary">
          点击左侧邮件列表中的邮件来查看详细内容
        </p>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div
      v-if="selectedEmail"
      class="p-4 border-t border-glass-border space-y-2"
    >
      <Button
        variant="primary"
        size="sm"
        class="w-full"
        @click="handleReply"
      >
        <ArrowUturnLeftIcon class="w-4 h-4 mr-2" />
        回复
      </Button>
      <div class="flex space-x-2">
        <Button
          variant="secondary"
          size="sm"
          class="flex-1"
          @click="handleReplyAll"
        >
          全部回复
        </Button>
        <Button
          variant="secondary"
          size="sm"
          class="flex-1"
          @click="handleForward"
        >
          转发
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Email, Attachment } from '@/types'
import { emailApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'
import Button from '@/components/ui/Button.vue'
import {
  EnvelopeIcon,
  PaperClipIcon,
  ArrowDownTrayIcon,
  ArrowUturnLeftIcon,
  StarIcon,
  TrashIcon,
  EllipsisVerticalIcon,
  PrinterIcon,
  DocumentDuplicateIcon,
  FlagIcon
} from '@heroicons/vue/24/outline'

// Props
interface Props {
  email?: Email | null
}

const props = withDefaults(defineProps<Props>(), {
  email: null
})

// Emits
interface Emits {
  (e: 'email-updated', email: Email): void
  (e: 'email-deleted', emailId: string): void
}

const emit = defineEmits<Emits>()

// 路由和通知
const router = useRouter()
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const showMoreActions = ref(false)
const loading = ref(false)

// 计算属性
const selectedEmail = computed(() => props.email)

const hasAttachments = computed(() => {
  return selectedEmail.value?.attachments && selectedEmail.value.attachments.length > 0
})

const attachmentCount = computed(() => {
  return selectedEmail.value?.attachments?.length || 0
})

const attachments = computed(() => {
  return selectedEmail.value?.attachments || []
})

// 获取姓名首字母
const getInitials = (email: string) => {
  if (!email) return 'U'
  const name = email.split('@')[0]
  return name.slice(0, 2).toUpperCase()
}

// 获取显示名称
const getDisplayName = (email: string) => {
  if (!email) return '未知发件人'
  return email.split('@')[0]
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''

  const date = new Date(dateString)
  if (isNaN(date.getTime())) return ''

  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 切换星标状态
const toggleStar = async () => {
  if (!selectedEmail.value) return

  try {
    loading.value = true
    const response = await emailApi.markAsStarred(selectedEmail.value.id, !selectedEmail.value.isStarred)

    if (response.success) {
      const updatedEmail = { ...selectedEmail.value, isStarred: !selectedEmail.value.isStarred }
      emit('email-updated', updatedEmail)
      showSuccess(selectedEmail.value.isStarred ? '已取消星标' : '已添加星标')
    } else {
      showError('操作失败')
    }
  } catch (error) {
    console.error('Toggle star error:', error)
    showError('操作失败')
  } finally {
    loading.value = false
  }
}

// 删除邮件
const handleDelete = async () => {
  if (!selectedEmail.value) return

  if (!confirm('确定要删除这封邮件吗？')) return

  try {
    loading.value = true
    const response = await emailApi.deleteEmail(selectedEmail.value.id)

    if (response.success) {
      emit('email-deleted', selectedEmail.value.id)
      showSuccess('邮件已删除')
    } else {
      showError('删除失败')
    }
  } catch (error) {
    console.error('Delete email error:', error)
    showError('删除失败')
  } finally {
    loading.value = false
  }
}

// 下载附件
const downloadAttachment = async (attachment: Attachment) => {
  try {
    loading.value = true
    // 这里应该调用附件下载API
    console.log('Download attachment:', attachment.filename)
    showSuccess('开始下载附件')
  } catch (error) {
    console.error('Download attachment error:', error)
    showError('下载失败')
  } finally {
    loading.value = false
  }
}

// 打印邮件
const handlePrint = () => {
  window.print()
}

// 标记为未读
const handleMarkAsUnread = async () => {
  if (!selectedEmail.value) return

  try {
    loading.value = true
    const response = await emailApi.markAsRead(selectedEmail.value.id, false)

    if (response.success) {
      const updatedEmail = { ...selectedEmail.value, isRead: false }
      emit('email-updated', updatedEmail)
      showSuccess('已标记为未读')
    } else {
      showError('操作失败')
    }
  } catch (error) {
    console.error('Mark as unread error:', error)
    showError('操作失败')
  } finally {
    loading.value = false
  }
}

// 标记为重要
const handleMarkAsImportant = async () => {
  if (!selectedEmail.value) return

  try {
    loading.value = true
    // 这里需要实现标记重要的API
    console.log('Mark as important')
    showSuccess('已标记为重要')
  } catch (error) {
    console.error('Mark as important error:', error)
    showError('操作失败')
  } finally {
    loading.value = false
  }
}

// 复制内容到剪贴板
const handleCopyToClipboard = async () => {
  if (!selectedEmail.value) return

  try {
    const content = `主题: ${selectedEmail.value.subject}\n发件人: ${selectedEmail.value.fromEmail}\n内容: ${selectedEmail.value.content}`
    await navigator.clipboard.writeText(content)
    showSuccess('内容已复制到剪贴板')
  } catch (error) {
    console.error('Copy to clipboard error:', error)
    showError('复制失败')
  }
}

// 回复邮件
const handleReply = () => {
  if (!selectedEmail.value) return

  router.push({
    path: '/compose',
    query: {
      type: 'reply',
      emailId: selectedEmail.value.id
    }
  })
}

// 全部回复
const handleReplyAll = () => {
  if (!selectedEmail.value) return

  router.push({
    path: '/compose',
    query: {
      type: 'replyAll',
      emailId: selectedEmail.value.id
    }
  })
}

// 转发邮件
const handleForward = () => {
  if (!selectedEmail.value) return

  router.push({
    path: '/compose',
    query: {
      type: 'forward',
      emailId: selectedEmail.value.id
    }
  })
}
</script>

<style scoped>
.glass-card {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  background: var(--color-glass-light);
  border: 1px solid var(--color-glass-border);
}

.prose {
  color: inherit;
  line-height: 1.6;
}

.prose p {
  margin-bottom: 1em;
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  color: inherit;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
}

.prose a {
  color: var(--color-primary);
  text-decoration: underline;
}

.prose a:hover {
  opacity: 0.8;
}

.prose blockquote {
  border-left: 4px solid var(--color-glass-border);
  padding-left: 1rem;
  margin: 1rem 0;
  font-style: italic;
  color: var(--color-text-secondary);
}

.prose ul, .prose ol {
  margin: 1rem 0;
  padding-left: 1.5rem;
}

.prose li {
  margin-bottom: 0.5rem;
}

.prose img {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 1rem 0;
}

.prose table {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
}

.prose th, .prose td {
  border: 1px solid var(--color-glass-border);
  padding: 0.5rem;
  text-align: left;
}

.prose th {
  background-color: var(--color-glass-light);
  font-weight: 600;
}

.email-content {
  /* 防止恶意样式影响页面 */
  max-width: 100%;
  overflow-wrap: break-word;
  word-wrap: break-word;
}

.email-content * {
  max-width: 100% !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .w-96 {
    width: 100%;
  }
}

/* 打印样式 */
@media print {
  .border-l, .border-b, .border-t {
    border: none !important;
  }

  .bg-background-secondary {
    background: white !important;
  }

  .text-text-primary, .text-text-secondary {
    color: black !important;
  }

  button {
    display: none !important;
  }
}
</style>
