<template>
  <nav class="flex items-center justify-between" aria-label="分页导航">
    <!-- 页面信息 -->
    <div class="flex-1 flex justify-between sm:hidden">
      <Button
        variant="ghost"
        size="sm"
        :disabled="currentPage <= 1"
        @click="$emit('pageChange', currentPage - 1)"
      >
        上一页
      </Button>
      <span class="text-sm text-text-secondary">
        第 {{ currentPage }} 页，共 {{ totalPages }} 页
      </span>
      <Button
        variant="ghost"
        size="sm"
        :disabled="currentPage >= totalPages"
        @click="$emit('pageChange', currentPage + 1)"
      >
        下一页
      </Button>
    </div>

    <!-- 桌面端分页 -->
    <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
      <div>
        <p class="text-sm text-text-secondary">
          第 <span class="font-medium">{{ currentPage }}</span> 页，共 <span class="font-medium">{{ totalPages }}</span> 页
        </p>
      </div>
      <div>
        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="分页">
          <!-- 上一页按钮 -->
          <button
            :disabled="currentPage <= 1"
            @click="$emit('pageChange', currentPage - 1)"
            class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-glass-border bg-glass-light text-sm font-medium text-text-secondary hover:bg-glass-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span class="sr-only">上一页</span>
            <ChevronLeftIcon class="h-5 w-5" />
          </button>

          <!-- 页码按钮 -->
          <template v-for="page in visiblePages" :key="page">
            <button
              v-if="page !== '...'"
              :class="[
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium transition-colors',
                page === currentPage
                  ? 'z-10 bg-primary-500 border-primary-500 text-white'
                  : 'bg-glass-light border-glass-border text-text-secondary hover:bg-glass-medium'
              ]"
              @click="$emit('pageChange', page)"
            >
              {{ page }}
            </button>
            <span
              v-else
              class="relative inline-flex items-center px-4 py-2 border border-glass-border bg-glass-light text-sm font-medium text-text-secondary"
            >
              ...
            </span>
          </template>

          <!-- 下一页按钮 -->
          <button
            :disabled="currentPage >= totalPages"
            @click="$emit('pageChange', currentPage + 1)"
            class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-glass-border bg-glass-light text-sm font-medium text-text-secondary hover:bg-glass-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <span class="sr-only">下一页</span>
            <ChevronRightIcon class="h-5 w-5" />
          </button>
        </nav>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
import Button from './Button.vue'

// Props
interface Props {
  currentPage: number
  totalPages: number
  maxVisible?: number
}

const props = withDefaults(defineProps<Props>(), {
  maxVisible: 7
})

// Emits
const emit = defineEmits<{
  pageChange: [page: number]
}>()

// 计算可见的页码
const visiblePages = computed(() => {
  const { currentPage, totalPages, maxVisible } = props
  const pages: (number | string)[] = []

  if (totalPages <= maxVisible) {
    // 如果总页数小于等于最大可见数，显示所有页码
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i)
    }
  } else {
    // 复杂的分页逻辑
    const halfVisible = Math.floor(maxVisible / 2)
    
    if (currentPage <= halfVisible + 1) {
      // 当前页在前半部分
      for (let i = 1; i <= maxVisible - 2; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages)
    } else if (currentPage >= totalPages - halfVisible) {
      // 当前页在后半部分
      pages.push(1)
      pages.push('...')
      for (let i = totalPages - (maxVisible - 3); i <= totalPages; i++) {
        pages.push(i)
      }
    } else {
      // 当前页在中间部分
      pages.push(1)
      pages.push('...')
      for (let i = currentPage - halfVisible + 1; i <= currentPage + halfVisible - 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPages)
    }
  }

  return pages
})
</script>
