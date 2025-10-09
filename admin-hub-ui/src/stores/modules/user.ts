import { ref } from 'vue'
import { defineStore } from 'pinia'

// 定义用户状态管理
export const useUserStore = defineStore('user', () => {

  const token = ref('')

  // 设置token
  const setToken = (newToken: string) => {
    if (newToken) {
      token.value = newToken
    }
  }

  // 获取token
  const getToken = () => {
    return token.value
  }

  // 清除token
  const removeToken = () => {
    token.value = ''
  }

  return { setToken, getToken, removeToken }
},{persist:true})
