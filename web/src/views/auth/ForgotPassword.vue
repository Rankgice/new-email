<template>
  <div class="min-h-screen flex items-center justify-center relative overflow-hidden">
    <!-- 动态背景 -->
    <div class="absolute inset-0 bg-gradient-to-br from-primary-900/20 via-background-primary to-secondary-900/20" />

    <!-- 忘记密码卡片 -->
    <div class="relative z-10 w-full max-w-md px-6">
      <GlassCard
        :level="3"
        :hover="false"
        padding="lg"
        border
        class="animate-scale-in"
      >
        <!-- 标题区域 -->
        <div class="text-center mb-8">
          <div class="flex justify-center mb-4">
            <div class="w-16 h-16 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-2xl flex items-center justify-center">
              <KeyIcon class="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 class="text-2xl font-bold text-text-primary mb-2">
            忘记密码
          </h1>
          <p class="text-text-secondary">
            输入您的邮箱地址，我们将发送重置密码的链接
          </p>
        </div>

        <!-- 表单 -->
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- 邮箱输入 -->
          <Input
            v-model="email"
            type="email"
            label="邮箱地址"
            placeholder="请输入您的邮箱"
            :left-icon="AtSymbolIcon"
            :error="error"
            required
            autocomplete="email"
          />

          <!-- 提交按钮 -->
          <Button
            type="submit"
            variant="primary"
            size="lg"
            :loading="isLoading"
            class="w-full"
          >
            <span v-if="!isLoading">📧 发送重置链接</span>
            <span v-else>发送中...</span>
          </Button>
        </form>

        <!-- 返回登录 -->
        <div class="mt-6 text-center">
          <router-link
            to="/auth/login"
            class="text-primary-400 hover:text-primary-300 font-medium transition-colors duration-200"
          >
            ← 返回登录
          </router-link>
        </div>
      </GlassCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import {
  KeyIcon,
  AtSymbolIcon
} from '@heroicons/vue/24/outline'

const authStore = useAuthStore()
const { showNotification } = useNotification()

const email = ref('')
const error = ref('')
const isLoading = ref(false)

const handleSubmit = async () => {
  error.value = ''

  if (!email.value) {
    error.value = '请输入邮箱地址'
    return
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    error.value = '请输入有效的邮箱地址'
    return
  }

  isLoading.value = true

  try {
    const result = await authStore.forgotPassword(email.value)

    if (result.success) {
      showNotification({
        type: 'success',
        title: '发送成功',
        message: '重置密码的邮件已发送到您的邮箱'
      })
    } else {
      showNotification({
        type: 'error',
        title: '发送失败',
        message: result.message || '发送失败，请稍后重试'
      })
    }
  } catch (error) {
    console.error('Forgot password error:', error)
    showNotification({
      type: 'error',
      title: '发送失败',
      message: '网络错误，请稍后重试'
    })
  } finally {
    isLoading.value = false
  }
}
</script>
