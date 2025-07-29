<template>
  <div class="min-h-screen flex items-center justify-center relative overflow-hidden">
    <!-- 动态背景 -->
    <div class="absolute inset-0 bg-gradient-to-br from-primary-900/20 via-background-primary to-secondary-900/20">
      <!-- 浮动粒子效果 -->
      <div class="absolute inset-0">
        <div
          v-for="i in 20"
          :key="i"
          class="absolute w-2 h-2 bg-primary-500/20 rounded-full animate-float"
          :style="{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 3}s`,
            animationDuration: `${3 + Math.random() * 2}s`
          }"
        />
      </div>
    </div>

    <!-- 注册卡片 -->
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
              <UserPlusIcon class="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 class="text-2xl font-bold text-text-primary mb-2">
            创建账户
          </h1>
          <p class="text-text-secondary">
            加入我们，开始您的邮件之旅
          </p>
        </div>

        <!-- 注册表单 -->
        <form @submit.prevent="handleRegister" class="space-y-6">
          <!-- 用户名输入 -->
          <Input
            v-model="form.username"
            type="text"
            label="用户名"
            placeholder="请输入用户名"
            :left-icon="UserIcon"
            :error="errors.username"
            required
            autocomplete="username"
          />

          <!-- 昵称输入 -->
          <Input
            v-model="form.nickname"
            type="text"
            label="昵称"
            placeholder="请输入昵称（可选）"
            :left-icon="IdentificationIcon"
            :error="errors.nickname"
            autocomplete="name"
          />

          <!-- 邮箱输入 -->
          <Input
            v-model="form.email"
            type="email"
            label="邮箱地址"
            placeholder="请输入您的邮箱"
            :left-icon="AtSymbolIcon"
            :error="errors.email"
            required
            autocomplete="email"
          />

          <!-- 密码输入 -->
          <Input
            v-model="form.password"
            type="password"
            label="密码"
            placeholder="请输入密码（至少6位）"
            :left-icon="LockClosedIcon"
            :error="errors.password"
            required
            autocomplete="new-password"
          />

          <!-- 确认密码输入 -->
          <Input
            v-model="form.confirmPassword"
            type="password"
            label="确认密码"
            placeholder="请再次输入密码"
            :left-icon="LockClosedIcon"
            :error="errors.confirmPassword"
            required
            autocomplete="new-password"
          />

          <!-- 服务条款 -->
          <div class="flex items-start">
            <input
              v-model="form.agreeToTerms"
              type="checkbox"
              class="mt-1 w-4 h-4 text-primary-600 bg-glass-light border-glass-border rounded focus:ring-primary-500 focus:ring-2"
            />
            <label class="ml-2 text-sm text-text-secondary">
              我已阅读并同意
              <a href="#" class="text-primary-400 hover:text-primary-300">服务条款</a>
              和
              <a href="#" class="text-primary-400 hover:text-primary-300">隐私政策</a>
            </label>
          </div>

          <!-- 注册按钮 -->
          <Button
            type="submit"
            variant="primary"
            size="lg"
            :loading="isLoading"
            :disabled="!form.agreeToTerms"
            class="w-full"
          >
            <span v-if="!isLoading">🚀 创建账户</span>
            <span v-else>创建中...</span>
          </Button>
        </form>

        <!-- 分割线 -->
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-glass-border" />
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-background-primary text-text-secondary">或</span>
          </div>
        </div>

        <!-- 登录链接 -->
        <div class="text-center">
          <p class="text-text-secondary">
            已有账户？
            <router-link
              to="/auth/login"
              class="text-primary-400 hover:text-primary-300 font-medium transition-colors duration-200"
            >
              立即登录
            </router-link>
          </p>
        </div>
      </GlassCard>

      <!-- 主题切换器 -->
      <div class="mt-6 flex justify-center">
        <ThemeSelector />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import ThemeSelector from '@/components/ui/ThemeSelector.vue'
import {
  UserPlusIcon,
  UserIcon,
  IdentificationIcon,
  AtSymbolIcon,
  LockClosedIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const authStore = useAuthStore()
const { showNotification } = useNotification()

// 表单数据
const form = reactive({
  username: '',
  nickname: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeToTerms: false
})

// 表单错误
const errors = reactive({
  username: '',
  nickname: '',
  email: '',
  password: '',
  confirmPassword: ''
})

// 加载状态
const isLoading = ref(false)

// 验证表单
const validateForm = () => {
  // 重置错误
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  // 验证用户名
  if (!form.username) {
    errors.username = '请输入用户名'
    isValid = false
  } else if (form.username.length < 3) {
    errors.username = '用户名长度至少3位'
    isValid = false
  }

  // 验证邮箱
  if (!form.email) {
    errors.email = '请输入邮箱地址'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = '请输入有效的邮箱地址'
    isValid = false
  }

  // 验证密码
  if (!form.password) {
    errors.password = '请输入密码'
    isValid = false
  } else if (form.password.length < 6) {
    errors.password = '密码长度至少6位'
    isValid = false
  }

  // 验证确认密码
  if (!form.confirmPassword) {
    errors.confirmPassword = '请确认密码'
    isValid = false
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = '两次输入的密码不一致'
    isValid = false
  }

  return isValid
}

// 处理注册
const handleRegister = async () => {
  if (!validateForm()) return

  isLoading.value = true

  try {
    const result = await authStore.register({
      username: form.username,
      email: form.email,
      password: form.password,
      nickname: form.nickname || undefined
    })

    if (result.success) {
      showNotification({
        type: 'success',
        title: '注册成功',
        message: '账户创建成功，请登录'
      })

      // 跳转到登录页面
      router.push('/auth/login')
    } else {
      showNotification({
        type: 'error',
        title: '注册失败',
        message: result.message || '注册失败，请稍后重试'
      })
    }
  } catch (error) {
    console.error('Register error:', error)
    showNotification({
      type: 'error',
      title: '注册失败',
      message: '网络错误，请稍后重试'
    })
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* 浮动动画 */
@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.5;
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 1;
  }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}
</style>
