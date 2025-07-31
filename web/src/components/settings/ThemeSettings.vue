<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h2 class="text-xl font-semibold text-text-primary mb-2">主题设置</h2>
      <p class="text-text-secondary">个性化您的界面外观和视觉体验</p>
    </div>

    <!-- 主题选择 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">预设主题</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="themeOption in availableThemes"
          :key="themeOption.id"
          :class="[
            'relative p-4 rounded-xl border-2 cursor-pointer transition-all duration-200',
            currentTheme.id === themeOption.id
              ? 'border-primary-500 bg-primary-500/10'
              : 'border-glass-border hover:border-primary-500/50'
          ]"
          @click="selectTheme(themeOption)"
        >
          <!-- 主题预览 -->
          <div class="mb-3">
            <div
              class="h-20 rounded-lg overflow-hidden"
              :style="{ background: themeOption.preview.background }"
            >
              <div class="h-full flex">
                <!-- 侧边栏预览 -->
                <div
                  class="w-1/3 h-full"
                  :style="{ background: themeOption.preview.sidebar }"
                />
                <!-- 主内容预览 -->
                <div class="flex-1 p-2 space-y-1">
                  <div
                    class="h-2 rounded"
                    :style="{ background: themeOption.preview.primary }"
                  />
                  <div
                    class="h-1 w-2/3 rounded"
                    :style="{ background: themeOption.preview.text }"
                  />
                  <div
                    class="h-1 w-1/2 rounded"
                    :style="{ background: themeOption.preview.secondary }"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- 主题信息 -->
          <div class="text-center">
            <h4 class="font-medium text-text-primary mb-1">
              {{ themeOption.name }}
            </h4>
            <p class="text-sm text-text-secondary">
              {{ themeOption.description }}
            </p>
          </div>

          <!-- 选中标识 -->
          <div
            v-if="currentTheme.id === themeOption.id"
            class="absolute top-2 right-2 w-6 h-6 bg-primary-500 rounded-full flex items-center justify-center"
          >
            <CheckIcon class="w-4 h-4 text-white" />
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 自定义设置 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">自定义设置</h3>
      <div class="space-y-6">
        <!-- 主色调 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-3">
            主色调
          </label>
          <div class="grid grid-cols-8 gap-3">
            <button
              v-for="color in primaryColors"
              :key="color.name"
              :class="[
                'w-10 h-10 rounded-lg border-2 transition-all duration-200',
                customSettings.primaryColor === color.value
                  ? 'border-white scale-110'
                  : 'border-transparent hover:scale-105'
              ]"
              :style="{ backgroundColor: color.value }"
              @click="updatePrimaryColor(color.value)"
            />
          </div>
        </div>

        <!-- 背景模糊度 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-3">
            背景模糊度: {{ customSettings.blurIntensity }}px
          </label>
          <input
            v-model="customSettings.blurIntensity"
            type="range"
            min="0"
            max="40"
            step="2"
            class="w-full h-2 bg-gray-600 rounded-lg appearance-none cursor-pointer slider"
            @input="updateBlurIntensity"
          />
        </div>

        <!-- 透明度 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-3">
            界面透明度: {{ Math.round(customSettings.opacity * 100) }}%
          </label>
          <input
            v-model="customSettings.opacity"
            type="range"
            min="0.3"
            max="1"
            step="0.05"
            class="w-full h-2 bg-gray-600 rounded-lg appearance-none cursor-pointer slider"
            @input="updateOpacity"
          />
        </div>

        <!-- 圆角大小 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-3">
            圆角大小: {{ customSettings.borderRadius }}px
          </label>
          <input
            v-model="customSettings.borderRadius"
            type="range"
            min="4"
            max="24"
            step="2"
            class="w-full h-2 bg-gray-600 rounded-lg appearance-none cursor-pointer slider"
            @input="updateBorderRadius"
          />
        </div>

        <!-- 字体大小 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-3">
            字体大小
          </label>
          <Select
            v-model="customSettings.fontSize"
            :options="fontSizeOptions"
            @change="updateFontSize"
          />
        </div>

        <!-- 动画效果 -->
        <Switch
          v-model="customSettings.animations"
          label="动画效果"
          description="启用界面过渡动画和交互效果"
          @change="updateAnimations"
        />

        <!-- 紧凑模式 -->
        <Switch
          v-model="customSettings.compactMode"
          label="紧凑模式"
          description="减少界面元素间距，显示更多内容"
          @change="updateCompactMode"
        />
      </div>
    </GlassCard>

    <!-- 高级设置 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">高级设置</h3>
      <div class="space-y-4">
        <!-- 自定义CSS -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            自定义CSS
          </label>
          <Textarea
            v-model="customSettings.customCss"
            placeholder="/* 在这里输入自定义CSS代码 */"
            :rows="6"
            help="高级用户可以通过CSS进一步自定义界面样式"
            @input="updateCustomCss"
          />
        </div>

        <!-- 重置按钮 -->
        <div class="flex justify-between items-center pt-4 border-t border-glass-border">
          <div>
            <h4 class="font-medium text-text-primary">重置设置</h4>
            <p class="text-sm text-text-secondary">
              将所有主题设置恢复为默认值
            </p>
          </div>
          <Button
            variant="secondary"
            @click="resetToDefault"
          >
            重置为默认
          </Button>
        </div>
      </div>
    </GlassCard>

    <!-- 导入导出 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">导入导出</h3>
      <div class="space-y-4">
        <div class="flex space-x-3">
          <Button
            variant="secondary"
            @click="exportTheme"
          >
            <ArrowDownTrayIcon class="w-4 h-4 mr-2" />
            导出主题
          </Button>
          <Button
            variant="secondary"
            @click="triggerImport"
          >
            <ArrowUpTrayIcon class="w-4 h-4 mr-2" />
            导入主题
          </Button>
        </div>
        <input
          ref="fileInput"
          type="file"
          accept=".json"
          class="hidden"
          @change="importTheme"
        />
        <p class="text-sm text-text-secondary">
          导出您的主题设置以便在其他设备上使用，或导入他人分享的主题配置
        </p>
      </div>
    </GlassCard>

    <!-- 保存按钮 -->
    <div class="flex justify-end">
      <Button
        variant="primary"
        @click="saveThemeSettings"
        :loading="loading"
      >
        保存设置
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme'
import { userSettingsApi } from '@/api/user-settings'
import { useNotification } from '@/composables/useNotification'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'
import Textarea from '@/components/ui/Textarea.vue'

