<template>
  <div class="anchor-list">
    <el-card class="box-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>主播管理</span>
          <el-button type="primary" size="small" class="cyber-btn" @click="openAddDialog">
            <el-icon><Plus /></el-icon>
            新增主播
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="filter-bar">
        <el-input 
          v-model="filter.keyword" 
          placeholder="搜索主播姓名/昵称/直播间ID" 
          style="width: 280px;" 
          clearable 
          class="cyber-input"
          @keyup.enter="getList" 
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="filter.platform" placeholder="平台" style="width: 120px; margin-left: 12px;" clearable @change="getList" class="cyber-select">
          <el-option label="抖音" :value="1" />
          <el-option label="快手" :value="2" />
          <el-option label="TikTok" :value="3" />
        </el-select>
        <el-select v-model="filter.status" placeholder="状态" style="width: 100px; margin-left: 12px;" clearable @change="getList" class="cyber-select">
          <el-option label="正常" :value="1" />
          <el-option label="封禁" :value="2" />
        </el-select>
        <el-button type="primary" size="small" class="cyber-btn" style="margin-left: 12px;" @click="getList">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
      </div>

      <!-- 操作栏 -->
      <div class="batch-bar">
        <el-button 
          v-if="selectedIds.length > 0" 
          type="danger" 
          size="small" 
          class="danger-btn"
          @click="batchUpdateStatus(2)"
        >批量封禁</el-button>
        <el-button 
          v-if="selectedIds.length > 0" 
          type="success" 
          size="small" 
          class="success-btn"
          style="margin-left: 8px;"
          @click="batchUpdateStatus(1)"
        >批量解封</el-button>
        <span v-if="selectedIds.length === 0" style="color: rgba(255,255,255,0.4); font-size: 14px;">
          已选择 0 条记录
        </span>
        <span v-else style="color: #00d4ff; font-size: 14px; margin-left: 12px;">
          已选择 {{ selectedIds.length }} 条记录
        </span>
      </div>

      <!-- 表格 -->
      <el-table 
        :data="tableData" 
        border 
        style="width: 100%;" 
        height="420"
        class="cyber-table"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="主播姓名" width="120" />
        <el-table-column prop="nickname" label="主播昵称" width="120" />
        <el-table-column prop="platform" label="平台" width="80">
          <template #default="{ row }">
            <el-tag :type="getPlatformType(row.platform)" effect="dark" size="small">
              {{ getPlatformName(row.platform) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="room_id" label="直播间ID" width="150" />
        <el-table-column prop="platform_uid" label="平台UID" width="150" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" effect="dark" size="small">
              {{ row.status === 1 ? '正常' : '封禁' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="filter.page"
          v-model:page-size="filter.page_size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          background
          @Size-change="getList"
          @Current-change="getList"
        />
      </div>
    </el-card>

    <!-- 弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增主播' : '编辑主播'" width="600px" class="cyber-dialog">
      <el-form :model="form" label-width="100px">
        <el-form-item label="主播姓名" prop="name" required>
          <el-input v-model="form.name" placeholder="请输入主播姓名" />
        </el-form-item>
        <el-form-item label="主播昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入主播昵称" />
        </el-form-item>
        <el-form-item label="平台" prop="platform" required>
          <el-select v-model="form.platform" placeholder="请选择平台" style="width: 100%;">
            <el-option label="抖音" :value="1" />
            <el-option label="快手" :value="2" />
            <el-option label="TikTok" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="直播间ID" prop="room_id" required>
          <el-input v-model="form.room_id" placeholder="请输入直播间ID" />
        </el-form-item>
        <el-form-item label="平台UID" prop="platform_uid">
          <el-input v-model="form.platform_uid" placeholder="请输入平台用户ID" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">正常</el-radio>
            <el-radio :label="2">封禁</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import request from '../utils/request'

const tableData = ref([])
const total = ref(0)
const selectedIds = ref([])
const dialogVisible = ref(false)
const dialogType = ref('add')

const filter = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  platform: '',
  status: ''
})

const form = reactive({
  id: '',
  name: '',
  nickname: '',
  platform: '',
  room_id: '',
  platform_uid: '',
  status: 1,
  remark: ''
})

const getList = async () => {
  try {
    const res = await request.get('/api/v1/anchor/list', { params: filter })
    if (res.code === 200) {
      tableData.value = res.data.list
      total.value = res.data.total
    }
  } catch (err) {
    ElMessage.error('获取列表失败')
  }
}

const getPlatformName = (platform) => {
  const map = { 1: '抖音', 2: '快手', 3: 'TikTok' }
  return map[platform] || '-'
}

const getPlatformType = (platform) => {
  const map = { 1: 'danger', 2: 'warning', 3: 'info' }
  return map[platform] || ''
}

const handleSelectionChange = (val) => {
  selectedIds.value = val.map(item => item.id)
}

const openAddDialog = () => {
  dialogType.value = 'add'
  Object.keys(form).forEach(key => {
    form[key] = key === 'status' ? 1 : ''
  })
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!form.name || !form.platform || !form.room_id) {
    ElMessage.warning('请填写必填项')
    return
  }
  try {
    if (dialogType.value === 'add') {
      await request.post('/api/v1/anchor/create', form)
      ElMessage.success('创建成功')
    } else {
      await request.post('/api/v1/anchor/update', form)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    getList()
  } catch (err) {
    ElMessage.error(err.response?.data?.msg || '操作失败')
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该主播吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.delete(`/api/v1/anchor/delete/${row.id}`)
      ElMessage.success('删除成功')
      getList()
    } catch (err) {
      ElMessage.error('删除失败')
    }
  })
}

