<template>
  <div class="settings-page">
    <el-row :gutter="16">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>游戏参数设置</span>
            </div>
          </template>
          
          <el-form label-width="120px">
            <el-form-item label="游戏名称">
              <el-input v-model="settings.game_name" />
            </el-form-item>
            <el-form-item label="每局时长(秒)">
              <el-input-number v-model="settings.game_duration" :min="30" :max="600" />
            </el-form-item>
            <el-form-item label="匹配等待时间">
              <el-input-number v-model="settings.match_wait_time" :min="5" :max="60" />
            </el-form-item>
            <el-form-item label="最大同时对局">
              <el-input-number v-model="settings.max_battles" :min="1" :max="100" />
            </el-form-item>
            <el-form-item label="弹幕间隔(ms)">
              <el-input-number v-model="settings.danmaku_interval" :min="100" :max="5000" :step="100" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>礼物配置</span>
            </div>
          </template>
          
          <el-table :data="giftList" border size="small">
            <el-table-column prop="name" label="礼物名称" />
            <el-table-column prop="price" label="价格(抖币)" width="100" />
            <el-table-column prop="value" label="游戏价值" width="100" />
            <el-table-column prop="dinosaur" label="恐龙类型" width="100" />
            <el-table-column label="操作" width="80">
              <template #default>
                <el-button size="small" type="primary" link>编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
        
        <el-card style="margin-top: 16px;">
          <template #header>
            <div class="card-header">
              <span>平台配置</span>
            </div>
          </template>
          
          <el-form label-width="100px">
            <el-form-item label="抖音AppID">
              <el-input v-model="settings.douyin_appid" />
            </el-form-item>
            <el-form-item label="快手AppID">
              <el-input v-model="settings.kuaishou_appid" />
            </el-form-item>
            <el-form-item label="TikTok AppID">
              <el-input v-model="settings.tiktok_appid" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="savePlatformSettings">保存平台配置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const settings = reactive({
  game_name: '恐龙攻守',
  game_duration: 180,
  match_wait_time: 30,
  max_battles: 50,
  danmaku_interval: 500,
  douyin_appid: '',
  kuaishou_appid: '',
  tiktok_appid: ''
})

const giftList = ref([
  { name: '小心心', price: 1, value: 1, dinosaur: '迅猛龙' },
  { name: '玫瑰花', price: 10, value: 10, dinosaur: '三角龙' },
  { name: '眼镜', price: 99, value: 99, dinosaur: '剑龙' },
  { name: '火箭', price: 1000, value: 1000, dinosaur: '霸王龙' },
  { name: '嘉年华', price: 30000, value: 30000, dinosaur: '史诗龙' }
])

const saveSettings = () => {
  ElMessage.success('游戏参数保存成功')
}

const savePlatformSettings = () => {
  ElMessage.success('平台配置保存成功')
}

onMounted(() => {
  // 加载配置
})
</script>

<style scoped>
.settings-page {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
