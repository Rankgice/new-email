<template>
  <div class="h-screen flex bg-background-primary">
    <!-- 侧边栏 -->
    <EmailSidebar />

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部工具栏 -->
      <EmailToolbar />

      <!-- 邮件列表 -->
      <div class="flex-1 overflow-hidden">
        <EmailList />
      </div>
    </div>

    <!-- 邮件预览面板 (桌面端) -->
    <EmailPreview v-if="!isMobile" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useResponsive } from '@/composables/useResponsive'
import { useAuthStore } from '@/stores/auth'
import EmailSidebar from '@/components/email/EmailSidebar.vue'
import EmailToolbar from '@/components/email/EmailToolbar.vue'
import EmailList from '@/components/email/EmailList.vue'
import EmailPreview from '@/components/email/EmailPreview.vue'

const { isMobile } = useResponsive()
const authStore = useAuthStore()

onMounted(() => {
  // 确保用户已登录
  if (!authStore.isAuthenticated) {
    authStore.setRedirectPath('/inbox')
  }
})
</script>
