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
                <!-- 邮箱项目 -->
                <div class="flex-1 overflow-y-auto">
                  <div
                    v-for="mailbox in (mailboxes || [])"
                    :key="mailbox?.id || Math.random()"
                    :class="[
                      'border-b border-glass-border transition-colors',
                      selectedMailbox?.id === mailbox?.id ? 'bg-primary-500/10' : 'hover:bg-white/5'
                    ]"
                  >
                    <!-- 邮箱信息 -->
                    <div
                      class="flex items-center justify-between p-4 cursor-pointer hover:bg-white/5 transition-colors"
                      @click="handleMailboxClick(mailbox)"
                    >
                      <div class="flex items-center space-x-3">
                        <!-- 邮箱图标 -->
                        <div class="flex items-center justify-center w-10 h-10 bg-primary-500/20 rounded-lg">
                          <PaperAirplaneIcon class="w-5 h-5 text-primary-400" />
                        </div>

                        <!-- 邮箱信息 -->
                        <div>
                          <h3 class="font-medium text-text-primary">
                            {{ mailbox?.email || '未知邮箱' }}
                          </h3>
                          <p class="text-sm text-text-secondary">
                            {{ mailbox?.name || mailbox?.email }}
                          </p>
                        </div>
                      </div>

                      <div class="flex items-center space-x-3">
                        <!-- 展开/收起图标 -->
                        <ChevronDownIcon
                          :class="[
                            'w-5 h-5 text-text-secondary transition-transform duration-200',
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
    console.log('Sent page - Mailboxes API response:', response)

    // 处理不同的响应格式
    if (response.data) {
      if (Array.isArray(response.data)) {
        mailboxes.value = response.data
      } else if (response.data.list && Array.isArray(response.data.list)) {
        mailboxes.value = response.data.list
      } else if (response.data.data && Array.isArray(response.data.data)) {
        mailboxes.value = response.data.data
      } else {
        console.warn('Unexpected response format:', response.data)
        mailboxes.value = []
      }
    } else {
      mailboxes.value = []
    }

    // 如果没有数据，添加一些模拟数据用于测试
    if (mailboxes.value.length === 0) {
      console.log('Using mock mailbox data for testing in sent page')
      mailboxes.value = [
        {
          id: 1,
          userId: 1,
          domainId: 1,
          email: 'test1@example.com',
          name: '测试邮箱1',
          autoReceive: true,
          status: 1,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        },
        {
          id: 2,
          userId: 1,
          domainId: 1,
          email: 'test2@gmail.com',
          name: '测试邮箱2',
          autoReceive: true,
          status: 1,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
      ]
    }

    console.log('Final mailboxes for sent page:', mailboxes.value)
  } catch (error) {
    console.error('Failed to load mailboxes for sent page:', error)
    showError('加载邮箱列表失败')

    // 在错误情况下也提供模拟数据
    mailboxes.value = [
      {
        id: 1,
        userId: 1,
        domainId: 1,
        email: 'test1@example.com',
        name: '测试邮箱1',
        autoReceive: true,
        status: 1,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
    ]
  } finally {
    loading.value = false
  }
}

const refreshMailboxes = () => {
  loadMailboxes()
}

const handleMailboxClick = (mailbox: Mailbox) => {
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


</script>
