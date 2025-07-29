<template>
  <div class="w-96 bg-background-secondary border-l border-glass-border flex flex-col">
    <!-- 预览头部 -->
    <div class="p-4 border-b border-glass-border">
      <h2 class="text-lg font-medium text-text-primary">
        邮件预览
      </h2>
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
                {{ getInitials(selectedEmail.from.name || selectedEmail.from.email) }}
              </span>
            </div>
            <div class="flex-1">
              <p class="text-sm font-medium text-text-primary">
                {{ selectedEmail.from.name || selectedEmail.from.email }}
              </p>
              <p class="text-xs text-text-secondary">
                {{ formatDate(selectedEmail.createdAt) }}
              </p>
            </div>
          </div>

          <!-- 收件人信息 -->
          <div class="text-sm text-text-secondary">
            <p>收件人: {{ selectedEmail.to.map(t => t.name || t.email).join(', ') }}</p>
            <p v-if="selectedEmail.cc && selectedEmail.cc.length > 0">
              抄送: {{ selectedEmail.cc.map(c => c.name || c.email).join(', ') }}
            </p>
          </div>
        </div>

        <!-- 邮件内容 -->
        <div class="prose prose-sm max-w-none text-text-primary">
          <div v-html="selectedEmail.htmlContent || selectedEmail.content.replace(/\n/g, '<br>')"></div>
        </div>

        <!-- 附件 -->
        <div
          v-if="selectedEmail.attachments && selectedEmail.attachments.length > 0"
          class="mt-6 pt-4 border-t border-glass-border"
        >
          <h4 class="text-sm font-medium text-text-primary mb-3">
            附件 ({{ selectedEmail.attachments.length }})
          </h4>
          <div class="space-y-2">
            <div
              v-for="attachment in selectedEmail.attachments"
              :key="attachment.id"
              class="flex items-center space-x-3 p-2 glass-card rounded-lg"
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
              >
                <ArrowDownTrayIcon class="w-4 h-4" />
              </Button>
            </div>
          </div>
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
import { ref } from 'vue'
import type { Email, Attachment } from '@/types'
import Button from '@/components/ui/Button.vue'
import {
  EnvelopeIcon,
  PaperClipIcon,
  ArrowDownTrayIcon,
  ArrowUturnLeftIcon
} from '@heroicons/vue/24/outline'

// 当前选中的邮件 (这里用模拟数据，实际应该从状态管理获取)
const selectedEmail = ref<Email | null>(null)

// 获取姓名首字母
const getInitials = (name: string) => {
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
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

// 下载附件
const downloadAttachment = (attachment: Attachment) => {
  console.log('Download attachment:', attachment.filename)
  // TODO: 实现附件下载
}

// 回复邮件
const handleReply = () => {
  console.log('Reply to email')
  // TODO: 跳转到写邮件页面，预填回复内容
}

// 全部回复
const handleReplyAll = () => {
  console.log('Reply all')
  // TODO: 跳转到写邮件页面，预填全部回复内容
}

// 转发邮件
const handleForward = () => {
  console.log('Forward email')
  // TODO: 跳转到写邮件页面，预填转发内容
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
}

.prose p {
  margin-bottom: 1em;
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  color: inherit;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
}

.prose a {
  color: var(--color-primary);
  text-decoration: underline;
}

.prose a:hover {
  opacity: 0.8;
}
</style>
