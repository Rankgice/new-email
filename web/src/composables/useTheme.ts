import { useThemeStore } from '@/stores/theme'

export const useTheme = () => {
  const themeStore = useThemeStore()
  
  return {
    // 状态
    theme: themeStore.theme,
    currentTheme: themeStore.currentTheme,
    availableThemes: themeStore.availableThemes,
    customTheme: themeStore.customTheme,
    
    // 方法
    setTheme: themeStore.setTheme,
    setCustomTheme: themeStore.setCustomTheme,
    nextTheme: themeStore.nextTheme,
    initTheme: themeStore.initTheme,
    resetTheme: themeStore.resetTheme,
    exportTheme: themeStore.exportTheme,
    importTheme: themeStore.importTheme
  }
}
