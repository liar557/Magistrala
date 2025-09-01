<template>
  <div class="channels-container">
    <div class="channels-header">
      <h2>Channels</h2>
      <div class="channels-actions">
        <button class="btn-create">+ Create</button>
        <button class="btn-upload">Upload</button>
      </div>
    </div>
    <div class="channels-toolbar">
      <input class="search-input" v-model="search" placeholder="Search Channel" />
    </div>
    <table class="channels-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Route</th>
          <th>Tags</th>
          <th>Status</th>
          <th>Created At</th>
          <th>Updated At</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="filteredChannels.length === 0">
          <td colspan="6" class="empty-tip">No channels found.</td>
        </tr>
        <tr
          v-for="channel in filteredChannels"
          :key="channel.id"
          class="channel-row"
          @click="goChannel(channel.id)"
          style="cursor:pointer;"
        >
          <td>{{ channel.name }}</td>
          <td>{{ channel.route ? channel.route : '-' }}</td>
          <td>-</td>
          <td>
            <span :class="channel.status === 'enabled' ? 'status-enabled' : 'status-disabled'">
              {{ channel.status === 'enabled' ? 'Enabled' : 'Disabled' }}
            </span>
          </td>
          <td>{{ formatDate(channel.created_at) }}</td>
          <td>{{ formatDate(channel.updated_at) }}</td>
        </tr>
      </tbody>
    </table>
    <div class="channels-pagination">
      Rows per page
      <select v-model="rowsPerPage">
        <option :value="10">10</option>
        <option :value="20">20</option>
      </select>
      Page {{ page }} of {{ totalPages }}
      <button @click="prevPage" :disabled="page === 1">&lt;</button>
      <button @click="nextPage" :disabled="page === totalPages">&gt;</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const domainId = route.params.id
const userToken = localStorage.getItem('token')

const channels = ref([])
const search = ref('')
const rowsPerPage = ref(10)
const page = ref(1)

function formatDate(dateStr) {
  if (!dateStr || dateStr.startsWith('0001-01-01')) return 'Not updated yet'
  const d = new Date(dateStr)
  return d.toLocaleString()
}

onMounted(async () => {
  const res = await fetch(`/Channels/${domainId}/channels`, {
    headers: {
      'Authorization': `Bearer ${userToken}`
    }
  })
  if (res.ok) {
    const data = await res.json()
    console.log('channels api 返回数据:', data)
    channels.value = data.channels || []
  } else {
    channels.value = []
  }
})

const filteredChannels = computed(() => {
  let result = channels.value
  if (search.value) {
    result = result.filter(c => c.name && c.name.toLowerCase().includes(search.value.toLowerCase()))
  }
  const start = (page.value - 1) * rowsPerPage.value
  return result.slice(start, start + rowsPerPage.value)
})

const totalPages = computed(() => Math.max(1, Math.ceil(channels.value.length / rowsPerPage.value)))

function prevPage() {
  if (page.value > 1) page.value--
}
function nextPage() {
  if (page.value < totalPages.value) page.value++
}

// 点击进入对应 channel 详情页
function goChannel(channelId) {
  router.push(`/domain/${domainId}/channels/${channelId}`)
}
</script>

<style scoped>
.channels-container {
  background: #f8fcfc;
  border-radius: 16px;
  box-shadow: 0 0 24px #eee;
  padding-bottom: 40px;
}
.channels-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px 0 32px;
}
.channels-actions button {
  margin-left: 8px;
  padding: 6px 16px;
  border-radius: 4px;
  border: none;
  background: #174e8a;
  color: #fff;
  font-weight: bold;
  cursor: pointer;
}
.channels-toolbar {
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
.channels-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 16px;
  background: #fff;
  box-shadow: 0 2px 8px #eee;
}
.channels-table th, .channels-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #eee;
  text-align: left;
}
.channel-row:hover {
  background: #eaf7f6;
}
.status-enabled {
  color: #0a3566;
  font-weight: bold;
}
.status-disabled {
  color: #e67c6b;
  font-weight: bold;
}
.action-tag {
  display: inline-block;
  background: #eaf7f6;
  color: #174e8a;
  border-radius: 4px;
  padding: 2px 6px;
  margin: 2px 2px 2px 0;
  font-size: 0.95em;
}
.channels-pagination {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}
.empty-tip {
  color: #888;
  font-size: 1.2em;
  text-align: center;
}
</style>