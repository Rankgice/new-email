# 🛠️ 前端技术实现指南 (Vue 3)

## 🎯 技术栈选择

### 核心框架
```
🟢 Vue 3 + TypeScript
├── Composition API
├── 类型安全
├── 响应式系统
└── SSR/SSG 支持 (Nuxt 3)

🎨 样式方案
├── Tailwind CSS (原子化CSS)
├── @vueuse/motion (动画库)
├── Headless UI Vue (无样式组件)
└── CSS Modules / SCSS

📱 响应式方案
├── CSS Grid + Flexbox
├── Container Queries
├── Viewport Units
└── Media Queries
```

### 状态管理
```
🗃️ 状态管理架构:
├── Pinia (Vue 官方状态管理)
├── VueQuery (服务端状态)
├── VeeValidate (表单验证)
└── Provide/Inject (主题状态)
```

---

## 🔮 毛玻璃效果实现

### CSS 实现
```css
/* 基础毛玻璃效果 */
.glass-card {
  position: relative;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

/* 悬浮效果 */
.glass-card:hover {
  transform: translateY(-4px);
  box-shadow: 
    0 12px 40px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
  transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
}

/* 渐变边框效果 */
.glass-card::before {
  content: '';
  position: absolute;
  inset: 0;
  padding: 1px;
  background: linear-gradient(
    135deg, 
    rgba(255, 255, 255, 0.2), 
    rgba(255, 255, 255, 0.05)
  );
  border-radius: inherit;
  mask: linear-gradient(#fff 0 0) content-box, 
        linear-gradient(#fff 0 0);
  mask-composite: exclude;
}
```

### Vue 3 组件实现
```vue
<!-- GlassCard.vue -->
<template>
  <div
    :class="[
      'relative rounded-2xl border transition-all duration-300 ease-out',
      'shadow-[0_8px_32px_rgba(0,0,0,0.3)]',
      glassLevels[level],
      {
        'hover:-translate-y-1 hover:shadow-2xl': hover
      },
      className
    ]"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
interface Props {
  className?: string;
  level?: 1 | 2 | 3; // 毛玻璃层级
  hover?: boolean;   // 是否启用悬浮效果
}

const props = withDefaults(defineProps<Props>(), {
  className: '',
  level: 1,
  hover: true
});

const glassLevels = {
  1: 'backdrop-blur-[20px] bg-white/5 border-white/10',
  2: 'backdrop-blur-[15px] bg-white/3 border-white/8',
  3: 'backdrop-blur-[25px] bg-white/8 border-white/15'
};
</script>
```

---

## 🎨 主题系统实现

### Pinia 主题 Store
```ts
// stores/theme.ts
import { defineStore } from 'pinia'

interface Theme {
  name: string;
  colors: {
    primary: string;
    secondary: string;
    background: string;
    surface: string;
    text: {
      primary: string;
      secondary: string;
      disabled: string;
    };
    glass: {
      light: string;
      medium: string;
      heavy: string;
      border: string;
    };
  };
}

const themes: Record<string, Theme> = {
  dark: {
    name: 'dark',
    colors: {
      primary: '#6366f1',
      secondary: '#8b5cf6',
      background: '#0a0a0a',
      surface: '#1a1a1a',
      text: {
        primary: '#ffffff',
        secondary: '#a0a0a0',
        disabled: '#666666'
      },
      glass: {
        light: 'rgba(255,255,255,0.05)',
        medium: 'rgba(255,255,255,0.08)',
        heavy: 'rgba(255,255,255,0.12)',
        border: 'rgba(255,255,255,0.1)'
      }
    }
  },
  ocean: {
    name: 'ocean',
    colors: {
      primary: '#0ea5e9',
      secondary: '#0284c7',
      background: '#0c1426',
      // ... 其他颜色
    }
  }
};

export const useThemeStore = defineStore('theme', () => {
  const currentTheme = ref('dark')

  const theme = computed(() => themes[currentTheme.value])

  const setTheme = (themeName: string) => {
    currentTheme.value = themeName

    // 更新 CSS 变量
    const themeData = themes[themeName]
    const root = document.documentElement

    Object.entries(themeData.colors).forEach(([key, value]) => {
      if (typeof value === 'string') {
        root.style.setProperty(`--color-${key}`, value)
      } else {
        Object.entries(value).forEach(([subKey, subValue]) => {
          root.style.setProperty(`--color-${key}-${subKey}`, subValue)
        })
      }
    })

    // 持久化到 localStorage
    localStorage.setItem('theme', themeName)
  }

  // 初始化主题
  const initTheme = () => {
    const savedTheme = localStorage.getItem('theme') || 'dark'
    setTheme(savedTheme)
  }

  return {
    currentTheme: readonly(currentTheme),
    theme,
    setTheme,
    initTheme
  }
})
```

