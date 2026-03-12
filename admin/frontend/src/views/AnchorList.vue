<template>
  <div class="anchor-list">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>主播管理</span>
          <el-button type="primary" size="small" @click="openAddDialog">新增主播</el-button>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="filter-bar">
        <el-input v-model="filter.keyword" placeholder="搜索主播姓名/昵称/直播间ID" style="width: 300px; margin-right: 20px;" clearable @keyup.enter="getList" />
        <el-select v-model="filter.platform" placeholder="平台筛选" style="width: 150px; margin-right: 20px;" clearable @change="getList">
          <el-option label="抖音" :value="1" />
          <el-option label="快手" :value="2" />
          <el-option label="TikTok" :value="3" />
        </el-select>
        <el-select v-model="filter.status" placeholder="状态筛选" style="width: 150px; margin-right: 20px;" clearable @change="getList">
          <el-option label="正常" :value="1" />
          <el-option label="封禁" :value="2" />
        </el-select>
        <el-button type="primary" size="small" @click="getList">搜索</el-button>
      </div>

      <!-- 批量操作 -->
      <div class="batch-bar" style="margin: 10px 0;">
        <el-button v-if="selectedIds.length > 0" type="danger" size="small" @click="batchUpdateStatus(2)">批量封禁</el-button>
        <el-button v-if="selectedIds.length > 0" type="success" size="small" @click="batchUpdateStatus(1)">批量解封</el-button>
      </div>

      <!-- 列表 -->
      <el-table :data="tableData" border style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="主播姓名" />
        <el-table-column prop="nickname" label="主播昵称" />
        <el-table-column prop="platform" label="平台" :formatter="formatPlatform" />
        <el-table-column prop="room_id" label="直播间ID" />
        <el-table-column prop="platform_uid" label="平台UID" />
        <el-table-column prop="status" label="状态" :formatter="formatStatus" />
        <el-table-column prop="remark" label="备注" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="openEditDialog(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination" style="margin-top: 20px; text-align: right;">
        <el-pagination
          v-model:current-page="filter.page"
          v-model:page-size="filter.page_size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="getList"
          @current-change="getList"
        />
      </div>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增主播' : '编辑主播'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="主播姓名" prop="name" required>
          <el-input v-model="form.name" placeholder="请输入主播姓名" />
        </el-form-item>
        <el-form-item label="主播昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入主播昵称" />
        </el-form-item>
        <el-form-item label="平台" prop="platform" required>
          <el-select v-model="form.platform" placeholder="请选择平台">
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

// 获取列表
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

// 格式化平台
const formatPlatform = (row) => {
  const map = { 1: '抖音', 2: '快手', 3: 'TikTok' }
  return map[row.platform] || '-'
}

// 格式化状态
const formatStatus = (row) => {
  return row.status === 1 ? '正常' : '封禁'
}

// 选择变化
const handleSelectionChange = (val) => {
  selectedIds.value = val.map(item => item.id)
}

// 打开新增弹窗
const openAddDialog = () => {
  dialogType.value = 'add'
  Object.keys(form).forEach(key => {
    form[key] = key === 'status' ? 1 : ''
  })
  dialogVisible.value = true
}

// 打开编辑弹窗
const openEditDialog = (row) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

// 提交表单
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

// 删除
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

// 批量更新状态
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
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter-bar {
  margin-bottom: 20px;
}
</style>
