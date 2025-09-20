<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 初始化认证状态
onMounted(async () => {
  await authStore.initAuth()
})
</script>

<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
/* 全局样式防止横向滚动 */
html, body {
  margin: 0;
  padding: 0;
  overflow-x: hidden;
  width: 100vw;
  height: 100vh;
}

* {
  box-sizing: border-box;
}

#app {
  height: 100vh;
  width: 100vw;
  overflow-x: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}
</style>
