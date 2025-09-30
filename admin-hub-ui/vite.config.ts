import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { loadEnv } from 'vite'
import type { EnvMeta } from './env'
import { Components } from 'ant-design-vue/es/date-picker/generatePicker'

const envDir = './'

// https://vite.dev/config/
export default defineConfig((config) => {
  const env = loadEnv(config.mode, envDir) as EnvMeta
  console.log('后端地址: ', env.VITE_APP_SERVER_URL)
  return {
    plugins: [vue(),
      vueDevTools(),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    envDir: envDir,
    server: {
      host: '0.0.0.0',
      port: env.VITE_APP_PORT,
      proxy: {
        '/api': {
          target: env.VITE_APP_SERVER_URL,
          changeOrigin: true,
        },
      },
    },
  }
})
