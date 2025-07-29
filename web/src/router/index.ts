import { createRouter, createWebHistory } from 'vue-router'
import type { RouteMeta } from '@/types'
import { useAuthStore } from '@/stores/auth'

// 路由组件懒加载
const Login = () => import('@/views/auth/Login.vue')
const Register = () => import('@/views/auth/Register.vue')
const ForgotPassword = () => import('@/views/auth/ForgotPassword.vue')

const Inbox = () => import('@/views/email/Inbox.vue')
const Compose = () => import('@/views/email/Compose.vue')
const EmailDetail = () => import('@/views/email/EmailDetail.vue')
const Sent = () => import('@/views/email/Sent.vue')
const Drafts = () => import('@/views/email/Drafts.vue')
const Trash = () => import('@/views/email/Trash.vue')

const Settings = () => import('@/views/user/Settings.vue')
const Profile = () => import('@/views/user/Profile.vue')

const AdminDashboard = () => import('@/views/admin/Dashboard.vue')
const AdminUsers = () => import('@/views/admin/Users.vue')
const AdminEmails = () => import('@/views/admin/Emails.vue')
const AdminSettings = () => import('@/views/admin/Settings.vue')

const NotFound = () => import('@/views/error/NotFound.vue')

declare module 'vue-router' {
  interface RouteMeta extends RouteMeta {}
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // 根路径重定向
    {
      path: '/',
      redirect: '/inbox'
    },
    
    // 认证相关路由
    {
      path: '/auth',
      children: [
        {
          path: 'login',
          name: 'Login',
          component: Login,
          meta: {
            title: '登录',
            layout: 'auth',
            requiresAuth: false
          }
        },
        {
          path: 'register',
          name: 'Register',
          component: Register,
          meta: {
            title: '注册',
            layout: 'auth',
            requiresAuth: false
          }
        },
        {
          path: 'forgot-password',
          name: 'ForgotPassword',
          component: ForgotPassword,
          meta: {
            title: '忘记密码',
            layout: 'auth',
            requiresAuth: false
          }
        }
      ]
    },
    
    // 邮件相关路由
    {
      path: '/inbox',
      name: 'Inbox',
      component: Inbox,
      meta: {
        title: '收件箱',
        requiresAuth: true,
        keepAlive: true
      }
    },
    {
      path: '/compose',
      name: 'Compose',
      component: Compose,
      meta: {
        title: '写邮件',
        requiresAuth: true
      }
    },
    {
      path: '/email/:id',
      name: 'EmailDetail',
      component: EmailDetail,
      meta: {
        title: '邮件详情',
        requiresAuth: true
      }
    },
    {
      path: '/sent',
      name: 'Sent',
      component: Sent,
      meta: {
        title: '已发送',
        requiresAuth: true,
        keepAlive: true
      }
    },
    {
      path: '/drafts',
      name: 'Drafts',
      component: Drafts,
      meta: {
        title: '草稿箱',
        requiresAuth: true,
        keepAlive: true
      }
    },
    {
      path: '/trash',
      name: 'Trash',
      component: Trash,
      meta: {
        title: '垃圾箱',
        requiresAuth: true,
        keepAlive: true
      }
    },
    
    // 用户相关路由
    {
      path: '/settings',
      name: 'Settings',
      component: Settings,
      meta: {
        title: '设置',
        requiresAuth: true
      }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      meta: {
        title: '个人资料',
        requiresAuth: true
      }
    },
    
    // 管理员路由
    {
      path: '/admin',
      meta: {
        requiresAuth: true,
        requiresAdmin: true,
        layout: 'admin'
      },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard'
        },
        {
          path: 'dashboard',
          name: 'AdminDashboard',
          component: AdminDashboard,
          meta: {
            title: '管理仪表板',
            requiresAuth: true,
            requiresAdmin: true,
            layout: 'admin'
          }
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: AdminUsers,
          meta: {
            title: '用户管理',
            requiresAuth: true,
            requiresAdmin: true,
            layout: 'admin'
          }
        },
        {
          path: 'emails',
          name: 'AdminEmails',
          component: AdminEmails,
          meta: {
            title: '邮件监控',
            requiresAuth: true,
            requiresAdmin: true,
            layout: 'admin'
          }
        },
        {
          path: 'settings',
          name: 'AdminSettings',
          component: AdminSettings,
          meta: {
            title: '系统设置',
            requiresAuth: true,
            requiresAdmin: true,
            layout: 'admin'
          }
        }
      ]
    },
    
    // 404 页面
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFound,
      meta: {
        title: '页面未找到'
      }
    }
  ],
  
  // 滚动行为
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 邮件系统`
  }
  
  // 检查认证状态
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      // 保存目标路由，登录后跳转
      authStore.setRedirectPath(to.fullPath)
      next('/auth/login')
      return
    }
    
    // 检查管理员权限
    if (to.meta.requiresAdmin && !authStore.isAdmin) {
      next('/inbox')
      return
    }
  }
  
  // 已登录用户访问认证页面，重定向到收件箱
  if (!to.meta.requiresAuth && authStore.isAuthenticated) {
    if (to.path.startsWith('/auth')) {
      next('/inbox')
      return
    }
  }
  
  next()
})

export default router
