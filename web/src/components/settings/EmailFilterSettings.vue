<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-text-primary mb-2">邮件过滤</h2>
        <p class="text-text-secondary">创建规则自动处理收到的邮件</p>
      </div>
      <Button
        variant="primary"
        @click="showCreateModal = true"
      >
        <PlusIcon class="w-4 h-4 mr-2" />
        新建规则
      </Button>
    </div>

    <!-- 过滤规则列表 -->
    <GlassCard padding="lg">
      <div v-if="filters.length === 0" class="text-center py-12">
        <FunnelIcon class="w-16 h-16 text-text-secondary mx-auto mb-4" />
        <h3 class="text-lg font-medium text-text-primary mb-2">暂无过滤规则</h3>
        <p class="text-text-secondary mb-4">
          创建过滤规则来自动整理您的邮件
        </p>
        <Button
          variant="primary"
          @click="showCreateModal = true"
        >
          创建第一个规则
        </Button>
      </div>

      <div v-else class="space-y-4">
        <div
          v-for="filter in filters"
          :key="filter.id"
          class="p-4 bg-white/5 rounded-lg border border-glass-border"
        >
          <!-- 规则头部 -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center space-x-3">
              <Switch
                v-model="filter.enabled"
                @change="toggleFilter(filter.id, filter.enabled)"
              />
              <div>
                <h4 class="text-md font-medium text-text-primary">
                  {{ filter.name }}
                </h4>
                <p class="text-sm text-text-secondary">
                  优先级: {{ filter.priority }}
                </p>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <Button
                variant="ghost"
                size="sm"
                @click="editFilter(filter)"
              >
                <PencilIcon class="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                @click="deleteFilter(filter.id)"
              >
                <TrashIcon class="w-4 h-4" />
              </Button>
            </div>
          </div>

          <!-- 规则条件 -->
          <div class="mb-3">
            <h5 class="text-sm font-medium text-text-primary mb-2">条件:</h5>
            <div class="space-y-1">
              <div
                v-for="(condition, index) in filter.conditions"
                :key="index"
                class="text-sm text-text-secondary"
              >
                {{ formatCondition(condition) }}
              </div>
            </div>
          </div>

          <!-- 规则动作 -->
          <div>
            <h5 class="text-sm font-medium text-text-primary mb-2">动作:</h5>
            <div class="space-y-1">
              <div
                v-for="(action, index) in filter.actions"
                :key="index"
                class="text-sm text-text-secondary"
              >
                {{ formatAction(action) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>

    <!-- 创建/编辑规则模态框 -->
    <Modal
      v-model="showCreateModal"
      :title="editingFilter ? '编辑过滤规则' : '创建过滤规则'"
      size="lg"
      @confirm="saveFilter"
      :loading="saving"
    >
      <div class="space-y-6">
        <!-- 规则名称 -->
        <Input
          v-model="filterForm.name"
          label="规则名称"
          placeholder="请输入规则名称"
          :error="errors.name"
          required
        />

        <!-- 优先级 -->
        <Select
          v-model="filterForm.priority"
          label="优先级"
          :options="priorityOptions"
          help="数字越小优先级越高"
        />

        <!-- 条件设置 -->
        <div>
          <h4 class="text-md font-medium text-text-primary mb-3">匹配条件</h4>
          <div class="space-y-3">
            <div
              v-for="(condition, index) in filterForm.conditions"
              :key="index"
              class="flex items-end space-x-3 p-3 bg-white/5 rounded-lg"
            >
              <!-- 字段 -->
              <Select
                v-model="condition.field"
                label="字段"
                :options="fieldOptions"
                class="flex-1"
              />

              <!-- 操作符 -->
              <Select
                v-model="condition.operator"
                label="条件"
                :options="operatorOptions"
                class="flex-1"
              />

              <!-- 值 -->
              <Input
                v-model="condition.value"
                label="值"
                placeholder="请输入匹配值"
                class="flex-1"
              />

              <!-- 删除按钮 -->
              <Button
                variant="ghost"
                size="sm"
                @click="removeCondition(index)"
                class="mb-1"
              >
                <TrashIcon class="w-4 h-4" />
              </Button>
            </div>

            <!-- 添加条件 -->
            <Button
              variant="ghost"
              @click="addCondition"
            >
              <PlusIcon class="w-4 h-4 mr-2" />
              添加条件
            </Button>
          </div>
        </div>

        <!-- 动作设置 -->
        <div>
          <h4 class="text-md font-medium text-text-primary mb-3">执行动作</h4>
          <div class="space-y-3">
            <div
              v-for="(action, index) in filterForm.actions"
              :key="index"
              class="flex items-end space-x-3 p-3 bg-white/5 rounded-lg"
            >
              <!-- 动作类型 -->
              <Select
                v-model="action.type"
                label="动作"
                :options="actionOptions"
                class="flex-1"
                @change="onActionTypeChange(action, index)"
              />

              <!-- 动作值 -->
              <Input
                v-if="needsActionValue(action.type)"
                v-model="action.value"
                label="目标"
                :placeholder="getActionValuePlaceholder(action.type)"
                class="flex-1"
              />

              <!-- 删除按钮 -->
              <Button
                variant="ghost"
                size="sm"
                @click="removeAction(index)"
                class="mb-1"
              >
                <TrashIcon class="w-4 h-4" />
              </Button>
            </div>

            <!-- 添加动作 -->
            <Button
              variant="ghost"
              @click="addAction"
            >
              <PlusIcon class="w-4 h-4 mr-2" />
              添加动作
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
import type { EmailFilter, FilterCondition, FilterAction } from '@/types'

import GlassCard from '@/components/ui/GlassCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Switch from '@/components/ui/Switch.vue'
import Modal from '@/components/ui/Modal.vue'

import {
  PlusIcon,
  FunnelIcon,
  PencilIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'

// 通知
const { showSuccess, showError } = useNotification()

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const showCreateModal = ref(false)
const editingFilter = ref<EmailFilter | null>(null)

const filters = ref<EmailFilter[]>([])

// 表单数据
const filterForm = reactive({
  name: '',
  priority: 1,
  enabled: true,
  conditions: [] as FilterCondition[],
  actions: [] as FilterAction[]
})

const errors = reactive({
  name: ''
})

// 选项数据
const priorityOptions = [
  { label: '最高 (1)', value: 1 },
  { label: '高 (2)', value: 2 },
  { label: '中 (3)', value: 3 },
  { label: '低 (4)', value: 4 },
  { label: '最低 (5)', value: 5 }
]

const fieldOptions = [
  { label: '发件人', value: 'from' },
  { label: '收件人', value: 'to' },
  { label: '主题', value: 'subject' },
  { label: '正文', value: 'body' },
  { label: '附件', value: 'attachment' }
]

const operatorOptions = [
  { label: '包含', value: 'contains' },
  { label: '等于', value: 'equals' },
  { label: '开始于', value: 'startsWith' },
  { label: '结束于', value: 'endsWith' },
  { label: '正则表达式', value: 'regex' }
]

const actionOptions = [
  { label: '移动到文件夹', value: 'move' },
  { label: '添加标签', value: 'label' },
  { label: '删除邮件', value: 'delete' },
  { label: '标记为已读', value: 'markRead' },
  { label: '标记为重要', value: 'markImportant' },
  { label: '转发到', value: 'forward' }
]

// 生命周期
onMounted(() => {
  loadFilters()
})

// 方法
const loadFilters = async () => {
  try {
    loading.value = true
    filters.value = await userSettingsApi.getEmailFilters()
  } catch (error) {
    showError('加载过滤规则失败')
  } finally {
    loading.value = false
  }
}

const addCondition = () => {
  filterForm.conditions.push({
    field: 'from',
    operator: 'contains',
    value: ''
  })
}

const removeCondition = (index: number) => {
  filterForm.conditions.splice(index, 1)
}

const addAction = () => {
  filterForm.actions.push({
    type: 'move',
    value: ''
  })
}

const removeAction = (index: number) => {
  filterForm.actions.splice(index, 1)
}

const needsActionValue = (actionType: string) => {
  return ['move', 'label', 'forward'].includes(actionType)
}

const getActionValuePlaceholder = (actionType: string) => {
  switch (actionType) {
    case 'move':
      return '文件夹名称'
    case 'label':
      return '标签名称'
    case 'forward':
      return '转发邮箱地址'
    default:
      return ''
  }
}

const onActionTypeChange = (action: FilterAction, index: number) => {
  if (!needsActionValue(action.type)) {
    action.value = undefined
  }
}

const formatCondition = (condition: FilterCondition) => {
  const fieldMap: Record<string, string> = {
    from: '发件人',
    to: '收件人',
    subject: '主题',
    body: '正文',
    attachment: '附件'
  }
  
  const operatorMap: Record<string, string> = {
    contains: '包含',
    equals: '等于',
    startsWith: '开始于',
    endsWith: '结束于',
    regex: '匹配正则'
  }

  return `${fieldMap[condition.field]} ${operatorMap[condition.operator]} "${condition.value}"`
}

const formatAction = (action: FilterAction) => {
  const actionMap: Record<string, string> = {
    move: '移动到',
    label: '添加标签',
    delete: '删除邮件',
    markRead: '标记为已读',
    markImportant: '标记为重要',
    forward: '转发到'
  }

  const actionText = actionMap[action.type]
  return action.value ? `${actionText} "${action.value}"` : actionText
}

const editFilter = (filter: EmailFilter) => {
  editingFilter.value = filter
  filterForm.name = filter.name
  filterForm.priority = filter.priority
  filterForm.enabled = filter.enabled
  filterForm.conditions = [...filter.conditions]
  filterForm.actions = [...filter.actions]
  showCreateModal.value = true
}

const resetForm = () => {
  editingFilter.value = null
  filterForm.name = ''
  filterForm.priority = 1
  filterForm.enabled = true
  filterForm.conditions = []
  filterForm.actions = []
  errors.name = ''
}

const validateForm = () => {
  errors.name = ''
  
  if (!filterForm.name.trim()) {
    errors.name = '请输入规则名称'
    return false
  }

  if (filterForm.conditions.length === 0) {
    showError('请至少添加一个匹配条件')
    return false
  }

  if (filterForm.actions.length === 0) {
    showError('请至少添加一个执行动作')
    return false
  }

  // 验证条件
  for (const condition of filterForm.conditions) {
    if (!condition.value.trim()) {
      showError('请填写所有条件的匹配值')
      return false
    }
  }

  // 验证动作
  for (const action of filterForm.actions) {
    if (needsActionValue(action.type) && !action.value?.trim()) {
      showError('请填写所有动作的目标值')
      return false
    }
  }

  return true
}

const saveFilter = async () => {
  if (!validateForm()) return

  try {
    saving.value = true

    const filterData = {
      name: filterForm.name,
      priority: filterForm.priority,
      enabled: filterForm.enabled,
      conditions: filterForm.conditions,
      actions: filterForm.actions
    }

    if (editingFilter.value) {
      // 更新现有规则
      await userSettingsApi.updateEmailFilter(editingFilter.value.id, filterData)
      showSuccess('过滤规则更新成功')
    } else {
      // 创建新规则
      await userSettingsApi.createEmailFilter(filterData)
      showSuccess('过滤规则创建成功')
    }

    showCreateModal.value = false
    resetForm()
    loadFilters()
  } catch (error) {
    showError('保存过滤规则失败')
  } finally {
    saving.value = false
  }
}

const toggleFilter = async (id: string, enabled: boolean) => {
  try {
    await userSettingsApi.toggleEmailFilter(id, enabled)
    showSuccess(enabled ? '规则已启用' : '规则已禁用')
  } catch (error) {
    showError('操作失败')
    // 恢复状态
    const filter = filters.value.find(f => f.id === id)
    if (filter) {
      filter.enabled = !enabled
    }
  }
}

const deleteFilter = async (id: string) => {
  if (!confirm('确定要删除这个过滤规则吗？')) return

  try {
    await userSettingsApi.deleteEmailFilter(id)
    filters.value = filters.value.filter(f => f.id !== id)
    showSuccess('过滤规则删除成功')
  } catch (error) {
    showError('删除过滤规则失败')
  }
}

// 监听模态框关闭
const handleModalClose = () => {
  resetForm()
}
</script>
