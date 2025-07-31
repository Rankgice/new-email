<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div>
      <h2 class="text-xl font-semibold text-text-primary mb-2">邮件签名</h2>
      <p class="text-text-secondary">设置您的邮件签名，将自动添加到发送的邮件末尾</p>
    </div>

    <!-- 签名设置 -->
    <GlassCard padding="lg">
      <div class="space-y-6">
        <!-- 启用签名 -->
        <Switch
          v-model="signatureEnabled"
          label="启用邮件签名"
          description="在发送的邮件中自动添加签名"
          @change="updateSignatureStatus"
        />

        <!-- 签名编辑器 -->
        <div v-if="signatureEnabled">
          <!-- 编辑模式切换 -->
          <div class="flex items-center space-x-4 mb-4">
            <label class="text-sm font-medium text-text-primary">编辑模式:</label>
            <div class="flex space-x-2">
              <button
                :class="[
                  'px-3 py-1 rounded-lg text-sm font-medium transition-colors',
                  editMode === 'text'
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-600 text-gray-300 hover:bg-gray-500'
                ]"
                @click="editMode = 'text'"
              >
                纯文本
              </button>
              <button
                :class="[
                  'px-3 py-1 rounded-lg text-sm font-medium transition-colors',
                  editMode === 'html'
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-600 text-gray-300 hover:bg-gray-500'
                ]"
                @click="editMode = 'html'"
              >
                富文本
              </button>
            </div>
          </div>

          <!-- 纯文本编辑器 -->
          <div v-if="editMode === 'text'">
            <Textarea
              v-model="textSignature"
              label="签名内容"
              placeholder="请输入您的邮件签名..."
              :rows="8"
              :max-length="1000"
              show-count
              help="支持换行，最多1000个字符"
            />
          </div>

          <!-- 富文本编辑器 -->
          <div v-else>
            <label class="block text-sm font-medium text-text-primary mb-2">
              签名内容 (HTML)
            </label>
            <div class="space-y-3">
              <!-- 工具栏 -->
              <div class="flex items-center space-x-2 p-2 bg-white/5 rounded-lg border border-glass-border">
                <button
                  v-for="tool in htmlTools"
                  :key="tool.name"
                  :class="[
                    'p-2 rounded hover:bg-white/10 transition-colors',
                    tool.active ? 'bg-primary-500 text-white' : 'text-text-secondary'
                  ]"
                  :title="tool.title"
                  @click="applyHtmlFormat(tool.command, tool.value)"
                >
                  <component :is="tool.icon" class="w-4 h-4" />
                </button>
              </div>

              <!-- HTML编辑区 -->
              <div
                ref="htmlEditor"
                contenteditable="true"
                class="min-h-[200px] p-4 bg-white/5 rounded-lg border border-glass-border focus:border-primary-500 focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-background-primary outline-none text-text-primary"
                @input="updateHtmlSignature"
                @blur="saveHtmlContent"
              />

              <!-- HTML源码编辑 -->
              <div class="flex items-center space-x-2">
                <Switch
                  v-model="showHtmlSource"
                  label="显示HTML源码"
                />
              </div>

              <Textarea
                v-if="showHtmlSource"
                v-model="htmlSignature"
                label="HTML源码"
                :rows="6"
                help="直接编辑HTML代码"
                @input="updateHtmlEditor"
              />
            </div>
          </div>

          <!-- 预览区域 -->
          <div class="mt-6">
            <h4 class="text-md font-medium text-text-primary mb-3">预览效果</h4>
            <div class="p-4 bg-white/5 rounded-lg border border-glass-border">
              <div class="text-text-primary">
                <p class="mb-4">这是一封测试邮件的内容...</p>
                <div class="border-t border-gray-600 pt-4">
                  <div
                    v-if="editMode === 'html'"
                    v-html="htmlSignature"
                    class="signature-preview"
                  />
                  <pre
                    v-else
                    class="whitespace-pre-wrap font-sans text-sm text-text-secondary"
                  >{{ textSignature }}</pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 签名模板 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">签名模板</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div
          v-for="template in signatureTemplates"
          :key="template.id"
          class="p-4 bg-white/5 rounded-lg border border-glass-border hover:border-primary-500 cursor-pointer transition-colors"
          @click="applyTemplate(template)"
        >
          <h4 class="text-md font-medium text-text-primary mb-2">
            {{ template.name }}
          </h4>
          <div
            class="text-sm text-text-secondary signature-preview"
            v-html="template.content"
          />
        </div>
      </div>
    </GlassCard>

    <!-- 保存按钮 -->
    <div class="flex justify-end space-x-3">
      <Button
        variant="ghost"
        @click="resetSignature"
        :disabled="loading"
      >
        重置
      </Button>
      <Button
        variant="primary"
        @click="saveSignature"
        :loading="loading"
      >
        保存签名
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { userSettingsApi } from '@/api/user-settings'
import { useNotification } from '@/composables/useNotification'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Textarea from '@/components/ui/Textarea.vue'
import Switch from '@/components/ui/Switch.vue'

