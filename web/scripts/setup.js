#!/usr/bin/env node

const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

console.log('🚀 开始设置邮件系统前端项目...\n')

// 检查 Node.js 版本
const nodeVersion = process.version
const majorVersion = parseInt(nodeVersion.slice(1).split('.')[0])

if (majorVersion < 18) {
  console.error('❌ 需要 Node.js 18.0 或更高版本')
  console.error(`   当前版本: ${nodeVersion}`)
  process.exit(1)
}

console.log(`✅ Node.js 版本检查通过: ${nodeVersion}`)

// 检查包管理器
let packageManager = 'npm'
try {
  execSync('yarn --version', { stdio: 'ignore' })
  packageManager = 'yarn'
  console.log('✅ 检测到 Yarn，将使用 Yarn 安装依赖')
} catch {
  console.log('✅ 将使用 npm 安装依赖')
}

// 安装依赖
console.log('\n📦 安装项目依赖...')
try {
  const installCommand = packageManager === 'yarn' ? 'yarn install' : 'npm install'
  execSync(installCommand, { stdio: 'inherit' })
  console.log('✅ 依赖安装完成')
} catch (error) {
  console.error('❌ 依赖安装失败:', error.message)
  process.exit(1)
}

// 检查环境文件
console.log('\n🔧 检查环境配置...')
const envFile = path.join(__dirname, '../.env')
if (fs.existsSync(envFile)) {
  console.log('✅ 环境配置文件已存在')
} else {
  console.log('⚠️  环境配置文件不存在，请检查 .env 文件')
}

// 运行类型检查
console.log('\n🔍 运行类型检查...')
try {
  execSync('npx vue-tsc --noEmit', { stdio: 'inherit' })
  console.log('✅ 类型检查通过')
} catch (error) {
  console.log('⚠️  类型检查发现问题，但不影响运行')
}

// 完成设置
console.log('\n🎉 项目设置完成！')
console.log('\n📋 可用命令:')
console.log(`   ${packageManager} run dev      - 启动开发服务器`)
console.log(`   ${packageManager} run build    - 构建生产版本`)
console.log(`   ${packageManager} run preview  - 预览生产版本`)
console.log(`   ${packageManager} run lint     - 代码检查`)
console.log(`   ${packageManager} run format   - 代码格式化`)

console.log('\n🚀 现在可以运行以下命令启动开发服务器:')
console.log(`   cd web`)
console.log(`   ${packageManager} run dev`)

console.log('\n📖 更多信息请查看 README.md 文件')
console.log('💡 如有问题，请查看项目文档或联系开发团队\n')
