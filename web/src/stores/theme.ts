import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import type { Theme } from '@/types'

// 预设主题
const themes: Record<string, Theme> = {
  dark: {
    name: 'dark',
    displayName: '经典深色',
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
      },
      status: {
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#3b82f6'
      }
    }
  },
  ocean: {
    name: 'ocean',
    displayName: '海洋蓝',
    colors: {
      primary: '#0ea5e9',
      secondary: '#0284c7',
      background: '#0c1426',
      surface: '#1e293b',
      text: {
        primary: '#f1f5f9',
        secondary: '#94a3b8',
        disabled: '#64748b'
      },
      glass: {
        light: 'rgba(14,165,233,0.05)',
        medium: 'rgba(14,165,233,0.08)',
        heavy: 'rgba(14,165,233,0.12)',
        border: 'rgba(14,165,233,0.1)'
      },
      status: {
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#0ea5e9'
      }
    }
  },
  sakura: {
    name: 'sakura',
    displayName: '樱花粉',
    colors: {
      primary: '#ec4899',
      secondary: '#db2777',
      background: '#1f0a14',
      surface: '#2d1b20',
      text: {
        primary: '#fdf2f8',
        secondary: '#f9a8d4',
        disabled: '#be185d'
      },
      glass: {
        light: 'rgba(236,72,153,0.05)',
        medium: 'rgba(236,72,153,0.08)',
        heavy: 'rgba(236,72,153,0.12)',
        border: 'rgba(236,72,153,0.1)'
      },
      status: {
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#ec4899'
      }
    }
  },
  forest: {
    name: 'forest',
    displayName: '森林绿',
    colors: {
      primary: '#059669',
      secondary: '#047857',
      background: '#0a1f17',
      surface: '#1a2e23',
      text: {
        primary: '#f0fdf4',
        secondary: '#86efac',
        disabled: '#166534'
      },
      glass: {
        light: 'rgba(5,150,105,0.05)',
        medium: 'rgba(5,150,105,0.08)',
        heavy: 'rgba(5,150,105,0.12)',
        border: 'rgba(5,150,105,0.1)'
      },
      status: {
        success: '#059669',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#3b82f6'
      }
    }
  },
  flame: {
    name: 'flame',
    displayName: '火焰橙',
    colors: {
      primary: '#ea580c',
      secondary: '#dc2626',
      background: '#1f0a05',
      surface: '#2d1b16',
      text: {
        primary: '#fff7ed',
        secondary: '#fed7aa',
        disabled: '#c2410c'
      },
      glass: {
        light: 'rgba(234,88,12,0.05)',
        medium: 'rgba(234,88,12,0.08)',
        heavy: 'rgba(234,88,12,0.12)',
        border: 'rgba(234,88,12,0.1)'
      },
      status: {
        success: '#10b981',
        warning: '#ea580c',
        error: '#dc2626',
        info: '#3b82f6'
      }
    }
  },
  purple: {
    name: 'purple',
    displayName: '神秘紫',
    colors: {
      primary: '#7c3aed',
      secondary: '#6d28d9',
      background: '#1a0b2e',
      surface: '#2d1b3d',
      text: {
        primary: '#faf5ff',
        secondary: '#c4b5fd',
        disabled: '#6b21a8'
      },
      glass: {
        light: 'rgba(124,58,237,0.05)',
        medium: 'rgba(124,58,237,0.08)',
        heavy: 'rgba(124,58,237,0.12)',
        border: 'rgba(124,58,237,0.1)'
      },
      status: {
        success: '#10b981',
        warning: '#f59e0b',
        error: '#ef4444',
        info: '#7c3aed'
      }
    }
  }
}

