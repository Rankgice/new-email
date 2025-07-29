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

        <!-- 设置页面 -->
        <div class="flex-1 grid grid-cols-1 lg:grid-cols-4 gap-6">
          <!-- 设置菜单 -->
          <div class="lg:col-span-1">
            <GlassCard padding="md">
              <nav class="space-y-1">
                <button
                  v-for="item in settingsMenu"
                  :key="item.id"
                  :class="[
                    'w-full flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors duration-200',
                    activeTab === item.id
                      ? 'bg-primary-500/20 text-primary-400'
                      : 'text-text-secondary hover:text-text-primary hover:bg-white/5'
                  ]"
                  @click="activeTab = item.id"
                >
                  <component :is="item.icon" class="w-5 h-5 mr-3" />
                  {{ item.name }}
                </button>
              </nav>
            </GlassCard>
          </div>

          <!-- 设置内容 -->
          <div class="lg:col-span-3">
            <GlassCard padding="lg" class="h-full">
              <div class="text-center py-20">
                <CogIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
                <h2 class="text-xl font-medium text-text-primary mb-2">
                  用户设置
                </h2>
                <p class="text-text-secondary">
                  此页面正在开发中...
                </p>
              </div>
            </GlassCard>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import {
  ArrowLeftIcon,
  CogIcon,
  UserIcon,
  EnvelopeIcon,
  LockClosedIcon,
  BellIcon,
  FunnelIcon,
  PencilIcon,
  KeyIcon,
  SwatchIcon
} from '@heroicons/vue/24/outline'

const activeTab = ref('profile')

const settingsMenu = [
  { id: 'profile', name: '个人资料', icon: UserIcon },
  { id: 'email', name: '邮箱管理', icon: EnvelopeIcon },
  { id: 'security', name: '安全设置', icon: LockClosedIcon },
  { id: 'notifications', name: '通知设置', icon: BellIcon },
  { id: 'filters', name: '邮件过滤', icon: FunnelIcon },
  { id: 'signature', name: '邮件签名', icon: PencilIcon },
  { id: 'api', name: 'API密钥', icon: KeyIcon },
  { id: 'theme', name: '主题设置', icon: SwatchIcon }
]
</script>
