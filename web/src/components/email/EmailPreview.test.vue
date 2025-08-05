<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">EmailPreview 组件测试</h1>
    
    <div class="flex space-x-4">
      <!-- 测试按钮 -->
      <div class="w-64 space-y-2">
        <Button @click="showTestEmail" variant="primary" class="w-full">
          显示测试邮件
        </Button>
        <Button @click="showHtmlEmail" variant="secondary" class="w-full">
          显示HTML邮件
        </Button>
        <Button @click="showEmailWithAttachments" variant="secondary" class="w-full">
          显示带附件邮件
        </Button>
        <Button @click="clearEmail" variant="ghost" class="w-full">
          清空选择
        </Button>
      </div>
      
      <!-- 邮件预览 -->
      <EmailPreview 
        :email="selectedEmail" 
        @email-updated="handleEmailUpdated"
        @email-deleted="handleEmailDeleted"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Email } from '@/types'
import EmailPreview from './EmailPreview.vue'
import Button from '@/components/ui/Button.vue'

const selectedEmail = ref<Email | null>(null)

const testEmail: Email = {
  id: '1',
  userId: 1,
  mailboxId: 1,
  subject: '测试邮件主题',
  fromEmail: 'sender@example.com',
  toEmails: 'recipient@example.com',
  content: '这是一封测试邮件的内容。\n\n包含多行文本和一些基本格式。',
  contentType: 'text',
  direction: 'received',
  isRead: false,
  isStarred: false,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString()
}

const htmlEmail: Email = {
  id: '2',
  userId: 1,
  mailboxId: 1,
  subject: 'HTML格式邮件测试',
  fromEmail: 'marketing@company.com',
  toEmails: 'user@example.com',
  ccEmail: 'manager@company.com',
  content: `
    <h2>欢迎使用我们的服务</h2>
    <p>这是一封<strong>HTML格式</strong>的邮件，包含：</p>
    <ul>
      <li>列表项目1</li>
      <li>列表项目2</li>
      <li><a href="https://example.com">链接示例</a></li>
    </ul>
    <blockquote>
      这是一个引用块的示例
    </blockquote>
    <p>感谢您的使用！</p>
  `,
  contentType: 'html',
  direction: 'received',
  isRead: true,
  isStarred: true,
  createdAt: new Date(Date.now() - 86400000).toISOString(),
  updatedAt: new Date(Date.now() - 86400000).toISOString()
}

const emailWithAttachments: Email = {
  id: '3',
  userId: 1,
  mailboxId: 1,
  subject: '带附件的邮件',
  fromEmail: 'files@company.com',
  toEmails: 'user@example.com',
  content: '请查看附件中的重要文档。',
  contentType: 'text',
  direction: 'received',
  isRead: false,
  isStarred: false,
  attachments: [
    {
      id: '1',
      filename: '重要文档.pdf',
      size: 1024000,
      contentType: 'application/pdf'
    },
    {
      id: '2',
      filename: '图片.jpg',
      size: 512000,
      contentType: 'image/jpeg'
    }
  ],
  createdAt: new Date(Date.now() - 3600000).toISOString(),
  updatedAt: new Date(Date.now() - 3600000).toISOString()
}

const showTestEmail = () => {
  selectedEmail.value = testEmail
}

const showHtmlEmail = () => {
  selectedEmail.value = htmlEmail
}

const showEmailWithAttachments = () => {
  selectedEmail.value = emailWithAttachments
}

const clearEmail = () => {
  selectedEmail.value = null
}

const handleEmailUpdated = (email: Email) => {
  console.log('Email updated:', email)
  selectedEmail.value = email
}

const handleEmailDeleted = (emailId: string) => {
  console.log('Email deleted:', emailId)
  selectedEmail.value = null
}
</script>
