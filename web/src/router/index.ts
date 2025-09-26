import { createRouter, createWebHistory, type Router } from 'vue-router'
import LoginView from '@/views/LoginView.vue'

// A instância é responsável por gerenciar a configuração das rotas da aplicação e o estado de navegação.
const router: Router = createRouter({
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
      component: (): Promise<typeof import('@/views/AdminView.vue')> => import('@/views/AdminView.vue'),
    },
    {
      path: '/user/:id',
      name: 'user',
      component: (): Promise<typeof import('@/views/UserView.vue')> => import('@/views/UserView.vue'),
    },
  ],
})

export default router
