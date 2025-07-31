<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
    <div class="max-w-md w-full space-y-8 p-8">
      <!-- Logo和标题 -->
      <div class="text-center">
        <div class="mx-auto h-16 w-16 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-full flex items-center justify-center">
          <i class="fas fa-shield-alt text-white text-2xl"></i>
        </div>
        <h2 class="mt-6 text-3xl font-bold text-gray-900 dark:text-white">
          管理员登录
        </h2>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          请使用管理员账户登录系统
        </p>
      </div>

      <!-- 登录表单 -->
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="space-y-4">
          <!-- 用户名 -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              用户名
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <i class="fas fa-user text-gray-400"></i>
              </div>
              <input
                id="username"
                v-model="form.username"
                type="text"
                required
                class="appearance-none relative block w-full pl-10 pr-3 py-3 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                placeholder="请输入用户名"
                :disabled="loading"
              />
            </div>
          </div>

          <!-- 密码 -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              密码
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <i class="fas fa-lock text-gray-400"></i>
              </div>
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                class="appearance-none relative block w-full pl-10 pr-10 py-3 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white bg-white dark:bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
                placeholder="请输入密码"
                :disabled="loading"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
                @click="showPassword = !showPassword"
              >
                <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'" class="text-gray-400 hover:text-gray-600"></i>
              </button>
            </div>
          </div>
        </div>

        <!-- 记住我 -->
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              v-model="form.rememberMe"
              type="checkbox"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
              记住我
            </label>
          </div>
        </div>

        <!-- 登录按钮 -->
        <div>
          <button
            type="submit"
            :disabled="loading || !form.username || !form.password"
            class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="fas fa-spinner fa-spin text-white"></i>
            </span>
            <span v-else class="absolute left-0 inset-y-0 flex items-center pl-3">
              <i class="fas fa-sign-in-alt text-white"></i>
            </span>
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>

        <!-- 错误信息 -->
        <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <i class="fas fa-exclamation-circle text-red-400"></i>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800 dark:text-red-200">
                登录失败
              </h3>
              <div class="mt-2 text-sm text-red-700 dark:text-red-300">
                {{ error }}
              </div>
            </div>
          </div>
        </div>
      </form>

      <!-- 返回用户登录 -->
      <div class="text-center">
        <router-link
          to="/login"
          class="text-sm text-blue-600 hover:text-blue-500 dark:text-blue-400 dark:hover:text-blue-300"
        >
          <i class="fas fa-arrow-left mr-1"></i>
          返回用户登录
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { adminApi } from '@/api/admin'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'

// 路由
const router = useRouter()

// 状态管理
const authStore = useAuthStore()

// 通知
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const showPassword = ref(false)
const error = ref('')

// 表单数据
const form = reactive({
  username: '',
  password: '',
  rememberMe: false
})

// 方法
async function handleLogin() {
  if (!form.username || !form.password) {
    error.value = '请输入用户名和密码'
    return
  }

  try {
    loading.value = true
    error.value = ''

    // 调用登录API
    const response = await adminApi.login({
      username: form.username,
      password: form.password
    })

    // 保存登录信息
    authStore.setAdminToken(response.token, response.admin)

    // 如果选择记住我，保存到localStorage
    if (form.rememberMe) {
      localStorage.setItem('admin_remember', 'true')
      localStorage.setItem('admin_username', form.username)
    } else {
      localStorage.removeItem('admin_remember')
      localStorage.removeItem('admin_username')
    }

    showSuccess('登录成功')

    // 跳转到管理员仪表板
    router.push('/admin/dashboard')
  } catch (err: any) {
    console.error('Admin login error:', err)
    error.value = err.response?.data?.msg || '登录失败，请检查用户名和密码'
    showError(error.value)
  } finally {
    loading.value = false
  }
}

// 初始化
function init() {
  // 如果之前选择了记住我，自动填充用户名
  const remember = localStorage.getItem('admin_remember')
  const savedUsername = localStorage.getItem('admin_username')
  
  if (remember === 'true' && savedUsername) {
    form.username = savedUsername
    form.rememberMe = true
  }

  // 如果已经登录，直接跳转到仪表板
  if (authStore.isAdminLoggedIn) {
    router.push('/admin/dashboard')
  }
}

// 页面加载时初始化
init()
</script>

<style scoped>
/* 自定义样式 */
.gradient-bg {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 输入框聚焦效果 */
input:focus {
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

/* 按钮悬停效果 */
button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

/* 动画效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .max-w-md {
    max-width: 90%;
  }
}
</style>
