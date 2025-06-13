import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const rootPath = new URL('.', import.meta.url).pathname

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': rootPath + 'src',
      'wailsjs': rootPath + 'wailsjs',
    },
  },
})
