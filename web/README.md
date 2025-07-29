# 📧 邮件系统前端

基于 Vue 3 + TypeScript 的现代化邮件系统前端，采用毛玻璃美学设计。

## ✨ 特性

### 🎨 设计特色
- **毛玻璃美学** - 现代化的毛玻璃效果和深度层次
- **深色主题** - 护眼的深色主题，支持多种配色方案
- **响应式设计** - 完美适配移动端、平板和桌面端
- **流畅动画** - 丰富的微交互和页面切换动画
- **现代 UI** - 基于最新设计趋势的用户界面

### 🚀 技术特性
- **Vue 3** - 使用最新的 Composition API
- **TypeScript** - 完整的类型安全支持
- **Vite** - 极速的开发构建工具
- **Tailwind CSS** - 原子化 CSS 框架
- **Pinia** - 现代化的状态管理
- **Vue Router 4** - 声明式路由
- **@vueuse/core** - 强大的组合式工具库

### 📱 功能模块
- **用户认证** - 登录、注册、忘记密码
- **邮件管理** - 收件箱、发件箱、草稿箱、垃圾箱
- **邮件编辑** - 富文本编辑器、附件上传
- **用户设置** - 个人资料、主题设置、通知配置
- **管理后台** - 用户管理、系统监控、数据统计

## 🛠️ 技术栈

### 核心框架
- **Vue 3.3+** - 渐进式 JavaScript 框架
- **TypeScript 5.0+** - JavaScript 的超集
- **Vite 5.0+** - 下一代前端构建工具

### 状态管理
- **Pinia 2.1+** - Vue 官方推荐的状态管理库
- **@tanstack/vue-query** - 强大的数据获取和缓存库

### 样式方案
- **Tailwind CSS 3.3+** - 原子化 CSS 框架
- **@vueuse/motion** - Vue 动画库
- **@headlessui/vue** - 无样式 UI 组件

### 开发工具
- **ESLint** - 代码质量检查
- **Prettier** - 代码格式化
- **TypeScript** - 类型检查

## 📦 安装和运行

### 环境要求
- Node.js 18.0+
- npm 9.0+ 或 yarn 1.22+

### 安装依赖
```bash
cd web
npm install
```

### 开发模式
```bash
npm run dev
```
访问 http://localhost:3000

### 构建生产版本
```bash
npm run build
```

### 预览生产版本
```bash
npm run preview
```

### 代码检查
```bash
npm run lint
```

### 代码格式化
```bash
npm run format
```

## 📁 项目结构

```
web/
├── public/                 # 静态资源
├── src/
│   ├── assets/            # 资源文件
│   │   └── styles/        # 样式文件
│   ├── components/        # 组件
│   │   ├── ui/           # UI 基础组件
│   │   ├── email/        # 邮件相关组件
│   │   └── layout/       # 布局组件
│   ├── composables/       # 组合式函数
│   ├── stores/           # Pinia 状态管理
│   ├── utils/            # 工具函数
│   ├── views/            # 页面组件
│   │   ├── auth/         # 认证页面
│   │   ├── email/        # 邮件页面
│   │   ├── user/         # 用户页面
│   │   ├── admin/        # 管理页面
│   │   └── error/        # 错误页面
│   ├── types/            # TypeScript 类型定义
│   ├── router/           # 路由配置
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── index.html            # HTML 模板
├── package.json          # 项目配置
├── tailwind.config.js    # Tailwind 配置
├── tsconfig.json         # TypeScript 配置
├── vite.config.ts        # Vite 配置
└── README.md            # 项目说明
```

## 🎨 主题系统

### 预设主题
- 🌙 **经典深色** - 默认深色主题
- 🌊 **海洋蓝** - 蓝色系主题
- 🌸 **樱花粉** - 粉色系主题
- 🍃 **森林绿** - 绿色系主题
- 🔥 **火焰橙** - 橙色系主题
- 💜 **神秘紫** - 紫色系主题

### 自定义主题
支持用户自定义颜色配置，可在设置页面进行个性化定制。

## 📱 响应式设计

### 断点系统
- **xs**: < 475px (超小屏)
- **sm**: 475px - 640px (小屏)
- **md**: 640px - 768px (中屏)
- **lg**: 768px - 1024px (大屏)
- **xl**: 1024px - 1280px (超大屏)
- **2xl**: > 1280px (极大屏)

### 布局适配
- **移动端**: 单栏布局 + 底部导航
- **平板端**: 两栏布局 (侧边栏 + 主内容)
- **桌面端**: 三栏布局 (侧边栏 + 邮件列表 + 预览面板)

## 🔧 开发指南

### 组件开发
- 使用 Composition API 编写组件
- 遵循 TypeScript 类型安全
- 采用 Props + Emits 模式
- 支持 v-model 双向绑定

### 状态管理
- 使用 Pinia 管理全局状态
- 按功能模块划分 Store
- 支持 TypeScript 类型推导
- 提供 Composables 封装

### 样式规范
- 使用 Tailwind CSS 原子类
- 自定义 CSS 变量支持主题
- 响应式设计优先
- 毛玻璃效果统一封装

## 🚀 部署

### 构建优化
- 代码分割和懒加载
- 资源压缩和优化
- Tree Shaking 去除无用代码
- 现代浏览器优化

### 部署方式
- **静态部署**: Nginx、Apache
- **CDN 部署**: 阿里云、腾讯云
- **容器部署**: Docker
- **云平台**: Vercel、Netlify

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系

如有问题，请联系开发团队。
