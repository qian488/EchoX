<template>
  <div class="tags-container">
    <el-card class="tags-card">
      <template #header>
        <div class="card-header">
          <h2>标签管理</h2>
          <el-button type="primary" @click="showCreateDialog">
            创建标签
          </el-button>
        </div>
      </template>

      <el-table :data="tags" style="width: 100%">
        <el-table-column prop="name" label="标签名" />
        <el-table-column label="颜色">
          <template #default="{ row }">
            <div class="color-preview" :style="{ backgroundColor: row.color }"></div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" @click="showEditDialog(row)">
                编辑
              </el-button>
              <el-button type="danger" @click="handleDelete(row)">
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑标签对话框 -->
    <el-dialog
      :title="dialogType === 'create' ? '创建标签' : '编辑标签'"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form
        ref="tagFormRef"
        :model="tagForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="标签名" prop="name">
          <el-input v-model="tagForm.name" />
        </el-form-item>
        <el-form-item label="颜色" prop="color">
          <el-color-picker v-model="tagForm.color" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import axios from 'axios'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const tags = ref([])
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const loading = ref(false)
const tagFormRef = ref<FormInstance>()

const tagForm = reactive({
  id: 0,
  name: '',
  color: '#409EFF'
})

const rules = {
  name: [
    { required: true, message: '请输入标签名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  color: [
    { required: true, message: '请选择颜色', trigger: 'change' }
  ]
}

const fetchTags = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/v1/tags', {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    })
    tags.value = response.data
  } catch (error) {
    ElMessage.error('获取标签列表失败')
  }
}

const showCreateDialog = () => {
  dialogType.value = 'create'
  tagForm.id = 0
  tagForm.name = ''
  tagForm.color = '#409EFF'
  dialogVisible.value = true
}

const showEditDialog = (tag: any) => {
  dialogType.value = 'edit'
  tagForm.id = tag.id
  tagForm.name = tag.name
  tagForm.color = tag.color
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!tagFormRef.value) return

  await tagFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (dialogType.value === 'create') {
          await axios.post('http://localhost:8080/api/v1/tags', tagForm, {
            headers: {
              Authorization: `Bearer ${userStore.token}`
            }
          })
          ElMessage.success('标签创建成功')
        } else {
          await axios.put(`http://localhost:8080/api/v1/tags/${tagForm.id}`, tagForm, {
            headers: {
              Authorization: `Bearer ${userStore.token}`
            }
          })
          ElMessage.success('标签更新成功')
        }
        dialogVisible.value = false
        fetchTags()
      } catch (error) {
        ElMessage.error(dialogType.value === 'create' ? '创建标签失败' : '更新标签失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const handleDelete = async (tag: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await axios.delete(`http://localhost:8080/api/v1/tags/${tag.id}`, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    })
    ElMessage.success('标签删除成功')
    fetchTags()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除标签失败')
    }
  }
}

onMounted(() => {
  fetchTags()
})
</script>

<style scoped>
.tags-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 