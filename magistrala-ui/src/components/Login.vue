<template>
  <div class="login-container">
    <div class="login-title">Abstract Machines</div>
    <div class="login-divider"></div>
    <div class="login-subtitle">Sign In</div>
    <form class="login-form" @submit.prevent="handleLogin">
      <label class="login-label" for="username">Email or Username</label>
      <input id="username" v-model="username" type="text" placeholder="Email or Username" required />

      <div class="login-row">
        <label class="login-label" for="password">Password</label>
        <div class="forgot-link" @click="forgotPwd">Forgot password?</div>
      </div>
      <div class="input-group">
        <input id="password" :type="showPwd ? 'text' : 'password'" v-model="password" placeholder="Password" required />
        <button type="button" class="eye-btn" @click="showPwd = !showPwd">👁️</button>
      </div>
      <button type="submit">Sign In</button>
    </form>
    <div class="login-links">
      Not registered? <router-link to="/signup">Sign Up</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const showPwd = ref(false)
const router = useRouter()

function handleLogin() {
  fetch('/Users/tokens/issue', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: username.value, password: password.value })
  })
    .then(res => res.json())
    .then(data => {
      console.log('登录接口返回:', data)
      // 支持 access_token 字段
      const token = data.token || data.access_token
      if (token) {
        alert('登录成功！')
        localStorage.setItem('token', token)
        router.push('/domains')
      } else {
        alert('登录失败：' + (data.message || '用户名或密码错误'))
      }
    })
    .catch(err => {
      console.error('登录请求异常:', err)
      alert('登录请求失败！')
    })
}

function forgotPwd() {
  alert('请联系管理员重置密码')
}
</script>

<style scoped>
.login-container {
  width: 400px;
  margin: 80px auto;
  background: #0a3566;
  border-radius: 16px;
  box-shadow: 0 0 24px #ccc;
  padding: 40px 32px;
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.login-title {
  font-size: 2.2em;
  font-family: 'Segoe UI', cursive;
  margin-bottom: 10px;
  letter-spacing: 2px;
  text-align: center;
}
.login-subtitle {
  font-size: 1.3em;
  margin-bottom: 24px;
  text-align: center;
}
.login-divider {
  width: 100%;
  height: 2px;
  background: #fff;
  opacity: 0.3;
  margin-bottom: 18px;
  border-radius: 1px;
}
.login-form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 18px;
  align-items: center;
}
.login-form input,
.login-form button[type="submit"] {
  width: 100%;
  box-sizing: border-box;
  display: block;
  margin: 0 auto;
}
.login-form input {
  padding: 12px;
  border-radius: 8px;
  border: none;
  font-size: 1.1em;
  margin-bottom: 4px;
  background: #f5f8ff;
  color: #222;
}
.login-form input:focus {
  outline: 2px solid #0a3566;
}
.input-group {
  width: 100%;
  position: relative;
}
.input-group input {
  width: 100%;
  padding-right: 40px; /* 给眼睛按钮留空间 */
}
.eye-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #0a3566;
  cursor: pointer;
  font-size: 1.2em;
}
.login-label {
  font-size: 1em;
  color: #fff;
  margin-bottom: 4px;
  display: block;
  text-align: left; /* 左对齐 */
  margin-left: 4px;
}

.login-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
  width: 100%;
}
.forgot-link {
  font-size: 0.95em;
  color: #fff;
  text-decoration: underline;
  cursor: pointer;
  margin-bottom: 0;
  text-align: right;
}
.login-form button[type="submit"] {
  padding: 12px;
  background: #fff;
  color: #0a3566;
  border: none;
  border-radius: 8px;
  font-size: 1.1em;
  cursor: pointer;
  font-weight: bold;
  margin-top: 8px;
  transition: background 0.2s;
}
.login-form button[type="submit"]:hover {
  background: #e0e0e0;
}
.login-links {
  margin-top: 18px;
  font-size: 1em;
  color: #fff;
  text-align: center;
}
.login-links a, .login-links .router-link-active {
  color: #fff;
  text-decoration: underline;
  margin-left: 8px;
}
</style>