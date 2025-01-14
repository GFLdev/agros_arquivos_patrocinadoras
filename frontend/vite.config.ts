import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import { resolve } from 'path'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

const pages = resolve(__dirname, 'src', 'pages')
const outDir = resolve(__dirname, 'dist')

// https://vite.dev/config/
export default defineConfig({
  base: './',
  plugins: [vue(), vueJsx(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    outDir: outDir,
    emptyOutDir: true,
    rollupOptions: {
      input: {
        login: resolve(pages, 'login', 'index.html'),
        admin: resolve(pages, 'admin', 'index.html'),
        user: resolve(pages, 'user', 'index.html'),
      },
    },
  },
})
