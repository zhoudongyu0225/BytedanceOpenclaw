<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #667eea, #764ba2);">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总主播数</div>
            <div class="stat-value">{{ anchorCount }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f093fb, #f5576c);">
            <el-icon><Coin /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">今日对战数</div>
            <div class="stat-value">{{ todayBattleCount }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #4facfe, #00f2fe);">
            <el-icon><Money /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">今日礼物收入</div>
            <div class="stat-value">¥{{ todayGiftIncome.toLocaleString() }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #43e97b, #38f9d7);">
            <el-icon><UserFilled /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">在线人数</div>
            <div class="stat-value">{{ onlineCount }}</div>
          </div>
        </div>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <div class="quick-card">
          <div class="card-title">快捷入口</div>
          <div class="quick-entry">
            <router-link to="/anchor" class="quick-btn">
              <el-icon><User /></el-icon>
              <span>主播管理</span>
            </router-link>
            <router-link to="/gift-simulator" class="quick-btn">
              <el-icon><Present /></el-icon>
              <span>礼物模拟</span>
            </router-link>
            <router-link to="/gm-command" class="quick-btn">
              <el-icon><Promotion /></el-icon>
              <span>GM命令</span>
            </router-link>
            <router-link to="/fans" class="quick-btn">
              <el-icon><UserFilled /></el-icon>
              <span>粉丝管理</span>
            </router-link>
            <router-link to="/announcement" class="quick-btn">
              <el-icon><Bell /></el-icon>
              <span>公告管理</span>
            </router-link>
            <router-link to="/settings" class="quick-btn">
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </router-link>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { User, Coin, Money, UserFilled, Present, Promotion, Bell, Setting } from '@element-plus/icons-vue'

const anchorCount = ref(128)
const todayBattleCount = ref(856)
const todayGiftIncome = ref(125680)
const onlineCount = ref(2341)
</script>

<style scoped>
.dashboard {
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  background: #f5f7fa;
  transition: background 0.3s;
}

/* 统计卡片 */
.stat-card {
  height: 100px;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
  transition: all 0.2s, background 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stat-icon .el-icon {
  font-size: 24px;
  color: #fff;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
  transition: color 0.3s;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  transition: color 0.3s;
}

/* 快捷入口 */
.quick-card {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
  transition: background 0.3s, box-shadow 0.3s;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
  transition: color 0.3s;
}

.quick-entry {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.quick-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 90px;
  background: #f5f7fa;
  border-radius: 8px;
  text-decoration: none;
  color: #606266;
  transition: all 0.2s;
}

.quick-btn:hover {
  background: #409eff;
  color: #fff;
}

.quick-btn .el-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.quick-btn span {
  font-size: 14px;
}

/* ===== 暗色主题 ===== */
:global(html.dark) .dashboard {
  background: #141414;
}

:global(html.dark) .stat-card {
  background: #1f1f1f;
  box-shadow: 0 2px 12px rgba(0,0,0,0.2);
}

:global(html.dark) .stat-label {
  color: #666;
}

:global(html.dark) .stat-value {
  color: #e5e5e5;
}

:global(html.dark) .quick-card {
  background: #1f1f1f;
  box-shadow: 0 2px 12px rgba(0,0,0,0.2);
}

:global(html.dark) .card-title {
  color: #e5e5e5;
}

:global(html.dark) .quick-btn {
  background: #2a2a2a;
  color: #999;
}

:global(html.dark) .quick-btn:hover {
  background: #409eff;
}
</style>
