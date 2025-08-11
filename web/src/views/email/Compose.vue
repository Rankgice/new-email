<template>
  <div class="h-screen bg-background-primary">
    <div class="container-responsive h-full py-6">
      <GlassCard padding="lg" class="h-full flex flex-col">
        <!-- 头部 -->
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center space-x-4">
            <Button
              variant="ghost"
              @click="$router.back()"
            >
              <ArrowLeftIcon class="w-4 h-4 mr-2" />
              返回
            </Button>
            <h1 class="text-xl font-semibold text-text-primary">
              ✏️ 写邮件
            </h1>
          </div>
          
          <div class="flex items-center space-x-2">
            <Button
              variant="secondary"
              :loading="isSavingDraft"
              @click="saveDraft"
            >
              💾 保存草稿
            </Button>
            <Button
              variant="primary"
              :loading="isSending"
              @click="sendEmail"
            >
              🚀 发送
            </Button>
          </div>
        </div>

        <!-- 邮件表单 -->
        <div class="flex-1 space-y-4">
          <!-- 发件人选择 -->
          <div>
            <label class="block text-sm font-medium text-text-primary mb-2">
              发件人
            </label>
            <select
              v-model="form.mailboxId"
              class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
              @change="handleMailboxChange"
            >
              <option value="">请选择发件邮箱</option>
              <option
                v-for="mailbox in activeMailboxes"
                :key="mailbox.id"
                :value="mailbox.id"
              >
                {{ mailbox.email }}{{ mailbox.name ? ` (${mailbox.name})` : '' }}
              </option>
            </select>
          </div>

          <!-- 收件人 -->
          <Input
            v-model="form.to"
            label="收件人"
            placeholder="输入收件人邮箱地址"
            :left-icon="UserIcon"
          />

          <!-- 抄送和密送 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              v-model="form.cc"
              label="抄送"
              placeholder="输入抄送邮箱地址（可选）"
              :left-icon="UserGroupIcon"
            />
            <Input
              v-model="form.bcc"
              label="密送"
              placeholder="输入密送邮箱地址（可选）"
              :left-icon="EyeSlashIcon"
            />
          </div>

          <!-- 主题 -->
          <Input
            v-model="form.subject"
            label="主题"
            placeholder="输入邮件主题"
            :left-icon="ChatBubbleLeftRightIcon"
          />

          <!-- 内容编辑器 -->
          <div>
            <label class="block text-sm font-medium text-text-primary mb-2">
              邮件内容
            </label>
            <textarea
              v-model="form.content"
              rows="12"
              class="w-full glass-card rounded-lg p-4 text-text-primary placeholder-text-secondary focus-ring resize-none"
              placeholder="输入邮件内容..."
            />
          </div>

          <!-- 附件 -->
          <div>
            <label class="block text-sm font-medium text-text-primary mb-2">
              附件
            </label>

            <!-- 文件上传区域 -->
            <div
              class="glass-card rounded-lg p-4 border-2 border-dashed border-glass-border transition-colors"
              :class="{
                'border-primary-400 bg-primary-50/10': isDragOver,
                'border-glass-border': !isDragOver
              }"
              @drop="handleDrop"
              @dragover="handleDragOver"
              @dragenter="handleDragEnter"
              @dragleave="handleDragLeave"
            >
              <div class="text-center">
                <PaperClipIcon class="w-8 h-8 text-text-secondary mx-auto mb-2" />
                <p class="text-sm text-text-secondary">
                  拖拽文件到此处或
                  <button
                    type="button"
                    class="text-primary-400 hover:text-primary-300"
                    @click="triggerFileInput"
                  >
                    点击选择文件
                  </button>
                </p>
                <p class="text-xs text-text-secondary mt-1">
                  支持的文件类型：{{ allowedTypes.join(', ') }}
                </p>
                <p class="text-xs text-text-secondary">
                  最大文件大小：{{ formatFileSize(maxFileSize) }}
                </p>
              </div>
            </div>

            <!-- 隐藏的文件输入 -->
            <input
              ref="fileInput"
              type="file"
              multiple
              :accept="acceptedFileTypes"
              class="hidden"
              @change="handleFileSelect"
            />

            <!-- 附件列表 -->
            <div v-if="attachments.length > 0" class="mt-4 space-y-2">
              <div
                v-for="(attachment, index) in attachments"
                :key="index"
                class="flex items-center justify-between p-3 glass-card rounded-lg"
              >
                <div class="flex items-center space-x-3">
                  <div class="flex-shrink-0">
                    <DocumentIcon class="w-5 h-5 text-text-secondary" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-text-primary truncate">
                      {{ attachment.name }}
                    </p>
                    <p class="text-xs text-text-secondary">
                      {{ formatFileSize(attachment.size) }}
                    </p>
                  </div>
                </div>
                <button
                  type="button"
                  class="flex-shrink-0 p-1 text-text-secondary hover:text-red-400 transition-colors"
                  @click="removeAttachment(index)"
                >
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </GlassCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useNotification } from '@/composables/useNotification'
import type { Mailbox, EmailSendRequest, AttachmentData } from '@/types'
import { emailApi } from '@/utils/api'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import {
  ArrowLeftIcon,
  UserIcon,
  UserGroupIcon,
  ChatBubbleLeftRightIcon,
  PaperClipIcon,
  EyeSlashIcon,
  DocumentIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const { showNotification } = useNotification()

// 表单数据
const form = reactive({
  mailboxId: '',
  to: '',
  cc: '',
  bcc: '',
  subject: '',
  content: '',
  contentType: 'html' as 'text' | 'html'
})

// 状态
const isSending = ref(false)
const isSavingDraft = ref(false)
const isLoading = ref(false)
const activeMailboxes = ref<Mailbox[]>([])

// 附件相关状态
const attachments = ref<File[]>([])
const isDragOver = ref(false)
const fileInput = ref<HTMLInputElement>()

// 附件配置
const maxFileSize = 10 * 1024 * 1024 // 10MB
const allowedTypes = ['jpg', 'jpeg', 'png', 'gif', 'pdf', 'doc', 'docx', 'xls', 'xlsx', 'txt', 'zip']
const acceptedFileTypes = allowedTypes.map(type => `.${type}`).join(',')

// 计算属性
const selectedMailbox = computed(() => {
  return activeMailboxes.value.find(m => m.id.toString() === form.mailboxId)
})

// 方法
const loadActiveMailboxes = async () => {
  isLoading.value = true
  try {
    const response = await emailApi.getActiveMailboxes()
    if (response.success && response.data) {
      activeMailboxes.value = response.data
      // 如果只有一个邮箱，自动选择
      if (activeMailboxes.value.length === 1) {
        form.mailboxId = activeMailboxes.value[0].id.toString()
        handleMailboxChange()
      }
    }
  } catch (error) {
    console.error('Failed to load mailboxes:', error)
    showNotification({
      type: 'error',
      title: '加载失败',
      message: '无法加载邮箱列表'
    })
  } finally {
    isLoading.value = false
  }
}

const handleMailboxChange = () => {
  // 当选择邮箱时，自动设置发件人邮箱地址
  if (selectedMailbox.value) {
    // 这里可以添加其他逻辑，比如加载签名等
  }
}

// 附件处理方法
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    addFiles(Array.from(target.files))
  }
  // 清空input值，允许重复选择同一文件
  target.value = ''
}

