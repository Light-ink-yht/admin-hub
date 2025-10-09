/// <reference types="vite/client" />

export interface EnvMeta extends Record<string, string>{
  NODE_ENV: string
  VITE_APP_TITLE: string
  VITE_API_BASE_URL: string
  VITE_API_TIMEOUT: number
  VITE_APP_PORT: number
  VITE_APP_SERVER_URL: string
}
