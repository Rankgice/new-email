<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-text-primary mb-2">API密钥管理</h2>
        <p class="text-text-secondary">管理您的API密钥，用于第三方应用集成</p>
      </div>
      <Button
        variant="primary"
        @click="showCreateModal = true"
      >
        <PlusIcon class="w-4 h-4 mr-2" />
        创建API密钥
      </Button>
    </div>

    <!-- API密钥列表 -->
    <GlassCard padding="lg">
      <div v-if="apiKeys.length === 0" class="text-center py-12">
        <KeyIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
        <h3 class="text-lg font-medium text-text-primary mb-2">暂无API密钥</h3>
        <p class="text-text-secondary mb-4">
          创建API密钥以便第三方应用访问您的邮件数据
        </p>
        <Button
          variant="primary"
          @click="showCreateModal = true"
        >
          创建第一个API密钥
        </Button>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="apiKey in apiKeys"
          :key="apiKey.id"
          class="p-4 bg-white/5 rounded-lg border border-glass-border"
        >
          <!-- API密钥头部 -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center space-x-3">
              <div
                :class="[
                  'w-3 h-3 rounded-full',
                  apiKey.enabled ? 'bg-green-400' : 'bg-gray-400'
                ]"
              />
              <div>
                <h4 class="text-md font-medium text-text-primary">
                  {{ apiKey.name }}
                </h4>
                <p class="text-sm text-text-secondary">
                  创建于 {{ formatDate(apiKey.createdAt) }}
                </p>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <Switch
                v-model="apiKey.enabled"
                @change="toggleApiKey(apiKey.id, apiKey.enabled)"
              />
              <Button
                variant="ghost"
                size="sm"
                @click="editApiKey(apiKey)"
              >
                <PencilIcon class="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="regenerateApiKey(apiKey.id)"
              >
                <ArrowPathIcon class="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="deleteApiKey(apiKey.id)"
              >
                <TrashIcon class="w-4 h-4" />
              </Button>
            </div>
          </div>

          <!-- API密钥信息 -->
          <div class="space-y-3">
            <!-- 密钥值 -->
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                API密钥
              </label>
              <div class="flex items-center space-x-2">
                <Input
                  :model-value="showKeys[apiKey.id] ? apiKey.key : '••••••••••••••••••••••••••••••••'"
                  readonly
                  class="flex-1"
                />
                <Button
                  variant="ghost"
                  size="sm"
                  @click="toggleKeyVisibility(apiKey.id)"
                >
                  <component :is="showKeys[apiKey.id] ? EyeSlashIcon : EyeIcon" class="w-4 h-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  @click="copyToClipboard(apiKey.key)"
                >
                  <ClipboardIcon class="w-4 h-4" />
                </Button>
              </div>
            </div>

            <!-- 权限 -->
            <div>
              <label class="block text-sm font-medium text-text-primary mb-1">
                权限
              </label>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="permission in apiKey.permissions"
                  :key="permission"
                  class="px-2 py-1 bg-primary-500/20 text-primary-400 rounded text-xs font-medium"
                >
                  {{ formatPermission(permission) }}
                </span>
              </div>
            </div>

            <!-- 使用信息 -->
            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span class="text-text-secondary">最后使用:</span>
                <span class="text-text-primary ml-2">
                  {{ apiKey.lastUsed ? formatDate(apiKey.lastUsed) : '从未使用' }}
                </span>
              </div>
              <div v-if="apiKey.expiresAt">
                <span class="text-text-secondary">过期时间:</span>
                <span
                  :class="[
                    'ml-2',
                    isExpiringSoon(apiKey.expiresAt) ? 'text-orange-400' : 'text-text-primary'
                  ]"
                >
                  {{ formatDate(apiKey.expiresAt) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 使用说明 -->
    <GlassCard padding="lg">
      <h3 class="text-lg font-medium text-text-primary mb-4">使用说明</h3>
      <div class="space-y-4 text-sm text-text-secondary">
        <div>
          <h4 class="font-medium text-text-primary mb-2">认证方式</h4>
          <p>在HTTP请求头中添加以下认证信息：</p>
          <div class="mt-2 p-3 bg-black/20 rounded-lg font-mono text-xs">
            Authorization: Bearer YOUR_API_KEY
          </div>
        </div>
        
        <div>
          <h4 class="font-medium text-text-primary mb-2">API端点</h4>
          <ul class="space-y-1">
            <li>• <code class="text-primary-400">GET /api/emails</code> - 获取邮件列表</li>
            <li>• <code class="text-primary-400">POST /api/emails</code> - 发送邮件</li>
            <li>• <code class="text-primary-400">GET /api/mailboxes</code> - 获取邮箱列表</li>
            <li>• <code class="text-primary-400">GET /api/user/profile</code> - 获取用户信息</li>
          </ul>
        </div>

        <div>
          <h4 class="font-medium text-text-primary mb-2">安全提醒</h4>
          <ul class="space-y-1">
            <li>• 请妥善保管您的API密钥，不要在公开场所分享</li>
            <li>• 建议定期更换API密钥以确保安全</li>
            <li>• 为不同的应用创建不同的API密钥</li>
            <li>• 及时删除不再使用的API密钥</li>
          </ul>
        </div>
      </div>
    </GlassCard>

    <!-- 创建/编辑API密钥模态框 -->
    <Modal
      v-model="showCreateModal"
      :title="editingApiKey ? '编辑API密钥' : '创建API密钥'"
      @confirm="saveApiKey"
      :loading="saving"
    >
      <div class="space-y-4">
        <!-- 密钥名称 -->
        <Input
          v-model="apiKeyForm.name"
          label="密钥名称"
          placeholder="请输入密钥名称"
          :error="errors.name"
          required
        />

        <!-- 权限选择 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            权限设置 <span class="text-red-400">*</span>
          </label>
          <div class="space-y-2">
            <div
              v-for="permission in availablePermissions"
              :key="permission.value"
              class="flex items-center space-x-3"
            >
              <input
                :id="permission.value"
                v-model="apiKeyForm.permissions"
                type="checkbox"
                :value="permission.value"
                class="rounded border-gray-600 bg-gray-700 text-primary-500 focus:ring-primary-500"
              />
              <label
                :for="permission.value"
                class="text-sm text-text-primary cursor-pointer"
              >
                {{ permission.label }}
              </label>
            </div>
          </div>
          <p v-if="errors.permissions" class="mt-1 text-sm text-red-400">
            {{ errors.permissions }}
          </p>
        </div>

        <!-- 过期时间 -->
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            过期时间
          </label>
          <Select
            v-model="apiKeyForm.expirationOption"
            :options="expirationOptions"
            @change="updateExpirationDate"
          />
          <Input
            v-if="apiKeyForm.expirationOption === 'custom'"
            v-model="apiKeyForm.customExpiration"
            type="datetime-local"
            class="mt-2"
          />
        </div>
      </div>
    </Modal>

    <!-- 新创建的API密钥显示模态框 -->
    <Modal
      v-model="showNewKeyModal"
      title="API密钥创建成功"
      :show-cancel="false"
      confirm-text="我已保存"
    >
      <div class="space-y-4">
        <div class="p-4 bg-green-500/10 rounded-lg border border-green-500/20">
          <p class="text-green-400 font-medium mb-2">
            ✅ API密钥创建成功！
          </p>
          <p class="text-text-secondary text-sm">
            请立即复制并保存您的API密钥，出于安全考虑，我们不会再次显示完整密钥。
          </p>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            您的API密钥
          </label>
          <div class="flex items-center space-x-2">
            <Input
              :model-value="newApiKey"
              readonly
              class="flex-1 font-mono"
            />
            <Button
              variant="secondary"
              @click="copyToClipboard(newApiKey)"
            >
              <ClipboardIcon class="w-4 h-4 mr-2" />
              复制
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { userSettingsApi } from '@/api/user-settings'
import { useNotification } from '@/composables/useNotification'
import type { ApiKey } from '@/types'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'
import Modal from '@/components/ui/Modal.vue'

import {
  PlusIcon,
  KeyIcon,
  PencilIcon,
  TrashIcon,
  ArrowPathIcon,
  EyeIcon,
  EyeSlashIcon,
  ClipboardIcon
} from '@heroicons/vue/24/outline'

// 通知
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const showCreateModal = ref(false)
const showNewKeyModal = ref(false)
const editingApiKey = ref<ApiKey | null>(null)
const newApiKey = ref('')

const apiKeys = ref<ApiKey[]>([])
const showKeys = reactive<Record<string, boolean>>({})

// 表单数据
const apiKeyForm = reactive({
  name: '',
  permissions: [] as string[],
  expirationOption: 'never',
  customExpiration: ''
})

const errors = reactive({
  name: '',
  permissions: ''
})

// 可用权限
const availablePermissions = [
  { label: '读取邮件', value: 'emails:read' },
  { label: '发送邮件', value: 'emails:send' },
  { label: '删除邮件', value: 'emails:delete' },
  { label: '管理邮箱', value: 'mailboxes:manage' },
  { label: '读取用户信息', value: 'user:read' },
  { label: '修改用户信息', value: 'user:write' }
]

// 过期时间选项
const expirationOptions = [
  { label: '永不过期', value: 'never' },
  { label: '30天', value: '30d' },
  { label: '90天', value: '90d' },
  { label: '1年', value: '1y' },
  { label: '自定义', value: 'custom' }
]

// 生命周期
onMounted(() => {
  loadApiKeys()
})

// 方法
const loadApiKeys = async () => {
  try {
    loading.value = true
    apiKeys.value = await userSettingsApi.getApiKeys()
  } catch (error) {
    showError('加载API密钥失败')
  } finally {
    loading.value = false
  }
}

const toggleKeyVisibility = (keyId: string) => {
  showKeys[keyId] = !showKeys[keyId]
}

const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    showSuccess('已复制到剪贴板')
  } catch (error) {
    showError('复制失败')
  }
}

