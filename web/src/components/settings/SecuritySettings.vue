<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h2 class="text-xl font-semibold text-text-primary mb-2">安全设置</h2>
      <p class="text-text-secondary">管理您的账户安全和隐私设置</p>
    </div>

    <!-- 密码设置 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">密码管理</h3>
      <div class="space-y-4">
        <!-- 当前密码 -->
        <Input
          v-model="passwordForm.currentPassword"
          label="当前密码"
          type="password"
          placeholder="请输入当前密码"
          :left-icon="LockClosedIcon"
          :error="passwordErrors.currentPassword"
          required
        />

        <!-- 新密码 -->
        <Input
          v-model="passwordForm.newPassword"
          label="新密码"
          type="password"
          placeholder="请输入新密码"
          :left-icon="KeyIcon"
          :error="passwordErrors.newPassword"
          help="密码至少8位，包含字母、数字和特殊字符"
        />

        <!-- 确认新密码 -->
        <Input
          v-model="passwordForm.confirmPassword"
          label="确认新密码"
          type="password"
          placeholder="请再次输入新密码"
          :left-icon="KeyIcon"
          :error="passwordErrors.confirmPassword"
        />

        <!-- 保存按钮 -->
        <div class="flex justify-end">
          <Button
            variant="primary"
            @click="changePassword"
            :loading="passwordLoading"
            :disabled="!canChangePassword"
          >
            更改密码
          </Button>
        </div>
      </div>
    </GlassCard>

    <!-- 两步验证 -->
    <GlassCard padding="lg">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h3 class="text-lg font-medium text-text-primary">两步验证</h3>
          <p class="text-sm text-text-secondary mt-1">
            为您的账户添加额外的安全保护
          </p>
        </div>
        <Switch
          v-model="securitySettings.twoFactorEnabled"
          @change="handleTwoFactorToggle"
        />
      </div>

      <!-- 两步验证设置 -->
      <div v-if="showTwoFactorSetup" class="mt-4 p-4 bg-blue-500/10 rounded-lg border border-blue-500/20">
        <h4 class="text-md font-medium text-text-primary mb-3">设置两步验证</h4>
        
        <!-- 二维码 -->
        <div v-if="twoFactorSetup.qrCode" class="text-center mb-4">
          <div class="inline-block p-4 bg-white rounded-lg">
            <img :src="twoFactorSetup.qrCode" alt="两步验证二维码" class="w-48 h-48" />
          </div>
          <p class="text-sm text-text-secondary mt-2">
            使用 Google Authenticator 或其他验证器应用扫描二维码
          </p>
        </div>

        <!-- 验证码输入 -->
        <div class="max-w-xs mx-auto">
          <Input
            v-model="twoFactorToken"
            label="验证码"
            placeholder="请输入6位验证码"
            :error="twoFactorError"
            maxlength="6"
          />
          <div class="flex space-x-2 mt-4">
            <Button
              variant="ghost"
              @click="cancelTwoFactorSetup"
              class="flex-1"
            >
              取消
            </Button>
            <Button
              variant="primary"
              @click="enableTwoFactor"
              :loading="twoFactorLoading"
              class="flex-1"
            >
              启用
            </Button>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 会话管理 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">会话管理</h3>
      
      <!-- 会话超时 -->
      <div class="mb-6">
        <label class="block text-sm font-medium text-text-primary mb-2">
          会话超时时间
        </label>
        <Select
          v-model="securitySettings.sessionTimeout"
          :options="sessionTimeoutOptions"
          @change="updateSessionTimeout"
        />
        <p class="text-sm text-text-secondary mt-1">
          超过指定时间未活动将自动退出登录
        </p>
      </div>

      <!-- 登录通知 -->
      <div class="mb-6">
        <Switch
          v-model="securitySettings.loginNotifications"
          label="登录通知"
          description="新设备登录时发送邮件通知"
          @change="updateLoginNotifications"
        />
      </div>

      <!-- 注销其他设备 -->
      <div class="p-4 bg-orange-500/10 rounded-lg border border-orange-500/20">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="text-md font-medium text-text-primary">注销其他设备</h4>
            <p class="text-sm text-text-secondary mt-1">
              强制注销所有其他设备上的登录会话
            </p>
          </div>
          <Button
            variant="secondary"
            @click="showLogoutModal = true"
          >
            注销其他设备
          </Button>
        </div>
      </div>
    </GlassCard>

    <!-- 登录历史 -->
    <GlassCard padding="lg">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-text-primary">登录历史</h3>
        <Button
          variant="ghost"
          size="sm"
          @click="loadLoginHistory"
          :loading="historyLoading"
        >
          <ArrowPathIcon class="w-4 h-4 mr-2" />
          刷新
        </Button>
      </div>

      <!-- 登录记录列表 -->
      <div class="space-y-3">
        <div
          v-for="record in loginHistory"
          :key="record.id"
          class="flex items-center justify-between p-3 bg-white/5 rounded-lg"
        >
          <div class="flex items-center space-x-3">
            <div
              :class="[
                'w-2 h-2 rounded-full',
                record.success ? 'bg-green-400' : 'bg-red-400'
              ]"
            />
            <div>
              <p class="text-sm font-medium text-text-primary">
                {{ record.location || '未知位置' }}
              </p>
              <p class="text-xs text-text-secondary">
                {{ record.ip }} • {{ formatDate(record.loginAt) }}
              </p>
            </div>
          </div>
          <div
            :class="[
              'px-2 py-1 rounded text-xs font-medium',
              record.success
                ? 'bg-green-500/20 text-green-400'
                : 'bg-red-500/20 text-red-400'
            ]"
          >
            {{ record.success ? '成功' : '失败' }}
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 注销确认模态框 -->
    <Modal
      v-model="showLogoutModal"
      title="注销其他设备"
      @confirm="logoutOtherDevices"
      :loading="logoutLoading"
    >
      <div class="space-y-4">
        <p class="text-text-secondary">
          此操作将强制注销您在所有其他设备上的登录会话。请输入密码确认：
        </p>
        <Input
          v-model="logoutPassword"
          type="password"
          placeholder="请输入密码"
          :error="logoutError"
        />
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { userSettingsApi } from '@/api/user-settings'
import { useNotification } from '@/composables/useNotification'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'
import Modal from '@/components/ui/Modal.vue'

