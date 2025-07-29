<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
      <TransitionGroup
        name="notification"
        tag="div"
        class="space-y-2"
      >
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-item"
        >
          <GlassCard
            :level="3"
            :hover="false"
            padding="md"
            :class="[
              'relative overflow-hidden',
              notificationStyles[notification.type]
            ]"
          >
            <!-- 进度条 -->
            <div
              v-if="notification.duration"
              class="absolute top-0 left-0 h-1 bg-current opacity-30 transition-all ease-linear"
              :style="{
                width: `${getProgress(notification)}%`,
                transitionDuration: `${notification.duration}ms`
              }"
            />

            <div class="flex items-start space-x-3">
              <!-- 图标 -->
              <div class="flex-shrink-0">
                <component
                  :is="notificationIcons[notification.type]"
                  class="w-5 h-5"
                />
              </div>

              <!-- 内容 -->
              <div class="flex-1 min-w-0">
                <h4 class="text-sm font-medium">
                  {{ notification.title }}
                </h4>
                <p
                  v-if="notification.message"
                  class="mt-1 text-sm opacity-90"
                >
                  {{ notification.message }}
                </p>

                <!-- 操作按钮 -->
                <div
                  v-if="notification.actions && notification.actions.length > 0"
                  class="mt-3 flex space-x-2"
                >
                  <button
                    v-for="action in notification.actions"
                    :key="action.label"
                    :class="[
                      'px-3 py-1 text-xs font-medium rounded transition-colors duration-200',
                      action.style === 'primary'
                        ? 'bg-white/20 hover:bg-white/30'
                        : 'bg-white/10 hover:bg-white/20'
                    ]"
                    @click="handleAction(action, notification.id)"
                  >
                    {{ action.label }}
                  </button>
                </div>
              </div>

              <!-- 关闭按钮 -->
              <button
                class="flex-shrink-0 p-1 rounded-full hover:bg-white/20 transition-colors duration-200"
                @click="removeNotification(notification.id)"
              >
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </GlassCard>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useNotification } from '@/composables/useNotification'
import GlassCard from './GlassCard.vue'
import {
  CheckCircleIcon,
  ExclamationTriangleIcon,
  XCircleIcon,
  InformationCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import type { Notification, NotificationAction } from '@/types'

const { notifications, removeNotification } = useNotification()

// 通知图标映射
const notificationIcons = {
  success: CheckCircleIcon,
  warning: ExclamationTriangleIcon,
  error: XCircleIcon,
  info: InformationCircleIcon
}

// 通知样式映射
const notificationStyles = {
  success: 'text-green-400 border-green-400/20',
  warning: 'text-yellow-400 border-yellow-400/20',
  error: 'text-red-400 border-red-400/20',
  info: 'text-blue-400 border-blue-400/20'
}

// 计算进度条进度
const getProgress = (notification: Notification) => {
  if (!notification.duration) return 100
  
  const elapsed = Date.now() - new Date(notification.createdAt).getTime()
  const progress = Math.max(0, 100 - (elapsed / notification.duration) * 100)
  return progress
}

// 处理操作按钮点击
const handleAction = (action: NotificationAction, notificationId: string) => {
  action.action()
  removeNotification(notificationId)
}
</script>

<style scoped>
/* 通知动画 */
.notification-enter-active {
  transition: all 0.3s ease-out;
}

.notification-leave-active {
  transition: all 0.3s ease-in;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.notification-move {
  transition: transform 0.3s ease;
}

/* 进度条动画 */
.notification-item {
  position: relative;
}

.notification-item .absolute {
  animation: progress-shrink linear;
}

@keyframes progress-shrink {
  from {
    width: 100%;
  }
  to {
    width: 0%;
  }
}
</style>
