// 用户相关类型
export interface User {
  id: string
  username: string
  email: string
  nickname?: string
  avatar?: string
  role: 'user' | 'admin'
  status: 'active' | 'inactive' | 'pending'
  createdAt: string
  lastLoginAt?: string
  settings?: UserSettings
}

export interface UserSettings {
  theme: string
  language: string
  emailSignature?: string
  notifications: {
    email: boolean
    desktop: boolean
    mobile: boolean
  }
  privacy: {
    showOnlineStatus: boolean
    allowDirectMessages: boolean
  }
}

// 邮件相关类型
export interface Email {
  id: string
  subject: string
  content: string
  htmlContent?: string
  from: EmailAddress
  to: EmailAddress[]
  cc?: EmailAddress[]
  bcc?: EmailAddress[]
  attachments?: Attachment[]
  isRead: boolean
  isStarred: boolean
  isImportant: boolean
  folder: EmailFolder
  labels?: string[]
  createdAt: string
  updatedAt: string
  size: number
}

export interface EmailAddress {
  email: string
  name?: string
}

export interface Attachment {
  id: string
  filename: string
  size: number
  mimeType: string
  url: string
  thumbnailUrl?: string
}

export interface EmailFolder {
  id: string
  name: string
  type: 'inbox' | 'sent' | 'drafts' | 'trash' | 'spam' | 'custom'
  unreadCount: number
  totalCount: number
  icon?: string
  color?: string
}

// 邮件编辑相关
export interface EmailDraft {
  id?: string
  to: EmailAddress[]
  cc?: EmailAddress[]
  bcc?: EmailAddress[]
  subject: string
  content: string
  htmlContent?: string
  attachments?: File[]
  isScheduled?: boolean
  scheduledAt?: string
  priority?: 'low' | 'normal' | 'high'
}

// API 响应类型
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
  code?: number
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// 邮件列表查询参数
export interface EmailListParams {
  folder?: string
  page?: number
  pageSize?: number
  search?: string
  isRead?: boolean
  isStarred?: boolean
  sortBy?: 'date' | 'subject' | 'from'
  sortOrder?: 'asc' | 'desc'
  labels?: string[]
}

// 主题相关类型
export interface Theme {
  name: string
  displayName: string
  colors: {
    primary: string
    secondary: string
    background: string
    surface: string
    text: {
      primary: string
      secondary: string
      disabled: string
    }
    glass: {
      light: string
      medium: string
      heavy: string
      border: string
    }
    status: {
      success: string
      warning: string
      error: string
      info: string
    }
  }
}

// 通知类型
export interface Notification {
  id: string
  type: 'success' | 'warning' | 'error' | 'info'
  title: string
  message?: string
  duration?: number
  actions?: NotificationAction[]
  createdAt: string
}

export interface NotificationAction {
  label: string
  action: () => void
  style?: 'primary' | 'secondary'
}

// 统计数据类型
export interface EmailStats {
  totalEmails: number
  unreadEmails: number
  todayEmails: number
  weeklyEmails: number
  monthlyEmails: number
  storageUsed: number
  storageLimit: number
}

export interface AdminStats {
  totalUsers: number
  activeUsers: number
  totalEmails: number
  systemStatus: 'healthy' | 'warning' | 'error'
  serverLoad: number
  memoryUsage: number
  diskUsage: number
}

// 表单验证类型
export interface ValidationRule {
  required?: boolean
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  custom?: (value: any) => boolean | string
}

export interface FormField {
  name: string
  label: string
  type: 'text' | 'email' | 'password' | 'textarea' | 'select' | 'checkbox'
  value: any
  rules?: ValidationRule[]
  options?: { label: string; value: any }[]
  placeholder?: string
  disabled?: boolean
}

// 路由元信息类型
export interface RouteMeta {
  title?: string
  requiresAuth?: boolean
  requiresAdmin?: boolean
  layout?: 'default' | 'auth' | 'admin'
  keepAlive?: boolean
}

// 组件 Props 类型
export interface GlassCardProps {
  level?: 1 | 2 | 3
  hover?: boolean
  className?: string
  padding?: 'none' | 'sm' | 'md' | 'lg'
}

export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  icon?: string
  iconPosition?: 'left' | 'right'
}

// 事件类型
export interface EmailEvent {
  type: 'read' | 'unread' | 'star' | 'unstar' | 'delete' | 'move' | 'label'
  emailId: string
  data?: any
}

export interface UserEvent {
  type: 'login' | 'logout' | 'register' | 'update' | 'delete'
  userId: string
  data?: any
}

// 邮箱管理相关类型
export interface Mailbox {
  id: number
  user_id: number
  domain_id: number
  email: string
  type: 'self' | 'third'
  provider: string
  imap_host: string
  imap_port: number
  imap_ssl: boolean
  smtp_host: string
  smtp_port: number
  smtp_ssl: boolean
  auto_receive: boolean
  status: number
  last_sync_at?: string
  created_at: string
  updated_at: string
}

export interface MailboxCreateRequest {
  domainId?: number
  email: string
  password: string
  type: 'self' | 'third'
  provider?: string
  imapHost?: string
  imapPort?: number
  imapSsl?: boolean
  smtpHost?: string
  smtpPort?: number
  smtpSsl?: boolean
  autoReceive?: boolean
  status?: number
}

export interface MailboxUpdateRequest {
  id: number
  domainId?: number
  email?: string
  password?: string
  type?: 'self' | 'third'
  provider?: string
  imapHost?: string
  imapPort?: number
  imapSsl?: boolean
  smtpHost?: string
  smtpPort?: number
  smtpSsl?: boolean
  autoReceive?: boolean
  status?: number
}

export interface MailboxListRequest {
  userId?: number
  domainId?: number
  email?: string
  type?: string
  provider?: string
  status?: number
  autoReceive?: boolean
  page?: number
  pageSize?: number
  startTime?: string
  endTime?: string
}

export interface MailboxProvider {
  provider: string
  imapHost: string
  imapPort: number
  imapSsl: boolean
  smtpHost: string
  smtpPort: number
  smtpSsl: boolean
}

export interface MailboxTestConnectionRequest {
  email: string
  password: string
  imapHost: string
  imapPort: number
  imapSsl: boolean
  smtpHost: string
  smtpPort: number
  smtpSsl: boolean
}

export interface MailboxTestConnectionResponse {
  imapSuccess: boolean
  smtpSuccess: boolean
  imapError: string
  smtpError: string
  message: string
}

export interface MailboxSyncRequest {
  id: number
  forceSync?: boolean
  syncDays?: number
}

export interface MailboxSyncResponse {
  success: boolean
  message: string
  syncCount: number
  errorCount: number
  lastSyncAt: string
}

export interface MailboxStats {
  totalMailboxes: number
  activeMailboxes: number
  selfMailboxes: number
  thirdMailboxes: number
}