const batchUpdateStatus = (status) => {
  const text = status === 2 ? '封禁' : '解封'
  ElMessageBox.confirm(`确定要批量${text}选中的主播吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await request.post('/api/v1/anchor/batch-update-status', {
        ids: selectedIds.value,
        status: status
      })
      ElMessage.success(`批量${text}成功`)
      getList()
    } catch (err) {
      ElMessage.error(`批量${text}失败`)
    }
  })
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.anchor-list {
  height: 100%;
  padding: 16px;
  box-sizing: border-box;
}

.box-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(255,255,255,0.03) !important;
  border: 1px solid rgba(255,255,255,0.1) !important;
  backdrop-filter: blur(10px);
}

.box-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.box-card :deep(.el-card__header) {
  background: rgba(0,0,0,0.2);
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 40px;
  color: #fff;
}

.card-header span {
  font-size: 16px;
  font-weight: 600;
}

/* 赛博按钮 */
.cyber-btn {
  background: linear-gradient(90deg, #00d4ff, #00ff88) !important;
  border: none !important;
  color: #0a0a0f !important;
  font-weight: 600;
}

.cyber-btn:hover {
  box-shadow: 0 0 20px rgba(0,212,255,0.5);
}

.danger-btn {
  background: linear-gradient(90deg, #ff4757, #ff6b81) !important;
  border: none !important;
  color: #fff !important;
}

.success-btn {
  background: linear-gradient(90deg, #00d4ff, #00ff88) !important;
  border: none !important;
  color: #0a0a0f !important;
}

/* 搜索栏 */
.filter-bar {
  display: flex;
  align-items: center;
  height: 50px;
  flex-shrink: 0;
}

.cyber-input :deep(.el-input__wrapper) {
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  box-shadow: none;
}

.cyber-input :deep(.el-input__wrapper:hover),
.cyber-input :deep(.el-input__wrapper.is-focus) {
  border-color: #00d4ff;
}

.cyber-input :deep(.el-input__inner) {
  color: #fff;
}

.cyber-input :deep(.el-input__inner::placeholder) {
  color: rgba(255,255,255,0.4);
}

.cyber-select :deep(.el-input__wrapper) {
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  box-shadow: none;
}

/* 操作栏 */
.batch-bar {
  display: flex;
  align-items: center;
  height: 40px;
  flex-shrink: 0;
  border-top: 1px solid rgba(255,255,255,0.1);
  border-bottom: 1px solid rgba(255,255,255,0.1);
  padding: 0 12px;
}

/* 表格 */
.cyber-table {
  background: transparent !important;
}

.cyber-table :deep(.el-table__header-wrapper th) {
  background: rgba(0,0,0,0.3) !important;
  color: #fff !important;
  border-bottom: 1px solid rgba(255,255,255,0.1) !important;
}

.cyber-table :deep(.el-table__body-wrapper) {
  background: transparent !important;
}

.cyber-table :deep(.el-table__row) {
  background: transparent !important;
  color: rgba(255,255,255,0.8) !important;
}

.cyber-table :deep(.el-table__row--hover) {
  background: rgba(0,212,255,0.1) !important;
}

.cyber-table :deep(.el-table__row td) {
  border-bottom: 1px solid rgba(255,255,255,0.05) !important;
}

/* 分页 */
.pagination {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-shrink: 0;
  border-top: 1px solid rgba(255,255,255,0.1);
  padding-top: 8px;
}

.pagination :deep(.el-pagination) {
  --el-pagination-bg-color: rgba(255,255,255,0.05);
  --el-pagination-text-color: rgba(255,255,255,0.8);
  --el-pagination-button-color: rgba(255,255,255,0.8);
  --el-pagination-hover-color: #00d4ff;
}
</style>
