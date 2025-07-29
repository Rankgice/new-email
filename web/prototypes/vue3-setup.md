# 🟢 Vue 3 邮件系统项目配置

## 📦 项目初始化

### 1. 创建 Vue 3 项目
```bash
# 使用 Vite 创建项目
npm create vue@latest email-system

# 选择配置
✅ TypeScript
✅ Router
✅ Pinia
✅ ESLint
✅ Prettier
❌ Vitest
❌ Playwright
❌ Cypress
```

### 2. 安装依赖
```bash
# 核心依赖
npm install @vueuse/core @vueuse/motion @headlessui/vue
npm install @tanstack/vue-query vee-validate
npm install tailwindcss @tailwindcss/forms @tailwindcss/typography

# 开发依赖
npm install -D @types/node autoprefixer postcss
```

---

## ⚙️ 配置文件

### Vite 配置 (vite.config.ts)
```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@/components': resolve(__dirname, 'src/components'),
      '@/composables': resolve(__dirname, 'src/composables'),
      '@/stores': resolve(__dirname, 'src/stores'),
      '@/utils': resolve(__dirname, 'src/utils'),
      '@/assets': resolve(__dirname, 'src/assets')
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/assets/styles/variables.scss";`
      }
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

### Tailwind CSS 配置 (tailwind.config.js)
```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f0f4ff',
          500: '#6366f1',
          600: '#5b5bd6',
          700: '#4f46e5',
          900: '#312e81'
        },
        glass: {
          light: 'rgba(255, 255, 255, 0.05)',
          medium: 'rgba(255, 255, 255, 0.08)',
          heavy: 'rgba(255, 255, 255, 0.12)',
          border: 'rgba(255, 255, 255, 0.1)'
        }
      },
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
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
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
      },
      fontFamily: {
        sans: ['Inter', 'SF Pro Display', 'system-ui', 'sans-serif'],
        mono: ['SF Mono', 'Monaco', 'Consolas', 'monospace']
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}
```

### TypeScript 配置 (tsconfig.json)
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"],
      "@/components/*": ["src/components/*"],
      "@/composables/*": ["src/composables/*"],
      "@/stores/*": ["src/stores/*"],
      "@/utils/*": ["src/utils/*"],
      "@/assets/*": ["src/assets/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

---

## 🎨 样式系统

### 主样式文件 (src/assets/styles/main.css)
```css
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

/* CSS 变量定义 */
:root {
  /* 深色主题 */
  --color-primary: #6366f1;
  --color-secondary: #8b5cf6;
  --color-background: #0a0a0a;
  --color-surface: #1a1a1a;
  --color-text-primary: #ffffff;
  --color-text-secondary: #a0a0a0;
  --color-text-disabled: #666666;
  --color-glass-light: rgba(255, 255, 255, 0.05);
  --color-glass-medium: rgba(255, 255, 255, 0.08);
  --color-glass-heavy: rgba(255, 255, 255, 0.12);
  --color-glass-border: rgba(255, 255, 255, 0.1);
}

/* 全局样式 */
* {
  box-sizing: border-box;
}

html {
  font-family: 'Inter', 'SF Pro Display', system-ui, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  margin: 0;
  background: var(--color-background);
  color: var(--color-text-primary);
  overflow-x: hidden;
}

/* 滚动条样式 */
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

/* 毛玻璃效果基类 */
.glass {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  background: var(--color-glass-light);
  border: 1px solid var(--color-glass-border);
}

.glass-medium {
  background: var(--color-glass-medium);
}

.glass-heavy {
  background: var(--color-glass-heavy);
}

/* 动画类 */
.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

.animate-slide-up {
  animation: slideUp 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { 
    opacity: 0; 
    transform: translateY(20px); 
  }
  to { 
    opacity: 1; 
    transform: translateY(0); 
  }
}

/* 响应式工具类 */
.container-responsive {
  @apply max-w-7xl mx-auto px-4 sm:px-6 lg:px-8;
}

.grid-responsive {
  @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6;
}
```

---

## 🔧 核心 Composables

### 主题管理 (src/composables/useTheme.ts)
```ts
import { useThemeStore } from '@/stores/theme'

export const useTheme = () => {
  const themeStore = useThemeStore()
  
  const toggleTheme = () => {
    const themes = ['dark', 'ocean', 'sakura', 'forest', 'flame', 'purple']
    const currentIndex = themes.indexOf(themeStore.currentTheme)
    const nextIndex = (currentIndex + 1) % themes.length
    themeStore.setTheme(themes[nextIndex])
  }
  
  const setCustomTheme = (colors: Partial<Theme['colors']>) => {
    themeStore.setCustomTheme(colors)
  }
  
  return {
    theme: themeStore.theme,
    currentTheme: themeStore.currentTheme,
    setTheme: themeStore.setTheme,
    toggleTheme,
    setCustomTheme,
    initTheme: themeStore.initTheme
  }
}
```

### API 请求 (src/composables/useApi.ts)
```ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { api } from '@/utils/api'

export const useEmailApi = () => {
  const queryClient = useQueryClient()
  
  // 获取邮件列表
  const useEmails = (params: EmailListParams) => {
    return useQuery({
      queryKey: ['emails', params],
      queryFn: () => api.getEmails(params),
      staleTime: 5 * 60 * 1000, // 5分钟
    })
  }
  
  // 发送邮件
  const useSendEmail = () => {
    return useMutation({
      mutationFn: api.sendEmail,
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ['emails'] })
      }
    })
  }
  
  // 删除邮件
  const useDeleteEmail = () => {
    return useMutation({
      mutationFn: api.deleteEmail,
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ['emails'] })
      }
    })
  }
  
  return {
    useEmails,
    useSendEmail,
    useDeleteEmail
  }
}
```

---

## 🚀 入口文件配置

### main.ts
```ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { MotionPlugin } from '@vueuse/motion'

import App from './App.vue'
import router from './router'

import './assets/styles/main.css'

const app = createApp(App)

// 状态管理
app.use(createPinia())

// 路由
app.use(router)

// Vue Query
app.use(VueQueryPlugin, {
  queryClientConfig: {
    defaultOptions: {
      queries: {
        staleTime: 5 * 60 * 1000,
        cacheTime: 10 * 60 * 1000,
      },
    },
  },
})

// 动画插件
app.use(MotionPlugin)

app.mount('#app')
```

### App.vue
```vue
<template>
  <div id="app" class="min-h-screen bg-background text-text-primary">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme'

const { initTheme } = useTheme()

onMounted(() => {
  initTheme()
})
</script>

<style>
#app {
  font-family: 'Inter', 'SF Pro Display', system-ui, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
```

这套 Vue 3 配置确保了项目的现代化架构和优秀的开发体验！🟢✨
