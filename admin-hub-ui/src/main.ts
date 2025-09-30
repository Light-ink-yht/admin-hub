import '@/styles/style.scss'

import { createApp } from 'vue'
import pinia from '@/stores/index.ts'

// 引入ant-design-vue
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(pinia)
app.use(router)

// 引入ant-design-vue
app.use(Antd);


app.mount('#app')