const handleDragEnter = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
}

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault()
  // 只有当离开整个拖拽区域时才设置为false
  if (!event.currentTarget?.contains(event.relatedTarget as Node)) {
    isDragOver.value = false
  }
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false

  if (event.dataTransfer?.files) {
    addFiles(Array.from(event.dataTransfer.files))
  }
}

const addFiles = (files: File[]) => {
  for (const file of files) {
    // 检查文件类型
    const fileExtension = file.name.split('.').pop()?.toLowerCase()
    if (!fileExtension || !allowedTypes.includes(fileExtension)) {
      showNotification({
        type: 'error',
        title: '文件类型不支持',
        message: `文件 ${file.name} 的类型不支持。支持的类型：${allowedTypes.join(', ')}`
      })
      continue
    }

    // 检查文件大小
    if (file.size > maxFileSize) {
      showNotification({
        type: 'error',
        title: '文件过大',
        message: `文件 ${file.name} 超过最大限制 ${formatFileSize(maxFileSize)}`
      })
      continue
    }

    // 检查是否已存在同名文件
    if (attachments.value.some(att => att.name === file.name)) {
      showNotification({
        type: 'warning',
        title: '文件已存在',
        message: `文件 ${file.name} 已经添加过了`
      })
      continue
    }

    attachments.value.push(file)
  }
}

