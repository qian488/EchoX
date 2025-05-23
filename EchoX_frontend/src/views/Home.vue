<template>
  <div class="home-container">
    <el-row justify="center" align="middle" class="hero-section">
      <el-col :span="16" class="text-center">
        <h1>Welcome to EchoX</h1>
        <p class="subtitle">记录你的游戏历程，创造独特的回忆</p>
        <div class="action-buttons">
          <el-button type="primary" size="large" @click="$router.push('/register')" v-if="!isAuthenticated">
            开始使用
          </el-button>
          <el-button type="success" size="large" @click="$router.push('/records')" v-else>
            查看记录
          </el-button>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="features-section">
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="feature-header">
              <el-icon><Document /></el-icon>
              <span>多样化记录</span>
            </div>
          </template>
          <p>支持文字、图片、视频等多种格式，全面记录游戏精彩瞬间</p>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="feature-header">
              <el-icon><Files /></el-icon>
              <span>个性化分类</span>
            </div>
          </template>
          <p>自定义标签和分类，打造专属于你的游戏回忆库</p>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="feature-header">
              <el-icon><Upload /></el-icon>
              <span>云同步</span>
            </div>
          </template>
          <p>数据自动同步，随时随地访问你的游戏记录</p>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Document, Files, Upload } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const isAuthenticated = computed(() => userStore.isAuthenticated)
</script>

<style scoped>
.home-container {
  padding: 40px 20px;
}

.hero-section {
  margin-bottom: 60px;
}

.text-center {
  text-align: center;
}

.subtitle {
  font-size: 1.2em;
  color: #666;
  margin: 20px 0;
}

.action-buttons {
  margin-top: 30px;
}

.features-section {
  margin-top: 40px;
}

.feature-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.1em;
  font-weight: bold;
}

.el-card {
  height: 100%;
  transition: transform 0.3s;
}

.el-card:hover {
  transform: translateY(-5px);
}
</style> 