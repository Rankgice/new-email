<template>
  <div class="h-screen bg-background-primary">
    <div class="container-responsive h-full py-6">
      <GlassCard padding="lg" class="h-full">
        <div class="text-center py-20">
          <EnvelopeIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
          <h2 class="text-xl font-medium text-text-primary mb-2">
            邮件监控
          </h2>
          <p class="text-text-secondary">
            此页面正在开发中...
          </p>
        </div>
      </GlassCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import GlassCard from '@/components/ui/GlassCard.vue'
import { EnvelopeIcon } from '@heroicons/vue/24/outline'
</script>
