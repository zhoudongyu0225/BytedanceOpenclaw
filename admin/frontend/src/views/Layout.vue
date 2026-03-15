<template>
  <div class="layout-container">
    <el-container style="height: 100vh; width: 100vw;">
      <!-- 左侧导航 - 游戏科技风 -->
      <el-aside width="240px" style="background: linear-gradient(180deg, #0f0c29 0%, #302b63 50%, #24243e 100%); overflow: hidden;">
        <div class="logo-area">
          <div class="logo-icon">🦖</div>
          <span class="logo-text">恐龙后台</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          background-color="transparent"
          text-color="rgba(255,255,255,0.7)"
          active-text-color="#00d4ff"
          :router="true"
        >
          <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
            <el-icon><component :is="item.icon" /></el-icon>
            <span>{{ item.title }}</span>
          </el-menu-item>
        </el-menu>
        <div class="sidebar-footer">
          <div class="version">v1.0.0</div>
        </div>
      </el-aside>
      
      <el-container style="flex: 1; overflow: hidden; background: #0a0a0f;">
        <!-- 顶部栏 -->
        <el-header class="top-header">
          <div class="header-left">
            <span class="page-title">{{ pageTitle }}</span>
          </div>
          <div class="header-right">
            <el-select v-model="currentEnv" placeholder="选择环境" size="small" class="env-select">
              <el-option label="🧪 测试环境" value="test" />
              <el-option label="🚀 正式环境" value="prod" />
            </el-select>
            <el-select v-model="currentPlatform" placeholder="选择平台" size="small" class="platform-select">
              <el-option label="🎵 抖音" :value="1" />
              <el-option label="📱 快手" :value="2" />
              <el-option label="🌍 TikTok" :value="3" />
            </el-select>
            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <div class="avatar">管</div>
                <span>{{ userInfo?.nickname || userInfo?.username || '管理员' }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">🚪 退出登录</el-dropdown-item>
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
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  DataAnalysis, User, Present, Promotion, 
  UserFilled, Bell, Setting, ArrowDown,
  House, ChatLineRound, Clock, Document
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

const userInfo = ref({})
const currentEnv = ref('test')
const currentPlatform = ref(1)

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
* {
  box-sizing: border-box;
}

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
.el-aside {
  position: relative;
  display: flex;
  flex-direction: column;
}

.logo-area {
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid rgba(255,255,255,0.1);
  background: rgba(0,0,0,0.2);
}

.logo-icon {
  font-size: 28px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(90deg, #00d4ff, #00ff88);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.el-menu-vertical {
  flex: 1;
  border-right: none;
  padding: 12px 0;
}

.el-menu-item {
  height: 48px;
  line-height: 48px;
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.el-menu-item:hover {
  background: rgba(0,212,255,0.1) !important;
}

.el-menu-item.is-active {
  background: linear-gradient(90deg, rgba(0,212,255,0.2) 0%, rgba(0,212,255,0.05) 100%) !important;
  border-left: 3px solid #00d4ff;
}

.el-menu-item .el-icon {
  margin-right: 10px;
  font-size: 18px;
}

.sidebar-footer {
  padding: 16px;
  text-align: center;
  border-top: 1px solid rgba(255,255,255,0.1);
}

.version {
  color: rgba(255,255,255,0.3);
  font-size: 12px;
}

/* 顶部栏 */
.top-header {
  background: rgba(15,15,25,0.95);
  border-bottom: 1px solid rgba(0,212,255,0.2);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
  backdrop-filter: blur(10px);
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #fff;
}

.page-title::before {
  content: '▶';
  margin-right: 10px;
  color: #00d4ff;
  font-size: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.env-select, .platform-select {
  width: 140px;
}

.env-select :deep(.el-input__wrapper),
.platform-select :deep(.el-input__wrapper) {
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  box-shadow: none;
}

.env-select :deep(.el-input__wrapper:hover),
.platform-select :deep(.el-input__wrapper:hover) {
  border-color: #00d4ff;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 20px;
  background: rgba(255,255,255,0.05);
  transition: all 0.3s;
  color: #fff;
}

.user-info:hover {
  background: rgba(0,212,255,0.2);
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #00d4ff, #00ff88);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}

/* 主内容区 */
.main-content {
  background: #0a0a0f;
  padding: 0;
  overflow: hidden;
  height: calc(100vh - 60px);
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: rgba(255,255,255,0.05);
}

::-webkit-scrollbar-thumb {
  background: rgba(0,212,255,0.3);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0,212,255,0.5);
}
</style>
