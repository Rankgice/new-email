<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h2 class="text-xl font-semibold text-text-primary mb-2">个人资料</h2>
      <p class="text-text-secondary">管理您的个人信息和偏好设置</p>
    </div>

    <!-- 头像设置 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">头像</h3>
      <div class="flex items-center space-x-6">
        <!-- 当前头像 -->
        <div class="relative">
          <div
            v-if="form.avatar"
            class="w-20 h-20 rounded-full overflow-hidden bg-gray-200"
          >
            <img
              :src="form.avatar"
              :alt="form.nickname || '用户头像'"
              class="w-full h-full object-cover"
            />
          </div>
          <div
            v-else
            class="w-20 h-20 rounded-full bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center text-white text-xl font-semibold"
          >
            {{ userInitials }}
          </div>
          
          <!-- 上传按钮 -->
          <label
            for="avatar-upload"
            class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 rounded-full opacity-0 hover:opacity-100 transition-opacity cursor-pointer"
          >
            <CameraIcon class="w-6 h-6 text-white" />
          </label>
          <input
            id="avatar-upload"
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleAvatarUpload"
          />
        </div>

        <!-- 头像操作 -->
        <div class="flex-1">
          <p class="text-sm text-text-secondary mb-2">
            点击头像上传新图片，支持 JPG、PNG 格式，建议尺寸 200x200 像素
          </p>
          <div class="flex space-x-2">
            <Button
              variant="secondary"
              size="sm"
              @click="triggerAvatarUpload"
            >
              <PhotoIcon class="w-4 h-4 mr-2" />
              上传头像
            </Button>
            <Button
              v-if="form.avatar"
              variant="ghost"
              size="sm"
              @click="removeAvatar"
            >
              <TrashIcon class="w-4 h-4 mr-2" />
              移除
            </Button>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 基本信息 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">基本信息</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- 用户名 -->
        <Input
          v-model="form.username"
          label="用户名"
          :left-icon="UserIcon"
          disabled
          readonly
          help="用户名不可修改"
        />

        <!-- 邮箱 -->
        <Input
          v-model="form.email"
          label="邮箱地址"
          type="email"
          :left-icon="EnvelopeIcon"
          disabled
          readonly
          help="邮箱地址不可修改"
        />

        <!-- 昵称 -->
        <Input
          v-model="form.nickname"
          label="昵称"
          placeholder="请输入昵称"
          :left-icon="UserIcon"
          :error="errors.nickname"
        />

        <!-- 语言 -->
        <Select
          v-model="form.language"
          label="语言"
          :options="languageOptions"
          :left-icon="LanguageIcon"
        />

        <!-- 时区 -->
        <Select
          v-model="form.timezone"
          label="时区"
          :options="timezoneOptions"
          :left-icon="ClockIcon"
        />
      </div>

      <!-- 个人简介 -->
      <div class="mt-6">
        <Textarea
          v-model="form.bio"
          label="个人简介"
          placeholder="介绍一下自己..."
          :rows="4"
          :max-length="500"
          show-count
          :error="errors.bio"
        />
      </div>
    </GlassCard>

    <!-- 保存按钮 -->
    <div class="flex justify-end space-x-3">
      <Button
        variant="ghost"
        @click="resetForm"
        :disabled="loading"
      >
        重置
      </Button>
      <Button
        variant="primary"
        @click="saveProfile"
        :loading="loading"
      >
        保存更改
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userSettingsApi } from '@/api/user-settings'
import { useNotification } from '@/composables/useNotification'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Textarea from '@/components/ui/Textarea.vue'

import {
  UserIcon,
  EnvelopeIcon,
  CameraIcon,
  PhotoIcon,
  TrashIcon,
  LanguageIcon,
  ClockIcon
} from '@heroicons/vue/24/outline'

// 状态管理
const authStore = useAuthStore()
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const errors = reactive({
  nickname: '',
  bio: ''
})

// 表单数据
const form = reactive({
  username: '',
  email: '',
  nickname: '',
  bio: '',
  language: 'zh-CN',
  timezone: 'Asia/Shanghai',
  avatar: ''
})

// 计算属性
const userInitials = computed(() => {
  const name = form.nickname || form.username
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
})

// 语言选项
const languageOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' },
  { label: '繁體中文', value: 'zh-TW' },
  { label: '日本語', value: 'ja-JP' }
]

// 时区选项
const timezoneOptions = [
  { label: '北京时间 (UTC+8)', value: 'Asia/Shanghai' },
  { label: '东京时间 (UTC+9)', value: 'Asia/Tokyo' },
  { label: '纽约时间 (UTC-5)', value: 'America/New_York' },
  { label: '伦敦时间 (UTC+0)', value: 'Europe/London' },
  { label: '洛杉矶时间 (UTC-8)', value: 'America/Los_Angeles' }
]

// 生命周期
onMounted(() => {
  loadUserProfile()
})

// 方法
const loadUserProfile = () => {
  if (authStore.user) {
    form.username = authStore.user.username
    form.email = authStore.user.email
    form.nickname = authStore.user.nickname || ''
    form.avatar = authStore.user.avatar || ''
    // 从用户设置中加载其他信息
    if (authStore.user.settings) {
      form.language = authStore.user.settings.language || 'zh-CN'
    }
  }
}

const triggerAvatarUpload = () => {
  document.getElementById('avatar-upload')?.click()
}

const handleAvatarUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    showError('请选择图片文件')
    return
  }

  // 验证文件大小 (5MB)
  if (file.size > 5 * 1024 * 1024) {
    showError('图片大小不能超过 5MB')
    return
  }

  try {
    loading.value = true
    const result = await userSettingsApi.uploadAvatar(file)
    form.avatar = result.url
    showSuccess('头像上传成功')
  } catch (error) {
    showError('头像上传失败')
  } finally {
    loading.value = false
  }
}

const removeAvatar = () => {
  form.avatar = ''
}

const validateForm = () => {
  // 清空错误
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  // 验证昵称
  if (form.nickname && form.nickname.length > 50) {
    errors.nickname = '昵称不能超过50个字符'
    isValid = false
  }

  // 验证个人简介
  if (form.bio && form.bio.length > 500) {
    errors.bio = '个人简介不能超过500个字符'
    isValid = false
  }

  return isValid
}

const saveProfile = async () => {
  if (!validateForm()) return

  try {
    loading.value = true
    
    await userSettingsApi.updateProfile({
      nickname: form.nickname,
      bio: form.bio,
      timezone: form.timezone,
      language: form.language,
      avatar: form.avatar
    })

    // 更新本地用户信息
    await authStore.updateUser({
      ...authStore.user!,
      nickname: form.nickname,
      avatar: form.avatar
    })

    showSuccess('个人资料保存成功')
  } catch (error) {
    showError('保存失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  loadUserProfile()
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })
}
</script>
