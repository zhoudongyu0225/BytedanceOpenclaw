<template>
  <div class="fans-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>粉丝管理</span>
          <el-button type="primary" size="small" @click="refreshData">刷新数据</el-button>
        </div>
      </template>
      
      <div class="filter-bar">
        <el-input 
          v-model="searchKeyword" 
          placeholder="搜索粉丝昵称/ID" 
          style="width: 240px;" 
          clearable 
          @keyup.enter="handleSearch"
        />
        <el-select v-model="filterPlatform" placeholder="平台筛选" clearable style="width: 120px; margin-left: 12px;">
          <el-option label="抖音" :value="1" />
          <el-option label="快手" :value="2" />
          <el-option label="TikTok" :value="3" />
        </el-select>
        <el-button type="primary" style="margin-left: 12px;" @click="handleSearch">搜索</el-button>
      </div>

      <el-table :data="tableData" border stripe height="500" style="margin-top: 16px;">
        <el-table-column prop="platform" label="平台" width="80">
          <template #default="{ row }">
            <el-tag :type="getPlatformType(row.platform)" size="small">
              {{ getPlatformName(row.platform) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="fans_name" label="粉丝昵称" min-width="120" />
        <el-table-column prop="fans_id" label="粉丝ID" width="150" />
        <el-table-column prop="anchor_name" label="关注主播" min-width="120" />
        <el-table-column prop="gift_total" label="礼物总额" width="100">
          <template #default="{ row }">
            ¥{{ row.gift_total || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="gift_count" label="礼物数" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁言' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link>详情</el-button>
            <el-button size="small" type="warning" link>禁言</el-button>
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
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const tableData = ref([])
const searchKeyword = ref('')
const filterPlatform = ref('')

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const getPlatformName = (platform) => {
  const map = { 1: '抖音', 2: '快手', 3: 'TikTok' }
  return map[platform] || '-'
}

const getPlatformType = (platform) => {
  const map = { 1: 'danger', 2: 'warning', 3: 'info' }
  return map[platform] || ''
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const refreshData = () => {
  ElMessage.success('刷新成功')
  loadData()
}

const loadData = () => {
  // 模拟数据
  tableData.value = [
    {
      platform: 1,
      fans_name: '粉丝用户A',
      fans_id: '2001',
      anchor_name: '测试主播',
      gift_total: 1000,
      gift_count: 50,
      status: 1,
      created_at: '2025-01-15 10:30:00'
    }
  ]
  pagination.total = 1
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.fans-page {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-bar {
  display: flex;
  align-items: center;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
