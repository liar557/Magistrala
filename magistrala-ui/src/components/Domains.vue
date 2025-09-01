<template>
  <div class="domains-container">
    <header class="domains-header">
      <div class="logo">Abstract Machines</div>
      <div class="header-actions">
        <span class="help-icon">?</span>
        <span class="user-icon">👤</span>
      </div>
    </header>
    <div class="domains-toolbar">
      <input class="search-input" type="text" v-model="search" placeholder="Search Domain" />
      <button class="create-btn" @click="showCreate = true">+ Create</button>
    </div>
    <div class="domains-list">
      <div v-if="domains.length === 0" class="empty-tip">
        No domains yet, click "Create" to add one.
      </div>
      <div
        v-for="domain in filteredDomains"
        :key="domain.id"
        class="domain-card"
        @click="goDomain(domain.id)"
        style="cursor:pointer;"
      >
        <div class="domain-avatar">{{ domain.name[0] }}</div>
        <div class="domain-info">
          <div class="domain-name">{{ domain.name }}</div>
          <div class="domain-admin">{{ domain.role_name }}</div>
          <div class="domain-status">
            <span :class="domain.status === 'enabled' ? 'enabled' : 'disabled'">
              {{ domain.status }}
            </span>
          </div>
          <div class="domain-route">{{ domain.route }}</div>
        </div>
      </div>
    </div>
    <!-- 创建弹窗可后续实现 -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userToken = localStorage.getItem('token')

const domains = ref([])
const search = ref('')
const showCreate = ref(false)

onMounted(async () => {
  const res = await fetch('/Domains', {
    headers: {
      'Authorization': `Bearer ${userToken}`
    }
  })
  if (res.ok) {
    const data = await res.json()
    domains.value = data.domains || []
  } else {
    domains.value = []
  }
})

const filteredDomains = computed(() => {
  if (!search.value) return domains.value
  return domains.value.filter(d => d.name.toLowerCase().includes(search.value.toLowerCase()))
})

function goDomain(id) {
  router.push(`/domain/${id}`)
}
</script>

<style scoped>
.domains-container {
  max-width: 1100px;
  margin: 32px auto;
  background: #f8fcfc;
  border-radius: 16px;
  min-height: 600px;
  box-shadow: 0 0 24px #eee;
  padding-bottom: 40px;
}
.domains-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px 0 32px;
}
.logo {
  font-size: 2em;
  font-family: 'Segoe UI', cursive;
  color: #0a3566;
}
.header-actions {
  display: flex;
  gap: 18px;
  font-size: 1.3em;
  color: #0a3566;
}
.domains-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 32px 0 32px;
}
.search-input {
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 1em;
}
.create-btn {
  background: #0a3566;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 8px 18px;
  font-size: 1em;
  cursor: pointer;
  font-weight: bold;
}
.domains-list {
  margin: 32px 32px;
  background: #eaf7f6;
  border-radius: 16px;
  min-height: 300px;
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
  align-items: flex-start;
  padding: 32px;
}
.empty-tip {
  color: #888;
  font-size: 1.2em;
  margin: 32px auto;
}
.domain-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px #ddd;
  padding: 24px 32px;
  min-width: 220px;
  max-width: 240px;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: box-shadow 0.2s;
}
.domain-card:hover {
  box-shadow: 0 4px 16px #bbb;
}
.domain-avatar {
  width: 48px;
  height: 48px;
  background: #eaf7f6;
  border-radius: 50%;
  font-size: 2em;
  color: #0a3566;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}
.domain-info {
  text-align: center;
}
.domain-name {
  font-size: 1.2em;
  font-weight: bold;
  margin-bottom: 4px;
}
.domain-admin {
  font-size: 1em;
  color: #666;
  margin-bottom: 4px;
}
.domain-status {
  margin-bottom: 4px;
}
.enabled {
  color: #0a3566;
  font-weight: bold;
}
.disabled {
  color: #aaa;
}
.domain-route {
  font-size: 0.95em;
  color: #888;
}
</style>