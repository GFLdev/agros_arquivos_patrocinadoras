import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      meta: {
        title: 'Login | Espaço Patrocinadora',
      },
      component: LoginView,
    },
    {
      path: '/admin',
      name: 'admin',
      meta: {
        title: 'Admin | Espaço Patrocinadora',
        requiredAuth: true,
      },
      component: () => import('../views/AdminView.vue'),
    },
    {
      path: '/user',
      name: 'user',
      meta: {
        title: 'Home | Espaço Patrocinadora',
        requiredAuth: true,
      },
      component: () => import('../views/UserView.vue'),
    },
  ],
})

export default router
