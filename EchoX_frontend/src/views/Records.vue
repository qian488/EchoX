<template>
  <div class="records-container">
    <el-row :gutter="20">
      <el-col :span="18">
        <el-card class="records-card">
          <template #header>
            <div class="card-header">
              <h2>游戏记录</h2>
              <div class="header-actions">
                <el-button-group class="filter-buttons">
                  <el-button
                    :type="selectedTags.length > 0 ? 'primary' : 'default'"
                    @click="showTagFilter = !showTagFilter"
                  >
                    <el-icon><Filter /></el-icon>
                    标签筛选
                    <el-badge
                      v-if="selectedTags.length > 0"
                      :value="selectedTags.length"
                      class="filter-badge"
                    />
                  </el-button>
                  <el-button
                    v-if="selectedTags.length > 0"
                    type="danger"
                    plain
                    @click="clearTagFilter"
                  >
                    清除筛选
                  </el-button>
                </el-button-group>
                <el-button type="primary" @click="showCreateDialog">
                  创建记录
                </el-button>
              </div>
            </div>
          </template>

          <!-- 标签筛选面板 -->
          <el-collapse-transition>
            <div v-show="showTagFilter" class="filter-panel">
              <div class="filter-header">
                <span>选择标签进行筛选</span>
                <el-button
                  v-if="selectedTags.length > 0"
                  type="text"
                  @click="clearTagFilter"
                >
                  清除已选
                </el-button>
              </div>
              <div class="selected-tags">
                <el-tag
                  v-for="tagId in selectedTags"
                  :key="tagId"
                  closable
                  :style="{
                    backgroundColor: getTagColor(tagId),
                    color: '#fff',
                    margin: '4px'
                  }"
                  @close="removeTag(tagId)"
                >
                  {{ getTagName(tagId) }}
                </el-tag>
              </div>
              <el-select
                v-model="selectedTags"
                multiple
                filterable
                collapse-tags
                placeholder="选择标签筛选"
                class="tag-select"
              >
                <el-option
                  v-for="tag in tags"
                  :key="tag.id"
                  :label="tag.name"
                  :value="tag.id"
                >
                  <div class="tag-option">
                    <div class="tag-color" :style="{ backgroundColor: tag.color }"></div>
                    <span>{{ tag.name }}</span>
                    <span class="tag-count">({{ getTagCount(tag.id) }})</span>
                  </div>
                </el-option>
              </el-select>
            </div>
          </el-collapse-transition>

          <el-table :data="filteredRecords" style="width: 100%">
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="type" label="类型">
              <template #default="{ row }">
                <el-tag :type="getTypeTag(row.type)">{{ getTypeLabel(row.type) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="标签">
              <template #default="{ row }">
                <el-tag
                  v-for="tag in row.tags"
                  :key="tag.id"
                  :style="{ backgroundColor: tag.color, color: '#fff', marginRight: '4px' }"
                >
                  {{ tag.name }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" />
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
      </el-col>

      <el-col :span="6">
        <TagCloud
          :tags="tagsWithCount"
          @select="handleTagCloudSelect"
        />
      </el-col>
    </el-row>

    <!-- 创建/编辑记录对话框 -->
    <el-dialog
      :title="dialogType === 'create' ? '创建记录' : '编辑记录'"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form
        ref="recordFormRef"
        :model="recordForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="recordForm.title" />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-select v-model="recordForm.type" placeholder="选择记录类型">
            <el-option label="文本" value="text" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
          </el-select>
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input
            v-model="recordForm.content"
            type="textarea"
            :rows="4"
            v-if="recordForm.type === 'text'"
          />
          <el-upload
            v-else
            class="upload-demo"
            action="http://localhost:8080/api/v1/upload"
            :headers="uploadHeaders"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
          >
            <el-button type="primary">点击上传</el-button>
          </el-upload>
        </el-form-item>

        <el-form-item label="标签" prop="tags">
          <el-select
            v-model="recordForm.tags"
            multiple
            collapse-tags
            placeholder="选择标签"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            >
              <div class="tag-option">
                <div class="tag-color" :style="{ backgroundColor: tag.color }"></div>
                <span>{{ tag.name }}</span>
              </div>
            </el-option>
          </el-select>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Filter } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import { recordService } from '../services/record'
import { tagService } from '../services/tag'
import type { Record, Tag } from '../types/record'
import TagCloud from '../components/TagCloud.vue'

const userStore = useUserStore()
const records = ref<Record[]>([])
const tags = ref<Tag[]>([])
const selectedTags = ref<number[]>([])
const dialogVisible = ref(false)
const dialogType = ref<'create' | 'edit'>('create')
const loading = ref(false)
const recordFormRef = ref<FormInstance>()
const showTagFilter = ref(false)

const recordForm = reactive({
  id: 0,
  title: '',
  content: '',
  type: 'text' as 'text' | 'image' | 'video',
  media_url: '' as string,
  tags: [] as number[]
})

const rules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择类型', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入内容', trigger: 'blur' }
  ]
}

