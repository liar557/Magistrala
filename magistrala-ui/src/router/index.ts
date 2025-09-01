import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue' 
import Domains from '../components/Domains.vue'
import MainLayout from '../components/MainLayout.vue'
import DomainHome from '../components/DomainHome.vue'
import Channels from '../components/Channels.vue'
import ChannelLayout from '../components/ChannelLayout.vue'
import ChannelDetail from '../components/ChannelDetail.vue'
import ChannelMessages from '../components/ChannelMessages.vue'
import ChannelChart from '../components/ChannelChart.vue'

const routes = [
  { path: '/login', component: Login }, 
  { path: '/domains', component: Domains },
  {
    path: '/domain/:id',
    component: MainLayout,
    children: [
      { path: '', component: DomainHome },
      { path: 'channels', component: Channels },
      {
        path: 'channels/:channelId',
        component: ChannelLayout,
        children: [
          { path: '', component: ChannelDetail },
          { path: 'messages', component: ChannelMessages },
          { path: 'chart', component: ChannelChart }
        ]
      }
    ]
  },
  { path: '/', redirect: '/login' }
] as any

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