import {
  CheckIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon
} from '@heroicons/vue/24/outline'

// 主题管理
const { currentTheme, availableThemes, setTheme } = useTheme()
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const fileInput = ref<HTMLInputElement>()

// 自定义设置
const customSettings = reactive({
  primaryColor: '#6366f1',
  blurIntensity: 20,
  opacity: 0.95,
  borderRadius: 12,
  fontSize: 'medium',
  animations: true,
  compactMode: false,
  customCss: ''
})

// 主色调选项
const primaryColors = [
  { name: '靛蓝', value: '#6366f1' },
  { name: '蓝色', value: '#3b82f6' },
  { name: '青色', value: '#06b6d4' },
  { name: '绿色', value: '#10b981' },
  { name: '黄色', value: '#f59e0b' },
  { name: '橙色', value: '#f97316' },
  { name: '红色', value: '#ef4444' },
  { name: '粉色', value: '#ec4899' },
  { name: '紫色', value: '#8b5cf6' },
  { name: '灰色', value: '#6b7280' }
]

// 字体大小选项
const fontSizeOptions = [
  { label: '小', value: 'small' },
  { label: '中等', value: 'medium' },
  { label: '大', value: 'large' },
  { label: '特大', value: 'extra-large' }
]

// 生命周期
onMounted(() => {
  loadThemeSettings()
})

// 方法
const loadThemeSettings = async () => {
  try {
    const settings = await userSettingsApi.getSettings()
    // 加载自定义主题设置
    if (settings.theme) {
      // 解析主题设置
    }
  } catch (error) {
    console.error('Failed to load theme settings:', error)
  }
}

const selectTheme = (theme: any) => {
  setTheme(theme.id)
}

