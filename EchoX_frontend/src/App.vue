<template>
  <el-container>
    <el-header>
      <el-menu
        :router="true"
        mode="horizontal"
        :ellipsis="false"
      >
        <el-menu-item index="/">Home</el-menu-item>
        <div class="flex-grow" />
        <template v-if="isAuthenticated">
          <el-menu-item index="/records">Records</el-menu-item>
          <el-menu-item index="/timeline">Timeline</el-menu-item>
          <el-menu-item index="/search">Search</el-menu-item>
          <el-menu-item index="/capture">Capture</el-menu-item>
          <el-menu-item index="/tags">Tags</el-menu-item>
          <el-menu-item index="/profile">Profile</el-menu-item>
          <el-menu-item @click="handleLogout">Logout</el-menu-item>
        </template>
        <template v-else>
          <el-menu-item index="/login">Login</el-menu-item>
          <el-menu-item index="/register">Register</el-menu-item>
        </template>
      </el-menu>
    </el-header>

    <el-main>
      <router-view />
    </el-main>

    <el-footer>
      <p>&copy; 2024 EchoX. All rights reserved.</p>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from './stores/user'

const router = useRouter()
const userStore = useUserStore()

const isAuthenticated = computed(() => userStore.isAuthenticated)

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('Logged out successfully')
  router.push('/login')
}
</script>

<style scoped>
.el-header {
  padding: 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.flex-grow {
  flex-grow: 1;
}

.el-main {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.el-footer {
  text-align: center;
  padding: 20px;
  background-color: #fff;
  color: #909399;
}
</style>
