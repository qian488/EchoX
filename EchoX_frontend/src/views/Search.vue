<template>
  <div class="search-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>搜索记录</h2>
          <el-button v-if="hasFilter" type="danger" plain @click="resetFilters">清除条件</el-button>
        </div>
      </template>

      <el-form :inline="true" class="search-form" @submit.prevent="handleSearch">
        <el-form-item label="关键词">
          <el-input
            v-model="filters.q"
            placeholder="标题或内容"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filters.type" placeholder="全部" clearable style="width: 120px">
            <el-option label="文字" value="text" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-select v-model="filters.tag" placeholder="全部" clearable filterable style="width: 140px">
            <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始">
          <el-date-picker
            v-model="filters.start"
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="开始日期"
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="结束">
          <el-date-picker
            v-model="filters.end"
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="结束日期"
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
        </el-form-item>
      </el-form>

      <div v-if="searched" class="result-summary">
        共找到 <strong>{{ results.length }}</strong> 条记录
      </div>

      <el-table v-if="results.length" :data="results" style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="标题" min-width="160" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="typeTagType(row.type)">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容摘要" min-width="200" show-overflow-tooltip />
        <el-table-column label="标签" min-width="140">
          <template #default="{ row }">
            <el-tag
              v-for="tag in row.tags"
              :key="tag.id"
              size="small"
              :style="{ backgroundColor: tag.color, color: '#fff', marginRight: '4px' }"
            >
              {{ tag.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
      </el-table>

      <el-empty v-else-if="searched && !loading" description="没有找到匹配的记录" />
      <el-empty v-else-if="!searched" description="输入条件开始搜索" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { discoveryService, type SearchParams } from '../services/discovery'
import { tagService } from '../services/tag'
import type { Record, Tag } from '../types/record'

const filters = reactive<SearchParams>({ q: '', type: '', tag: '', start: '', end: '' })
const results = ref<Record[]>([])
const tags = ref<Tag[]>([])
const loading = ref(false)
const searched = ref(false)

const hasFilter = computed(() =>
  !!(filters.q || filters.type || filters.tag || filters.start || filters.end)
)

const typeLabel = (t: string) => ({ text: '文字', image: '图片', video: '视频' }[t] || t)
const typeTagType = (t: string) =>
  ({ text: 'primary', image: 'success', video: 'warning' }[t] || 'info') as
    | 'primary' | 'success' | 'warning' | 'info'

const handleSearch = async () => {
  loading.value = true
  try {
    const params: SearchParams = {}
    if (filters.q) params.q = filters.q
    if (filters.type) params.type = filters.type
    if (filters.tag) params.tag = filters.tag
    if (filters.start) params.start = filters.start
    if (filters.end) params.end = filters.end
    results.value = await discoveryService.search(params)
    searched.value = true
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '搜索失败')
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.q = ''
  filters.type = ''
  filters.tag = ''
  filters.start = ''
  filters.end = ''
  results.value = []
  searched.value = false
}

onMounted(async () => {
  try {
    tags.value = await tagService.getTags()
  } catch {
    // 标签加载失败不影响搜索
  }
})
</script>

<style scoped>
.search-container { max-width: 1100px; margin: 0 auto; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.card-header h2 { margin: 0; font-size: 18px; }
.search-form { margin-bottom: 10px; }
.result-summary { margin: 10px 0; color: #606266; }
</style>
