<template>
  <div class="bg-background-secondary border-b border-glass-border p-4">
    <div class="flex items-center justify-between">
      <!-- 左侧操作 -->
      <div class="flex items-center space-x-3">
        <!-- 全选复选框 -->
        <input
          v-model="selectAll"
          type="checkbox"
          class="w-4 h-4 text-primary-600 bg-glass-light border-glass-border rounded focus:ring-primary-500 focus:ring-2"
          @change="handleSelectAll"
        />

        <!-- 批量操作按钮 -->
        <div class="flex items-center space-x-2">
          <Button
            variant="ghost"
            size="sm"
            :disabled="selectedCount === 0"
            @click="handleBatchDelete"
          >
            <TrashIcon class="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            :disabled="selectedCount === 0"
            @click="handleBatchStar"
          >
            <StarIcon class="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            :disabled="selectedCount === 0"
            @click="handleBatchMarkRead"
          >
            <EnvelopeOpenIcon class="w-4 h-4" />
          </Button>
        </div>

        <!-- 选中计数 -->
        <span
          v-if="selectedCount > 0"
          class="text-sm text-text-secondary"
        >
          已选择 {{ selectedCount }} 封邮件
        </span>
      </div>

      <!-- 右侧操作 -->
      <div class="flex items-center space-x-3">
        <!-- 搜索框 -->
        <div class="relative">
          <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-secondary" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索邮件..."
            class="pl-10 pr-4 py-2 w-64 glass-card rounded-lg text-sm text-text-primary placeholder-text-secondary focus-ring"
            @input="handleSearch"
          />
        </div>

        <!-- 刷新按钮 -->
        <Button
          variant="ghost"
          size="sm"
          :loading="isRefreshing"
          @click="handleRefresh"
        >
          <ArrowPathIcon class="w-4 h-4" />
        </Button>

        <!-- 视图切换 -->
        <div class="flex items-center bg-glass-light rounded-lg p-1">
          <button
            :class="[
              'p-1 rounded transition-colors duration-200',
              viewMode === 'list'
                ? 'bg-primary-500/20 text-primary-400'
                : 'text-text-secondary hover:text-text-primary'
            ]"
            @click="viewMode = 'list'"
          >
            <Bars3Icon class="w-4 h-4" />
          </button>
          <button
            :class="[
              'p-1 rounded transition-colors duration-200',
              viewMode === 'grid'
                ? 'bg-primary-500/20 text-primary-400'
                : 'text-text-secondary hover:text-text-primary'
            ]"
            @click="viewMode = 'grid'"
          >
            <Squares2X2Icon class="w-4 h-4" />
          </button>
        </div>

        <!-- 排序选择 -->
        <select
          v-model="sortBy"
          class="glass-card rounded-lg px-3 py-2 text-sm text-text-primary focus-ring"
          @change="handleSort"
        >
          <option value="date">按日期排序</option>
          <option value="subject">按主题排序</option>
          <option value="from">按发件人排序</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Button from '@/components/ui/Button.vue'
import {
  TrashIcon,
  StarIcon,
  EnvelopeOpenIcon,
  MagnifyingGlassIcon,
  ArrowPathIcon,
  Bars3Icon,
  Squares2X2Icon
} from '@heroicons/vue/24/outline'

// 状态
const selectAll = ref(false)
const selectedCount = ref(0)
const searchQuery = ref('')
const isRefreshing = ref(false)
const viewMode = ref<'list' | 'grid'>('list')
const sortBy = ref('date')

// 处理全选
const handleSelectAll = () => {
  // TODO: 实现全选逻辑
  console.log('Select all:', selectAll.value)
}

// 处理批量删除
const handleBatchDelete = () => {
  // TODO: 实现批量删除逻辑
  console.log('Batch delete')
}

// 处理批量加星
const handleBatchStar = () => {
  // TODO: 实现批量加星逻辑
  console.log('Batch star')
}

// 处理批量标记已读
const handleBatchMarkRead = () => {
  // TODO: 实现批量标记已读逻辑
  console.log('Batch mark read')
}

// 处理搜索
const handleSearch = () => {
  // TODO: 实现搜索逻辑
  console.log('Search:', searchQuery.value)
}

// 处理刷新
const handleRefresh = async () => {
  isRefreshing.value = true
  // TODO: 实现刷新逻辑
  setTimeout(() => {
    isRefreshing.value = false
  }, 1000)
}

// 处理排序
const handleSort = () => {
  // TODO: 实现排序逻辑
  console.log('Sort by:', sortBy.value)
}
</script>

<style scoped>
.focus-ring {
  @apply focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-background-primary;
}

.glass-card {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  background: var(--color-glass-light);
  border: 1px solid var(--color-glass-border);
}
</style>