import {
  LockClosedIcon,
  KeyIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline'

// 通知
const { showSuccess, showError } = useNotification()

// 响应式数据
const passwordLoading = ref(false)
const twoFactorLoading = ref(false)
const historyLoading = ref(false)
const logoutLoading = ref(false)

const showTwoFactorSetup = ref(false)
const showLogoutModal = ref(false)

const twoFactorToken = ref('')
const twoFactorError = ref('')
const logoutPassword = ref('')
const logoutError = ref('')

// 密码表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordErrors = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 安全设置
const securitySettings = reactive({
  twoFactorEnabled: false,
  sessionTimeout: 30,
  loginNotifications: true
})

// 两步验证设置
const twoFactorSetup = reactive({
  qrCode: '',
  secret: ''
})

// 登录历史
const loginHistory = ref<Array<{
  id: string
  ip: string
  userAgent: string
  location?: string
  loginAt: string
  success: boolean
}>>([])

// 计算属性
const canChangePassword = computed(() => {
  return passwordForm.currentPassword && passwordForm.newPassword && passwordForm.confirmPassword
})

// 会话超时选项
const sessionTimeoutOptions = [
  { label: '15分钟', value: 15 },
  { label: '30分钟', value: 30 },
  { label: '1小时', value: 60 },
  { label: '2小时', value: 120 },
  { label: '4小时', value: 240 },
  { label: '8小时', value: 480 }
]

// 生命周期
onMounted(() => {
  loadSecuritySettings()
  loadLoginHistory()
})

// 方法
const loadSecuritySettings = async () => {
  try {
    const settings = await userSettingsApi.getSettings()
    if (settings.security) {
      Object.assign(securitySettings, settings.security)
    }
  } catch (error) {
    console.error('Failed to load security settings:', error)
  }
}

const validatePassword = () => {
  // 清空错误
  Object.keys(passwordErrors).forEach(key => {
    passwordErrors[key as keyof typeof passwordErrors] = ''
  })

  let isValid = true

  // 验证当前密码
  if (!passwordForm.currentPassword) {
    passwordErrors.currentPassword = '请输入当前密码'
    isValid = false
  }

  // 验证新密码
  if (passwordForm.newPassword) {
    if (passwordForm.newPassword.length < 8) {
      passwordErrors.newPassword = '密码至少8位'
      isValid = false
    } else if (!/(?=.*[a-zA-Z])(?=.*\d)(?=.*[!@#$%^&*])/.test(passwordForm.newPassword)) {
      passwordErrors.newPassword = '密码必须包含字母、数字和特殊字符'
      isValid = false
    }
  }

  // 验证确认密码
  if (passwordForm.newPassword && passwordForm.confirmPassword !== passwordForm.newPassword) {
    passwordErrors.confirmPassword = '两次输入的密码不一致'
    isValid = false
  }

  return isValid
}

const changePassword = async () => {
  if (!validatePassword()) return

  try {
    passwordLoading.value = true
    
    await userSettingsApi.updateSecurity({
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword
    })

    // 清空表单
    Object.keys(passwordForm).forEach(key => {
      passwordForm[key as keyof typeof passwordForm] = ''
    })

    showSuccess('密码修改成功')
  } catch (error) {
    showError('密码修改失败，请检查当前密码是否正确')
  } finally {
    passwordLoading.value = false
  }
}

const handleTwoFactorToggle = async (enabled: boolean) => {
  if (enabled) {
    // 启用两步验证
    try {
      const setup = await userSettingsApi.setupTwoFactor()
      twoFactorSetup.qrCode = setup.qrCode
      twoFactorSetup.secret = setup.secret
      showTwoFactorSetup.value = true
    } catch (error) {
      securitySettings.twoFactorEnabled = false
      showError('无法设置两步验证')
    }
  } else {
    // 禁用两步验证
    try {
      const password = prompt('请输入密码以禁用两步验证：')
      if (password) {
        await userSettingsApi.disableTwoFactor(password)
        showSuccess('两步验证已禁用')
      } else {
        securitySettings.twoFactorEnabled = true
      }
    } catch (error) {
      securitySettings.twoFactorEnabled = true
      showError('禁用两步验证失败')
    }
  }
}

const enableTwoFactor = async () => {
  if (!twoFactorToken.value || twoFactorToken.value.length !== 6) {
    twoFactorError.value = '请输入6位验证码'
    return
  }

  try {
    twoFactorLoading.value = true
    const result = await userSettingsApi.enableTwoFactor(twoFactorToken.value)
    
    showTwoFactorSetup.value = false
    twoFactorToken.value = ''
    twoFactorError.value = ''
    
    showSuccess('两步验证已启用')
    
    // 显示恢复代码
    alert(`请保存以下恢复代码：\n${result.recoveryCodes.join('\n')}`)
  } catch (error) {
    twoFactorError.value = '验证码错误，请重试'
  } finally {
    twoFactorLoading.value = false
  }
}

const cancelTwoFactorSetup = () => {
  showTwoFactorSetup.value = false
  securitySettings.twoFactorEnabled = false
  twoFactorToken.value = ''
  twoFactorError.value = ''
}

const updateSessionTimeout = async () => {
  try {
    await userSettingsApi.updateSecurity({
      currentPassword: '', // 这里可能需要验证密码
      sessionTimeout: securitySettings.sessionTimeout
    })
    showSuccess('会话超时设置已更新')
  } catch (error) {
    showError('设置更新失败')
  }
}

const updateLoginNotifications = async () => {
  try {
    await userSettingsApi.updateSecurity({
      currentPassword: '',
      loginNotifications: securitySettings.loginNotifications
    })
    showSuccess('登录通知设置已更新')
  } catch (error) {
    showError('设置更新失败')
  }
}

const loadLoginHistory = async () => {
  try {
    historyLoading.value = true
    loginHistory.value = await userSettingsApi.getLoginHistory()
  } catch (error) {
    showError('加载登录历史失败')
  } finally {
    historyLoading.value = false
  }
}

const logoutOtherDevices = async () => {
  if (!logoutPassword.value) {
    logoutError.value = '请输入密码'
    return
  }

  try {
    logoutLoading.value = true
    await userSettingsApi.logoutOtherDevices(logoutPassword.value)
    
    showLogoutModal.value = false
    logoutPassword.value = ''
    logoutError.value = ''
    
    showSuccess('已注销其他设备')
  } catch (error) {
    logoutError.value = '密码错误'
  } finally {
    logoutLoading.value = false
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}
</script>