const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const filteredRecords = computed(() => {
  if (selectedTags.value.length === 0) {
    return records.value
  }
  return records.value.filter(record => 
    record.tags.some(tag => selectedTags.value.includes(tag.id))
  )
})

const getTypeLabel = (type: 'text' | 'image' | 'video') => {
  const types = {
    text: '文本',
    image: '图片',
    video: '视频'
  } as const
  return types[type]
}

const getTypeTag = (type: 'text' | 'image' | 'video') => {
  const types = {
    text: 'info',
    image: 'success',
    video: 'warning'
  } as const
  return types[type]
}

const fetchRecords = async () => {
  try {
    records.value = await recordService.getRecords()
  } catch (error) {
    ElMessage.error('获取记录列表失败')
  }
}

const fetchTags = async () => {
  try {
    tags.value = await tagService.getTags()
  } catch (error) {
    ElMessage.error('获取标签列表失败')
  }
}

const showCreateDialog = () => {
  dialogType.value = 'create'
  recordForm.id = 0
  recordForm.title = ''
  recordForm.content = ''
  recordForm.type = 'text'
  recordForm.media_url = ''
  recordForm.tags = []
  dialogVisible.value = true
}

const showEditDialog = (record: Record) => {
  dialogType.value = 'edit'
  recordForm.id = record.id
  recordForm.title = record.title
  recordForm.content = record.content
  recordForm.type = record.type
  recordForm.media_url = record.media_url || ''
  recordForm.tags = record.tags.map((tag: Tag) => tag.id)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!recordFormRef.value) return

  await recordFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          title: recordForm.title,
          content: recordForm.content,
          type: recordForm.type,
          media_url: recordForm.media_url,
          tag_ids: recordForm.tags
        }

        if (dialogType.value === 'create') {
          await recordService.createRecord(data)
          ElMessage.success('记录创建成功')
        } else {
          await recordService.updateRecord(recordForm.id, data)
          ElMessage.success('记录更新成功')
        }
        dialogVisible.value = false
        fetchRecords()
      } catch (error) {
        ElMessage.error(dialogType.value === 'create' ? '创建记录失败' : '更新记录失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const handleDelete = async (record: Record) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await recordService.deleteRecord(record.id)
    ElMessage.success('记录删除成功')
    fetchRecords()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除记录失败')
    }
  }
}

const handleUploadSuccess = (response: any) => {
  recordForm.media_url = response.url
  ElMessage.success('上传成功')
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isVideo = file.type.startsWith('video/')
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isImage && !isVideo) {
    ElMessage.error('只能上传图片或视频文件！')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB！')
    return false
  }
  return true
}

// 计算标签使用次数
const tagsWithCount = computed(() => {
  const tagCounts = new Map<number, number>()
  
  // 统计每个标签的使用次数
  records.value.forEach(record => {
    record.tags.forEach(tag => {
      tagCounts.set(tag.id, (tagCounts.get(tag.id) || 0) + 1)
    })
  })
  
  // 将统计结果添加到标签数据中
  return tags.value.map(tag => ({
    ...tag,
    count: tagCounts.get(tag.id) || 0
  }))
})

// 处理标签云点击事件
const handleTagCloudSelect = (tag: Tag) => {
  if (selectedTags.value.includes(tag.id)) {
    selectedTags.value = selectedTags.value.filter(id => id !== tag.id)
  } else {
    selectedTags.value = [...selectedTags.value, tag.id]
  }
}

const getTagColor = (tagId: number) => {
  const tag = tags.value.find(t => t.id === tagId)
  return tag?.color || '#909399'
}

const getTagName = (tagId: number) => {
  const tag = tags.value.find(t => t.id === tagId)
  return tag?.name || '未知标签'
}

const getTagCount = (tagId: number) => {
  return records.value.filter(record => 
    record.tags.some(tag => tag.id === tagId)
  ).length
}

const removeTag = (tagId: number) => {
  selectedTags.value = selectedTags.value.filter(id => id !== tagId)
}

const clearTagFilter = () => {
  selectedTags.value = []
  showTagFilter.value = false
}

onMounted(() => {
  fetchRecords()
  fetchTags()
})
</script>

<style scoped>
.records-container {
  max-width: 1200px;
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

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.filter-buttons {
  display: flex;
  align-items: center;
}

.filter-badge {
  margin-left: 8px;
}

.filter-panel {
  margin-bottom: 20px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  color: #606266;
}

.selected-tags {
  margin-bottom: 12px;
}

.tag-select {
  width: 100%;
}

.tag-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-color {
  width: 16px;
  height: 16px;
  border-radius: 4px;
}

.tag-count {
  margin-left: auto;
  color: #909399;
  font-size: 0.9em;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.upload-demo {
  width: 100%;
}
</style> 