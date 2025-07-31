<template>
  <div class="flex items-center">
    <!-- 标签 -->
    <label
      v-if="label"
      :for="switchId"
      class="text-sm font-medium text-text-primary mr-3 cursor-pointer"
    >
      {{ label }}
      <span v-if="required" class="text-red-400 ml-1">*</span>
    </label>

    <!-- 开关容器 -->
    <div class="relative">
      <!-- 隐藏的原生 checkbox -->
      <input
        :id="switchId"
        type="checkbox"
        :checked="modelValue"
        :disabled="disabled"
        class="sr-only"
        @change="handleChange"
      />

      <!-- 开关背景 -->
      <div
        :class="[
          'relative inline-flex h-6 w-11 items-center rounded-full transition-all duration-200 cursor-pointer',
          'focus-within:ring-2 focus-within:ring-primary-500 focus-within:ring-offset-2 focus-within:ring-offset-background-primary',
          {
            'bg-primary-500': modelValue,
            'bg-gray-600': !modelValue,
            'opacity-50 cursor-not-allowed': disabled
          }
        ]"
        @click="toggle"
      >
        <!-- 开关滑块 -->
        <span
          :class="[
            'inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-200 shadow-lg',
            {
              'translate-x-6': modelValue,
              'translate-x-1': !modelValue
            }
          ]"
        />
      </div>
    </div>

    <!-- 描述文字 -->
    <div v-if="description" class="ml-3 flex-1">
      <p class="text-sm text-text-secondary">{{ description }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  modelValue: boolean
  label?: string
  description?: string
  disabled?: boolean
  required?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'change', value: boolean): void
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: false
})

const emit = defineEmits<Emits>()

// 生成唯一 ID
const switchId = `switch-${Math.random().toString(36).substr(2, 9)}`

// 切换开关
const toggle = () => {
  if (!props.disabled) {
    const newValue = !props.modelValue
    emit('update:modelValue', newValue)
    emit('change', newValue)
  }
}

// 处理原生 change 事件
const handleChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.checked)
  emit('change', target.checked)
}
</script>

<style scoped>
.focus-ring {
  @apply focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-background-primary;
}
</style>
