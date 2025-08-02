<template>
  <Modal
    :visible="visible"
    :title="isEdit ? '编辑邮箱' : '创建邮箱'"
    size="lg"
    @close="$emit('close')"
  >
    <div class="space-y-6">




      <!-- 调试信息 -->
      <div style="background: #f0f0f0; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
        <p>调试信息:</p>
        <p>Email: {{ formData.email }}</p>
        <p>Password: {{ formData.password ? '***' : '(空)' }}</p>
        <p>Saving: {{ saving }}</p>
        <button @click="testClick" style="background: red; color: white; padding: 4px 8px; border: none; border-radius: 2px; margin: 4px;">
          调试按钮
        </button>
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
      <div class="flex justify-end pt-6 border-t border-glass-border" style="position: relative; z-index: 1000;">
        <div class="flex space-x-3" style="position: relative; z-index: 1001;">
          <Button
            type="button"
            variant="ghost"
            @click="$emit('close')"
          >
            取消
          </Button>
          <button
            type="button"
            @click="testClick"
            style="background: #28a745; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer; margin-right: 8px;"
          >
            测试
          </button>
          <button
            type="button"
            @click="handleButtonClick"
            :disabled="saving"
            style="background: #007bff; color: white; padding: 8px 16px; border: none; border-radius: 4px; cursor: pointer;"
          >
            {{ saving ? '创建中...' : (isEdit ? '更新' : '创建') }}
          </button>
        </div>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type {
  Mailbox,
  MailboxCreateRequest,
  MailboxUpdateRequest
} from '@/types'
import { mailboxApi } from '@/utils/api'
import { useNotification } from '@/composables/useNotification'

// Icons
import {
  EyeIcon,
  EyeSlashIcon
} from '@heroicons/vue/24/outline'

// Components
import Modal from '@/components/ui/Modal.vue'
import Button from '@/components/ui/Button.vue'

// Props
interface Props {
  visible: boolean
  mailbox?: Mailbox | null
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
  email: '',
  password: '',
  autoReceive: true,
  status: 1
})

// 方法
const testClick = () => {
  console.log('Test button clicked!')
  alert('Test button works!')
}

const handleButtonClick = async () => {
  console.log('Button clicked!')
  console.log('Form data:', formData.value)

  // 检查必填字段
  if (!formData.value.email || !formData.value.password) {
    console.log('Missing required fields')
    return
  }

  console.log('Saving mailbox...')
  saving.value = true
  try {
    const data = { ...formData.value }
    console.log('Emitting save event with data:', data)
    emit('save', data)
  } finally {
    saving.value = false
  }
}





// 组件挂载时的调试信息
onMounted(() => {
  console.log('MailboxModal mounted')
  console.log('Props:', props)
  console.log('Form data:', formData.value)
})

// 监听邮箱数据变化
watch(() => props.mailbox, (newMailbox) => {
  console.log('Mailbox data changed:', newMailbox)
  if (newMailbox) {
    formData.value = {
      id: newMailbox.id,
      domainId: newMailbox.domainId,
      email: newMailbox.email,
      password: '', // 密码不回显
      autoReceive: newMailbox.autoReceive,
      status: newMailbox.status
    }
  } else {
    // 重置表单
    formData.value = {
      email: '',
      password: '',
      autoReceive: true,
      status: 1
    }
  }
  console.log('Updated form data:', formData.value)
}, { immediate: true })
</script>
