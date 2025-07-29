<template>
  <div class="min-h-screen flex items-center justify-center relative overflow-hidden">
    <!-- 动态背景 -->
    <div class="absolute inset-0 bg-gradient-to-br from-primary-900/20 via-background-primary to-secondary-900/20 pointer-events-none z-0">
      <!-- 浮动粒子效果 -->
      <div class="absolute inset-0 pointer-events-none z-0">
        <div
          v-for="i in 20"
          :key="i"
          class="absolute w-2 h-2 bg-primary-500/20 rounded-full animate-float pointer-events-none"
          :style="{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 3}s`,
            animationDuration: `${3 + Math.random() * 2}s`
          }"
        />
      </div>
    </div>

    <!-- 登录卡片 -->
    <div class="relative z-50 w-full max-w-md px-6">
      <GlassCard
        :level="3"
        :hover="false"
        padding="lg"
        border
        class="animate-scale-in relative z-50"
      >
        <!-- 标题区域 -->
        <div class="text-center mb-8">
          <div class="flex justify-center mb-4">
            <div class="w-16 h-16 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-2xl flex items-center justify-center">
              <EnvelopeIcon class="w-8 h-8 text-white" />
            </div>
          </div>
          <h1 class="text-2xl font-bold text-text-primary mb-2">
            📧 邮件系统
          </h1>
          <p class="text-text-secondary">
            欢迎回来，请登录您的账户
          </p>
        </div>

        <!-- 登录表单 -->
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- 用户名输入 -->
          <Input
            v-model="form.username"
            type="text"
            label="用户名"
            placeholder="请输入您的用户名或邮箱"
            :left-icon="AtSymbolIcon"
            :error="errors.username"
            required
            autocomplete="username"
          />

          <!-- 密码输入 -->
          <Input
            v-model="form.password"
            type="password"
            label="密码"
            placeholder="请输入您的密码"
            :left-icon="LockClosedIcon"
            :error="errors.password"
            required
            autocomplete="current-password"
          />

          <!-- 记住我 -->
          <div class="flex items-center justify-between">
            <label class="flex items-center">
              <input
                v-model="form.rememberMe"
                type="checkbox"
                class="w-4 h-4 text-primary-600 bg-glass-light border-glass-border rounded focus:ring-primary-500 focus:ring-2"
              />
              <span class="ml-2 text-sm text-text-secondary">记住我</span>
            </label>

            <router-link
              to="/auth/forgot-password"
              class="text-sm text-primary-400 hover:text-primary-300 transition-colors duration-200"
            >
              忘记密码？
            </router-link>
          </div>

          <!-- 登录按钮 -->
          <Button
            type="submit"
            variant="primary"
            size="lg"
            :loading="isLoading"
            class="w-full"
          >
            <span v-if="!isLoading">🚀 登录</span>
            <span v-else>登录中...</span>
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

        <!-- 注册链接 -->
        <div class="text-center relative z-50">
          <p class="text-text-secondary">
            还没有账户？
            <button
              type="button"
              class="text-primary-400 hover:text-primary-300 font-medium transition-colors duration-200 cursor-pointer bg-transparent border-none underline relative z-50"
              @click="handleRegisterClick"
              style="pointer-events: auto !important; position: relative; z-index: 9999;"
            >
              立即注册
            </button>
          </p>
        </div>
      </GlassCard>

      <!-- 主题切换器 -->
      <div class="mt-6 flex justify-center relative z-50">
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
  EnvelopeIcon,
  AtSymbolIcon,
  LockClosedIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()
const authStore = useAuthStore()
const { showNotification } = useNotification()

// 表单数据
const form = reactive({
  username: '',
  password: '',
  rememberMe: false
})

// 表单错误
const errors = reactive({
  username: '',
  password: ''
})

// 加载状态
const isLoading = ref(false)

// 验证表单
const validateForm = () => {
  errors.username = ''
  errors.password = ''

  if (!form.username) {
    errors.username = '请输入用户名或邮箱'
    return false
  }

  if (form.username.length < 3) {
    errors.username = '用户名长度至少3位'
    return false
  }

  if (!form.password) {
    errors.password = '请输入密码'
    return false
  }

  if (form.password.length < 6) {
    errors.password = '密码长度至少6位'
    return false
  }

  return true
}

// 处理登录
const handleLogin = async () => {
  if (!validateForm()) return

  isLoading.value = true

  try {
    const result = await authStore.login({
      username: form.username,
      password: form.password
    })

    if (result.success) {
      showNotification({
        type: 'success',
        title: '登录成功',
        message: '欢迎回来！'
      })

      // 跳转到目标页面
      const redirectPath = authStore.getAndClearRedirectPath()
      router.push(redirectPath)
    } else {
      showNotification({
        type: 'error',
        title: '登录失败',
        message: result.message || '用户名或密码错误'
      })
    }
  } catch (error) {
    console.error('Login error:', error)
    showNotification({
      type: 'error',
      title: '登录失败',
      message: '网络错误，请稍后重试'
    })
  } finally {
    isLoading.value = false
  }
}

// 处理注册链接点击
const handleRegisterClick = () => {
  console.log('注册链接被点击')
  router.push('/auth/register')
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
  pointer-events: none; /* 确保浮动粒子不会阻止点击事件 */
}
</style>
