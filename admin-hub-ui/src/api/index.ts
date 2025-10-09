import axios from "axios";
import { message } from 'ant-design-vue';

export const useAxios = axios.create({
  timeout: import.meta.env.VITE_API_TIMEOUT,
  baseURL: import.meta.env.VITE_API_BASE_URL, // 在使用前端代理的情况下，这里必须留空，不然会跨域
})

useAxios.interceptors.request.use((config) => {
  config.headers.set("token", "xxx")
  return config
})

useAxios.interceptors.response.use((res) => {
  if (res.status === 200){
    return res.data
  }
  return res
}, (res) => {
  message.error(res.message).then(() => Promise.reject(res))
})
