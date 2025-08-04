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
    sound: boolean
    newEmail: boolean
    importantEmail: boolean
    securityAlerts: boolean
  }
  privacy: {
    showOnlineStatus: boolean
    allowDirectMessages: boolean
    showReadReceipts: boolean
  }
  security: {
    twoFactorEnabled: boolean
    sessionTimeout: number
    loginNotifications: boolean
  }
  emailFilters: EmailFilter[]
}

// 邮件过滤规则
export interface EmailFilter {
  id: string
  name: string
  enabled: boolean
  conditions: FilterCondition[]
  actions: FilterAction[]
  priority: number
  createdAt: string
  updatedAt: string
}

export interface FilterCondition {
  field: 'from' | 'to' | 'subject' | 'body' | 'attachment'
  operator: 'contains' | 'equals' | 'startsWith' | 'endsWith' | 'regex'
  value: string
  caseSensitive?: boolean
}

export interface FilterAction {
  type: 'move' | 'label' | 'delete' | 'markRead' | 'markImportant' | 'forward'
  value?: string
}

// API密钥
export interface ApiKey {
  id: string
  name: string
  key: string
  permissions: string[]
  lastUsed?: string
  expiresAt?: string
  createdAt: string
  enabled: boolean
}

// 用户资料更新请求
export interface UserProfileUpdateRequest {
  nickname?: string
  avatar?: string
  bio?: string
  timezone?: string
  language?: string
}

// 安全设置更新请求
export interface SecuritySettingsUpdateRequest {
  currentPassword: string
  newPassword?: string
  twoFactorEnabled?: boolean
  sessionTimeout?: number
  loginNotifications?: boolean
}

// 通知设置更新请求
export interface NotificationSettingsUpdateRequest {
  email: boolean
  desktop: boolean
  mobile: boolean
  sound: boolean
  newEmail: boolean
  importantEmail: boolean
  securityAlerts: boolean
}

// 邮件相关类型
export interface Email {
  id: string
  userId: number
  mailboxId: number
  subject: string
  fromEmail: string
  toEmails: string
  ccEmail?: string
  bccEmail?: string
  content: string
  contentType: 'text' | 'html'
  attachments?: Attachment[]
  status: number
  type: 'inbox' | 'sent' | 'draft'
  isRead: boolean
  isStarred: boolean
  isImportant?: boolean
  labels?: string[]
  createdAt: string
  updatedAt: string
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
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 邮件列表查询参数
export interface EmailListParams {
  userId?: number
  mailboxId?: number
  messageId?: string
  subject?: string
  fromEmail?: string
  toEmails?: string
  direction?: 'sent' | 'received'
  isRead?: boolean
  isStarred?: boolean
  contentType?: string
  page?: number
  pageSize?: number
  createdAtStart?: string
  createdAtEnd?: string
  updatedAtStart?: string
  updatedAtEnd?: string
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
  userId: number
  domainId: number
  email: string
  name?: string
  autoReceive: boolean
  status: number
  unreadCount?: number
  lastSyncAt?: string
  createdAt: string
  updatedAt: string
}

export interface MailboxCreateRequest {
  domainId?: number
  email: string
  password: string
  autoReceive?: boolean
  status?: number
}

export interface MailboxUpdateRequest {
  id: number
  domainId?: number
  email?: string
  password?: string
  autoReceive?: boolean
  status?: number
}

export interface MailboxListRequest {
  userId?: number
  domainId?: number
  email?: string
  status?: number
  autoReceive?: boolean
  page?: number
  pageSize?: number
  startTime?: string
  endTime?: string
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
}

// 邮件发送相关类型
export interface EmailSendRequest {
  mailboxId: number
  subject: string
  fromEmail: string
  toEmail: string
  ccEmail?: string
  bccEmail?: string
  content: string
  contentType: 'text' | 'html'
  attachments?: string
}

export interface EmailSendResponse {
  success: boolean
  message: string
  emailId: number
  sentAt: string
}

// 验证码相关类型
export interface VerificationCode {
  id: number
  userId: number
  emailId: number
  code: string
  source: string
  type?: string
  context?: string
  confidence?: number
  pattern?: string
  description?: string
  isUsed: boolean
  isExpired: boolean
  usedAt?: string
  expiresAt: string
  createdAt: string
}

export interface VerificationCodeExtractRequest {
  emailId: number
}

export interface VerificationCodeResult {
  code: string
  type: string
  context: string
  confidence: number
  position: number
  length: number
  pattern: string
  description: string
}

export interface VerificationCodeExtractResponse {
  emailId: number
  subject: string
  fromEmail: string
  extractedAt: string
  codes: VerificationCodeResult[]
}

export interface VerificationCodeBatchExtractRequest {
  emailIds: number[]
}

export interface VerificationCodeBatchExtractResponse {
  totalEmails: number
  processedEmails: number
  extractedCodes: number
  results: VerificationCodeExtractResponse[]
  errors: string[]
}

export interface VerificationCodeStats {
  totalCodes: number
  usedCodes: number
  unusedCodes: number
  todayCodes: number
  typeStats: Array<{
    type: string
    count: number
  }>
  sourceStats: Array<{
    fromEmail: string
    count: number
  }>
}

export interface VerificationCodeListParams {
  page?: number
  pageSize?: number
  emailId?: number
  code?: string
  source?: string
  isUsed?: boolean
  isExpired?: boolean
  createdAtStart?: string
  createdAtEnd?: string
}

// 管理员相关类型
export interface AdminLoginRequest {
  username: string
  password: string
}

export interface AdminLoginResponse {
  token: string
  expiresAt: string
  admin: AdminInfo
}

export interface AdminInfo {
  id: number
  username: string
  nickname: string
  email: string
  role: string
  status: number
}

export interface AdminSystemStats {
  userStats: AdminUserStats
  emailStats: AdminEmailStats
  mailboxStats: AdminMailboxStats
  systemStats: AdminSystemInfo
}

export interface AdminUserStats {
  totalUsers: number
  activeUsers: number
  newUsers: number
  onlineUsers: number
}

export interface AdminEmailStats {
  totalEmails: number
  todayEmails: number
  sentEmails: number
  receivedEmails: number
}

export interface AdminMailboxStats {
  totalMailboxes: number
  activeMailboxes: number
  imapMailboxes: number
  pop3Mailboxes: number
}

export interface AdminSystemInfo {
  version: string
  startTime: string
  uptime: string
  goVersion: string
  platform: string
  cpuUsage: number
  memUsage: number
  diskUsage: number
}

export interface AdminDashboard {
  stats: AdminSystemStats
  recentUsers: User[]
  recentLogs: OperationLog[]
  systemAlerts: SystemAlert[]
}

export interface SystemAlert {
  id: number
  type: string
  level: string
  title: string
  message: string
  status: string
  createdAt: string
}

export interface AdminSystemSettings {
  siteName: string
  siteDescription: string
  siteLogo: string
  allowRegister: boolean
  requireInvite: boolean
  defaultQuota: number
  maxMailboxes: number
  updatedAt?: string
}

export interface DailyStats {
  date: string
  count: number
}

export interface UserActivityStats {
  userId: number
  username: string
  emailCount: number
  mailboxCount: number
  lastLoginAt?: string
}

export interface SystemHealth {
  status: string
  score: number
  cpuUsage: number
  memUsage: number
  diskUsage: number
  uptime: string
  checkedAt: string
}