const updatePrimaryColor = (color: string) => {
  customSettings.primaryColor = color
  applyCustomSettings()
}

const updateBlurIntensity = () => {
  applyCustomSettings()
}

const updateOpacity = () => {
  applyCustomSettings()
}

const updateBorderRadius = () => {
  applyCustomSettings()
}

const updateFontSize = () => {
  applyCustomSettings()
}

const updateAnimations = () => {
  applyCustomSettings()
}

const updateCompactMode = () => {
  applyCustomSettings()
}

const updateCustomCss = () => {
  applyCustomSettings()
}

const applyCustomSettings = () => {
  // 应用自定义设置到CSS变量
  const root = document.documentElement
  
  root.style.setProperty('--color-primary', customSettings.primaryColor)
  root.style.setProperty('--blur-intensity', `${customSettings.blurIntensity}px`)
  root.style.setProperty('--glass-opacity', customSettings.opacity.toString())
  root.style.setProperty('--border-radius', `${customSettings.borderRadius}px`)
  
  // 字体大小
  const fontSizeMap = {
    small: '14px',
    medium: '16px',
    large: '18px',
    'extra-large': '20px'
  }
  root.style.setProperty('--font-size-base', fontSizeMap[customSettings.fontSize as keyof typeof fontSizeMap])
  
  // 动画
  root.style.setProperty('--animation-duration', customSettings.animations ? '200ms' : '0ms')
  
  // 紧凑模式
  root.style.setProperty('--spacing-scale', customSettings.compactMode ? '0.8' : '1')
  
  // 自定义CSS
  let customStyleElement = document.getElementById('custom-theme-styles')
  if (!customStyleElement) {
    customStyleElement = document.createElement('style')
    customStyleElement.id = 'custom-theme-styles'
    document.head.appendChild(customStyleElement)
  }
  customStyleElement.textContent = customSettings.customCss
}

const resetToDefault = () => {
  if (!confirm('确定要重置所有主题设置为默认值吗？')) return
  
  customSettings.primaryColor = '#6366f1'
  customSettings.blurIntensity = 20
  customSettings.opacity = 0.95
  customSettings.borderRadius = 12
  customSettings.fontSize = 'medium'
  customSettings.animations = true
  customSettings.compactMode = false
  customSettings.customCss = ''
  
  applyCustomSettings()
  setTheme('dark') // 重置为默认主题
  showSuccess('主题设置已重置为默认值')
}

const exportTheme = () => {
  const themeConfig = {
    theme: currentTheme.value.id,
    customSettings: { ...customSettings },
    exportDate: new Date().toISOString(),
    version: '1.0'
  }
  
  const blob = new Blob([JSON.stringify(themeConfig, null, 2)], {
    type: 'application/json'
  })
  
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `theme-${currentTheme.value.id}-${Date.now()}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  showSuccess('主题配置已导出')
}

const triggerImport = () => {
  fileInput.value?.click()
}

const importTheme = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const themeConfig = JSON.parse(e.target?.result as string)
      
      // 验证配置格式
      if (!themeConfig.theme || !themeConfig.customSettings) {
        throw new Error('Invalid theme configuration')
      }
      
      // 应用导入的设置
      setTheme(themeConfig.theme)
      Object.assign(customSettings, themeConfig.customSettings)
      applyCustomSettings()
      
      showSuccess('主题配置导入成功')
    } catch (error) {
      showError('主题配置文件格式错误')
    }
  }
  
  reader.readAsText(file)
  
  // 清空文件输入
  target.value = ''
}

const saveThemeSettings = async () => {
  try {
    loading.value = true
    
    const themeSettings = {
      theme: currentTheme.value.id,
      customSettings: { ...customSettings }
    }
    
    await userSettingsApi.updateTheme(JSON.stringify(themeSettings))
    showSuccess('主题设置保存成功')
  } catch (error) {
    showError('保存主题设置失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 滑块样式 */
.slider::-webkit-slider-thumb {
  appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: var(--color-primary);
  cursor: pointer;
  border: 2px solid #ffffff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}

.slider::-moz-range-thumb {
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: var(--color-primary);
  cursor: pointer;
  border: 2px solid #ffffff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
}
</style>
