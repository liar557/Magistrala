import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      '/Domains': {
        target: 'http://localhost:9003',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/Domains/, '/domains'),
      },
      '/Channels': {
        target: 'http://localhost:9005',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/Channels/, ''),
      },
      '/Users': {
        target: 'http://localhost:80', // 后端实际端口
        changeOrigin: true,
        rewrite: path => path.replace(/^\/Users/, '/users'),
      },
      '/Messages': {
        target: 'http://localhost:9011',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/Messages/, ''),
      }
    }
  }
})
