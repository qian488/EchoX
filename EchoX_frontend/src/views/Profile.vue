<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h2>个人资料</h2>
          <el-button type="primary" @click="isEditing = true" v-if="!isEditing">
            编辑资料
          </el-button>
        </div>
      </template>

      <el-form
        :model="profileForm"
        :rules="rules"
        ref="profileFormRef"
        label-position="top"
        :disabled="!isEditing"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="profileForm.username" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="profileForm.email" />
        </el-form-item>

        <el-form-item v-if="isEditing">
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            保存
          </el-button>
          <el-button @click="cancelEdit">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="stats-card">
      <template #header>
        <div class="card-header">
          <h3>统计信息</h3>
        </div>
      </template>

      <div class="stats-content">
        <div class="stat-item">
          <div class="stat-value">{{ stats.totalRecords }}</div>
          <div class="stat-label">总记录数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.textRecords }}</div>
          <div class="stat-label">文本记录</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.imageRecords }}</div>
          <div class="stat-label">图片记录</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ stats.videoRecords }}</div>
          <div class="stat-label">视频记录</div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { useUserStore } from '../stores/user'
import axios from 'axios'

const userStore = useUserStore()
const profileFormRef = ref<FormInstance>()
const isEditing = ref(false)
const loading = ref(false)

const profileForm = reactive({
  username: '',
  email: ''
})

const stats = reactive({
  totalRecords: 0,
  textRecords: 0,
  imageRecords: 0,
  videoRecords: 0
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const fetchProfile = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/v1/users/profile', {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    })
    const { username, email } = response.data
    profileForm.username = username
    profileForm.email = email
  } catch (error) {
    ElMessage.error('获取个人资料失败')
  }
}

const fetchStats = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/v1/records', {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    })
    const records = response.data
    stats.totalRecords = records.length
    stats.textRecords = records.filter((r: any) => r.type === 'text').length
    stats.imageRecords = records.filter((r: any) => r.type === 'image').length
    stats.videoRecords = records.filter((r: any) => r.type === 'video').length
  } catch (error) {
    ElMessage.error('获取统计信息失败')
  }
}

const handleSubmit = async () => {
  if (!profileFormRef.value) return

  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await axios.put('http://localhost:8080/api/v1/users/profile', profileForm, {
          headers: {
            Authorization: `Bearer ${userStore.token}`
          }
        })
        ElMessage.success('个人资料更新成功')
        isEditing.value = false
      } catch (error) {
        ElMessage.error('更新个人资料失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const cancelEdit = () => {
  isEditing.value = false
  fetchProfile()
}

onMounted(() => {
  fetchProfile()
  fetchStats()
})
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
}

.stats-content {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  text-align: center;
}

.stat-item {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 8px;
}

.stat-label {
  color: #606266;
}
</style> 