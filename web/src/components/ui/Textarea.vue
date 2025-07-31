<template>
  <div>
    <!-- 标签 -->
    <label
      v-if="label"
      :for="textareaId"
      class="block text-sm font-medium text-text-primary mb-2"
    >
      {{ label }}
      <span v-if="required" class="text-red-400 ml-1">*</span>
    </label>

    <!-- 文本域容器 -->
    <div class="relative">
      <!-- 文本域 -->
      <textarea
        :id="textareaId"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :rows="rows"
        :maxlength="maxLength"
        :class="[
          'w-full rounded-lg border transition-all duration-200 focus-ring glass-card resize-none',
          'text-text-primary placeholder-text-secondary',
          {
            'p-4': size === 'lg',
            'p-3': size === 'md',
            'p-2': size === 'sm',
            'border-glass-border': !error,
            'border-red-400 focus:border-red-400 focus:ring-red-400': error,
            'opacity-50 cursor-not-allowed': disabled
          }
        ]"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
      />

      <!-- 字符计数 -->
      <div
        v-if="maxLength && showCount"
        class="absolute bottom-2 right-2 text-xs text-text-secondary"
      >
        {{ currentLength }}/{{ maxLength }}
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
import { computed } from 'vue'

interface Props {
  modelValue: string
  label?: string
  placeholder?: string
  error?: string
  help?: string
  disabled?: boolean
  readonly?: boolean
  required?: boolean
  rows?: number
  maxLength?: number
  showCount?: boolean
  size?: 'sm' | 'md' | 'lg'
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'blur', event: FocusEvent): void
  (e: 'focus', event: FocusEvent): void
}

const props = withDefaults(defineProps<Props>(), {
  rows: 4,
  size: 'md',
  disabled: false,
  readonly: false,
  required: false,
  showCount: false
})

const emit = defineEmits<Emits>()

// 生成唯一 ID
const textareaId = `textarea-${Math.random().toString(36).substr(2, 9)}`

// 当前字符长度
const currentLength = computed(() => props.modelValue?.length || 0)

// 处理输入
const handleInput = (event: Event) => {
  const target = event.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
}

// 处理失焦
const handleBlur = (event: FocusEvent) => {
  emit('blur', event)
}

// 处理聚焦
const handleFocus = (event: FocusEvent) => {
  emit('focus', event)
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
