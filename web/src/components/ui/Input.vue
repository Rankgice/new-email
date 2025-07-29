<template>
  <div class="relative">
    <!-- 标签 -->
    <label
      v-if="label"
      :for="inputId"
      :class="[
        'block text-sm font-medium mb-2 transition-colors duration-200',
        {
          'text-text-primary': !error,
          'text-red-400': error
        }
      ]"
    >
      {{ label }}
      <span v-if="required" class="text-red-400 ml-1">*</span>
    </label>

    <!-- 输入框容器 -->
    <div class="relative">
      <!-- 左侧图标 -->
      <div
        v-if="leftIcon"
        class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
      >
        <component
          :is="leftIcon"
          class="h-5 w-5 text-text-secondary"
        />
      </div>

      <!-- 输入框 -->
      <input
        :id="inputId"
        :type="inputType"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :autocomplete="autocomplete"
        :class="[
          'w-full rounded-lg border transition-all duration-200 focus-ring glass-card',
          'text-text-primary placeholder-text-secondary',
          {
            'pl-10': leftIcon,
            'pr-10': rightIcon || type === 'password',
            'pl-4': !leftIcon,
            'pr-4': !rightIcon && type !== 'password',
            'py-3': size === 'md',
            'py-2': size === 'sm',
            'py-4': size === 'lg',
            'border-glass-border': !error,
            'border-red-400 focus:border-red-400 focus:ring-red-400': error,
            'opacity-50 cursor-not-allowed': disabled
          }
        ]"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
        @keydown.enter="handleEnter"
      />

      <!-- 右侧图标/密码切换 -->
      <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
        <!-- 密码显示/隐藏按钮 -->
        <button
          v-if="type === 'password'"
          type="button"
          class="text-text-secondary hover:text-text-primary transition-colors duration-200"
          @click="togglePasswordVisibility"
        >
          <EyeIcon v-if="showPassword" class="h-5 w-5" />
          <EyeSlashIcon v-else class="h-5 w-5" />
        </button>

        <!-- 自定义右侧图标 -->
        <component
          v-else-if="rightIcon"
          :is="rightIcon"
          class="h-5 w-5 text-text-secondary"
        />
      </div>
    </div>

    <!-- 错误信息 -->
    <p
      v-if="error"
      class="mt-2 text-sm text-red-400 animate-slide-down"
    >
      {{ error }}
    </p>

    <!-- 帮助文本 -->
    <p
      v-else-if="help"
      class="mt-2 text-sm text-text-secondary"
    >
      {{ help }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { EyeIcon, EyeSlashIcon } from '@heroicons/vue/24/outline'

interface Props {
  modelValue: string
  type?: 'text' | 'email' | 'password' | 'search' | 'tel' | 'url'
  label?: string
  placeholder?: string
  error?: string
  help?: string
  disabled?: boolean
  readonly?: boolean
  required?: boolean
  size?: 'sm' | 'md' | 'lg'
  leftIcon?: any
  rightIcon?: any
  autocomplete?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'blur', event: FocusEvent): void
  (e: 'focus', event: FocusEvent): void
  (e: 'enter', event: KeyboardEvent): void
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  size: 'md',
  disabled: false,
  readonly: false,
  required: false
})

const emit = defineEmits<Emits>()

// 生成唯一 ID
const inputId = `input-${Math.random().toString(36).substr(2, 9)}`

// 密码显示状态
const showPassword = ref(false)

// 计算实际输入框类型
const inputType = computed(() => {
  if (props.type === 'password') {
    return showPassword.value ? 'text' : 'password'
  }
  return props.type
})

// 切换密码显示
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

// 处理输入
const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
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

// 处理回车
const handleEnter = (event: KeyboardEvent) => {
  emit('enter', event)
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

.glass-card:focus {
  background: var(--color-glass-medium);
  border-color: var(--color-primary);
  box-shadow: 0 0 0 1px var(--color-primary);
}
</style>
