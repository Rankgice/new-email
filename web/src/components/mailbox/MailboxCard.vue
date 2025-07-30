<template>
  <GlassCard :level="2" hover padding="lg" class="transition-all duration-200">
    <div class="flex items-start justify-between">
      <!-- 邮箱信息 -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center space-x-3 mb-3">
          <!-- 邮箱图标 -->
          <div class="flex-shrink-0">
            <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-primary-500 to-secondary-500 flex items-center justify-center">
              <EnvelopeIcon class="w-5 h-5 text-white" />
            </div>
          </div>
          
          <!-- 邮箱地址和状态 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2">
              <h3 class="text-lg font-semibold text-text-primary truncate">
                {{ mailbox.email }}
              </h3>
              <StatusBadge :status="mailbox.status" />
            </div>
            <div class="flex items-center space-x-4 mt-1 text-sm text-text-secondary">
              <span class="flex items-center">
                <TagIcon class="w-4 h-4 mr-1" />
                {{ typeLabel }}
              </span>
              <span v-if="mailbox.provider" class="flex items-center">
                <ServerIcon class="w-4 h-4 mr-1" />
                {{ providerLabel }}
              </span>
              <span v-if="mailbox.lastSyncAt" class="flex items-center">
                <ClockIcon class="w-4 h-4 mr-1" />
                {{ formatDate(mailbox.lastSyncAt) }}
              </span>
            </div>
          </div>
        </div>

        <!-- 配置信息 -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <!-- IMAP 配置 -->
          <div class="bg-glass-light rounded-lg p-3">
            <h4 class="text-sm font-medium text-text-primary mb-2 flex items-center">
              <ArrowDownTrayIcon class="w-4 h-4 mr-1" />
              IMAP 接收
            </h4>
            <div class="space-y-1 text-xs text-text-secondary">
              <div>{{ mailbox.imapHost }}:{{ mailbox.imapPort }}</div>
              <div class="flex items-center">
                <ShieldCheckIcon class="w-3 h-3 mr-1" />
                {{ mailbox.imapSsl ? 'SSL/TLS' : '无加密' }}
              </div>
            </div>
          </div>

          <!-- SMTP 配置 -->
          <div class="bg-glass-light rounded-lg p-3">
            <h4 class="text-sm font-medium text-text-primary mb-2 flex items-center">
              <ArrowUpTrayIcon class="w-4 h-4 mr-1" />
              SMTP 发送
            </h4>
            <div class="space-y-1 text-xs text-text-secondary">
              <div>{{ mailbox.smtpHost }}:{{ mailbox.smtpPort }}</div>
              <div class="flex items-center">
                <ShieldCheckIcon class="w-3 h-3 mr-1" />
                {{ mailbox.smtpSsl ? 'SSL/TLS' : '无加密' }}
              </div>
            </div>
          </div>
        </div>

        <!-- 功能选项 -->
        <div class="flex items-center space-x-4 text-sm">
          <div class="flex items-center">
            <input
              type="checkbox"
              :checked="mailbox.autoReceive"
              disabled
              class="w-4 h-4 text-primary-600 bg-glass-light border-glass-border rounded focus:ring-primary-500"
            />
            <label class="ml-2 text-text-secondary">自动收信</label>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex-shrink-0 ml-4">
        <div class="flex items-center space-x-2">
          <!-- 同步按钮 -->
          <Button
            variant="ghost"
            size="sm"
            @click="$emit('sync', mailbox)"
            :loading="syncLoading"
            class="text-primary-400 hover:text-primary-300"
            title="同步邮件"
          >
            <ArrowPathIcon class="w-4 h-4" />
          </Button>

          <!-- 测试连接按钮 -->
          <Button
            variant="ghost"
            size="sm"
            @click="$emit('test', mailbox)"
            :loading="testLoading"
            class="text-blue-400 hover:text-blue-300"
            title="测试连接"
          >
            <WifiIcon class="w-4 h-4" />
          </Button>

          <!-- 更多操作菜单 -->
          <div class="relative" ref="menuRef">
            <Button
              variant="ghost"
              size="sm"
              @click="showMenu = !showMenu"
              class="text-text-secondary hover:text-text-primary"
            >
              <EllipsisVerticalIcon class="w-4 h-4" />
            </Button>

            <!-- 下拉菜单 -->
            <Transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div
                v-if="showMenu"
                class="absolute right-0 top-full mt-1 w-48 bg-glass-medium backdrop-blur-md border border-glass-border rounded-lg shadow-lg z-10"
              >
                <div class="py-1">
                  <button
                    @click="handleEdit"
                    class="flex items-center w-full px-4 py-2 text-sm text-text-primary hover:bg-glass-light transition-colors"
                  >
                    <PencilIcon class="w-4 h-4 mr-3" />
                    编辑配置
                  </button>
                  <button
                    @click="handleToggleStatus"
                    class="flex items-center w-full px-4 py-2 text-sm text-text-primary hover:bg-glass-light transition-colors"
                  >
                    <PowerIcon class="w-4 h-4 mr-3" />
                    {{ mailbox.status === 1 ? '禁用' : '启用' }}
                  </button>
                  <div class="border-t border-glass-border my-1"></div>
                  <button
                    @click="handleDelete"
                    class="flex items-center w-full px-4 py-2 text-sm text-red-400 hover:bg-glass-light transition-colors"
                  >
                    <TrashIcon class="w-4 h-4 mr-3" />
                    删除邮箱
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </div>
  </GlassCard>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { Mailbox } from '@/types'

// Icons
import {
  EnvelopeIcon,
  TagIcon,
  ServerIcon,
  ClockIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  ShieldCheckIcon,
  ArrowPathIcon,
  WifiIcon,
  EllipsisVerticalIcon,
  PencilIcon,
  PowerIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'

// Components
import Button from '@/components/ui/Button.vue'
import GlassCard from '@/components/ui/GlassCard.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'

// Props
interface Props {
  mailbox: Mailbox
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  edit: [mailbox: Mailbox]
  delete: [mailbox: Mailbox]
  sync: [mailbox: Mailbox]
  test: [mailbox: Mailbox]
  toggleStatus: [mailbox: Mailbox]
}>()

// 响应式数据
const showMenu = ref(false)
const menuRef = ref<HTMLElement>()
const syncLoading = ref(false)
const testLoading = ref(false)

// 计算属性
const typeLabel = computed(() => {
  return props.mailbox.type === 'self' ? '自建邮箱' : '第三方邮箱'
})

const providerLabel = computed(() => {
  const providerMap: Record<string, string> = {
    gmail: 'Gmail',
    outlook: 'Outlook',
    qq: 'QQ邮箱',
    163: '163邮箱',
    imap: '自定义IMAP'
  }
  return providerMap[props.mailbox.provider] || props.mailbox.provider
})

// 方法
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    return '今天'
  } else if (days === 1) {
    return '昨天'
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString()
  }
}

const handleEdit = () => {
  showMenu.value = false
  emit('edit', props.mailbox)
}

const handleDelete = () => {
  showMenu.value = false
  emit('delete', props.mailbox)
}

const handleToggleStatus = () => {
  showMenu.value = false
  emit('toggleStatus', props.mailbox)
}

const handleClickOutside = (event: Event) => {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    showMenu.value = false
  }
}

// 生命周期
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
