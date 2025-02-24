import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import dotenv from 'dotenv'
import * as fs from 'node:fs'

dotenv.config()

// https://vite.dev/config/
export default defineConfig({
  server: {
    https: {
      key: fs.readFileSync(process.env.SSL_KEY ?? './key.pem', 'utf8'),
      cert: fs.readFileSync(process.env.SSL_CERT ?? './cert.pem', 'utf8'),
    },
    host: '0.0.0.0',
  },
  plugins: [vue(), vueJsx(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
