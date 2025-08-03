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
      '@/assets': resolve(__dirname, 'src/assets'),
      '@/views': resolve(__dirname, 'src/views'),
      '@/types': resolve(__dirname, 'src/types')
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
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          ui: ['@headlessui/vue', '@heroicons/vue'],
          utils: ['@vueuse/core', '@vueuse/motion']
        }
      }
    }
  }
})
