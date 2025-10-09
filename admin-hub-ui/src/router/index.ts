import { createRouter, createWebHistory } from 'vue-router'

// 导入各模块路由
import authRoutes from './modules/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...authRoutes,   // 认证模块路由
    {
      path: '/',
      component: () => import('@/views/layout/index.vue'),
      children: [

      ]
    }
  ],
})

// 路由守卫可以在这里统一配置
router.beforeEach((to, from, next) => {

  // 未登录用户访问需要授权的页面时跳转登录
  if (to.meta.requiresAuth) {
    next('/login')
  } else {
    next()
  }
})

export default router