### 主题 Composable
```ts
// composables/useTheme.ts
export const useTheme = () => {
  const themeStore = useThemeStore()

  return {
    theme: themeStore.theme,
    currentTheme: themeStore.currentTheme,
    setTheme: themeStore.setTheme,
    initTheme: themeStore.initTheme
  }
}
```

---

## ⚡ 动画系统实现

### @vueuse/motion 动画配置
```vue
<!-- AnimatedCard.vue -->
<template>
  <div
    v-motion
    :initial="{ opacity: 0, scale: 0.95, y: 20 }"
    :enter="{ opacity: 1, scale: 1, y: 0 }"
    :hovered="{ y: -4 }"
    :transition="{ duration: 200, ease: 'easeOut' }"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
// 无需额外逻辑，@vueuse/motion 自动处理
</script>
```

```vue
<!-- PageTransition.vue -->
<template>
  <Transition
    name="page"
    mode="out-in"
    @enter="onEnter"
    @leave="onLeave"
  >
    <slot />
  </Transition>
</template>

<script setup lang="ts">
import { useMotion } from '@vueuse/motion'

const onEnter = (el: Element) => {
  const { apply } = useMotion(el as HTMLElement, {
    initial: { opacity: 0, y: 20 },
    enter: { opacity: 1, y: 0, transition: { duration: 300 } }
  })
  apply('enter')
}

const onLeave = (el: Element) => {
  const { apply } = useMotion(el as HTMLElement, {
    leave: { opacity: 0, y: -20, transition: { duration: 300 } }
  })
  apply('leave')
}
</script>

<style scoped>
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
</style>
```

### 动画 Composable
```ts
// composables/useAnimations.ts
export const useAnimations = () => {
  const cardEntry = {
    initial: { opacity: 0, scale: 0.95, y: 20 },
    enter: { opacity: 1, scale: 1, y: 0 },
    transition: { duration: 200, ease: 'easeOut' }
  }

  const listItem = {
    initial: { opacity: 0, x: -20 },
    enter: { opacity: 1, x: 0 },
    leave: { opacity: 0, x: 20, height: 0 },
    transition: { duration: 200 }
  }

  const buttonHover = {
    hovered: { scale: 1.02 },
    tapped: { scale: 0.98 }
  }

  return {
    cardEntry,
    listItem,
    buttonHover
  }
}
```

---

## 📱 响应式实现

### Tailwind CSS 配置
```js
// tailwind.config.js
module.exports = {
  theme: {
    screens: {
      'xs': '475px',
      'sm': '640px',
      'md': '768px',
      'lg': '1024px',
      'xl': '1280px',
      '2xl': '1536px',
    },
    extend: {
      backdropBlur: {
        'xs': '2px',
        'sm': '4px',
        'md': '8px',
        'lg': '12px',
        'xl': '16px',
        '2xl': '24px',
        '3xl': '32px',
      },
      animation: {
        'glass-shimmer': 'shimmer 2s linear infinite',
        'float': 'float 3s ease-in-out infinite',
      },
      keyframes: {
        shimmer: {
          '0%': { transform: 'translateX(-100%)' },
          '100%': { transform: 'translateX(100%)' }
        },
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' }
        }
      }
    }
  }
}
```

### 响应式 Composable
```ts
// composables/useResponsive.ts
import { useWindowSize } from '@vueuse/core'

export const useResponsive = () => {
  const { width, height } = useWindowSize()

  const isMobile = computed(() => width.value < 768)
  const isTablet = computed(() => width.value >= 768 && width.value < 1024)
  const isDesktop = computed(() => width.value >= 1024)

  const breakpoint = computed(() => {
    if (width.value < 576) return 'xs'
    if (width.value < 768) return 'sm'
    if (width.value < 1024) return 'md'
    if (width.value < 1280) return 'lg'
    return 'xl'
  })

  return {
    width: readonly(width),
    height: readonly(height),
    isMobile,
    isTablet,
    isDesktop,
    breakpoint
  }
}
```

