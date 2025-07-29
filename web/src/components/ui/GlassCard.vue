<template>
  <div
    :class="[
      'relative rounded-2xl border transition-all duration-300 ease-out',
      glassLevels[level],
      paddingClasses[padding],
      {
        'hover:-translate-y-1 hover:shadow-glass-lg cursor-pointer': hover,
        'glass-border': border
      },
      className
    ]"
    v-motion
    :initial="{ opacity: 0, scale: 0.95 }"
    :enter="{ opacity: 1, scale: 1 }"
    :transition="{ duration: 200, ease: 'easeOut' }"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import type { GlassCardProps } from '@/types'

interface Props extends GlassCardProps {
  border?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  level: 1,
  hover: false,
  className: '',
  padding: 'md',
  border: false
})

// 毛玻璃层级样式
const glassLevels = {
  1: 'backdrop-blur-[20px] bg-glass-light border-glass-border shadow-glass-md',
  2: 'backdrop-blur-[15px] border-glass-border shadow-glass-sm glass-medium',
  3: 'backdrop-blur-[25px] border-glass-border shadow-glass-lg glass-heavy'
}

// 内边距样式
const paddingClasses = {
  none: '',
  sm: 'p-3',
  md: 'p-4 sm:p-6',
  lg: 'p-6 sm:p-8'
}
</script>

<style scoped>
/* 毛玻璃效果 */
.glass-medium {
  background: var(--color-glass-medium);
}

.glass-heavy {
  background: var(--color-glass-heavy);
}

/* 渐变边框效果 */
.glass-border::before {
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
  -webkit-mask-composite: xor;
}

/* 悬浮时的额外效果 */
.hover\:shadow-glass-lg:hover {
  box-shadow:
    0 12px 40px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}
</style>