const removeAttachment = (index: number) => {
  attachments.value.splice(index, 1)
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 处理附件，将文件转换为Base64
const processAttachments = async (): Promise<AttachmentData[]> => {
  const attachmentData: AttachmentData[] = []

  for (const file of attachments.value) {
    try {
      const base64Data = await fileToBase64(file)
      attachmentData.push({
        filename: file.name,
        contentType: file.type || 'application/octet-stream',
        data: base64Data,
        size: file.size
      })
    } catch (error) {
      console.error('Failed to process attachment:', file.name, error)
      showNotification({
        type: 'error',
        title: '附件处理失败',
        message: `无法处理附件 ${file.name}`
      })
    }
  }

  return attachmentData
}

// 将文件转换为Base64字符串
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      const result = reader.result as string
      // 移除data:xxx;base64,前缀，只保留Base64数据
      const base64Data = result.split(',')[1]
      resolve(base64Data)
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}



const validateForm = () => {
  if (!form.mailboxId) {
    showNotification({
      type: 'warning',
      title: '请选择发件邮箱',
      message: '请先选择一个发件邮箱'
    })
    return false
  }

  if (!form.to.trim()) {
    showNotification({
      type: 'warning',
      title: '请输入收件人',
      message: '收件人不能为空'
    })
    return false
  }

  if (!form.subject.trim()) {
    showNotification({
      type: 'warning',
      title: '请输入主题',
      message: '邮件主题不能为空'
    })
    return false
  }

  if (!form.content.trim()) {
    showNotification({
      type: 'warning',
      title: '请输入内容',
      message: '邮件内容不能为空'
    })
    return false
  }

  return true
}

// 解析邮件地址字符串为数组
const parseEmailAddresses = (emailString: string): string[] => {
  if (!emailString.trim()) return []

  // 支持多种分隔符：逗号、分号、空格
  return emailString
    .split(/[,;，；\s]+/)
    .map(email => email.trim())
    .filter(email => email.length > 0)
}

const sendEmail = async () => {
  if (!validateForm()) return

  isSending.value = true
  try {
    // 处理附件
    const attachmentData = await processAttachments()

    const sendData: EmailSendRequest = {
      mailboxId: parseInt(form.mailboxId),
      subject: form.subject,
      fromEmail: selectedMailbox.value?.email || '',
      toEmail: parseEmailAddresses(form.to),
      ccEmail: form.cc ? parseEmailAddresses(form.cc) : undefined,
      bccEmail: form.bcc ? parseEmailAddresses(form.bcc) : undefined,
      content: form.content,
      contentType: form.contentType,
      attachments: attachmentData.length > 0 ? attachmentData : undefined
    }

    console.log('Sending email with data:', sendData)

    const response = await emailApi.sendEmail(sendData)
    if (response.success) {
      showNotification({
        type: 'success',
        title: '邮件已发送',
        message: '您的邮件已成功发送'
      })
      router.push('/inbox')
    }
  } catch (error) {
    console.error('Failed to send email:', error)
    showNotification({
      type: 'error',
      title: '发送失败',
      message: '邮件发送失败，请重试'
    })
  } finally {
    isSending.value = false
  }
}

const saveDraft = async () => {
  isSavingDraft.value = true
  try {
    // TODO: 实现保存草稿逻辑
    await new Promise(resolve => setTimeout(resolve, 1000))
    showNotification({
      type: 'success',
      title: '草稿已保存'
    })
  } catch (error) {
    console.error('Failed to save draft:', error)
    showNotification({
      type: 'error',
      title: '保存失败',
      message: '草稿保存失败，请重试'
    })
  } finally {
    isSavingDraft.value = false
  }
}

// 生命周期
onMounted(() => {
  loadActiveMailboxes()
})
</script>

<style scoped>
.focus-ring {
  @apply focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-background-primary;
}

.glass-card {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  background: var(--color-glass-light);
  border: 1px solid var(--color-glass-border);
}
</style>