```vue
<!-- ResponsiveLayout.vue -->
<template>
  <div :class="layoutClass">
    <!-- 移动端布局 -->
    <template v-if="isMobile">
      <main class="flex-1 overflow-hidden">
        <slot />
      </main>
      <MobileNavigation />
    </template>

    <!-- 平板端布局 -->
    <template v-else-if="isTablet">
      <Sidebar />
      <main class="overflow-hidden">
        <slot />
      </main>
    </template>

    <!-- 桌面端布局 -->
    <template v-else>
      <Sidebar />
      <main class="overflow-hidden">
        <slot />
      </main>
      <PreviewPanel />
    </template>
  </div>
</template>

<script setup lang="ts">
const { isMobile, isTablet, isDesktop } = useResponsive()

const layoutClass = computed(() => {
  if (isMobile.value) return 'flex flex-col h-screen'
  if (isTablet.value) return 'grid grid-cols-[250px_1fr] h-screen'
  return 'grid grid-cols-[250px_1fr_300px] h-screen'
})
</script>
```

---

## 🎪 交互组件实现

### 手势支持
```ts
// composables/useSwipeGesture.ts
import { useSwipe } from '@vueuse/gesture'

export const useSwipeGesture = (
  onSwipeLeft?: () => void,
  onSwipeRight?: () => void
) => {
  const target = ref<HTMLElement>()

  const { lengthX, direction } = useSwipe(target, {
    threshold: 50,
    onSwipeEnd() {
      if (direction.value === 'left' && onSwipeLeft) {
        onSwipeLeft()
      } else if (direction.value === 'right' && onSwipeRight) {
        onSwipeRight()
      }
    }
  })

  return {
    target,
    lengthX,
    direction
  }
}
```

```vue
<!-- SwipeableEmailItem.vue -->
<template>
  <div
    ref="swipeTarget"
    class="relative overflow-hidden"
    v-motion
    :style="{ transform: `translateX(${swipeOffset}px)` }"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
  >
    <GlassCard class="p-4 relative z-10">
      <EmailContent :email="email" />
    </GlassCard>

    <!-- 左滑删除背景 -->
    <div
      class="absolute inset-y-0 right-0 w-20 bg-red-500 flex items-center justify-center"
      :class="{ 'opacity-100': swipeOffset < -20, 'opacity-0': swipeOffset >= -20 }"
    >
      <TrashIcon class="w-6 h-6 text-white" />
    </div>

    <!-- 右滑标记背景 -->
    <div
      class="absolute inset-y-0 left-0 w-20 bg-green-500 flex items-center justify-center"
      :class="{ 'opacity-100': swipeOffset > 20, 'opacity-0': swipeOffset <= 20 }"
    >
      <CheckIcon class="w-6 h-6 text-white" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  email: Email
}

interface Emits {
  (e: 'delete'): void
  (e: 'markRead'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const swipeTarget = ref<HTMLElement>()
const swipeOffset = ref(0)
const startX = ref(0)
const isDragging = ref(false)

const onTouchStart = (e: TouchEvent) => {
  startX.value = e.touches[0].clientX
  isDragging.value = true
}

const onTouchMove = (e: TouchEvent) => {
  if (!isDragging.value) return

  const currentX = e.touches[0].clientX
  const diff = currentX - startX.value

  // 限制滑动距离
  swipeOffset.value = Math.max(-100, Math.min(100, diff))
}

const onTouchEnd = () => {
  isDragging.value = false

  if (swipeOffset.value < -50) {
    emit('delete')
  } else if (swipeOffset.value > 50) {
    emit('markRead')
  }

  // 重置位置
  swipeOffset.value = 0
}
</script>
```

---

## 🚀 性能优化

