<template>
  <div class="h-screen bg-background-primary">
    <div class="container-responsive h-full py-6">
      <div class="h-full flex flex-col">
        <!-- 返回按钮 -->
        <div class="mb-4">
          <Button
            variant="ghost"
            @click="$router.back()"
          >
            <ArrowLeftIcon class="w-4 h-4 mr-2" />
            返回
          </Button>
        </div>

        <!-- 邮件详情 -->
        <GlassCard padding="lg" class="flex-1">
          <div class="text-center py-20">
            <EnvelopeIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
            <h2 class="text-xl font-medium text-text-primary mb-2">
              邮件详情页面
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
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import {
  ArrowLeftIcon,
  EnvelopeIcon
} from '@heroicons/vue/24/outline'
</script>
