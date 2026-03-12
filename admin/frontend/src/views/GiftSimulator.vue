<template>
  <div class="gift-simulator">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>礼物模拟工具</span>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form label-width="100px">
            <el-form-item label="平台" required>
              <el-select v-model="form.platform" placeholder="请选择平台" style="width: 100%;">
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
            <el-form-item label="礼物类型" required>
              <el-select v-model="form.gift_id" placeholder="请选择礼物" style="width: 100%;">
                <el-option label="小心心(1抖币)" :value="1" />
                <el-option label="玫瑰花(10抖币)" :value="2" />
                <el-option label="眼镜(99抖币)" :value="3" />
                <el-option label="火箭(1000抖币)" :value="4" />
                <el-option label="嘉年华(30000抖币)" :value="5" />
              </el-select>
            </el-form-item>
            <el-form-item label="礼物数量" required>
              <el-input-number v-model="form.count" :min="1" :max="10000" style="width: 100%;" />
            </el-form-item>
            <el-form-item label="发送用户" required>
              <el-input v-model="form.user_name" placeholder="请输入发送用户昵称" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="sendGift" style="width: 100%;">发送模拟礼物</el-button>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="16">
          <div class="log-box">
            <div class="log-header">
              <span>发送日志</span>
              <el-button size="small" @click="clearLog">清空日志</el-button>
            </div>
            <div class="log-content">
              <div v-for="(log, index) in logs" :key="index" class="log-item" :class="log.type">
                <span class="log-time">[{{ log.time }}]</span>
                <span class="log-text">{{ log.text }}</span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const anchorList = ref([])
const logs = ref([])

const form = reactive({
  platform: '',
  anchor_id: '',
  gift_id: '',
  count: 1,
  user_name: '模拟用户'
})

// 获取主播列表
const getAnchorList = async () => {
  try {
    const res = await request.get('/api/v1/anchor/list', { params: { page: 1, page_size: 100 } })
    if (res.code === 200) {
      anchorList.value = res.data.list
    }
  } catch (err) {
    ElMessage.error('获取主播列表失败')
  }
}

// 发送礼物
const sendGift = async () => {
  if (!form.platform || !form.anchor_id || !form.gift_id) {
    ElMessage.warning('请填写必填项')
    return
  }
  // 先打印日志，后续对接真实发送逻辑
  const anchor = anchorList.value.find(item => item.id === form.anchor_id)
  const platformMap = { 1: '抖音', 2: '快手', 3: 'TikTok' }
  const giftMap = { 1: '小心心', 2: '玫瑰花', 3: '眼镜', 4: '火箭', 5: '嘉年华' }
  const log = {
    time: new Date().toLocaleString(),
    text: `向【${platformMap[form.platform]}】主播【${anchor?.name || '未知'}】(${anchor?.room_id || ''}) 发送礼物【${giftMap[form.gift_id]}】x ${form.count}，发送用户：${form.user_name}`,
    type: 'success'
  }
  logs.value.unshift(log)
  ElMessage.success('模拟礼物发送成功（功能开发中，暂未对接真实回调）')
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
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.log-box {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  height: 500px;
  display: flex;
  flex-direction: column;
}
.log-header {
  padding: 10px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.log-content {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
}
.log-item {
  margin-bottom: 8px;
  font-size: 14px;
  line-height: 1.5;
}
.log-item.success {
  color: #67c23a;
}
.log-time {
  color: #909399;
  margin-right: 10px;
}
</style>
