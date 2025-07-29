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
        <GlassCard padding="lg" class="h-full">
          <div class="text-center py-20">
            <PaperAirplaneIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
            <h2 class="text-xl font-medium text-text-primary mb-2">
              已发送邮件
            </h2>
            <p class="text-text-secondary">
              此页面正在开发中...
            </p>
          </div>
        </GlassCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import EmailSidebar from '@/components/email/EmailSidebar.vue'
import EmailToolbar from '@/components/email/EmailToolbar.vue'
import GlassCard from '@/components/ui/GlassCard.vue'
import { PaperAirplaneIcon } from '@heroicons/vue/24/outline'
</script>
