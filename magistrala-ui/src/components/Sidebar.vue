<template>
  <aside class="sidebar">
    <div class="sidebar-logo">Abstract Machines</div>
    <div class="sidebar-domain">
      <span>{{ domainName }}</span>
    </div>
    <div class="sidebar-scroll">
      <nav class="sidebar-menu">
        <router-link
          :to="`/domain/${domainId}`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}`, true) }"
        >Home</router-link>
        <router-link
          :to="`/domain/${domainId}/dashboards`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/dashboards`) }"
        >Dashboards</router-link>
        <div class="menu-section">Domain Management</div>
        <router-link
          :to="`/domain/${domainId}/domain`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/domain`) }"
        >Domain</router-link>
        <router-link
          :to="`/domain/${domainId}/members`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/members`) }"
        >Members</router-link>
        <router-link
          :to="`/domain/${domainId}/roles`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/roles`) }"
        >Roles</router-link>
        <router-link
          :to="`/domain/${domainId}/invitations`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/invitations`) }"
        >Invitations</router-link>
        <div class="menu-section">Clients Management</div>
        <router-link
          :to="`/domain/${domainId}/groups`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/groups`) }"
        >Groups</router-link>
        <router-link
          :to="`/domain/${domainId}/channels`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/channels`) }"
        >Channels</router-link>
        <router-link
          :to="`/domain/${domainId}/clients`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/clients`) }"
        >Clients</router-link>
        <router-link
          :to="`/domain/${domainId}/rules`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/rules`) }"
        >Rules</router-link>
        <router-link
          :to="`/domain/${domainId}/alarms`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/alarms`) }"
        >Alarms</router-link>
        <router-link
          :to="`/domain/${domainId}/reports`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/reports`) }"
        >Reports</router-link>
        <div class="menu-section">Smart Agriculture</div>
        <router-link
          :to="`/domain/${domainId}/agri-integration`"
          class="menu-item"
          :class="{ active: isActive(`/domain/${domainId}/agri-integration`) }"
        >🌾 农业数据集成</router-link>
      </nav>
    </div>
    <div class="sidebar-support">Support</div>
  </aside>
</template>

<script setup>
import { useRoute } from 'vue-router'
const route = useRoute()
const domainId = route.params.id || ''
const domainName = 'testDomain'

// 判断当前菜单是否选中
function isActive(path, exact = false) {
  if (exact) {
    return route.path === path
  }
  return route.path.startsWith(path)
}
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  width: 220px;
  height: 100%;
  background: #0a3566;
  color: #fff;
  display: flex;
  flex-direction: column;
}
.sidebar-logo {
  font-size: 1.4em;
  font-family: 'Segoe UI', cursive;
  text-align: center;
  margin: 18px 0 8px 0;
  letter-spacing: 2px;
}
.sidebar-domain {
  font-size: 1.1em;
  font-weight: bold;
  text-align: center;
  margin-bottom: 18px;
}
.sidebar-scroll {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.sidebar-scroll::-webkit-scrollbar {
  display: none;
}
.sidebar-menu {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 18px 18px 0 18px; /* 顶部增加18px内边距 */
}
.menu-item {
  color: #fff;
  text-decoration: none;
  padding: 8px 0 8px 8px;
  border-radius: 4px;
  font-size: 1em;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background 0.2s, color 0.2s;
}
.menu-item.active {
  background: #fff;
  color: #174e8a;
  font-weight: bold;
  border-radius: 8px;
  box-shadow: 0 0 0 2px #fff;
}
.menu-section {
  margin: 14px 0 4px 0;
  font-size: 0.95em;
  color: #b0c4de;
  font-weight: bold;
}
.sidebar-support {
  text-align: center;
  padding: 16px 0;
  color: #b0c4de;
  font-size: 1em;
}
</style>