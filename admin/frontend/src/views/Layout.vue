<template>
  <div class="layout-container" :class="{ 'dark': isDark }">
    <el-container style="height: 100vh; width: 100vw;">
      <!-- 左侧导航 -->
      <el-aside width="220px" class="sidebar">
        <div class="logo-area">
          <span class="logo-text">后台管理</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          :router="true"
        >
          <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.title }}</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-container style="flex: 1; overflow: hidden;">
        <!-- 顶部栏 -->
        <el-header class="top-header">
          <div class="header-left">
            <span class="page-title">{{ pageTitle }}</span>
          </div>
          <div class="header-right">
            <el-select v-model="currentEnv" placeholder="选择环境" size="small" style="width: 120px; margin-right: 12px;">
              <el-option label="测试环境" value="test" />
              <el-option label="正式环境" value="prod" />
            </el-select>
            <el-select v-model="currentPlatform" placeholder="选择平台" size="small" style="width: 120px; margin-right: 12px;">
              <el-option label="抖音" :value="1" />
              <el-option label="快手" :value="2" />
              <el-option label="TikTok" :value="3" />
            </el-select>
            <!-- 主题切换 -->
            <el-switch
              v-model="isDark"
              inline-prompt
              active-text="暗"
              inactive-text="亮"
              style="margin-right: 16px;"
            />
            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <el-icon><UserFilled /></el-icon>
                <span>{{ userInfo?.nickname || userInfo?.username || '管理员' }}</span>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        
        <!-- 主内容区 -->
        <el-main class="main-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  DataAnalysis, User, Present, Promotion, 
  UserFilled, Bell, Setting
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const userInfo = ref({})
const currentEnv = ref('test')
const currentPlatform = ref(1)
const isDark = ref(false)

// 监听主题变化，广播给所有组件
watch(isDark, (val) => {
  const theme = val ? 'dark' : 'light'
  localStorage.setItem('theme', theme)
  // 触发自定义事件
  window.dispatchEvent(new CustomEvent('theme-change', { detail: val }))
  // 同时修改 html class
  if (val) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
})

const menuItems = [
  { path: '/dashboard', title: '数据概览', icon: DataAnalysis },
  { path: '/anchor', title: '主播管理', icon: User },
  { path: '/gift-simulator', title: '礼物模拟', icon: Present },
  { path: '/gm-command', title: 'GM命令', icon: Promotion },
  { path: '/fans', title: '粉丝管理', icon: UserFilled },
  { path: '/announcement', title: '公告管理', icon: Bell },
  { path: '/settings', title: '系统设置', icon: Setting },
]

const activeMenu = computed(() => route.path)

const pageTitle = computed(() => {
  const titleMap = {
    '/dashboard': '数据概览',
    '/anchor': '主播管理',
    '/gift-simulator': '礼物模拟工具',
    '/gm-command': 'GM命令',
    '/fans': '粉丝管理',
    '/announcement': '公告管理',
    '/settings': '系统设置'
  }
  return titleMap[route.path] || '游戏后台管理'
})

onMounted(() => {
  // 读取主题设置
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
  
  const info = localStorage.getItem('userInfo')
  if (info) {
    userInfo.value = JSON.parse(info)
  }
})

const handleCommand = (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      ElMessage.success('退出成功')
      router.push('/login')
    })
  }
}
</script>

<style scoped>
.layout-container {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}

.el-container {
  width: 100%;
  height: 100%;
}

/* 侧边栏 */
.sidebar {
  background: #fff;
  border-right: 1px solid #e4e7ed;
  transition: background 0.3s, border-color 0.3s;
}

.dark .sidebar {
  background: #1f1f1f;
  border-color: #333;
}

.logo-area {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #409eff, #3375d9);
}

.logo-text {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
}

.el-menu-vertical {
  border-right: none;
  background: #fff;
  transition: background 0.3s;
}

.dark .el-menu-vertical {
  background: #1f1f1f;
}

.el-menu-item {
  height: 48px;
  line-height: 48px;
  margin: 4px 8px;
  border-radius: 8px;
  color: #606266;
  transition: all 0.3s;
}

.dark .el-menu-item {
  color: #999;
}

.el-menu-item:hover {
  background: #f5f7fa;
}

.dark .el-menu-item:hover {
  background: #2a2a2a;
}

.el-menu-item.is-active {
  background: #ecf5ff;
  color: #409eff;
}

.dark .el-menu-item.is-active {
  background: #1a3a5c;
}

.el-menu-item .el-icon {
  margin-right: 10px;
}

/* 顶部栏 */
.top-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
  transition: background 0.3s, border-color 0.3s;
}

.dark .top-header {
  background: #1f1f1f;
  border-color: #333;
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  transition: color 0.3s;
}

.dark .page-title {
  color: #e5e5e5;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  color: #303133;
  transition: all 0.3s;
}

.dark .user-info {
  color: #e5e5e5;
}

.user-info:hover {
  background: #f5f7fa;
}

.dark .user-info:hover {
  background: #2a2a2a;
}

/* 主内容区 */
.main-content {
  background: #f5f7fa;
  padding: 0;
  overflow: hidden;
  height: calc(100vh - 60px);
  transition: background 0.3s;
}

.dark .main-content {
  background: #141414;
}
</style>
