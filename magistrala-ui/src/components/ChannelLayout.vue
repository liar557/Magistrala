<template>
  <div class="channel-layout">
    <div class="channel-menu-area">
      <ChannelMenu
        :name="channel.name"
        :id="channel.id"
        :menu="menuList"
        :active="activeMenu"
        @menu-click="handleMenuClick"
      />
    </div>
    <div class="channel-content-area">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ChannelMenu from './ChannelMenu.vue'

const route = useRoute()
const router = useRouter()
const domainId = route.params.id
const channelId = route.params.channelId
const userToken = localStorage.getItem('token')

const channel = ref({
  name: '',
  id: '',
})

const menuList = [
  { key: 'settings', label: 'Settings' },
  { key: 'connections', label: 'Connections' },
  { key: 'roles', label: 'Roles' },
  { key: 'members', label: 'Members' },
  { key: 'messages', label: 'Messages' },
  { key: 'chart', label: '图表展示' }, 
  { key: 'audit', label: 'Audit Logs' }
]

const activeMenu = ref('settings')

watch(
  () => route.path,
  (newPath) => {
    const matched = menuList.find(item =>
      newPath.endsWith(item.key) || (item.key === 'settings' && newPath === `/domain/${domainId}/channels/${channelId}`)
    )
    activeMenu.value = matched ? matched.key : 'settings'
  },
  { immediate: true }
)

function handleMenuClick(key) {
  activeMenu.value = key
  if (key === 'settings') {
    router.push(`/domain/${domainId}/channels/${channelId}`)
  } else {
    router.push(`/domain/${domainId}/channels/${channelId}/${key}`)
  }
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
    }
  }
})
</script>

<style scoped>
.channel-layout {
  display: flex;
  width: 100%;
  min-height: 100%;
}
.channel-menu-area {
  flex: 0 0 20%;
  max-width: 20%;
  min-width: 200px;
}
.channel-content-area {
  flex: 1;
  background: #fff;
  padding: 32px 32px 24px 32px;
  border-radius: 18px;
  box-shadow: 0 0 24px #eee;
  min-width: 0;
}
</style>