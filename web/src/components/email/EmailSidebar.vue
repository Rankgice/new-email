<template>
  <div class="w-64 bg-background-secondary border-r border-glass-border flex flex-col">
    <!-- 用户信息 -->
    <div class="p-4 border-b border-glass-border">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full flex items-center justify-center">
          <span class="text-white font-medium text-sm">{{ userInitials }}</span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-text-primary truncate">
            {{ user?.nickname || user?.username }}
          </p>
          <p class="text-xs text-text-secondary truncate">
            {{ user?.email }}
          </p>
        </div>
      </div>
    </div>

    <!-- 写邮件按钮 -->
    <div class="p-4">
      <Button
        variant="primary"
        size="lg"
        class="w-full"
        @click="$router.push('/compose')"
      >
        <PencilIcon class="w-4 h-4 mr-2" />
        写邮件
      </Button>
    </div>

    <!-- 邮件文件夹 -->
    <nav class="flex-1 px-4 pb-4 space-y-1">
      <router-link
        v-for="folder in folders"
        :key="folder.path"
        :to="folder.path"
        :class="[
          'flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors duration-200',
          $route.path === folder.path
            ? 'bg-primary-500/20 text-primary-400'
            : 'text-text-secondary hover:text-text-primary hover:bg-white/5'
        ]"
      >
        <component :is="folder.icon" class="w-5 h-5 mr-3" />
        <span class="flex-1">{{ folder.name }}</span>
        <span
          v-if="folder.count > 0"
          class="ml-2 px-2 py-0.5 text-xs bg-primary-500/20 text-primary-400 rounded-full"
        >
          {{ folder.count }}
        </span>
      </router-link>
    </nav>

    <!-- 底部操作 -->
    <div class="p-4 border-t border-glass-border space-y-2">
      <Button
        variant="ghost"
        size="sm"
        class="w-full justify-start"
        @click="$router.push('/settings')"
      >
        <CogIcon class="w-4 h-4 mr-2" />
        设置
      </Button>
      <Button
        variant="ghost"
        size="sm"
        class="w-full justify-start"
        @click="handleLogout"
      >
        <ArrowRightOnRectangleIcon class="w-4 h-4 mr-2" />
        退出登录
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'
import Button from '@/components/ui/Button.vue'
import {
  InboxIcon,
  PaperAirplaneIcon,
  DocumentIcon,
  TrashIcon,
  StarIcon,
  PencilIcon,
  CogIcon,
  ArrowRightOnRectangleIcon
} from '@heroicons/vue/24/outline'

const authStore = useAuthStore()
const { showNotification } = useNotification()

const user = computed(() => authStore.user)
const userInitials = computed(() => authStore.userInitials)

// 邮件文件夹配置
const folders = [
  {
    name: '收件箱',
    path: '/inbox',
    icon: InboxIcon,
    count: 12
  },
  {
    name: '已发送',
    path: '/sent',
    icon: PaperAirplaneIcon,
    count: 0
  },
  {
    name: '草稿箱',
    path: '/drafts',
    icon: DocumentIcon,
    count: 3
  },
  {
    name: '已加星标',
    path: '/starred',
    icon: StarIcon,
    count: 5
  },
  {
    name: '垃圾箱',
    path: '/trash',
    icon: TrashIcon,
    count: 0
  }
]

// 处理退出登录
const handleLogout = async () => {
  try {
    await authStore.logout()
    showNotification({
      type: 'success',
      title: '已退出登录',
      message: '感谢使用邮件系统'
    })
  } catch (error) {
    console.error('Logout error:', error)
  }
}
</script>
