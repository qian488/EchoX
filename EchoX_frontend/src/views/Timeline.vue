<template>
  <div class="timeline-container">
    <el-card class="timeline-card">
      <template #header>
        <div class="card-header">
          <h2>游戏历程时间轴</h2>
          <el-radio-group :model-value="group" @update:model-value="handleGroupChange">
            <el-radio-button value="day">按天</el-radio-button>
            <el-radio-button value="month">按月</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <div v-loading="loading" class="timeline-body">
        <el-empty v-if="!loading && timeline.length === 0" description="暂无记录" />

        <el-timeline v-else>
          <el-timeline-item
            v-for="item in timeline"
            :key="item.date"
            :timestamp="formatDate(item.date)"
            placement="top"
            color="#409EFF"
          >
            <div class="group-header">
              <span class="group-date">{{ formatDate(item.date) }}</span>
              <el-tag size="small" type="info">{{ item.count }} 条记录</el-tag>
            </div>

            <el-collapse class="record-collapse">
              <el-collapse-item :title="`展开查看 ${item.count} 条记录`" :name="item.date">
                <div v-for="record in item.records" :key="record.id" class="record-item">
                  <div class="record-title-row">
                    <span class="record-title">{{ record.title }}</span>
                    <el-tag :type="getTypeTag(record.type)" size="small">
                      {{ getTypeLabel(record.type) }}
                    </el-tag>
                  </div>

                  <div class="record-content">{{ summarize(record.content) }}</div>

                  <div v-if="record.tags && record.tags.length > 0" class="record-tags">
                    <el-tag
                      v-for="tag in record.tags"
                      :key="tag.id"
                      size="small"
                      :style="{
                        backgroundColor: tag.color,
                        color: '#fff',
                        borderColor: tag.color
                      }"
                    >
                      {{ tag.name }}
                    </el-tag>
                  </div>

                  <div v-if="record.type === 'image' && record.media_url" class="record-media">
                    <img :src="getMediaUrl(record.media_url)" :alt="record.title" />
                  </div>
                </div>
              </el-collapse-item>
            </el-collapse>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { discoveryService } from '../services/discovery'
import type { TimelineGroup } from '../services/discovery'
import type { Record } from '../types/record'

const BACKEND_URL = 'http://localhost:8080'

const group = ref<'day' | 'month'>('day')
const timeline = ref<TimelineGroup[]>([])
const loading = ref(false)

const fetchTimeline = async () => {
  loading.value = true
  try {
    timeline.value = await discoveryService.getTimeline(group.value)
  } catch {
    timeline.value = []
    ElMessage.error('获取时间轴失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const handleGroupChange = (value: unknown) => {
  const next: 'day' | 'month' = value === 'month' ? 'month' : 'day'
  if (next !== group.value) {
    group.value = next
    fetchTimeline()
  }
}

const formatDate = (date: string) => {
  const parts = date.split('-')
  if (parts.length === 2) return `${parts[0]}年${parts[1]}月`
  if (parts.length === 3) return `${parts[0]}年${parts[1]}月${parts[2]}日`
  return date
}

const getTypeLabel = (type: Record['type']) => {
  const labels = { text: '文本', image: '图片', video: '视频' } as const
  return labels[type]
}

const getTypeTag = (type: Record['type']) => {
  const tags = { text: 'info', image: 'success', video: 'warning' } as const
  return tags[type]
}

const summarize = (content: string) => {
  if (!content) return '（无内容）'
  return content.length > 120 ? `${content.slice(0, 120)}…` : content
}

const getMediaUrl = (url?: string) => {
  if (!url) return ''
  return url.startsWith('http') ? url : `${BACKEND_URL}${url}`
}

onMounted(() => {
  fetchTimeline()
})
</script>

<style scoped>
.timeline-container {
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

.timeline-body {
  min-height: 200px;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.group-date {
  font-weight: 600;
  color: #303133;
}

.record-collapse {
  border: none;
}

.record-item {
  padding: 8px 0;
  border-bottom: 1px dashed #ebeef5;
}

.record-item:last-child {
  border-bottom: none;
}

.record-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.record-title {
  font-weight: 600;
  color: #303133;
}

.record-content {
  color: #606266;
  font-size: 0.95em;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.record-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.record-media {
  margin-top: 10px;
}

.record-media img {
  max-width: 100%;
  max-height: 260px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  object-fit: contain;
}
</style>
