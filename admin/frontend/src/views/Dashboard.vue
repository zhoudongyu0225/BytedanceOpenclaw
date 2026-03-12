<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>总主播数</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">{{ anchorCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>今日对战数</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">{{ todayBattleCount }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>今日礼物收入</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">¥{{ todayGiftIncome }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>在线人数</span>
            </div>
          </template>
          <div class="card-content">
            <span class="number">{{ onlineCount }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>功能进度</span>
            </div>
          </template>
          <el-steps :active="2" finish-status="success">
            <el-step title="基础框架" description="已完成" />
            <el-step title="登录/权限" description="已完成" />
            <el-step title="主播管理" description="开发中" />
            <el-step title="礼物模拟工具" description="待开发" />
            <el-step title="GM功能" description="待开发" />
          </el-steps>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const anchorCount = ref(0)
const todayBattleCount = ref(0)
const todayGiftIncome = ref(0)
const onlineCount = ref(0)
</script>

<style scoped>
.dashboard {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 500;
}
.card-content {
  text-align: center;
  padding: 20px 0;
}
.number {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
}
</style>