export const useThemeStore = defineStore('theme', () => {
  // 状态
  const currentTheme = ref('dark')
  const customTheme = ref<Partial<Theme['colors']> | null>(null)

  // 计算属性
  const theme = computed(() => {
    const baseTheme = themes[currentTheme.value]
    if (customTheme.value) {
      return {
        ...baseTheme,
        colors: {
          ...baseTheme.colors,
          ...customTheme.value
        }
      }
    }
    return baseTheme
  })

  const availableThemes = computed(() => Object.values(themes))

  // 设置主题
  const setTheme = (themeName: string) => {
    if (!themes[themeName]) {
      console.warn(`Theme "${themeName}" not found`)
      return
    }

    currentTheme.value = themeName
    customTheme.value = null
    applyTheme(themes[themeName])
    
    // 持久化到本地存储
    localStorage.setItem('theme', themeName)
    localStorage.removeItem('custom_theme')
  }

  // 设置自定义主题
  const setCustomTheme = (colors: Partial<Theme['colors']>) => {
    customTheme.value = colors
    
    const mergedTheme = {
      ...themes[currentTheme.value],
      colors: {
        ...themes[currentTheme.value].colors,
        ...colors
      }
    }
    
    applyTheme(mergedTheme)
    
    // 持久化到本地存储
    localStorage.setItem('custom_theme', JSON.stringify(colors))
  }

  // 应用主题到 DOM
  const applyTheme = (themeData: Theme) => {
    const root = document.documentElement

    // 设置基础颜色变量
    root.style.setProperty('--color-primary', themeData.colors.primary)
    root.style.setProperty('--color-secondary', themeData.colors.secondary)
    root.style.setProperty('--color-background', themeData.colors.background)
    root.style.setProperty('--color-surface', themeData.colors.surface)

    // 设置文字颜色变量
    root.style.setProperty('--color-text-primary', themeData.colors.text.primary)
    root.style.setProperty('--color-text-secondary', themeData.colors.text.secondary)
    root.style.setProperty('--color-text-disabled', themeData.colors.text.disabled)

    // 设置毛玻璃效果变量
    root.style.setProperty('--color-glass-light', themeData.colors.glass.light)
    root.style.setProperty('--color-glass-medium', themeData.colors.glass.medium)
    root.style.setProperty('--color-glass-heavy', themeData.colors.glass.heavy)
    root.style.setProperty('--color-glass-border', themeData.colors.glass.border)

    // 设置状态颜色变量
    root.style.setProperty('--color-success', themeData.colors.status.success)
    root.style.setProperty('--color-warning', themeData.colors.status.warning)
    root.style.setProperty('--color-error', themeData.colors.status.error)
    root.style.setProperty('--color-info', themeData.colors.status.info)

    // 更新 body 背景色
    document.body.style.background = themeData.colors.background
    document.body.style.color = themeData.colors.text.primary
  }

  // 切换到下一个主题
  const nextTheme = () => {
    const themeNames = Object.keys(themes)
    const currentIndex = themeNames.indexOf(currentTheme.value)
    const nextIndex = (currentIndex + 1) % themeNames.length
    setTheme(themeNames[nextIndex])
  }

  // 初始化主题
  const initTheme = () => {
    // 从本地存储恢复主题
    const savedTheme = localStorage.getItem('theme')
    const savedCustomTheme = localStorage.getItem('custom_theme')

    if (savedTheme && themes[savedTheme]) {
      currentTheme.value = savedTheme
    }

    if (savedCustomTheme) {
      try {
        customTheme.value = JSON.parse(savedCustomTheme)
      } catch (error) {
        console.error('Failed to parse custom theme:', error)
        localStorage.removeItem('custom_theme')
      }
    }

    // 应用主题
    applyTheme(theme.value)
  }

  // 重置为默认主题
  const resetTheme = () => {
    setTheme('dark')
  }

  // 导出主题配置
  const exportTheme = () => {
    return {
      theme: currentTheme.value,
      customColors: customTheme.value
    }
  }

  // 导入主题配置
  const importTheme = (config: { theme: string; customColors?: Partial<Theme['colors']> }) => {
    if (config.theme && themes[config.theme]) {
      currentTheme.value = config.theme
    }
    
    if (config.customColors) {
      customTheme.value = config.customColors
    }
    
    applyTheme(theme.value)
    
    // 保存到本地存储
    localStorage.setItem('theme', currentTheme.value)
    if (customTheme.value) {
      localStorage.setItem('custom_theme', JSON.stringify(customTheme.value))
    }
  }

  return {
    // 状态
    currentTheme: readonly(currentTheme),
    customTheme: readonly(customTheme),
    
    // 计算属性
    theme,
    availableThemes,
    
    // 方法
    setTheme,
    setCustomTheme,
    nextTheme,
    initTheme,
    resetTheme,
    exportTheme,
    importTheme
  }
})
