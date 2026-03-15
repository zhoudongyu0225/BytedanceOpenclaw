<template>
  <div class="gift-simulator" :class="{ 'dark': isDark }">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>礼物模拟工具</span>
          <el-button size="small" @click="refreshGifts">刷新礼物配置</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <!-- 左侧：礼物发送表单 -->
        <el-col :xs="24" :sm="24" :md="10" :lg="8">
          <el-form label-width="80px" label-position="top" class="gift-form">
            <el-form-item label="平台" required>
              <el-select v-model="form.platform" placeholder="请选择平台" style="width: 100%;" @change="onPlatformChange">
                <el-option label="抖音" :value="1" />
                <el-option label="快手" :value="2" />
                <el-option label="TikTok" :value="3" />
              </el-select>
            </el-form-item>
            <el-form-item label="主播" required>
              <el-select v-model="form.anchor_id" placeholder="请选择主播" style="width: 100%;" filterable>
                <el-option v-for="anchor in anchorList" :key="anchor.id" :label="`${anchor.name} (${anchor.room_id})`" :value="anchor.id" />
              </el-select>
            </el-form-item>
            
            <!-- 礼物选择 -->
            <el-form-item label="选择礼物" required>
              <div class="gift-grid">
                <div 
                  v-for="gift in giftList" 
                  :key="gift.id"
                  class="gift-item"
                  :class="{ active: form.gift_id === gift.id }"
                  @click="selectGift(gift)"
                >
                  <div class="gift-icon">{{ gift.icon }}</div>
                  <div class="gift-name">{{ gift.name }}</div>
                  <div class="gift-price">{{ gift.price }}币</div>
                  <div class="gift-effect">+{{ gift.soldier }}兵</div>
                </div>
              </div>
            </el-form-item>
            
            <el-form-item label="礼物数量">
              <el-input-number v-model="form.count" :min="1" :max="100" style="width: 100%;" />
            </el-form-item>
            <el-form-item label="发送用户">
              <el-input v-model="form.user_name" placeholder="请输入发送用户昵称" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="sendGift" style="width: 100%;">发送模拟礼物</el-button>
            </el-form-item>
          </el-form>
        </el-col>
        
        <!-- 右侧：发送日志 -->
        <el-col :xs="24" :sm="24" :md="14" :lg="16">
          <div class="log-box">
            <div class="log-header">
              <span>发送日志</span>
              <div>
                <el-button size="small" @click="clearLog">清空</el-button>
              </div>
            </div>
            <div class="log-content" ref="logContentRef">
              <div v-for="(log, index) in logs" :key="index" class="log-item" :class="log.type">
                <span class="log-time">[{{ log.time }}]</span>
                <span class="log-platform">{{ log.platform }}</span>
                <span class="log-text">{{ log.text }}</span>
              </div>
              <div v-if="logs.length === 0" class="log-empty">暂无发送记录</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const isDark = ref(localStorage.getItem('theme') === 'dark')

const anchorList = ref([])
const logs = ref([])
const logContentRef = ref(null)

// 真实礼物配置 - 抖音/快手礼物
const giftConfig = {
  1: [ // 抖音
    { id: 1, name: '小心心', icon: '❤️', price: 1, soldier: 10, desc: '基础礼物' },
    { id: 2, name: '玫瑰花', icon: '🌹', price: 10, soldier: 100, desc: '浪漫礼物' },
    { id: 3, name: '大啤酒', icon: '🍺', price: 30, soldier: 280, desc: '畅饮礼物' },
    { id: 4, name: '仙女棒', icon: '🪄', price: 52, soldier: 500, desc: '魔法礼物' },
    { id: 5, name: '大喇叭', icon: '📢', price: 100, soldier: 950, desc: '公告礼物' },
    { id: 6, name: '跑车', icon: '🚗', price: 200, soldier: 1900, desc: '豪气礼物' },
    { id: 7, name: '城堡', icon: '🏰', price: 500, soldier: 4800, desc: '梦幻礼物' },
    { id: 8, name: '火箭', icon: '🚀', price: 1000, soldier: 9900, desc: '超级火箭' },
    { id: 9, name: '嘉年华', icon: '🎆', price: 3000, soldier: 29000, desc: '顶级烟花' },
    { id: 10, name: '为你点亮', icon: '💖', price: 5200, soldier: 52000, desc: '真爱特效' },
    { id: 11, name: '求佛', icon: '🙏', price: 10000, soldier: 100000, desc: '许愿特效' },
  ],
  2: [ // 快手
    { id: 1, name: '小心心', icon: '❤️', price: 1, soldier: 10, desc: '基础礼物' },
    { id: 2, name: '玫瑰', icon: '🌹', price: 10, soldier: 100, desc: '鲜花礼物' },
    { id: 3, name: '棒棒糖', icon: '🍭', price: 20, soldier: 200, desc: '甜蜜礼物' },
    { id: 4, name: '大啤酒', icon: '🍺', price: 50, soldier: 500, desc: '畅饮礼物' },
    { id: 5, name: '仙女棒', icon: '🪄', price: 100, soldier: 1000, desc: '魔法礼物' },
    { id: 6, name: '跑车', icon: '🚗', price: 200, soldier: 2000, desc: '豪气礼物' },
    { id: 7, name: '城堡', icon: '🏰', price: 500, soldier: 5000, desc: '梦幻礼物' },
    { id: 8, name: '火箭', icon: '🚀', price: 1000, soldier: 10000, desc: '超级火箭' },
    { id: 9, name: '为你点亮', icon: '💖', price: 3000, soldier: 30000, desc: '真爱特效' },
  ],
  3: [ // TikTok
    { id: 1, name: 'Rose', icon: '🌹', price: 1, soldier: 10, desc: 'Basic Gift' },
    { id: 2, name: 'Coffee', icon: '☕', price: 5, soldier: 50, desc: 'Coffee Gift' },
    { id: 3, name: 'Dragon', icon: '🐉', price: 30, soldier: 300, desc: 'Dragon Gift' },
    { id: 4, name: 'Lambo', icon: '🏎️', price: 100, soldier: 1000, desc: 'Luxury Car' },
    { id: 5, name: 'Castle', icon: '🏰', price: 500, soldier: 5000, desc: 'Dream Castle' },
    { id: 6, name: 'Crown', icon: '👑', price: 1000, soldier: 10000, desc: 'Royal Crown' },
  ]
}

