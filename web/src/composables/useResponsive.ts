import { useWindowSize } from '@vueuse/core'
import { computed } from 'vue'

export const useResponsive = () => {
  const { width, height } = useWindowSize()

  // 断点定义
  const breakpoints = {
    xs: 475,
    sm: 640,
    md: 768,
    lg: 1024,
    xl: 1280,
    '2xl': 1536
  }

  // 设备类型判断
  const isMobile = computed(() => width.value < breakpoints.md)
  const isTablet = computed(() => width.value >= breakpoints.md && width.value < breakpoints.lg)
  const isDesktop = computed(() => width.value >= breakpoints.lg)
  const isLargeDesktop = computed(() => width.value >= breakpoints.xl)

  // 当前断点
  const breakpoint = computed(() => {
    if (width.value < breakpoints.xs) return 'xs'
    if (width.value < breakpoints.sm) return 'sm'
    if (width.value < breakpoints.md) return 'md'
    if (width.value < breakpoints.lg) return 'lg'
    if (width.value < breakpoints.xl) return 'xl'
    return '2xl'
  })

  // 屏幕方向
  const isPortrait = computed(() => height.value > width.value)
  const isLandscape = computed(() => width.value > height.value)

  // 是否为小屏设备
  const isSmallScreen = computed(() => width.value < breakpoints.lg)

  // 是否支持悬浮效果 (非触摸设备)
  const supportsHover = computed(() => {
    if (typeof window === 'undefined') return false
    return window.matchMedia('(hover: hover)').matches
  })

  // 获取容器类名
  const getContainerClass = () => {
    if (isMobile.value) return 'container-mobile'
    if (isTablet.value) return 'container-tablet'
    if (isDesktop.value) return 'container-desktop'
    return 'container-large'
  }

  // 获取网格列数
  const getGridCols = (options: {
    mobile?: number
    tablet?: number
    desktop?: number
    large?: number
  } = {}) => {
    const defaults = {
      mobile: 1,
      tablet: 2,
      desktop: 3,
      large: 4
    }
    
    const config = { ...defaults, ...options }
    
    if (isMobile.value) return config.mobile
    if (isTablet.value) return config.tablet
    if (isDesktop.value) return config.desktop
    return config.large
  }

  // 获取间距大小
  const getSpacing = (options: {
    mobile?: string
    tablet?: string
    desktop?: string
  } = {}) => {
    const defaults = {
      mobile: '4',
      tablet: '6',
      desktop: '8'
    }
    
    const config = { ...defaults, ...options }
    
    if (isMobile.value) return config.mobile
    if (isTablet.value) return config.tablet
    return config.desktop
  }

  return {
    // 尺寸信息
    width: readonly(width),
    height: readonly(height),
    
    // 设备类型
    isMobile,
    isTablet,
    isDesktop,
    isLargeDesktop,
    isSmallScreen,
    
    // 屏幕方向
    isPortrait,
    isLandscape,
    
    // 断点信息
    breakpoint,
    breakpoints,
    
    // 功能支持
    supportsHover,
    
    // 工具方法
    getContainerClass,
    getGridCols,
    getSpacing
  }
}
