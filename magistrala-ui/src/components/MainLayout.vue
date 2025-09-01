<template>
  <div class="main-layout">
    <Sidebar />
    <div class="main-content">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import Sidebar from './Sidebar.vue'
</script>

<style scoped>
.main-layout {
  min-height: 100%;
}
.main-content {
  margin-left: 220px;      /* 适配 sidebar 宽度 */
  min-height: 100vh;
  background: #fff;
  padding: 30px 20px;      /* 左右 padding 一致 */
  box-sizing: border-box;
}
</style>