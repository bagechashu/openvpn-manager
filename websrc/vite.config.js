import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    // http://localhost:5173/api -> http://localhost:8080/api
    proxy: {
      '/api': {
        target: 'http://192.168.56.4:8080',
        changeOrigin: true,
      },
    }
  },
})
