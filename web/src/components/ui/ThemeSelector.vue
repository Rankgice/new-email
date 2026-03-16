<template>
  <div class="relative">
    <!-- 主题切换按钮 -->
    <button
      class="p-3 glass-card rounded-full transition-all duration-200 group hover:glass-medium"
      @click="toggleDropdown"
    >
      <SwatchIcon class="w-5 h-5 text-text-secondary group-hover:text-text-primary transition-colors duration-200" />
    </button>

    <!-- 主题选择下拉菜单 -->
    <Transition
      name="dropdown"
      @enter="onEnter"
      @leave="onLeave"
    >
      <div
        v-if="isOpen"
        class="absolute bottom-full mb-2 right-0 w-64 glass-card rounded-xl p-4 shadow-glass-lg"
        @click.stop
      >
        <!-- 标题 -->
        <h3 class="text-sm font-medium text-text-primary mb-3">
          🎨 选择主题
        </h3>

        <!-- 主题网格 -->
        <div class="grid grid-cols-2 gap-2 mb-4">
          <button
            v-for="theme in availableThemes"
            :key="theme.id"
            :class="[
              'relative p-3 rounded-lg border-2 transition-all duration-200 group',
              currentTheme.id === theme.id
                ? 'border-primary-500 bg-primary-500/10'
                : 'border-glass-border hover:border-primary-400/50 hover:bg-white/5'
            ]"
            @click="setTheme(theme.id)"
          >
            <!-- 主题预览色块 -->
            <div class="flex space-x-1 mb-2">
              <div
                class="w-3 h-3 rounded-full"
                :style="{ backgroundColor: theme.colors.primary }"
              />
              <div
                class="w-3 h-3 rounded-full"
                :style="{ backgroundColor: theme.colors.secondary }"
              />
              <div
                class="w-3 h-3 rounded-full"
                :style="{ backgroundColor: theme.colors.status.success }"
              />
            </div>

            <!-- 主题名称 -->
            <div class="text-xs text-text-secondary group-hover:text-text-primary transition-colors duration-200">
              {{ theme.displayName }}
            </div>

            <!-- 选中指示器 -->
            <div
              v-if="currentTheme.id === theme.id"
              class="absolute top-1 right-1"
            >
              <CheckIcon class="w-3 h-3 text-primary-500" />
            </div>
          </button>
        </div>

        <!-- 快速切换按钮 -->
        <div class="flex space-x-2">
          <button
            class="flex-1 px-3 py-2 text-xs font-medium rounded-lg glass-card transition-all duration-200 hover:glass-medium"
            @click="nextTheme"
          >
            🔄 下一个
          </button>
          <button
            class="flex-1 px-3 py-2 text-xs font-medium rounded-lg glass-card transition-all duration-200 hover:glass-medium"
            @click="resetTheme"
          >
            🔙 重置
          </button>
        </div>

        <!-- 自定义主题提示 -->
        <div class="mt-3 pt-3 border-t border-glass-border">
          <p class="text-xs text-text-secondary text-center">
            💡 更多自定义选项请前往设置页面
          </p>
        </div>
      </div>
    </Transition>

    <!-- 点击外部关闭 -->
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[-1]"
      @click="closeDropdown"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useTheme } from '@/composables/useTheme'
import { useMotion } from '@vueuse/motion'
import {
  SwatchIcon,
  CheckIcon
} from '@heroicons/vue/24/outline'

const {
  currentTheme,
  availableThemes,
  setTheme,
  nextTheme,
  resetTheme
} = useTheme()

// 下拉菜单状态
const isOpen = ref(false)

// 切换下拉菜单
const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

// 关闭下拉菜单
const closeDropdown = () => {
  isOpen.value = false
}

// 进入动画
const onEnter = (el: Element) => {
  const { apply } = useMotion(el as HTMLElement, {
    initial: { opacity: 0, scale: 0.95, y: 10 },
    enter: { 
      opacity: 1, 
      scale: 1, 
      y: 0,
      transition: { 
        duration: 200,
        ease: 'easeOut'
      } 
    }
  })
  apply('enter')
}

// 离开动画
const onLeave = (el: Element) => {
  const { apply } = useMotion(el as HTMLElement, {
    leave: { 
      opacity: 0, 
      scale: 0.95, 
      y: 10,
      transition: { 
        duration: 150,
        ease: 'easeIn'
      } 
    }
  })
  apply('leave')
}
</script>

<style scoped>
.glass-card {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  background: var(--color-glass-light);
  border: 1px solid var(--color-glass-border);
}

.glass-medium {
  background: var(--color-glass-medium);
}

/* 下拉菜单动画 */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(10px);
}
</style>