import {
  BoldIcon,
  ItalicIcon,
  UnderlineIcon,
  LinkIcon,
  PhotoIcon
} from '@heroicons/vue/24/outline'

// 通知
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const signatureEnabled = ref(true)
const editMode = ref<'text' | 'html'>('text')
const showHtmlSource = ref(false)

const textSignature = ref('')
const htmlSignature = ref('')
const htmlEditor = ref<HTMLElement>()

// HTML编辑工具
const htmlTools = [
  { name: 'bold', icon: BoldIcon, command: 'bold', title: '粗体', active: false },
  { name: 'italic', icon: ItalicIcon, command: 'italic', title: '斜体', active: false },
  { name: 'underline', icon: UnderlineIcon, command: 'underline', title: '下划线', active: false },
  { name: 'link', icon: LinkIcon, command: 'createLink', title: '链接', active: false },
  { name: 'image', icon: PhotoIcon, command: 'insertImage', title: '图片', active: false }
]

// 签名模板
const signatureTemplates = [
  {
    id: 1,
    name: '简约商务',
    content: `
      <div style="font-family: Arial, sans-serif; color: #333;">
        <p><strong>张三</strong><br>
        产品经理<br>
        ABC科技有限公司</p>
        <p>📧 zhangsan@example.com<br>
        📱 +86 138 0000 0000<br>
        🌐 www.example.com</p>
      </div>
    `
  },
  {
    id: 2,
    name: '创意设计',
    content: `
      <div style="font-family: 'Helvetica Neue', sans-serif; color: #2c3e50;">
        <table style="border: none;">
          <tr>
            <td style="padding-right: 20px; border-right: 3px solid #3498db;">
              <h3 style="margin: 0; color: #3498db;">李四</h3>
              <p style="margin: 5px 0; color: #7f8c8d;">UI/UX设计师</p>
            </td>
            <td style="padding-left: 20px;">
              <p style="margin: 0; font-size: 14px;">
                📧 lisi@example.com<br>
                💼 XYZ设计工作室<br>
                🎨 让设计改变世界
              </p>
            </td>
          </tr>
        </table>
      </div>
    `
  },
  {
    id: 3,
    name: '极简风格',
    content: `
      <div style="font-family: 'SF Pro Text', -apple-system, sans-serif; color: #1d1d1f;">
        <p style="margin: 0; font-size: 16px; font-weight: 600;">王五</p>
        <p style="margin: 5px 0 0 0; font-size: 14px; color: #86868b;">
          软件工程师 | wangwu@example.com
        </p>
      </div>
    `
  },
  {
    id: 4,
    name: '社交媒体',
    content: `
      <div style="font-family: Arial, sans-serif;">
        <p><strong>赵六</strong><br>
        市场营销专家</p>
        <p>
          <a href="mailto:zhaoliu@example.com" style="color: #1da1f2; text-decoration: none;">📧 邮箱</a> |
          <a href="#" style="color: #1da1f2; text-decoration: none;">🐦 Twitter</a> |
          <a href="#" style="color: #0077b5; text-decoration: none;">💼 LinkedIn</a> |
          <a href="#" style="color: #25d366; text-decoration: none;">💬 微信</a>
        </p>
        <p style="font-style: italic; color: #666;">
          "创新营销，连接未来"
        </p>
      </div>
    `
  }
]

