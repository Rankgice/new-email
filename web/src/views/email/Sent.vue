<template>
  <div class="h-screen flex bg-background-primary">
    <!-- 侧边栏 -->
    <EmailSidebar />

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部工具栏 -->
      <EmailToolbar />

      <!-- 已发送邮件列表 -->
      <div class="flex-1 overflow-hidden p-4">
        <GlassCard padding="none" class="h-full">
          <EmailListView
            title="已发送"
            email-type="sent"
            empty-title="暂无已发送邮件"
            empty-message="您还没有发送过任何邮件"
          />
        </GlassCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import EmailSidebar from '@/components/email/EmailSidebar.vue'
import EmailToolbar from '@/components/email/EmailToolbar.vue'
import EmailListView from '@/components/email/EmailListView.vue'
import GlassCard from '@/components/ui/GlassCard.vue'
</script>
