<template>
  <Modal
    :visible="visible"
    :title="isEdit ? '编辑邮箱' : '添加邮箱'"
    size="lg"
    @close="$emit('close')"
  >
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- 邮箱类型选择 -->
      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">
          邮箱类型
        </label>
        <div class="grid grid-cols-2 gap-4">
          <button
            type="button"
            :class="[
              'p-4 border-2 rounded-lg transition-all duration-200',
              formData.type === 'third'
                ? 'border-primary-500 bg-primary-500/10'
                : 'border-glass-border hover:border-primary-500/50'
            ]"
            @click="formData.type = 'third'"
          >
            <ServerIcon class="w-8 h-8 mx-auto mb-2 text-primary-400" />
            <div class="text-sm font-medium text-text-primary">第三方邮箱</div>
            <div class="text-xs text-text-secondary mt-1">Gmail, Outlook, QQ等</div>
          </button>
          <button
            type="button"
            :class="[
              'p-4 border-2 rounded-lg transition-all duration-200',
              formData.type === 'self'
                ? 'border-primary-500 bg-primary-500/10'
                : 'border-glass-border hover:border-primary-500/50'
            ]"
            @click="formData.type = 'self'"
          >
            <BuildingOfficeIcon class="w-8 h-8 mx-auto mb-2 text-primary-400" />
            <div class="text-sm font-medium text-text-primary">自建邮箱</div>
            <div class="text-xs text-text-secondary mt-1">企业邮箱服务器</div>
          </button>
        </div>
      </div>

      <!-- 邮箱提供商选择（第三方邮箱） -->
      <div v-if="formData.type === 'third'">
        <label class="block text-sm font-medium text-text-primary mb-2">
          邮箱提供商
        </label>
        <select
          v-model="formData.provider"
          class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
          @change="handleProviderChange"
        >
          <option value="">请选择提供商</option>
          <option
            v-for="provider in providers"
            :key="provider.provider"
            :value="provider.provider"
          >
            {{ getProviderLabel(provider.provider) }}
          </option>
        </select>
      </div>

      <!-- 域名选择（自建邮箱） -->
      <div v-if="formData.type === 'self'">
        <label class="block text-sm font-medium text-text-primary mb-2">
          域名
        </label>
        <select
          v-model="formData.domainId"
          class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary"
        >
          <option value="">请选择域名</option>
          <!-- TODO: 从API获取域名列表 -->
        </select>
      </div>

      <!-- 邮箱地址 -->
      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">
          邮箱地址 *
        </label>
        <input
          v-model="formData.email"
          type="email"
          required
          class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
          placeholder="example@domain.com"
        />
      </div>

      <!-- 邮箱密码 -->
      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">
          邮箱密码 *
        </label>
        <div class="relative">
          <input
            v-model="formData.password"
            :type="showPassword ? 'text' : 'password'"
            required
            class="w-full px-3 py-2 pr-10 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
            placeholder="请输入邮箱密码或应用专用密码"
          />
          <button
            type="button"
            class="absolute right-3 top-1/2 transform -translate-y-1/2 text-text-secondary hover:text-text-primary"
            @click="showPassword = !showPassword"
          >
            <EyeIcon v-if="!showPassword" class="w-4 h-4" />
            <EyeSlashIcon v-else class="w-4 h-4" />
          </button>
        </div>
        <p class="text-xs text-text-secondary mt-1">
          对于Gmail等邮箱，请使用应用专用密码而非登录密码
        </p>
      </div>

      <!-- 服务器配置 -->
      <div class="space-y-4">
        <h3 class="text-lg font-medium text-text-primary">服务器配置</h3>
        
        <!-- IMAP配置 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-text-primary mb-2">
              IMAP服务器 *
            </label>
            <input
              v-model="formData.imapHost"
              type="text"
              required
              class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
              placeholder="imap.example.com"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-primary mb-2">
              IMAP端口 *
            </label>
            <input
              v-model.number="formData.imapPort"
              type="number"
              required
              class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
              placeholder="993"
            />
          </div>
        </div>

        <!-- SMTP配置 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-text-primary mb-2">
              SMTP服务器 *
            </label>
            <input
              v-model="formData.smtpHost"
              type="text"
              required
              class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
              placeholder="smtp.example.com"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-primary mb-2">
              SMTP端口 *
            </label>
            <input
              v-model.number="formData.smtpPort"
              type="number"
              required
              class="w-full px-3 py-2 bg-glass-light border border-glass-border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent text-text-primary placeholder-text-secondary"
              placeholder="587"
            />
          </div>
        </div>

        <!-- SSL选项 -->
        <div class="grid grid-cols-2 gap-4">
          <div class="flex items-center">
            <input
              id="imapSsl"
              v-model="formData.imapSsl"
              type="checkbox"
              class="w-4 h-4 text-primary-500 bg-glass-light border-glass-border rounded focus:ring-primary-500"
            />
            <label for="imapSsl" class="ml-2 text-sm text-text-primary">
              IMAP使用SSL/TLS
            </label>
          </div>
          <div class="flex items-center">
            <input
              id="smtpSsl"
              v-model="formData.smtpSsl"
              type="checkbox"
              class="w-4 h-4 text-primary-500 bg-glass-light border-glass-border rounded focus:ring-primary-500"
            />
            <label for="smtpSsl" class="ml-2 text-sm text-text-primary">
              SMTP使用SSL/TLS
            </label>
          </div>
        </div>
      </div>

      <!-- 其他选项 -->
      <div class="space-y-4">
        <div class="flex items-center">
          <input
            id="autoReceive"
            v-model="formData.autoReceive"
            type="checkbox"
            class="w-4 h-4 text-primary-500 bg-glass-light border-glass-border rounded focus:ring-primary-500"
          />
          <label for="autoReceive" class="ml-2 text-sm text-text-primary">
            自动收取邮件
          </label>
        </div>
        <div class="flex items-center">
          <input
            id="status"
            v-model="statusEnabled"
            type="checkbox"
            class="w-4 h-4 text-primary-500 bg-glass-light border-glass-border rounded focus:ring-primary-500"
          />
          <label for="status" class="ml-2 text-sm text-text-primary">
            启用邮箱
          </label>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-between pt-6 border-t border-glass-border">
        <div>
          <Button
            v-if="!isEdit"
            type="button"
            variant="ghost"
            @click="testConnection"
            :loading="testing"
          >
            <WifiIcon class="w-4 h-4 mr-2" />
            测试连接
          </Button>
        </div>
        <div class="flex space-x-3">
          <Button
            type="button"
            variant="ghost"
            @click="$emit('close')"
          >
            取消
          </Button>
          <Button
            type="submit"
            variant="primary"
            :loading="saving"
          >
            {{ isEdit ? '更新' : '添加' }}
          </Button>
        </div>
      </div>
    </form>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { 
  Mailbox, 
  MailboxCreateRequest, 
  MailboxUpdateRequest, 
  MailboxProvider,
  MailboxTestConnectionRequest 
} from '@/types'
import { mailboxApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'

// Icons
import {
  ServerIcon,
  BuildingOfficeIcon,
  EyeIcon,
  EyeSlashIcon,
  WifiIcon
} from '@heroicons/vue/24/outline'

// Components
import Modal from '@/components/ui/Modal.vue'
import Button from '@/components/ui/Button.vue'

// Props
interface Props {
  visible: boolean
  mailbox?: Mailbox | null
  providers: MailboxProvider[]
}

const props = withDefaults(defineProps<Props>(), {
  mailbox: null
})

// Emits
const emit = defineEmits<{
  close: []
  save: [data: MailboxCreateRequest | MailboxUpdateRequest]
}>()

// Composables
const { showNotification } = useNotification()

// State
const showPassword = ref(false)
const saving = ref(false)
const testing = ref(false)

// 计算属性
const isEdit = computed(() => !!props.mailbox)
const statusEnabled = computed({
  get: () => formData.value.status === 1,
  set: (value: boolean) => {
    formData.value.status = value ? 1 : 2
  }
})

// 表单数据
const formData = ref<MailboxCreateRequest & { id?: number }>({
  type: 'third',
  email: '',
  password: '',
  provider: '',
  domainId: 0,
  imapHost: '',
  imapPort: 993,
  imapSsl: true,
  smtpHost: '',
  smtpPort: 587,
  smtpSsl: true,
  autoReceive: true,
  status: 1
})

// 方法
const getProviderLabel = (provider: string) => {
  const labels: Record<string, string> = {
    gmail: 'Gmail',
    outlook: 'Outlook',
    qq: 'QQ邮箱',
    163: '163邮箱',
    126: '126邮箱',
    imap: '其他IMAP'
  }
  return labels[provider] || provider
}

const handleProviderChange = () => {
  const provider = props.providers.find(p => p.provider === formData.value.provider)
  if (provider) {
    formData.value.imapHost = provider.imapHost
    formData.value.imapPort = provider.imapPort
    formData.value.imapSsl = provider.imapSsl
    formData.value.smtpHost = provider.smtpHost
    formData.value.smtpPort = provider.smtpPort
    formData.value.smtpSsl = provider.smtpSsl
  }
}

const testConnection = async () => {
  if (!formData.value.email || !formData.value.password || !formData.value.imapHost || !formData.value.smtpHost) {
    showNotification({
      type: 'warning',
      title: '信息不完整',
      message: '请填写完整的邮箱信息后再测试连接'
    })
    return
  }

  testing.value = true
  try {
    const testData: MailboxTestConnectionRequest = {
      email: formData.value.email,
      password: formData.value.password,
      imapHost: formData.value.imapHost,
      imapPort: formData.value.imapPort,
      imapSsl: formData.value.imapSsl,
      smtpHost: formData.value.smtpHost,
      smtpPort: formData.value.smtpPort,
      smtpSsl: formData.value.smtpSsl
    }

    const response = await mailboxApi.testConnection(testData)
    if (response.success && response.data) {
      const result = response.data
      if (result.imapSuccess && result.smtpSuccess) {
        showNotification({
          type: 'success',
          title: '连接测试成功',
          message: result.message
        })
      } else {
        showNotification({
          type: 'warning',
          title: '连接测试部分失败',
          message: `IMAP: ${result.imapSuccess ? '成功' : result.imapError}, SMTP: ${result.smtpSuccess ? '成功' : result.smtpError}`
        })
      }
    }
  } catch (error) {
    console.error('Failed to test connection:', error)
    showNotification({
      type: 'error',
      title: '连接测试失败',
      message: '无法测试邮箱连接'
    })
  } finally {
    testing.value = false
  }
}

const handleSubmit = async () => {
  saving.value = true
  try {
    const data = { ...formData.value }
    
    // 清理不需要的字段
    if (data.type === 'third') {
      delete data.domainId
    } else {
      delete data.provider
    }

    emit('save', data)
  } finally {
    saving.value = false
  }
}

// 监听邮箱数据变化
watch(() => props.mailbox, (newMailbox) => {
  if (newMailbox) {
    formData.value = {
      id: newMailbox.id,
      type: newMailbox.type,
      email: newMailbox.email,
      password: '', // 密码不回显
      provider: newMailbox.provider,
      domainId: newMailbox.domain_id,
      imapHost: newMailbox.imap_host,
      imapPort: newMailbox.imap_port,
      imapSsl: newMailbox.imap_ssl,
      smtpHost: newMailbox.smtp_host,
      smtpPort: newMailbox.smtp_port,
      smtpSsl: newMailbox.smtp_ssl,
      autoReceive: newMailbox.auto_receive,
      status: newMailbox.status
    }
  } else {
    // 重置表单
    formData.value = {
      type: 'third',
      email: '',
      password: '',
      provider: '',
      domainId: 0,
      imapHost: '',
      imapPort: 993,
      imapSsl: true,
      smtpHost: '',
      smtpPort: 587,
      smtpSsl: true,
      autoReceive: true,
      status: 1
    }
  }
}, { immediate: true })
</script>
