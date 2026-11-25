<template>
  <div class="capture-container">
    <el-card class="capture-card">
      <template #header>
        <div class="card-header">
          <h2>数据采集 / 快捷记录</h2>
        </div>
      </template>

      <el-alert
        class="capture-tip"
        type="info"
        :closable="false"
        show-icon
        title="可直接按 Ctrl+V 粘贴截图，自动作为图片文件上传；标签也可在「游戏记录」页面管理。"
      />

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="90px"
        class="capture-form"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="给这次游戏历程起个标题" maxlength="100" />
        </el-form-item>

        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="选择记录类型" class="full-width">
            <el-option label="文本" value="text" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
          </el-select>
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="4"
            placeholder="记录玩了什么、有什么感想…"
          />
        </el-form-item>

        <el-form-item label="标签">
          <el-select
            v-model="form.tag_ids"
            multiple
            filterable
            collapse-tags
            placeholder="选择标签（可多选）"
            class="full-width"
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

        <el-form-item label="附件">
          <div class="file-area">
            <div class="file-actions">
              <input
                ref="fileInputRef"
                type="file"
                accept="image/*,video/*"
                class="hidden-file-input"
                @change="handleFileInputChange"
              />
              <el-button type="primary" plain @click="openFilePicker">
                选择图片/视频
              </el-button>
              <span class="file-tip">支持 png/jpg/jpeg/gif/webp/mp4/webm/mov/avi，最大 100MB</span>
            </div>

            <div v-if="selectedFile" class="selected-file">
              <span class="file-name">{{ selectedFile.name }}</span>
              <span class="file-size">({{ formatSize(selectedFile.size) }})</span>
              <el-button type="danger" link @click="removeFile">移除</el-button>
            </div>

            <el-progress
              v-if="submitting && selectedFile"
              :percentage="uploadPercent"
              :stroke-width="8"
              class="upload-progress"
            />
          </div>
        </el-form-item>
      </el-form>

      <div class="form-footer">
        <el-button @click="resetForm">重置</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          提交记录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { recordService } from '../services/record'
import { tagService } from '../services/tag'
import { mediaService } from '../services/media'
import type { Tag } from '../types/record'

const MAX_FILE_SIZE = 100 * 1024 * 1024

const formRef = ref<FormInstance>()
const fileInputRef = ref<HTMLInputElement | null>(null)
const tags = ref<Tag[]>([])
const selectedFile = ref<File | null>(null)
const uploadPercent = ref(0)
const submitting = ref(false)

const form = reactive({
  title: '',
  content: '',
  type: 'text' as 'text' | 'image' | 'video',
  media_url: '',
  tag_ids: [] as number[]
})

const rules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择类型', trigger: 'change' }
  ]
}

const fetchTags = async () => {
  try {
    tags.value = await tagService.getTags()
  } catch {
    ElMessage.error('获取标签列表失败')
  }
}

const formatSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

const setSelectedFile = (file: File) => {
  const isImage = file.type.startsWith('image/')
  const isVideo = file.type.startsWith('video/')
  if (!isImage && !isVideo) {
    ElMessage.error('仅支持上传图片或视频文件')
    return
  }
  if (file.size > MAX_FILE_SIZE) {
    ElMessage.error('文件大小不能超过 100MB')
    return
  }
  selectedFile.value = file
  form.type = isVideo ? 'video' : 'image'
  form.media_url = ''
  uploadPercent.value = 0
}

const openFilePicker = () => {
  fileInputRef.value?.click()
}

const handleFileInputChange = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) setSelectedFile(file)
  input.value = ''
}

const removeFile = () => {
  selectedFile.value = null
  form.media_url = ''
  uploadPercent.value = 0
  if (fileInputRef.value) fileInputRef.value.value = ''
}

// 粘贴截图：从剪贴板中提取图片文件并作为待上传附件
const handlePaste = (event: ClipboardEvent) => {
  const clipboardData = event.clipboardData
  if (!clipboardData) return

  let imageFile: File | null = Array.from(clipboardData.files).find((file) =>
    file.type.startsWith('image/')
  ) ?? null

  if (!imageFile) {
    for (const item of Array.from(clipboardData.items)) {
      if (item.kind === 'file' && item.type.startsWith('image/')) {
        imageFile = item.getAsFile()
        if (imageFile) break
      }
    }
  }

  if (imageFile) {
    setSelectedFile(imageFile)
    ElMessage.success('已粘贴截图，提交时将自动上传')
    event.preventDefault()
  }
}

const resetForm = () => {
  form.title = ''
  form.content = ''
  form.type = 'text'
  form.media_url = ''
  form.tag_ids = []
  removeFile()
  formRef.value?.clearValidate()
}

const handleSubmit = async () => {
  if (!formRef.value || submitting.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    uploadPercent.value = 0
    try {
      let mediaUrl = ''
      if (selectedFile.value) {
        const fileType: 'image' | 'video' = selectedFile.value.type.startsWith('video/')
          ? 'video'
          : 'image'
        const uploaded = await mediaService.upload(selectedFile.value, fileType, (percent) => {
          uploadPercent.value = percent
        })
        mediaUrl = uploaded.url
        form.media_url = uploaded.url
      }

      await recordService.createRecord({
        title: form.title,
        content: form.content,
        type: form.type,
        media_url: mediaUrl,
        tag_ids: form.tag_ids
      })

      ElMessage.success('记录创建成功')
      resetForm()
    } catch {
      ElMessage.error('创建记录失败，请检查网络或稍后重试')
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  fetchTags()
  document.addEventListener('paste', handlePaste)
})

onUnmounted(() => {
  document.removeEventListener('paste', handlePaste)
})
</script>

<style scoped>
.capture-container {
  max-width: 900px;
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

.capture-tip {
  margin-bottom: 20px;
}

.full-width {
  width: 100%;
}

.file-area {
  width: 100%;
}

.file-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.hidden-file-input {
  display: none;
}

.file-tip {
  color: #909399;
  font-size: 0.85em;
}

.selected-file {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 10px;
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.file-name {
  color: #303133;
  font-weight: 600;
}

.file-size {
  color: #909399;
  font-size: 0.9em;
}

.upload-progress {
  margin-top: 10px;
}

.form-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
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
</style>
