<template>
  <div class="forgot-container">
    <el-card class="forgot-card">
      <template #header>
        <h2>重置密码</h2>
      </template>

      <el-steps :active="step" finish-status="success" simple style="margin-bottom: 24px">
        <el-step title="验证邮箱" />
        <el-step title="设置新密码" />
      </el-steps>

      <template v-if="step === 0">
        <el-form ref="emailFormRef" :model="emailForm" :rules="emailRules" label-width="80px">
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="emailForm.email" placeholder="注册时使用的邮箱" @keyup.enter="sendToken" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="sendToken">获取重置令牌</el-button>
            <el-button @click="$router.push('/login')">返回登录</el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-else>
        <el-alert
          type="warning"
          :closable="false"
          title="开发模式提示"
          description="本项目未配置邮件服务，重置令牌已直接显示在下方，请复制填入。生产环境应改为邮件发送且不展示令牌。"
          show-icon
          style="margin-bottom: 20px"
        />
        <el-form ref="resetFormRef" :model="resetForm" :rules="resetRules" label-width="90px">
          <el-form-item label="重置令牌" prop="token">
            <el-input v-model="resetForm.token" type="textarea" :rows="3" placeholder="粘贴收到的令牌" />
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input v-model="resetForm.newPassword" type="password" show-password placeholder="至少 6 位" />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input v-model="resetForm.confirmPassword" type="password" show-password placeholder="再次输入新密码" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="doReset">确认重置</el-button>
            <el-button @click="step = 0">上一步</el-button>
          </el-form-item>
        </el-form>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'

const router = useRouter()
const API_URL = 'http://localhost:8080/api/v1'
const step = ref(0)
const loading = ref(false)
const emailFormRef = ref<FormInstance>()
const resetFormRef = ref<FormInstance>()

const emailForm = reactive({ email: '' })
const resetForm = reactive({ token: '', newPassword: '', confirmPassword: '' })

const emailRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ]
}
const resetRules = {
  token: [{ required: true, message: '请输入重置令牌', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: any) => {
        if (value !== resetForm.newPassword) callback(new Error('两次密码不一致'))
        else callback()
      },
      trigger: 'blur'
    }
  ]
}

const sendToken = async () => {
  if (!emailFormRef.value) return
  await emailFormRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      const resp = await axios.post(`${API_URL}/auth/forgot-password`, { email: emailForm.email })
      if (resp.data.reset_token) resetForm.token = resp.data.reset_token
      ElMessage.success(resp.data.message || '请求已提交')
      step.value = 1
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '请求失败')
    } finally {
      loading.value = false
    }
  })
}

const doReset = async () => {
  if (!resetFormRef.value) return
  await resetFormRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    try {
      await axios.post(`${API_URL}/auth/reset-password`, {
        token: resetForm.token,
        new_password: resetForm.newPassword
      })
      ElMessage.success('密码重置成功，请重新登录')
      router.push('/login')
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '重置失败')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.forgot-container { display: flex; justify-content: center; padding: 40px 20px; }
.forgot-card { width: 100%; max-width: 560px; }
.card-header { margin: 0; font-size: 18px; }
</style>