const form = reactive({
  platform: 1,
  anchor_id: '',
  gift_id: 1,
  count: 1,
  user_name: '模拟用户'
})

// 监听主题变化
onMounted(() => {
  window.addEventListener('storage', (e) => {
    if (e.key === 'theme') {
      isDark.value = e.newValue === 'dark'
    }
  })
})

// 获取礼物列表
const getGiftList = () => {
  return giftConfig[form.platform] || giftConfig[1]
}

// 平台变化时重置礼物选择
const onPlatformChange = () => {
  form.gift_id = 1
}

// 选择礼物
const selectGift = (gift) => {
  form.gift_id = gift.id
}

// 获取当前礼物列表
const giftList = computed(() => getGiftList())

import { computed } from 'vue'

// 获取主播列表
const getAnchorList = async () => {
  try {
    const res = await request.get('/api/v1/anchor/list', { params: { page: 1, page_size: 100, status: 1 } })
    if (res.code === 200) {
      anchorList.value = res.data.list || []
    }
  } catch (err) {
    console.error('获取主播列表失败', err)
  }
}

// 刷新礼物配置
const refreshGifts = () => {
  ElMessage.success('礼物配置已刷新')
}

// 发送礼物
const sendGift = async () => {
  if (!form.platform || !form.anchor_id || !form.gift_id) {
    ElMessage.warning('请填写必填项')
    return
  }
  
  const anchor = anchorList.value.find(item => item.id === form.anchor_id)
  const gifts = getGiftList()
  const gift = gifts.find(g => g.id === form.gift_id)
  const platformMap = { 1: '抖音', 2: '快手', 3: 'TikTok' }
  
  const log = {
    time: new Date().toLocaleTimeString(),
    platform: platformMap[form.platform],
    text: `向主播【${anchor?.name || '未知'}】发送 ${gift?.name || ''} x${form.count}，获得${(gift?.soldier || 0) * form.count}士兵`,
    type: 'success'
  }
  logs.value.unshift(log)
  
  // 滚动到顶部
  nextTick(() => {
    if (logContentRef.value) {
      logContentRef.value.scrollTop = 0
    }
  })
  
  ElMessage.success(`模拟礼物发送成功！${gift?.name} x${form.count}`)
}

// 清空日志
const clearLog = () => {
  logs.value = []
}

onMounted(() => {
  getAnchorList()
})
</script>

<style scoped>
.gift-simulator {
  height: 100%;
  padding: 16px;
  box-sizing: border-box;
  background: #f5f7fa;
  transition: background 0.3s;
}

.gift-simulator.dark {
  background: #141414;
}

.box-card {
  background: #fff;
  border-radius: 8px;
  transition: background 0.3s, border-color 0.3s;
}

.gift-simulator.dark .box-card {
  background: #1f1f1f;
  border-color: #333;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 礼物网格 */
.gift-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.gift-item {
  padding: 10px 5px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.dark .gift-item {
  border-color: #333;
}

.gift-item:hover {
  border-color: #409eff;
  background: #ecf5ff;
}

.dark .gift-item:hover {
  background: #1a3a5c;
}

.gift-item.active {
  border-color: #409eff;
  background: #ecf5ff;
  box-shadow: 0 0 8px rgba(64, 158, 255, 0.3);
}

.dark .gift-item.active {
  background: #1a3a5c;
}

.gift-icon {
  font-size: 24px;
  margin-bottom: 4px;
}

.gift-name {
  font-size: 12px;
  color: #303133;
  font-weight: 500;
}

.dark .gift-name {
  color: #e5e5e5;
}

.gift-price {
  font-size: 11px;
  color: #909399;
}

.dark .gift-price {
  color: #666;
}

.gift-effect {
  font-size: 10px;
  color: #67c23a;
  margin-top: 2px;
}

/* 日志框 */
.log-box {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  height: 500px;
  display: flex;
  flex-direction: column;
  transition: border-color 0.3s, background 0.3s;
}

.dark .log-box {
  border-color: #333;
  background: #1a1a1a;
}

.log-header {
  padding: 10px;
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: background 0.3s, border-color 0.3s;
}

.dark .log-header {
  background: #2a2a2a;
  border-bottom-color: #333;
}

.log-content {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
  background: #fff;
  transition: background 0.3s;
}

.dark .log-content {
  background: #1a1a1a;
}

.log-item {
  margin-bottom: 8px;
  font-size: 13px;
  line-height: 1.5;
}

.log-time {
  color: #909399;
  margin-right: 8px;
}

.dark .log-time {
  color: #666;
}

.log-platform {
  background: #409eff;
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  margin-right: 8px;
}

.log-text {
  color: #303133;
}

.dark .log-text {
  color: #ccc;
}

.log-empty {
  text-align: center;
  color: #909399;
  padding: 50px 0;
}

.dark .log-empty {
  color: #666;
}
</style>
