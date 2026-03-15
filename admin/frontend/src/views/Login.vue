<template>
  <div class="login-page">
    <div class="login-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>
    <div class="login-card">
      <div class="login-header">
        <div class="logo">🦖</div>
        <h1>恐龙攻守</h1>
        <p>游戏后台管理系统</p>
      </div>
      
      <el-form :model="form" class="login-form">
        <el-form-item>
          <div class="input-wrapper">
            <el-icon class="input-icon"><User /></el-icon>
            <input 
              v-model="form.username" 
              type="text" 
              placeholder="请输入用户名"
              class="cyber-input"
            />
          </div>
        </el-form-item>
        <el-form-item>
          <div class="input-wrapper">
            <el-icon class="input-icon"><Lock /></el-icon>
            <input 
              v-model="form.password" 
              type="password" 
              placeholder="请输入密码"
              class="cyber-input"
              @keyup.enter="handleLogin"
            />
          </div>
        </el-form-item>
        <el-form-item>
          <button 
            type="button"
            class="login-btn"
            @click="handleLogin"
          >
            <span>登 录</span>
            <div class="btn-glow"></div>
          </button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <span>© 2026 恐龙攻守游戏</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import request from '../utils/request'

const router = useRouter()

const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
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
  position: relative;
  overflow: hidden;
  background: #0a0a0f;
}

/* 背景动画 */
.login-bg {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: float 20s infinite;
}

.shape-1 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #00d4ff, #0066ff);
  top: -200px;
  left: -200px;
  animation-delay: 0s;
}

.shape-2 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #ff00ff, #6600ff);
  bottom: -150px;
  right: -150px;
  animation-delay: -7s;
}

.shape-3 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #00ff88, #00d4ff);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(50px, -50px) scale(1.1); }
  50% { transform: translate(-30px, 30px) scale(0.9); }
  75% { transform: translate(-50px, -30px) scale(1.05); }
}

/* 登录卡片 */
.login-card {
  width: 420px;
  padding: 50px 40px;
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 20px;
  backdrop-filter: blur(20px);
  box-shadow: 0 25px 50px rgba(0,0,0,0.5);
  position: relative;
  z-index: 10;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  font-size: 60px;
  margin-bottom: 16px;
  animation: logo-pulse 2s infinite;
}

@keyframes logo-pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
}

.login-header h1 {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(90deg, #00d4ff, #00ff88);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.login-header p {
  color: rgba(255,255,255,0.5);
  font-size: 14px;
}

/* 输入框 */
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  color: rgba(255,255,255,0.4);
  font-size: 18px;
  z-index: 1;
}

.cyber-input {
  width: 100%;
  height: 50px;
  padding: 0 16px 0 48px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 10px;
  color: #fff;
  font-size: 15px;
  outline: none;
  transition: all 0.3s;
}

.cyber-input::placeholder {
  color: rgba(255,255,255,0.3);
}

.cyber-input:focus {
  border-color: #00d4ff;
  box-shadow: 0 0 20px rgba(0,212,255,0.2);
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 50px;
  background: linear-gradient(90deg, #00d4ff, #00ff88);
  border: none;
  border-radius: 10px;
  color: #0a0a0f;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(0,212,255,0.4);
}

.login-btn span {
  position: relative;
  z-index: 1;
}

.btn-glow {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(
    45deg,
    transparent 30%,
    rgba(255,255,255,0.3) 50%,
    transparent 70%
  );
  animation: shine 3s infinite;
}

@keyframes shine {
  0% { transform: translateX(-100%) rotate(45deg); }
  100% { transform: translateX(100%) rotate(45deg); }
}

.login-footer {
  text-align: center;
  margin-top: 30px;
  color: rgba(255,255,255,0.3);
  font-size: 12px;
}
</style>
