<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <h2>Create an Account</h2>
      </template>
      
      <el-form :model="registerForm" :rules="rules" ref="registerFormRef" label-position="top">
        <el-form-item label="Username" prop="username">
          <el-input v-model="registerForm.username" placeholder="Choose a username" />
        </el-form-item>
        
        <el-form-item label="Email" prop="email">
          <el-input v-model="registerForm.email" placeholder="Enter your email" />
        </el-form-item>
        
        <el-form-item label="Password" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="Choose a password" />
        </el-form-item>
        
        <el-form-item label="Confirm Password" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="Confirm your password" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleRegister" :loading="loading" class="submit-btn">
            Register
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-link">
        Already have an account? <router-link to="/login">Login here</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validatePass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('Please enter your password'))
  } else {
    if (registerForm.confirmPassword !== '') {
      if (registerFormRef.value) {
        registerFormRef.value.validateField('confirmPassword')
      }
    }
    callback()
  }
}

const validatePass2 = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('Please enter your password again'))
  } else if (value !== registerForm.password) {
    callback(new Error('Two passwords do not match!'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: 'Please enter your username', trigger: 'blur' },
    { min: 3, max: 20, message: 'Length should be 3 to 20 characters', trigger: 'blur' }
  ],
  email: [
    { required: true, message: 'Please enter your email', trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email address', trigger: 'blur' }
  ],
  password: [
    { required: true, validator: validatePass, trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass2, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const success = await userStore.register(
          registerForm.username,
          registerForm.password,
          registerForm.email
        )
        if (success) {
          ElMessage.success('Registration successful! Please login.')
          router.push('/login')
        } else {
          ElMessage.error('Registration failed. Please try again.')
        }
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
}

.register-card {
  width: 100%;
  max-width: 400px;
}

.submit-btn {
  width: 100%;
}

.login-link {
  text-align: center;
  margin-top: 20px;
  color: #606266;
}

.login-link a {
  color: #409EFF;
  text-decoration: none;
}

.login-link a:hover {
  text-decoration: underline;
}
</style> 