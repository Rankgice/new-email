<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="visible"
        class="fixed inset-0 z-50 overflow-y-auto"
        @click="handleBackdropClick"
      >
        <!-- 背景遮罩 -->
        <div class="fixed inset-0 bg-black/50 backdrop-blur-sm"></div>
        
        <!-- 模态框容器 -->
        <div class="flex min-h-full items-center justify-center p-4">
          <Transition
            enter-active-class="transition ease-out duration-300"
            enter-from-class="opacity-0 scale-95"
            enter-to-class="opacity-100 scale-100"
            leave-active-class="transition ease-in duration-200"
            leave-from-class="opacity-100 scale-100"
            leave-to-class="opacity-0 scale-95"
          >
            <div
              v-if="visible"
              :class="[
                'relative w-full bg-glass-medium backdrop-blur-md border border-glass-border rounded-xl shadow-xl',
                sizeClasses
              ]"
              @click.stop
            >
              <!-- 模态框头部 -->
              <div v-if="title || $slots.header" class="flex items-center justify-between p-6 border-b border-glass-border">
                <div class="flex items-center space-x-3">
                  <slot name="header">
                    <h3 class="text-lg font-semibold text-text-primary">
                      {{ title }}
                    </h3>
                  </slot>
                </div>
                <button
                  v-if="closable"
                  @click="$emit('close')"
                  class="text-text-secondary hover:text-text-primary transition-colors p-1 rounded-lg hover:bg-glass-light"
                >
                  <XMarkIcon class="w-5 h-5" />
                </button>
              </div>

              <!-- 模态框内容 -->
              <div class="p-6">
                <slot />
              </div>

              <!-- 模态框底部 -->
              <div v-if="$slots.footer" class="flex items-center justify-end space-x-3 p-6 border-t border-glass-border">
                <slot name="footer" />
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'

// Props
interface Props {
  visible: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closable?: boolean
  closeOnBackdrop?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closable: true,
  closeOnBackdrop: true
})

// Emits
const emit = defineEmits<{
  close: []
}>()

// 计算属性
const sizeClasses = computed(() => {
  const sizes = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
    full: 'max-w-7xl'
  }
  return sizes[props.size]
})

// 方法
const handleBackdropClick = () => {
  if (props.closeOnBackdrop && props.closable) {
    emit('close')
  }
}

const handleEscapeKey = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && props.visible && props.closable) {
    emit('close')
  }
}

// 生命周期
onMounted(() => {
  document.addEventListener('keydown', handleEscapeKey)
  // 防止背景滚动
  if (props.visible) {
    document.body.style.overflow = 'hidden'
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscapeKey)
  // 恢复背景滚动
  document.body.style.overflow = ''
})

// 监听visible变化
import { watch } from 'vue'
watch(() => props.visible, (newVisible) => {
  if (newVisible) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>
