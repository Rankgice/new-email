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

        <!-- 个人资料内容 -->
        <div class="flex-1">
          <ProfileSettings />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import ProfileSettings from '@/components/settings/ProfileSettings.vue'
import { ArrowLeftIcon } from '@heroicons/vue/24/outline'
</script>
