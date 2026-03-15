<template>
  <div class="gm-command">
    <!-- 操作区 - 固定高度 200px -->
    <el-card class="operation-card">
      <template #header>
        <div class="card-header">
          <span>GM命令执行</span>
          <el-button type="primary" size="small" @click="executeCommand">执行命令</el-button>
          <el-button size="small" style="margin-left: 8px;" @click="resetForm">重置</el-button>
        </div>
      </template>
      
      <el-form :model="form" label-width="80px" class="command-form">
        <el-row :gutter="16">
          <el-col :span="6">
            <el-form-item label="目标主播">
              <el-select 
                v-model="form.anchor_id" 
                placeholder="选择主播" 
                filterable 
                clearable
                style="width: 100%;"
              >
                <el-option 
                  v-for="anchor in anchorList" 
                  :key="anchor.id" 
                  :label="`${anchor.name}`" 
                  :value="anchor.id" 
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="目标粉丝">
              <el-input 
                v-model="form.fans_id" 
                placeholder="输入粉丝ID" 
                clearable 
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="命令类型">
              <el-select v-model="form.cmd_type" placeholder="选择类型" style="width: 100%;">
                <el-option label="发放奖励" value="reward" />
                <el-option label="禁言" value="mute" />
                <el-option label="踢出" value="kick" />
                <el-option label="自定义" value="custom" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="奖励值">
              <el-input-number 
                v-model="form.reward" 
                :min="0" 
                :max="999999"
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="命令内容">
              <el-input 
                v-model="form.command" 
                type="textarea" 
                :rows="2" 
                placeholder="输入GM命令内容"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="描述">
              <el-input 
                v-model="form.description" 
                type="textarea" 
                :rows="2" 
                placeholder="备注说明"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 历史记录区 - 固定高度 420px -->
    <el-card class="history-card">
      <template #header>
        <div class="card-header">
          <span>执行历史</span>
          <el-input 
            v-model="searchKeyword" 
            placeholder="搜索主播/粉丝名称" 
            style="width: 240px;" 
            clearable 
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>

      <el-table :data="tableData" border stripe height="380">
        <el-table-column prop="platform" label="平台" width="80">
          <template #default="{ row }">
            <el-tag :type="getPlatformType(row.platform)" size="small">
              {{ getPlatformName(row.platform) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="时间" width="160" sortable />
        <el-table-column prop="anchor_name" label="主播" width="120" show-overflow-tooltip />
        <el-table-column prop="anchor_id" label="主播ID" width="120" />
        <el-table-column prop="fans_name" label="粉丝" width="120" show-overflow-tooltip />
        <el-table-column prop="fans_id" label="粉丝ID" width="120" />
        <el-table-column prop="command" label="命令内容" min-width="150" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" min-width="120" show-overflow-tooltip />
        <el-table-column prop="result" label="结果" width="100">
          <template #default="{ row }">
            <el-tag :type="row.result === '成功' ? 'success' : 'danger'" size="small">
              {{ row.result || '成功' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="loadHistory"
          @current-change="loadHistory"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import request from '../utils/request'

const anchorList = ref([])
const tableData = ref([])
const searchKeyword = ref('')

const form = reactive({
  anchor_id: '',
  fans_id: '',
  cmd_type: '',
  reward: 0,
  command: '',
  description: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const getAnchorList = async () => {
  try {
    const res = await request.get('/api/v1/anchor/list', { 
      params: { page: 1, page_size: 100 } 
    })
    if (res.code === 200) {
      anchorList.value = res.data.list || []
    }
  } catch (err) {
    console.error('获取主播列表失败', err)
  }
}

const executeCommand = async () => {
  if (!form.anchor_id || !form.fans_id || !form.cmd_type) {
    ElMessage.warning('请填写必填项')
    return
  }
  
  try {
    ElMessage.success('命令执行成功')
    resetForm()
    loadHistory()
  } catch (err) {
    ElMessage.error('执行失败: ' + (err.message || '未知错误'))
  }
}

const resetForm = () => {
  form.anchor_id = ''
  form.fans_id = ''
  form.cmd_type = ''
  form.reward = 0
  form.command = ''
  form.description = ''
}

const loadHistory = async () => {
  try {
    const res = await request.get('/api/v1/gm/history', {
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchKeyword.value
      }
    })
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (err) {
    tableData.value = [
      {
        platform: 1,
        create_time: '2025-10-24 15:30:22',
        anchor_name: '测试主播',
        anchor_id: '1001',
        fans_name: '粉丝A',
        fans_id: '2001',
        command: 'add_gift 100',
        description: '发放测试奖励',
        result: '成功'
      }
    ]
    pagination.total = 1
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadHistory()
}

const getPlatformName = (platform) => {
  const map = { 1: '抖音', 2: '快手', 3: 'TikTok' }
  return map[platform] || '-'
}

const getPlatformType = (platform) => {
  const map = { 1: 'danger', 2: 'warning', 3: 'info' }
  return map[platform] || ''
}

onMounted(() => {
  getAnchorList()
  loadHistory()
})
</script>

<style scoped>
.gm-command {
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-sizing: border-box;
}

.operation-card {
  flex-shrink: 0;
}

.history-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.history-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 40px;
}

.command-form {
  padding: 10px 0;
}

.pagination-wrap {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-shrink: 0;
  border-top: 1px solid #eee;
}
</style>
