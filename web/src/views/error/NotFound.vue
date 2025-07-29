<template>
  <div class="min-h-screen flex items-center justify-center bg-background-primary">
    <div class="text-center">
      <GlassCard padding="lg" class="max-w-md mx-auto">
        <!-- 404 图标 -->
        <div class="mb-6">
          <div class="w-24 h-24 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full flex items-center justify-center mx-auto mb-4">
            <span class="text-4xl font-bold text-white">404</span>
          </div>
        </div>

        <!-- 错误信息 -->
        <h1 class="text-2xl font-bold text-text-primary mb-2">
          页面未找到
        </h1>
        <p class="text-text-secondary mb-6">
          抱歉，您访问的页面不存在或已被移动
        </p>

        <!-- 操作按钮 -->
        <div class="space-y-3">
          <Button
            variant="primary"
            class="w-full"
            @click="$router.push('/inbox')"
          >
            <HomeIcon class="w-4 h-4 mr-2" />
            返回首页
          </Button>
          <Button
            variant="ghost"
            class="w-full"
            @click="$router.back()"
          >
            <ArrowLeftIcon class="w-4 h-4 mr-2" />
            返回上页
          </Button>
        </div>

        <!-- 帮助链接 -->
        <div class="mt-6 pt-4 border-t border-glass-border">
          <p class="text-sm text-text-secondary">
            需要帮助？
            <a href="#" class="text-primary-400 hover:text-primary-300">
              联系客服
            </a>
          </p>
        </div>
      </GlassCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import {
  HomeIcon,
  ArrowLeftIcon
} from '@heroicons/vue/24/outline'
</script>