### 虚拟滚动
```vue
<!-- VirtualEmailList.vue -->
<template>
  <div
    ref="containerRef"
    class="overflow-auto"
    :style="{ height: `${containerHeight}px` }"
    @scroll="onScroll"
  >
    <div
      :style="{
        height: `${totalHeight}px`,
        position: 'relative'
      }"
    >
      <div
        v-for="(email, index) in visibleEmails"
        :key="email.id"
        :style="{
          position: 'absolute',
          top: `${(visibleStart + index) * itemHeight}px`,
          height: `${itemHeight}px`,
          width: '100%'
        }"
      >
        <EmailItem :email="email" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  emails: Email[]
  itemHeight: number
  containerHeight?: number
}

const props = withDefaults(defineProps<Props>(), {
  containerHeight: 600
})

const containerRef = ref<HTMLElement>()
const scrollTop = ref(0)

const totalHeight = computed(() => props.emails.length * props.itemHeight)

const visibleStart = computed(() =>
  Math.floor(scrollTop.value / props.itemHeight)
)

const visibleEnd = computed(() =>
  Math.min(
    visibleStart.value + Math.ceil(props.containerHeight / props.itemHeight) + 1,
    props.emails.length
  )
)

const visibleEmails = computed(() =>
  props.emails.slice(visibleStart.value, visibleEnd.value)
)

const onScroll = (e: Event) => {
  const target = e.target as HTMLElement
  scrollTop.value = target.scrollTop
}
</script>
```

### 懒加载和预加载
```vue
<!-- LazyImage.vue -->
<template>
  <div ref="containerRef" :class="['relative', className]">
    <img
      v-if="isInView"
      :src="src"
      :alt="alt"
      :class="[
        'transition-opacity duration-300',
        { 'opacity-100': isLoaded, 'opacity-0': !isLoaded }
      ]"
      @load="isLoaded = true"
    />
    <div
      v-if="!isLoaded"
      class="absolute inset-0 bg-gray-200 animate-pulse rounded"
    />
  </div>
</template>

<script setup lang="ts">
import { useIntersectionObserver } from '@vueuse/core'

interface Props {
  src: string
  alt: string
  className?: string
}

const props = withDefaults(defineProps<Props>(), {
  className: ''
})

const containerRef = ref<HTMLElement>()
const isLoaded = ref(false)
const isInView = ref(false)

// 使用 VueUse 的 Intersection Observer
useIntersectionObserver(
  containerRef,
  ([{ isIntersecting }]) => {
    if (isIntersecting) {
      isInView.value = true
    }
  },
  { threshold: 0.1 }
)
</script>
```

## 🚀 Vue 3 项目结构

### 推荐项目结构
```
📁 src/
├── 📁 components/          # 通用组件
│   ├── 📁 ui/             # UI 基础组件
│   │   ├── GlassCard.vue
│   │   ├── Button.vue
│   │   └── Input.vue
│   ├── 📁 email/          # 邮件相关组件
│   │   ├── EmailItem.vue
│   │   ├── EmailList.vue
│   │   └── EmailComposer.vue
│   └── 📁 layout/         # 布局组件
│       ├── Sidebar.vue
│       ├── Header.vue
│       └── ResponsiveLayout.vue
├── 📁 composables/        # 组合式函数
│   ├── useTheme.ts
│   ├── useResponsive.ts
│   ├── useAnimations.ts
│   └── useSwipeGesture.ts
├── 📁 stores/             # Pinia 状态管理
│   ├── theme.ts
│   ├── auth.ts
│   ├── email.ts
│   └── user.ts
├── 📁 views/              # 页面组件
│   ├── Login.vue
│   ├── Inbox.vue
│   ├── Compose.vue
│   └── Settings.vue
├── 📁 assets/             # 静态资源
│   ├── 📁 styles/
│   │   ├── main.css
│   │   └── themes.css
│   └── 📁 images/
├── 📁 utils/              # 工具函数
│   ├── api.ts
│   ├── helpers.ts
│   └── constants.ts
└── main.ts                # 入口文件
```

### 主要依赖包
```json
{
  "dependencies": {
    "vue": "^3.3.0",
    "@vue/typescript": "^1.8.0",
    "vue-router": "^4.2.0",
    "pinia": "^2.1.0",
    "@vueuse/core": "^10.0.0",
    "@vueuse/motion": "^2.0.0",
    "@vueuse/gesture": "^2.0.0",
    "@headlessui/vue": "^1.7.0",
    "tailwindcss": "^3.3.0",
    "vee-validate": "^4.9.0",
    "@tanstack/vue-query": "^4.29.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^4.2.0",
    "typescript": "^5.0.0",
    "vite": "^4.3.0",
    "@types/node": "^20.0.0"
  }
}
```

这套基于 Vue 3 的技术实现指南确保了设计系统的完美落地和优秀的用户体验！🚀✨
