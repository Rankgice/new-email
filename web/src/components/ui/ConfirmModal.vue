<template>
  <Modal
    :visible="visible"
    :title="title"
    size="sm"
    @close="$emit('close')"
  >
    <!-- 内容 -->
    <div class="text-center">
      <!-- 图标 -->
      <div class="mx-auto flex items-center justify-center w-12 h-12 rounded-full mb-4" :class="iconBgClass">
        <component :is="iconComponent" class="w-6 h-6" :class="iconClass" />
      </div>
      
      <!-- 消息 -->
      <div class="mb-6">
        <p class="text-text-primary text-sm">
          {{ message }}
        </p>
      </div>
    </div>

    <!-- 操作按钮 -->
    <template #footer>
      <Button
        variant="ghost"
        @click="$emit('close')"
      >
        {{ cancelText }}
      </Button>
      <Button
        :variant="confirmVariant"
        :loading="loading"
        @click="$emit('confirm')"
      >
        {{ confirmText }}
      </Button>
    </template>
  </Modal>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  ExclamationTriangleIcon,
  InformationCircleIcon,
  CheckCircleIcon,
  XCircleIcon
} from '@heroicons/vue/24/outline'

import Modal from './Modal.vue'
import Button from './Button.vue'

// Props
interface Props {
  visible: boolean
  title?: string
  message: string
  type?: 'warning' | 'info' | 'success' | 'danger'
  confirmText?: string
  cancelText?: string
  confirmVariant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'ghost'
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '确认操作',
  type: 'warning',
  confirmText: '确认',
  cancelText: '取消',
  confirmVariant: 'primary',
  loading: false
})

// Emits
const emit = defineEmits<{
  close: []
  confirm: []
}>()

// 计算属性
const iconComponent = computed(() => {
  const icons = {
    warning: ExclamationTriangleIcon,
    info: InformationCircleIcon,
    success: CheckCircleIcon,
    danger: XCircleIcon
  }
  return icons[props.type]
})

const iconClass = computed(() => {
  const classes = {
    warning: 'text-yellow-400',
    info: 'text-blue-400',
    success: 'text-green-400',
    danger: 'text-red-400'
  }
  return classes[props.type]
})

const iconBgClass = computed(() => {
  const classes = {
    warning: 'bg-yellow-400/10',
    info: 'bg-blue-400/10',
    success: 'bg-green-400/10',
    danger: 'bg-red-400/10'
  }
  return classes[props.type]
})
</script>
