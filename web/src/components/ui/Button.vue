<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200 focus-ring',
      sizeClasses[size],
      variantClasses[variant],
      {
        'opacity-50 cursor-not-allowed': disabled,
        'cursor-wait': loading
      },
      className
    ]"
    v-motion
    :hovered="{ scale: disabled || loading ? 1 : 1.02 }"
    :tapped="{ scale: disabled || loading ? 1 : 0.98 }"
    @click="handleClick"
  >
    <!-- 加载图标 -->
    <svg
      v-if="loading"
      class="animate-spin -ml-1 mr-2 h-4 w-4"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>

    <!-- 左侧图标 -->
    <component
      v-else-if="icon && iconPosition === 'left'"
      :is="icon"
      :class="[
        'flex-shrink-0',
        $slots.default ? 'mr-2' : '',
        iconSizes[size]
      ]"
    />

    <!-- 按钮文字 -->
    <span v-if="$slots.default">
      <slot />
    </span>

    <!-- 右侧图标 -->
    <component
      v-if="icon && iconPosition === 'right' && !loading"
      :is="icon"
      :class="[
        'flex-shrink-0',
        $slots.default ? 'ml-2' : '',
        iconSizes[size]
      ]"
    />
  </button>
</template>

<script setup lang="ts">
import type { ButtonProps } from '@/types'

interface Props extends ButtonProps {
  type?: 'button' | 'submit' | 'reset'
  className?: string
}

interface Emits {
  (e: 'click', event: MouseEvent): void
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  loading: false,
  disabled: false,
  iconPosition: 'left',
  className: ''
})

const emit = defineEmits<Emits>()

// 尺寸样式
const sizeClasses = {
  sm: 'px-3 py-2 text-sm',
  md: 'px-4 py-2.5 text-sm',
  lg: 'px-6 py-3 text-base'
}

// 图标尺寸
const iconSizes = {
  sm: 'h-4 w-4',
  md: 'h-4 w-4',
  lg: 'h-5 w-5'
}

// 变体样式
const variantClasses = {
  primary: 'bg-gradient-to-r from-primary-600 to-primary-500 text-white hover:from-primary-700 hover:to-primary-600 shadow-lg hover:shadow-xl',
  secondary: 'glass-card text-text-primary border-glass-border hover:glass-medium',
  ghost: 'text-text-secondary hover:text-text-primary hover:bg-white/5',
  danger: 'bg-gradient-to-r from-red-600 to-red-500 text-white hover:from-red-700 hover:to-red-600 shadow-lg hover:shadow-xl'
}

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
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

.glass-medium {
  background: var(--color-glass-medium);
}
</style>
