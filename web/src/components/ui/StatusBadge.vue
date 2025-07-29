<template>
  <span :class="badgeClasses">
    <span :class="dotClasses"></span>
    {{ statusText }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  status: number | string
  type?: 'default' | 'mailbox' | 'user'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default'
})

const statusConfig = computed(() => {
  const status = Number(props.status)
  
  if (props.type === 'mailbox') {
    return status === 1
      ? { text: '启用', color: 'green' }
      : { text: '禁用', color: 'red' }
  }
  
  if (props.type === 'user') {
    return status === 1
      ? { text: '活跃', color: 'green' }
      : { text: '禁用', color: 'red' }
  }
  
  // 默认状态配置
  return status === 1
    ? { text: '正常', color: 'green' }
    : { text: '异常', color: 'red' }
})

const statusText = computed(() => statusConfig.value.text)

const badgeClasses = computed(() => {
  const baseClasses = 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium'
  const colorClasses = {
    green: 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400',
    red: 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400',
    yellow: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400',
    blue: 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400',
    gray: 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
  }
  
  return `${baseClasses} ${colorClasses[statusConfig.value.color as keyof typeof colorClasses]}`
})

const dotClasses = computed(() => {
  const baseClasses = 'w-1.5 h-1.5 rounded-full mr-1.5'
  const colorClasses = {
    green: 'bg-green-400',
    red: 'bg-red-400',
    yellow: 'bg-yellow-400',
    blue: 'bg-blue-400',
    gray: 'bg-gray-400'
  }
  
  return `${baseClasses} ${colorClasses[statusConfig.value.color as keyof typeof colorClasses]}`
})
</script>
</script>
