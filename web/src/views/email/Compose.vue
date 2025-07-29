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
          <!-- 收件人 -->
          <Input
            v-model="form.to"
            label="收件人"
            placeholder="输入收件人邮箱地址"
            :left-icon="UserIcon"
          />

          <!-- 抄送 -->
          <Input
            v-model="form.cc"
            label="抄送"
            placeholder="输入抄送邮箱地址（可选）"
            :left-icon="UserGroupIcon"
          />

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
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useNotification } from '@/composables/useNotification'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import {
  ArrowLeftIcon,
  UserIcon,
  UserGroupIcon,
  ChatBubbleLeftRightIcon,
  PaperClipIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const { showNotification } = useNotification()

const form = reactive({
  to: '',
  cc: '',
  subject: '',
  content: ''
})

const isSending = ref(false)
const isSavingDraft = ref(false)

const sendEmail = async () => {
  isSending.value = true
  // TODO: 实现发送邮件逻辑
  setTimeout(() => {
    isSending.value = false
    showNotification({
      type: 'success',
      title: '邮件已发送',
      message: '您的邮件已成功发送'
    })
    router.push('/inbox')
  }, 2000)
}

const saveDraft = async () => {
  isSavingDraft.value = true
  // TODO: 实现保存草稿逻辑
  setTimeout(() => {
    isSavingDraft.value = false
    showNotification({
      type: 'success',
      title: '草稿已保存'
    })
  }, 1000)
}
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