// 生命周期
onMounted(() => {
  loadSignature()
})

// 方法
const loadSignature = async () => {
  try {
    const settings = await userSettingsApi.getSettings()
    if (settings.emailSignature) {
      // 检测是否为HTML格式
      if (settings.emailSignature.includes('<') && settings.emailSignature.includes('>')) {
        editMode.value = 'html'
        htmlSignature.value = settings.emailSignature
        updateHtmlEditor()
      } else {
        editMode.value = 'text'
        textSignature.value = settings.emailSignature
      }
    }
  } catch (error) {
    console.error('Failed to load signature:', error)
  }
}

const updateSignatureStatus = async () => {
  // 这里可以立即保存启用状态
  if (!signatureEnabled.value) {
    try {
      await userSettingsApi.updateEmailSignature('')
      showSuccess('邮件签名已禁用')
    } catch (error) {
      showError('更新失败')
    }
  }
}

const applyHtmlFormat = (command: string, value?: string) => {
  if (!htmlEditor.value) return

  htmlEditor.value.focus()

  if (command === 'createLink') {
    const url = prompt('请输入链接地址:')
    if (url) {
      document.execCommand(command, false, url)
    }
  } else if (command === 'insertImage') {
    const url = prompt('请输入图片地址:')
    if (url) {
      document.execCommand(command, false, url)
    }
  } else {
    document.execCommand(command, false, value)
  }

  updateHtmlSignature()
}

const updateHtmlSignature = () => {
  if (htmlEditor.value) {
    htmlSignature.value = htmlEditor.value.innerHTML
  }
}

const updateHtmlEditor = () => {
  if (htmlEditor.value) {
    htmlEditor.value.innerHTML = htmlSignature.value
  }
}

const saveHtmlContent = () => {
  updateHtmlSignature()
}

const applyTemplate = (template: any) => {
  if (editMode.value === 'html') {
    htmlSignature.value = template.content.trim()
    nextTick(() => {
      updateHtmlEditor()
    })
  } else {
    // 将HTML转换为纯文本
    const tempDiv = document.createElement('div')
    tempDiv.innerHTML = template.content
    textSignature.value = tempDiv.textContent || tempDiv.innerText || ''
  }
}

const resetSignature = () => {
  textSignature.value = ''
  htmlSignature.value = ''
  if (htmlEditor.value) {
    htmlEditor.value.innerHTML = ''
  }
}

const saveSignature = async () => {
  if (!signatureEnabled.value) {
    showError('请先启用邮件签名')
    return
  }

  try {
    loading.value = true
    
    const signature = editMode.value === 'html' ? htmlSignature.value : textSignature.value
    
    if (!signature.trim()) {
      showError('签名内容不能为空')
      return
    }

    await userSettingsApi.updateEmailSignature(signature)
    showSuccess('邮件签名保存成功')
  } catch (error) {
    showError('保存签名失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.signature-preview {
  font-size: 14px;
  line-height: 1.5;
}

.signature-preview p {
  margin: 0.5em 0;
}

.signature-preview a {
  color: #3b82f6;
  text-decoration: none;
}

.signature-preview a:hover {
  text-decoration: underline;
}

/* HTML编辑器样式 */
[contenteditable="true"]:focus {
  outline: none;
}

[contenteditable="true"] p {
  margin: 0.5em 0;
}

[contenteditable="true"] strong {
  font-weight: 600;
}

[contenteditable="true"] em {
  font-style: italic;
}

[contenteditable="true"] u {
  text-decoration: underline;
}
</style>
