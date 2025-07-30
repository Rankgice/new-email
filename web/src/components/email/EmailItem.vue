<template>
  <GlassCard
    :hover="true"
    padding="md"
    :class="[
      'cursor-pointer transition-all duration-200',
      {
        'border-l-4 border-l-primary-500': !email.isRead,
        'opacity-75': email.isRead
      }
    ]"
    @click="$emit('click', email)"
  >
    <div class="flex items-start space-x-3">
      <!-- 选择框 -->
      <input
        :checked="isSelected"
        type="checkbox"
        class="mt-1 w-4 h-4 text-primary-600 bg-glass-light border-glass-border rounded focus:ring-primary-500 focus:ring-2"
        @click.stop
        @change="handleSelect"
      />

      <!-- 发件人/收件人头像 -->
      <div class="flex-shrink-0">
        <div
          v-if="getDisplayName()"
          class="w-10 h-10 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full flex items-center justify-center"
        >
          <span class="text-white font-medium text-sm">
            {{ getInitials(getDisplayName()) }}
          </span>
        </div>
        <div
          v-else
          class="w-10 h-10 rounded-full flex items-center justify-center glass-medium"
        >
          <UserIcon class="w-5 h-5 text-text-secondary" />
        </div>
      </div>

      <!-- 邮件内容 -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center justify-between mb-1">
          <!-- 发件人和主题 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-2">
              <p
                :class="[
                  'text-sm font-medium truncate',
                  email.isRead ? 'text-text-secondary' : 'text-text-primary'
                ]"
              >
                {{ getDisplayName() }}
              </p>
              <p class="text-xs text-text-secondary truncate">
                {{ getDisplayEmail() }}
              </p>
              <StarIcon
                v-if="email.isStarred"
                class="w-4 h-4 text-yellow-400 flex-shrink-0"
              />
              <ExclamationTriangleIcon
                v-if="email.isImportant"
                class="w-4 h-4 text-red-400 flex-shrink-0"
              />
            </div>
            <h3
              :class="[
                'text-sm truncate mt-1',
                email.isRead ? 'text-text-secondary font-normal' : 'text-text-primary font-medium'
              ]"
            >
              {{ email.subject }}
            </h3>
          </div>

          <!-- 时间和操作 -->
          <div class="flex items-center space-x-2 ml-4">
            <span class="text-xs text-text-secondary whitespace-nowrap">
              {{ formatDate(email.createdAt) }}
            </span>
            
            <!-- 快速操作按钮 -->
            <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              <button
                class="p-1 rounded hover:bg-white/10 transition-colors duration-200"
                @click.stop="$emit('star', email)"
              >
                <StarIcon
                  :class="[
                    'w-4 h-4',
                    email.isStarred ? 'text-yellow-400' : 'text-text-secondary'
                  ]"
                />
              </button>
              <button
                class="p-1 rounded hover:bg-white/10 transition-colors duration-200"
                @click.stop="$emit('delete', email)"
              >
                <TrashIcon class="w-4 h-4 text-text-secondary hover:text-red-400" />
              </button>
            </div>
          </div>
        </div>

        <!-- 邮件预览 -->
        <p class="text-sm text-text-secondary line-clamp-2 mt-1">
          {{ email.content }}
        </p>

        <!-- 附件指示器 -->
        <div
          v-if="email.attachments && email.attachments.length > 0"
          class="flex items-center mt-2"
        >
          <PaperClipIcon class="w-4 h-4 text-text-secondary mr-1" />
          <span class="text-xs text-text-secondary">
            {{ email.attachments.length }} 个附件
          </span>
        </div>
      </div>
    </div>
  </GlassCard>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Email } from '@/types'
import GlassCard from '@/components/ui/GlassCard.vue'
import {
  UserIcon,
  StarIcon,
  ExclamationTriangleIcon,
  TrashIcon,
  PaperClipIcon
} from '@heroicons/vue/24/outline'

interface Props {
  email: Email
  selected?: boolean
  emailType?: 'inbox' | 'sent'
}

interface Emits {
  (e: 'click', email: Email): void
  (e: 'select', email: Email, selected: boolean): void
  (e: 'star', email: Email): void
  (e: 'delete', email: Email): void
}

const props = withDefaults(defineProps<Props>(), {
  selected: false,
  emailType: 'inbox'
})

const emit = defineEmits<Emits>()

const isSelected = ref(props.selected)

// 获取显示名称
const getDisplayName = () => {
  if (props.emailType === 'sent') {
    // 已发送邮件显示收件人
    const firstRecipient = props.email.toEmails?.split(',')[0]?.trim()
    return firstRecipient?.split('@')[0] || '未知收件人'
  } else {
    // 收件箱邮件显示发件人
    return props.email.fromEmail?.split('@')[0] || '未知发件人'
  }
}

// 获取显示邮箱
const getDisplayEmail = () => {
  if (props.emailType === 'sent') {
    return props.email.toEmails?.split(',')[0]?.trim() || ''
  } else {
    return props.email.fromEmail || ''
  }
}

// 获取姓名首字母
const getInitials = (name: string) => {
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = (now.getTime() - date.getTime()) / (1000 * 60 * 60)

  if (diffInHours < 24) {
    return date.toLocaleTimeString('zh-CN', { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  } else if (diffInHours < 24 * 7) {
    return date.toLocaleDateString('zh-CN', { 
      weekday: 'short' 
    })
  } else {
    return date.toLocaleDateString('zh-CN', { 
      month: 'short', 
      day: 'numeric' 
    })
  }
}

// 处理选择
const handleSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  isSelected.value = target.checked
  emit('select', props.email, target.checked)
}
</script>

<style scoped>
.glass-medium {
  background: var(--color-glass-medium);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.group:hover .group-hover\:opacity-100 {
  opacity: 1;
}
</style>
