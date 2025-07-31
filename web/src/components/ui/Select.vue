<template>
  <div class="relative">
    <!-- 标签 -->
    <label
      v-if="label"
      :for="selectId"
      class="block text-sm font-medium text-text-primary mb-2"
    >
      {{ label }}
      <span v-if="required" class="text-red-400 ml-1">*</span>
    </label>

    <!-- 选择框容器 -->
    <div class="relative">
      <!-- 左侧图标 -->
      <div
        v-if="leftIcon"
        class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
      >
        <component :is="leftIcon" class="h-5 w-5 text-text-secondary" />
      </div>

      <!-- 选择框 -->
      <select
        :id="selectId"
        :value="modelValue"
        :disabled="disabled"
        :class="[
          'w-full rounded-lg border transition-all duration-200 focus-ring glass-card',
          'text-text-primary bg-transparent appearance-none cursor-pointer',
          {
            'pl-10': leftIcon,
            'pl-4': !leftIcon,
            'pr-10': true, // 为下拉箭头留空间
            'py-3': size === 'md',
            'py-2': size === 'sm',
            'py-4': size === 'lg',
            'border-glass-border': !error,
            'border-red-400 focus:border-red-400 focus:ring-red-400': error,
            'opacity-50 cursor-not-allowed': disabled
          }
        ]"
        @change="handleChange"
      >
        <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
        <option
          v-for="option in options"
          :key="option.value"
          :value="option.value"
          :disabled="option.disabled"
        >
          {{ option.label }}
        </option>
      </select>

      <!-- 下拉箭头 -->
      <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
        <ChevronDownIcon class="h-5 w-5 text-text-secondary" />
      </div>
    </div>

    <!-- 错误信息 -->
    <p v-if="error" class="mt-1 text-sm text-red-400">
      {{ error }}
    </p>

    <!-- 帮助文字 -->
    <p v-if="help && !error" class="mt-1 text-sm text-text-secondary">
      {{ help }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ChevronDownIcon } from '@heroicons/vue/24/outline'

interface Option {
  label: string
  value: string | number
  disabled?: boolean
}

interface Props {
  modelValue: string | number
  options: Option[]
  label?: string
  placeholder?: string
  error?: string
  help?: string
  disabled?: boolean
  required?: boolean
  size?: 'sm' | 'md' | 'lg'
  leftIcon?: any
}

interface Emits {
  (e: 'update:modelValue', value: string | number): void
  (e: 'change', value: string | number): void
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  disabled: false,
  required: false
})

const emit = defineEmits<Emits>()

// 生成唯一 ID
const selectId = `select-${Math.random().toString(36).substr(2, 9)}`

// 处理选择变化
const handleChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  const value = target.value
  
  // 尝试转换为数字（如果原始值是数字类型）
  const convertedValue = isNaN(Number(value)) ? value : Number(value)
  
  emit('update:modelValue', convertedValue)
  emit('change', convertedValue)
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

/* 自定义选择框样式 */
select {
  background-image: none;
}

select option {
  background-color: var(--color-surface);
  color: var(--color-text-primary);
}
</style>
