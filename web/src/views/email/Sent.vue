<template>
  <div class="h-screen flex bg-background-primary">
    <!-- 侧边栏 -->
    <EmailSidebar />

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部工具栏 -->
      <EmailToolbar />

      <!-- 已发送邮件内容 -->
      <div class="flex-1 overflow-hidden p-4">
        <div class="h-full flex flex-col">
          <!-- 页面标题 -->
          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-3">
              <PaperAirplaneIcon class="w-6 h-6 text-text-primary" />
              <h1 class="text-2xl font-bold text-text-primary">已发送</h1>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center space-x-3">
              <Button
                variant="ghost"
                size="sm"
                @click="refreshMailboxes"
                :loading="loading"
              >
                <ArrowPathIcon class="w-4 h-4 mr-2" />
                刷新
              </Button>
            </div>
          </div>

          <!-- 邮箱列表 -->
          <div class="flex-1 overflow-hidden">
            <GlassCard padding="none" class="h-full">
              <!-- 空状态 -->
              <div v-if="!loading && (!mailboxes || mailboxes.length === 0)" class="h-full flex items-center justify-center">
                <div class="text-center">
                  <PaperAirplaneIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
                  <h3 class="text-lg font-medium text-text-primary mb-2">
                    暂无邮箱
                  </h3>
                  <p class="text-text-secondary mb-4">
                    请先添加邮箱账户以发送邮件
                  </p>
                  <Button
                    variant="primary"
                    @click="$router.push('/settings')"
                  >
                    添加邮箱
                  </Button>
                </div>
              </div>

              <!-- 邮箱列表 -->
              <div v-else class="h-full flex flex-col">
                <!-- 列表头部 -->
                <div class="flex items-center px-4 py-3 border-b border-glass-border bg-white/5">
                  <div class="flex-1 grid grid-cols-12 gap-4 text-sm font-medium text-text-secondary">
                    <div class="col-span-1">状态</div>
                    <div class="col-span-6">邮箱地址</div>
                    <div class="col-span-3">邮箱类型</div>
                    <div class="col-span-2">操作</div>
                  </div>
                </div>

                <!-- 邮箱项目 -->
                <div class="flex-1 overflow-y-auto">
                  <div
                    v-for="mailbox in (mailboxes || [])"
                    :key="mailbox?.id || Math.random()"
                    :class="[
                      'flex items-center px-4 py-3 border-b border-glass-border hover:bg-white/5 cursor-pointer transition-colors',
                      selectedMailbox?.id === mailbox?.id ? 'bg-primary-500/10' : ''
                    ]"
                    @click="mailbox && selectMailbox(mailbox)"
                  >
                    <div class="flex-1 grid grid-cols-12 gap-4 items-center">
                      <!-- 状态 -->
                      <div class="col-span-1">
                        <div
                          :class="[
                            'w-3 h-3 rounded-full',
                            mailbox?.isActive ? 'bg-green-400' : 'bg-gray-400'
                          ]"
                        />
                      </div>

                      <!-- 邮箱地址 -->
                      <div class="col-span-6">
                        <div class="text-text-primary font-medium">
                          {{ mailbox?.email || '未知邮箱' }}
                        </div>
                        <div class="text-xs text-text-secondary">
                          {{ mailbox?.displayName || mailbox?.email }}
                        </div>
                      </div>

                      <!-- 邮箱类型 -->
                      <div class="col-span-3">
                        <span
                          :class="[
                            'px-2 py-1 rounded text-xs font-medium',
                            getMailboxTypeClass(mailbox?.provider)
                          ]"
                        >
                          {{ getMailboxTypeLabel(mailbox?.provider) }}
                        </span>
                      </div>

                      <!-- 操作 -->
                      <div class="col-span-2 flex items-center space-x-2">
                        <ChevronDownIcon
                          :class="[
                            'w-5 h-5 text-text-secondary transition-transform duration-200 cursor-pointer',
                            selectedMailbox?.id === mailbox?.id ? 'rotate-180' : ''
                          ]"
                        />
                      </div>
                    </div>

                    <!-- 已发送邮件列表 (展开状态) -->
                    <div
                      v-if="selectedMailbox?.id === mailbox?.id"
                      class="border-t border-glass-border bg-white/5"
                    >
                      <SentEmailList
                        :mailbox-id="mailbox.id"
                        @email-selected="handleEmailSelected"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </GlassCard>
          </div>
        </div>
      </div>
    </div>

    <!-- 邮件预览面板 (桌面端) -->
    <EmailPreview v-if="!isMobile && selectedEmail" :email="selectedEmail" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'
import { mailboxApi } from '@/utils/api'
import type { Mailbox, Email } from '@/types'

import EmailSidebar from '@/components/email/EmailSidebar.vue'
import EmailToolbar from '@/components/email/EmailToolbar.vue'
import EmailPreview from '@/components/email/EmailPreview.vue'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import SentEmailList from '@/components/email/SentEmailList.vue'

import {
  PaperAirplaneIcon,
  ArrowPathIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/outline'

// Composables
const { isMobile } = useResponsive()
const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()

// 响应式数据
const loading = ref(false)
const mailboxes = ref<Mailbox[]>([])
const selectedMailbox = ref<Mailbox | null>(null)
const selectedEmail = ref<Email | null>(null)

// 生命周期
onMounted(() => {
  loadMailboxes()
})

// 方法
const loadMailboxes = async () => {
  try {
    loading.value = true
    const response = await mailboxApi.getMailboxes()
    mailboxes.value = response.data || []
  } catch (error) {
    showError('加载邮箱列表失败')
  } finally {
    loading.value = false
  }
}

const refreshMailboxes = () => {
  loadMailboxes()
}

const selectMailbox = (mailbox: Mailbox) => {
  if (selectedMailbox.value?.id === mailbox.id) {
    // 如果点击的是已选中的邮箱，则收起
    selectedMailbox.value = null
    selectedEmail.value = null
  } else {
    // 否则展开新的邮箱
    selectedMailbox.value = mailbox
    selectedEmail.value = null
  }
}

const handleEmailSelected = (email: Email) => {
  selectedEmail.value = email
}

const getMailboxTypeClass = (provider?: string) => {
  switch (provider?.toLowerCase()) {
    case 'gmail':
      return 'bg-red-500/20 text-red-400'
    case 'outlook':
    case 'hotmail':
      return 'bg-blue-500/20 text-blue-400'
    case 'yahoo':
      return 'bg-purple-500/20 text-purple-400'
    case 'qq':
      return 'bg-green-500/20 text-green-400'
    case '163':
    case '126':
      return 'bg-orange-500/20 text-orange-400'
    default:
      return 'bg-gray-500/20 text-gray-400'
  }
}

const getMailboxTypeLabel = (provider?: string) => {
  switch (provider?.toLowerCase()) {
    case 'gmail':
      return 'Gmail'
    case 'outlook':
      return 'Outlook'
    case 'hotmail':
      return 'Hotmail'
    case 'yahoo':
      return 'Yahoo'
    case 'qq':
      return 'QQ邮箱'
    case '163':
      return '163邮箱'
    case '126':
      return '126邮箱'
    default:
      return provider || '其他'
  }
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
</script>
