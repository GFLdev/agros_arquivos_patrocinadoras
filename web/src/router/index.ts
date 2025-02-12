import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/admin',
      name: 'admin',
      // meta: {
      //   requiredAuth: true,
      // },
      component: () => import('@/views/AdminView.vue'),
    },
    {
      path: '/user/:id',
      name: 'user',
      // meta: {
      //   requiredAuth: true,
      // },
      component: () => import('@/views/UserView.vue'),
    },
  ],
})

export default router
