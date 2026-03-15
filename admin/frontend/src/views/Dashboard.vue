<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card card-1">
          <div class="stat-icon">👑</div>
          <div class="stat-info">
            <div class="stat-label">总主播数</div>
            <div class="stat-value">{{ anchorCount }}</div>
          </div>
          <div class="stat-trend up">+12%</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-2">
          <div class="stat-icon">⚔️</div>
          <div class="stat-info">
            <div class="stat-label">今日对战数</div>
            <div class="stat-value">{{ todayBattleCount }}</div>
          </div>
          <div class="stat-trend up">+8%</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-3">
          <div class="stat-icon">💎</div>
          <div class="stat-info">
            <div class="stat-label">今日礼物收入</div>
            <div class="stat-value">¥{{ todayGiftIncome }}</div>
          </div>
          <div class="stat-trend up">+25%</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-4">
          <div class="stat-icon">👥</div>
          <div class="stat-info">
            <div class="stat-label">在线人数</div>
            <div class="stat-value">{{ onlineCount }}</div>
          </div>
          <div class="stat-trend">在线</div>
        </div>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <div class="quick-entry-card">
          <div class="card-title">
            <span>🚀</span> 快捷入口
          </div>
          <div class="quick-entry">
            <router-link to="/anchor" class="quick-btn">
              <div class="quick-icon">👤</div>
              <span>主播管理</span>
            </router-link>
            <router-link to="/gift-simulator" class="quick-btn">
              <div class="quick-icon">🎁</div>
              <span>礼物模拟</span>
            </router-link>
            <router-link to="/gm-command" class="quick-btn">
              <div class="quick-icon">⚡</div>
              <span>GM命令</span>
            </router-link>
            <router-link to="/fans" class="quick-btn">
              <div class="quick-icon">❤️</div>
              <span>粉丝管理</span>
            </router-link>
            <router-link to="/announcement" class="quick-btn">
              <div class="quick-icon">📢</div>
              <span>公告管理</span>
            </router-link>
            <router-link to="/settings" class="quick-btn">
              <div class="quick-icon">⚙️</div>
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
}

/* 统计卡片 */
.stat-card {
  height: 120px;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  position: relative;
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0,0,0,0.3);
}

.card-1 {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.card-2 {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.card-3 {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.card-4 {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-icon {
  font-size: 40px;
  margin-right: 16px;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: rgba(255,255,255,0.8);
  margin-bottom: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
}

.stat-trend {
  position: absolute;
  top: 16px;
  right: 16px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  background: rgba(255,255,255,0.2);
  color: #fff;
}

.stat-trend.up {
  background: rgba(0,255,128,0.3);
}

/* 快捷入口 */
.quick-entry-card {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 16px;
  padding: 24px;
  backdrop-filter: blur(10px);
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 20px;
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
  width: 140px;
  height: 100px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 12px;
  text-decoration: none;
  color: rgba(255,255,255,0.8);
  transition: all 0.3s;
}

.quick-btn:hover {
  background: linear-gradient(135deg, rgba(0,212,255,0.2) 0%, rgba(0,255,136,0.2) 100%);
  border-color: #00d4ff;
  transform: translateY(-3px);
  box-shadow: 0 10px 30px rgba(0,212,255,0.2);
}

.quick-icon {
  font-size: 32px;
  margin-bottom: 8px;
}

.quick-btn span {
  font-size: 14px;
  font-weight: 500;
}
</style>
