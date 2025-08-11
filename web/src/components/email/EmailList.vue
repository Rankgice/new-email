<template>
  <div class="h-full overflow-auto">
    <!-- 邮件列表 -->
    <div class="p-4 space-y-2">
      <EmailItem
        v-for="email in emails"
        :key="email.id"
        :email="email"
        @click="handleEmailClick"
        @select="handleEmailSelect"
        @star="handleEmailStar"
        @delete="handleEmailDelete"
      />
    </div>

    <!-- 空状态 -->
    <div
      v-if="emails.length === 0"
      class="flex flex-col items-center justify-center h-64 text-center"
    >
      <InboxIcon class="w-16 h-16 text-text-secondary mb-4" />
      <h3 class="text-lg font-medium text-text-primary mb-2">
        收件箱为空
      </h3>
      <p class="text-text-secondary">
        您的收件箱中暂时没有邮件
      </p>
    </div>

    <!-- 加载更多 -->
    <div
      v-if="hasMore"
      class="p-4 text-center"
    >
      <Button
        variant="ghost"
        :loading="isLoadingMore"
        @click="loadMore"
      >
        加载更多
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Email } from '@/types'
import EmailItem from './EmailItem.vue'
import Button from '@/components/ui/Button.vue'
import { InboxIcon } from '@heroicons/vue/24/outline'

// 邮件数据（不使用模拟数据）
const emails = ref<Email[]>([])

const hasMore = ref(true)
const isLoadingMore = ref(false)

// 处理邮件点击
const handleEmailClick = (email: Email) => {
  console.log('Email clicked:', email.id)
  // TODO: 跳转到邮件详情页面
}

// 处理邮件选择
const handleEmailSelect = (email: Email, selected: boolean) => {
  console.log('Email selected:', email.id, selected)
  // TODO: 更新选择状态
}

// 处理邮件加星
const handleEmailStar = (email: Email) => {
  console.log('Email starred:', email.id)
  // TODO: 更新加星状态
  email.isStarred = !email.isStarred
}

// 处理邮件删除
const handleEmailDelete = (email: Email) => {
  console.log('Email deleted:', email.id)
  // TODO: 删除邮件
  const index = emails.value.findIndex(e => e.id === email.id)
  if (index > -1) {
    emails.value.splice(index, 1)
  }
}

// 加载更多邮件
const loadMore = async () => {
  isLoadingMore.value = true
  // TODO: 实现加载更多逻辑
  setTimeout(() => {
    isLoadingMore.value = false
    hasMore.value = false
  }, 1000)
}

onMounted(() => {
  // TODO: 加载邮件数据
})
</script>
