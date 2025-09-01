<template>
  <div class="channel-detail-card">
    <table class="detail-table">
      <tbody>
        <tr>
          <td>Name</td>
          <td>{{ channel.name }}</td>
        </tr>
        <tr>
          <td>ID</td>
          <td>
            {{ channel.id }}
            <button @click="copyId">复制</button>
          </td>
        </tr>
        <tr>
          <td>Route</td>
          <td>{{ channel.route }}</td>
        </tr>
        <tr>
          <td>Tags</td>
          <td>{{ channel.tags.join(', ') }}</td>
        </tr>
        <tr>
          <td>Status</td>
          <td>{{ channel.status }}</td>
        </tr>
      </tbody>
    </table>
    <button class="delete-btn">Delete Channel</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const domainId = route.params.id
const channelId = route.params.channelId
const userToken = localStorage.getItem('token')

const channel = ref({
  name: '',
  id: '',
  route: '',
  tags: [],
  status: '',
})

function copyId() {
  navigator.clipboard.writeText(channel.value.id)
  alert('Channel ID copied!')
}

onMounted(async () => {
  const res = await fetch(`/Channels/${domainId}/channels/${channelId}`, {
    headers: {
      'Authorization': `Bearer ${userToken}`
    }
  })
  if (res.ok) {
    const data = await res.json()
    channel.value = {
      name: data.name || '',
      id: data.id || '',
      route: data.route || '',
      tags: data.tags || [],
      status: data.status || '',
    }
  }
})
</script>

<style scoped>
.channel-detail-card {
  /* width: 100%; */
  background: #fff;
  border-radius: 18px;
  box-shadow: 0 0 24px #eee;
  padding: 32px 32px 24px 32px;
  align-self: flex-start;
}
.detail-table {
  width: 100%;
  margin-bottom: 24px;
  border-collapse: collapse;
}
.detail-table td {
  padding: 8px 12px;
  border-bottom: 1px solid #eee;
}
.delete-btn {
  background: #e74c3c;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 24px;
  font-size: 1em;
  cursor: pointer;
  font-weight: bold;
}
</style>