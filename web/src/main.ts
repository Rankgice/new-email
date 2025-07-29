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

// Vue Query - 服务端状态管理
app.use(VueQueryPlugin, {
  queryClientConfig: {
    defaultOptions: {
      queries: {
        staleTime: 5 * 60 * 1000, // 5分钟
        cacheTime: 10 * 60 * 1000, // 10分钟
        retry: 3,
        refetchOnWindowFocus: false,
      },
      mutations: {
        retry: 1,
      },
    },
  },
})

// 动画插件
app.use(MotionPlugin)

// 全局错误处理
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue Error:', err, info)
}

// 移除加载屏幕
const removeLoadingScreen = () => {
  const loadingScreen = document.querySelector('.loading-screen')
  if (loadingScreen) {
    loadingScreen.style.opacity = '0'
    setTimeout(() => {
      loadingScreen.remove()
    }, 300)
  }
}

app.mount('#app')

// 应用挂载后移除加载屏幕
setTimeout(removeLoadingScreen, 1000)
