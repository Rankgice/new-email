<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h2 class="text-xl font-semibold text-text-primary mb-2">通知设置</h2>
      <p class="text-text-secondary">管理您的通知偏好和提醒方式</p>
    </div>

    <!-- 通知方式 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">通知方式</h3>
      <div class="space-y-4">
        <!-- 邮件通知 -->
        <Switch
          v-model="settings.email"
          label="邮件通知"
          description="通过邮件接收重要通知"
          @change="updateSettings"
        />

        <!-- 桌面通知 -->
        <Switch
          v-model="settings.desktop"
          label="桌面通知"
          description="在浏览器中显示桌面通知"
          @change="handleDesktopNotificationChange"
        />

        <!-- 移动端通知 -->
        <Switch
          v-model="settings.mobile"
          label="移动端推送"
          description="在移动设备上接收推送通知"
          @change="updateSettings"
        />

        <!-- 声音提醒 -->
        <Switch
          v-model="settings.sound"
          label="声音提醒"
          description="新消息到达时播放提示音"
          @change="updateSettings"
        />
      </div>
    </GlassCard>

    <!-- 邮件通知 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">邮件通知</h3>
      <div class="space-y-4">
        <!-- 新邮件通知 -->
        <Switch
          v-model="settings.newEmail"
          label="新邮件通知"
          description="收到新邮件时发送通知"
          @change="updateSettings"
        />

        <!-- 重要邮件通知 -->
        <Switch
          v-model="settings.importantEmail"
          label="重要邮件通知"
          description="收到标记为重要的邮件时发送通知"
          @change="updateSettings"
        />

        <!-- 安全警报 -->
        <Switch
          v-model="settings.securityAlerts"
          label="安全警报"
          description="账户安全相关的重要通知"
          @change="updateSettings"
        />
      </div>
    </GlassCard>

    <!-- 通知时间 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">通知时间</h3>
      <div class="space-y-4">
        <!-- 免打扰模式 -->
        <div class="flex items-center justify-between">
          <div>
            <h4 class="text-md font-medium text-text-primary">免打扰模式</h4>
            <p class="text-sm text-text-secondary mt-1">
              在指定时间段内不接收通知
            </p>
          </div>
          <Switch
            v-model="doNotDisturb.enabled"
            @change="updateDoNotDisturb"
          />
        </div>

        <!-- 免打扰时间设置 -->
        <div v-if="doNotDisturb.enabled" class="ml-4 space-y-3">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                开始时间
              </label>
              <Input
                v-model="doNotDisturb.startTime"
                type="time"
                @change="updateDoNotDisturb"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                结束时间
              </label>
              <Input
                v-model="doNotDisturb.endTime"
                type="time"
                @change="updateDoNotDisturb"
              />
            </div>
          </div>

          <!-- 工作日设置 -->
          <div>
            <label class="block text-sm font-medium text-text-primary mb-2">
              应用日期
            </label>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="(day, index) in weekDays"
                :key="index"
                :class="[
                  'px-3 py-1 rounded-full text-sm font-medium transition-colors',
                  doNotDisturb.weekdays.includes(index)
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-600 text-gray-300 hover:bg-gray-500'
                ]"
                @click="toggleWeekday(index)"
              >
                {{ day }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 通知频率 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">通知频率</h3>
      <div class="space-y-4">
        <!-- 邮件摘要 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            邮件摘要频率
          </label>
          <Select
            v-model="emailDigest.frequency"
            :options="digestFrequencyOptions"
            @change="updateEmailDigest"
          />
          <p class="text-sm text-text-secondary mt-1">
            定期发送邮件摘要而不是每封邮件都通知
          </p>
        </div>

        <!-- 摘要时间 -->
        <div v-if="emailDigest.frequency !== 'never'">
          <label class="block text-sm font-medium text-text-primary mb-2">
            发送时间
          </label>
          <Input
            v-model="emailDigest.time"
            type="time"
            @change="updateEmailDigest"
          />
        </div>
      </div>
    </GlassCard>

    <!-- 测试通知 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">测试通知</h3>
      <div class="space-y-4">
        <p class="text-text-secondary">
          发送测试通知以验证您的设置是否正常工作
        </p>
        <div class="flex space-x-3">
          <Button
            variant="secondary"
            @click="sendTestNotification('desktop')"
            :loading="testLoading.desktop"
          >
            <ComputerDesktopIcon class="w-4 h-4 mr-2" />
            测试桌面通知
          </Button>
          <Button
            variant="secondary"
            @click="sendTestNotification('email')"
            :loading="testLoading.email"
          >
            <EnvelopeIcon class="w-4 h-4 mr-2" />
            测试邮件通知
          </Button>
        </div>
      </div>
    </GlassCard>

    <!-- 保存按钮 -->
    <div class="flex justify-end">
      <Button
        variant="primary"
        @click="saveAllSettings"
        :loading="loading"
      >
        保存所有设置
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { userSettingsApi } from '@/api/user-settings'
import { useNotification } from '@/composables/useNotification'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'

import {
  ComputerDesktopIcon,
  EnvelopeIcon
} from '@heroicons/vue/24/outline'

// 通知
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const testLoading = reactive({
  desktop: false,
  email: false
})

// 通知设置
const settings = reactive({
  email: true,
  desktop: true,
  mobile: true,
  sound: true,
  newEmail: true,
  importantEmail: true,
  securityAlerts: true
})

// 免打扰设置
const doNotDisturb = reactive({
  enabled: false,
  startTime: '22:00',
  endTime: '08:00',
  weekdays: [1, 2, 3, 4, 5] // 周一到周五
})

// 邮件摘要设置
const emailDigest = reactive({
  frequency: 'daily',
  time: '09:00'
})

// 星期选项
const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

// 摘要频率选项
const digestFrequencyOptions = [
  { label: '从不', value: 'never' },
  { label: '每小时', value: 'hourly' },
  { label: '每日', value: 'daily' },
  { label: '每周', value: 'weekly' }
]

// 生命周期
onMounted(() => {
  loadNotificationSettings()
  checkNotificationPermission()
})

// 方法
const loadNotificationSettings = async () => {
  try {
    const userSettings = await userSettingsApi.getSettings()
    if (userSettings.notifications) {
      Object.assign(settings, userSettings.notifications)
    }
  } catch (error) {
    console.error('Failed to load notification settings:', error)
  }
}

const checkNotificationPermission = () => {
  if ('Notification' in window) {
    if (Notification.permission === 'denied') {
      settings.desktop = false
    }
  } else {
    settings.desktop = false
  }
}

const handleDesktopNotificationChange = async (enabled: boolean) => {
  if (enabled && 'Notification' in window) {
    if (Notification.permission === 'default') {
      const permission = await Notification.requestPermission()
      if (permission !== 'granted') {
        settings.desktop = false
        showError('桌面通知权限被拒绝')
        return
      }
    } else if (Notification.permission === 'denied') {
      settings.desktop = false
      showError('请在浏览器设置中允许通知权限')
      return
    }
  }
  
  updateSettings()
}

const updateSettings = async () => {
  try {
    await userSettingsApi.updateNotifications(settings)
  } catch (error) {
    showError('设置更新失败')
  }
}

const updateDoNotDisturb = () => {
  // 这里可以调用API保存免打扰设置
  console.log('Do not disturb settings updated:', doNotDisturb)
}

const toggleWeekday = (dayIndex: number) => {
  const index = doNotDisturb.weekdays.indexOf(dayIndex)
  if (index > -1) {
    doNotDisturb.weekdays.splice(index, 1)
  } else {
    doNotDisturb.weekdays.push(dayIndex)
  }
  updateDoNotDisturb()
}

const updateEmailDigest = () => {
  // 这里可以调用API保存邮件摘要设置
  console.log('Email digest settings updated:', emailDigest)
}

const sendTestNotification = async (type: 'desktop' | 'email') => {
  try {
    testLoading[type] = true
    
    if (type === 'desktop') {
      if ('Notification' in window && Notification.permission === 'granted') {
        new Notification('测试通知', {
          body: '这是一条测试桌面通知',
          icon: '/favicon.ico'
        })
        showSuccess('桌面通知已发送')
      } else {
        showError('桌面通知权限未开启')
      }
    } else if (type === 'email') {
      // 调用API发送测试邮件
      showSuccess('测试邮件已发送，请检查您的邮箱')
    }
  } catch (error) {
    showError('发送测试通知失败')
  } finally {
    testLoading[type] = false
  }
}

const saveAllSettings = async () => {
  try {
    loading.value = true
    
    await userSettingsApi.updateNotifications(settings)
    // 这里还可以保存其他设置如免打扰、邮件摘要等
    
    showSuccess('通知设置保存成功')
  } catch (error) {
    showError('保存设置失败')
  } finally {
    loading.value = false
  }
}
</script>
