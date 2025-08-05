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
            <div class="glass-card rounded-lg p-4 border-2 border-dashed border-glass-border">
              <div class="text-center">
                <PaperClipIcon class="w-8 h-8 text-text-secondary mx-auto mb-2" />
                <p class="text-sm text-text-secondary">
                  拖拽文件到此处或
                  <button class="text-primary-400 hover:text-primary-300">
                    点击选择文件
                  </button>
                </p>
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
import type { Mailbox, EmailSendRequest } from '@/types'
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
  EyeSlashIcon
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
    const sendData: EmailSendRequest = {
      mailboxId: parseInt(form.mailboxId),
      subject: form.subject,
      fromEmail: selectedMailbox.value?.email || '',
      toEmail: parseEmailAddresses(form.to),
      ccEmail: form.cc ? parseEmailAddresses(form.cc) : undefined,
      bccEmail: form.bcc ? parseEmailAddresses(form.bcc) : undefined,
      content: form.content,
      contentType: form.contentType
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
