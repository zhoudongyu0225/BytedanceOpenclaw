<template>
  <div style="width: 100vw; height: 100vh; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); display: flex; align-items: center; justify-content: center; overflow: hidden;">
    <div style="width: 500px; background: white; padding: 50px; border-radius: 12px; box-shadow: 0 0 30px rgba(0,0,0,0.2);">
      <h2 style="text-align: center; margin-bottom: 40px; font-size: 28px; color: #333;">游戏后台管理系统</h2>
      
      <div style="margin-bottom: 20px;">
        <label style="display: block; margin-bottom: 8px; font-size: 14px; color: #666;">用户名</label>
        <input 
          v-model="username" 
          type="text" 
          placeholder="请输入用户名"
          style="width: 100%; height: 48px; padding: 0 15px; border: 1px solid #dcdfe6; border-radius: 4px; font-size: 16px; box-sizing: border-box;"
        />
      </div>

      <div style="margin-bottom: 30px;">
        <label style="display: block; margin-bottom: 8px; font-size: 14px; color: #666;">密码</label>
        <input 
          v-model="password" 
          type="password" 
          placeholder="请输入密码"
          @keyup.enter="handleLogin"
          style="width: 100%; height: 48px; padding: 0 15px; border: 1px solid #dcdfe6; border-radius: 4px; font-size: 16px; box-sizing: border-box;"
        />
      </div>

      <button 
        @click="handleLogin"
        style="width: 100%; height: 48px; background: #409eff; color: white; border: none; border-radius: 4px; font-size: 16px; cursor: pointer;"
      >
        登录
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const router = useRouter()
const username = ref('')
const password = ref('')

const handleLogin = async () => {
  if (!username.value || !password.value) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  try {
    const res = await request.post('/api/v1/login', {
      username: username.value,
      password: password.value
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
  }
}
</script>

<style scoped>
/* 移除所有复杂样式，纯内联样式确保兼容性 */
</style>
