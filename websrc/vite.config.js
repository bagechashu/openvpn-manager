import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  base: './', //for embedded deployment
  publicDir: 'public',
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: "../web",
    emptyOutDir: true,
    chunkSizeWarningLimit: 500,
    rollupOptions: {
      output: {
        entryFileNames: 'js/[name].[hash].js',
        assetFileNames: '[ext]/[name].[ext]',
        chunkFileNames: 'js/chunks/[name].[hash].js',
        // chunkFileNames: (chunkInfo) => {
        //   const facadeModuleId = chunkInfo.facadeModuleId ? chunkInfo.facadeModuleId.split('/') : [];
        //   const fileName = facadeModuleId[facadeModuleId.length - 2] || '[name]';
        //   return `js/${fileName}/[name].[hash].js`;
        // },
        manualChunks(id) {
          // if (id.includes("node_modules")) {
          //   return id.toString().split("node_modules/")[1].split("/")[0].toString()
          // }
          if (id.includes("naive-ui")) {
            return "navie-ui";
          }
          if (id.includes("node_modules") && !id.includes("naive-ui")) {
            return "vendor";
          }
        },
      }
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
