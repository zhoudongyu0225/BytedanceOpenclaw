<template>
  <div class="login-page" :class="{ 'dark-mode': isDark }">
    <div class="login-card">
      <div class="login-header">
        <h1>后台管理系统</h1>
        <p>Management System</p>
      </div>
      
      <el-form :model="form" class="login-form">
        <el-form-item>
          <el-input 
            v-model="form.username" 
            placeholder="请输入用户名"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item>
          <el-input 
            v-model="form.password" 
            type="password" 
            placeholder="请输入密码"
            prefix-icon="Lock"
            size="large"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            size="large" 
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="theme-toggle">
        <el-switch
          v-model="isDark"
          inline-prompt
          active-text="暗色"
          inactive-text="亮色"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const form = reactive({
  username: '',
  password: ''
})

const loading = ref(false)
const isDark = ref(false)

// 监听主题变化
watch(isDark, (val) => {
  const theme = val ? 'dark' : 'light'
  localStorage.setItem('theme', theme)
  if (val) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
})

onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
})

const handleLogin = async () => {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const res = await request.post('/api/v1/login', {
      username: form.username,
      password: form.password
    })
    if (res.code === 200) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('userInfo', JSON.stringify(res.data.user_info))
      ElMessage.success('登录成功')
      router.push('/')
    } else {
      ElMessage.error(res.msg || '登录失败')
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.msg || '登录失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
  transition: background 0.3s;
}

.login-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
  transition: background 0.3s, box-shadow 0.3s;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  transition: color 0.3s;
}

.login-header p {
  font-size: 14px;
  color: #909399;
  transition: color 0.3s;
}

.login-form {
  margin-bottom: 24px;
}

.login-btn {
  width: 100%;
  font-weight: 500;
}

.theme-toggle {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid #eee;
  transition: border-color 0.3s;
}

/* 暗色主题 */
.dark-mode.login-page {
  background: #141414;
}

.dark-mode .login-card {
  background: #1f1f1f;
  box-shadow: 0 2px 12px rgba(0,0,0,0.3);
}

.dark-mode .login-header h1 {
  color: #e5e5e5;
}

.dark-mode .login-header p {
  color: #666;
}

.dark-mode .theme-toggle {
  border-top-color: #333;
}
</style>
