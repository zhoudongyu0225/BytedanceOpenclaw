<template>
  <div class="announcement-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>公告管理</span>
          <el-button type="primary" size="small" @click="openAddDialog">发布新公告</el-button>
        </div>
      </template>
      
      <el-table :data="tableData" border stripe height="500">
        <el-table-column prop="title" label="公告标题" min-width="180" />
        <el-table-column prop="content" label="公告内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="platform" label="平台" width="100">
          <template #default="{ row }">
            <el-tag :type="getPlatformType(row.platform)" size="small">
              {{ getPlatformName(row.platform) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            {{ row.type === 1 ? '系统' : '活动' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="start_time" label="开始时间" width="160" />
        <el-table-column prop="end_time" label="结束时间" width="160" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button size="small" type="warning" link @click="toggleStatus(scope.row)">
              {{ scope.row.status === 1 ? '隐藏' : '显示' }}
            </el-button>
            <el-button size="small" type="danger" link @click="handleDelete(scope.row)">删除</el-button>
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

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '发布公告' : '编辑公告'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="公告标题" required>
          <el-input v-model="form.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-form-item label="公告内容" required>
          <el-input v-model="form.content" type="textarea" :rows="4" placeholder="请输入公告内容" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="平台">
              <el-select v-model="form.platform" placeholder="选择平台" style="width: 100%;">
                <el-option label="全部" :value="0" />
                <el-option label="抖音" :value="1" />
                <el-option label="快手" :value="2" />
                <el-option label="TikTok" :value="3" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型">
              <el-select v-model="form.type" placeholder="选择类型" style="width: 100%;">
                <el-option label="系统公告" :value="1" />
                <el-option label="活动公告" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间">
              <el-date-picker 
                v-model="form.start_time" 
                type="datetime" 
                placeholder="选择开始时间" 
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间">
              <el-date-picker 
                v-model="form.end_time" 
                type="datetime" 
                placeholder="选择结束时间" 
                style="width: 100%;"
              />
            </el-form-item>
          </el-col>
        </el-row>
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

const tableData = ref([])
const dialogVisible = ref(false)
const dialogType = ref('add')

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const form = reactive({
  id: '',
  title: '',
  content: '',
  platform: 0,
  type: 1,
  start_time: '',
  end_time: '',
  status: 1
})

const getPlatformName = (platform) => {
  const map = { 0: '全部', 1: '抖音', 2: '快手', 3: 'TikTok' }
  return map[platform] || '-'
}

const getPlatformType = (platform) => {
  const map = { 1: 'danger', 2: 'warning', 3: 'info' }
  return map[platform] || ''
}

const openAddDialog = () => {
  dialogType.value = 'add'
  Object.assign(form, {
    id: '', title: '', content: '', platform: 0, type: 1,
    start_time: '', end_time: '', status: 1
  })
  dialogVisible.value = true
}

const openEditDialog = (row) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

const toggleStatus = (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  ElMessage.success(newStatus === 1 ? '已显示' : '已隐藏')
  loadData()
}

const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该公告吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('删除成功')
    loadData()
  })
}

const submitForm = () => {
  if (!form.title || !form.content) {
    ElMessage.warning('请填写必填项')
    return
  }
  ElMessage.success(dialogType.value === 'add' ? '发布成功' : '更新成功')
  dialogVisible.value = false
  loadData()
}

const loadData = () => {
  tableData.value = [
    {
      title: '测试公告',
      content: '这是一个测试公告内容',
      platform: 1,
      type: 1,
      status: 1,
      start_time: '2025-01-01 00:00:00',
      end_time: '2025-12-31 23:59:59'
    }
  ]
  pagination.total = 1
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.announcement-page {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