const formatPermission = (permission: string) => {
  const permissionMap: Record<string, string> = {
    'emails:read': '读取邮件',
    'emails:send': '发送邮件',
    'emails:delete': '删除邮件',
    'mailboxes:manage': '管理邮箱',
    'user:read': '读取用户信息',
    'user:write': '修改用户信息'
  }
  return permissionMap[permission] || permission
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const isExpiringSoon = (expiresAt: string) => {
  const expireDate = new Date(expiresAt)
  const now = new Date()
  const diffDays = (expireDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
  return diffDays <= 7
}

const editApiKey = (apiKey: ApiKey) => {
  editingApiKey.value = apiKey
  apiKeyForm.name = apiKey.name
  apiKeyForm.permissions = [...apiKey.permissions]
  
  if (apiKey.expiresAt) {
    apiKeyForm.expirationOption = 'custom'
    apiKeyForm.customExpiration = new Date(apiKey.expiresAt).toISOString().slice(0, 16)
  } else {
    apiKeyForm.expirationOption = 'never'
  }
  
  showCreateModal.value = true
}

const resetForm = () => {
  editingApiKey.value = null
  apiKeyForm.name = ''
  apiKeyForm.permissions = []
  apiKeyForm.expirationOption = 'never'
  apiKeyForm.customExpiration = ''
  errors.name = ''
  errors.permissions = ''
}

const updateExpirationDate = () => {
  if (apiKeyForm.expirationOption !== 'custom') {
    apiKeyForm.customExpiration = ''
  }
}

const validateForm = () => {
  errors.name = ''
  errors.permissions = ''
  
  if (!apiKeyForm.name.trim()) {
    errors.name = '请输入密钥名称'
    return false
  }

  if (apiKeyForm.permissions.length === 0) {
    errors.permissions = '请至少选择一个权限'
    return false
  }

  return true
}

const saveApiKey = async () => {
  if (!validateForm()) return

  try {
    saving.value = true
    
    let expiresAt: string | undefined
    
    if (apiKeyForm.expirationOption !== 'never') {
      if (apiKeyForm.expirationOption === 'custom') {
        expiresAt = new Date(apiKeyForm.customExpiration).toISOString()
      } else {
        const now = new Date()
        const days = {
          '30d': 30,
          '90d': 90,
          '1y': 365
        }[apiKeyForm.expirationOption] || 30
        
        expiresAt = new Date(now.getTime() + days * 24 * 60 * 60 * 1000).toISOString()
      }
    }

    if (editingApiKey.value) {
      // 更新现有密钥
      await userSettingsApi.updateApiKey(editingApiKey.value.id, {
        name: apiKeyForm.name,
        permissions: apiKeyForm.permissions
      })
      showSuccess('API密钥更新成功')
    } else {
      // 创建新密钥
      const result = await userSettingsApi.createApiKey({
        name: apiKeyForm.name,
        permissions: apiKeyForm.permissions,
        expiresAt
      })
      
      newApiKey.value = result.key
      showNewKeyModal.value = true
    }

    showCreateModal.value = false
    resetForm()
    loadApiKeys()
  } catch (error) {
    showError('保存API密钥失败')
  } finally {
    saving.value = false
  }
}

const toggleApiKey = async (id: string, enabled: boolean) => {
  try {
    await userSettingsApi.updateApiKey(id, { enabled })
    showSuccess(enabled ? 'API密钥已启用' : 'API密钥已禁用')
  } catch (error) {
    showError('操作失败')
    // 恢复状态
    const apiKey = apiKeys.value.find(k => k.id === id)
    if (apiKey) {
      apiKey.enabled = !enabled
    }
  }
}

const regenerateApiKey = async (id: string) => {
  if (!confirm('确定要重新生成API密钥吗？旧密钥将立即失效。')) return

  try {
    const result = await userSettingsApi.regenerateApiKey(id)
    newApiKey.value = result.key
    showNewKeyModal.value = true
    loadApiKeys()
    showSuccess('API密钥重新生成成功')
  } catch (error) {
    showError('重新生成API密钥失败')
  }
}

const deleteApiKey = async (id: string) => {
  if (!confirm('确定要删除这个API密钥吗？此操作不可撤销。')) return

  try {
    await userSettingsApi.deleteApiKey(id)
    apiKeys.value = apiKeys.value.filter(k => k.id !== id)
    showSuccess('API密钥删除成功')
  } catch (error) {
    showError('删除API密钥失败')
  }
}
</script>
