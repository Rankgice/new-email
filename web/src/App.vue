<template>
  <div id="app" class="min-h-screen bg-background-primary text-text-primary">
    <!-- 全局背景效果 -->
    <div class="fixed inset-0 bg-gradient-to-br from-primary-900/20 via-background-primary to-secondary-900/20 pointer-events-none" />
    
    <!-- 路由视图 -->
    <RouterView v-slot="{ Component, route }">
      <Transition
        name="page"
        mode="out-in"
        @enter="onPageEnter"
        @leave="onPageLeave"
      >
        <component :is="Component" :key="route.path" />
      </Transition>
    </RouterView>
    
    <!-- 全局通知 -->
    <NotificationContainer />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme'
import { useMotion } from '@vueuse/motion'
import NotificationContainer from '@/components/ui/NotificationContainer.vue'

const { initTheme } = useTheme()

// 页面切换动画
const onPageEnter = (el: Element) => {
  const { apply } = useMotion(el as HTMLElement, {
    initial: { opacity: 0, y: 20 },
    enter: { 
      opacity: 1, 
      y: 0, 
      transition: { 
        duration: 300,
        ease: 'easeOut'
      } 
    }
  })
  apply('enter')
}

const onPageLeave = (el: Element) => {
  const { apply } = useMotion(el as HTMLElement, {
    leave: { 
      opacity: 0, 
      y: -20, 
      transition: { 
        duration: 300,
        ease: 'easeIn'
      } 
    }
  })
  apply('leave')
}

onMounted(() => {
  // 初始化主题
  initTheme()
  
  // 设置全局CSS变量
  document.documentElement.style.setProperty('--vh', `${window.innerHeight * 0.01}px`)
  
  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    document.documentElement.style.setProperty('--vh', `${window.innerHeight * 0.01}px`)
  })
})
</script>

<style>
#app {
  font-family: 'Inter', 'SF Pro Display', system-ui, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 页面切换动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* Firefox 滚动条 */
* {
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
}
</style>
