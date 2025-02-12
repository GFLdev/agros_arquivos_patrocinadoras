import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import type { UserData } from '@/@types/User.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosError } from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()

  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<UserData | null>(null)

  // Função para salvar o token no estado e no localStorage
  function setToken(newToken: string | null) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  // Função para restaurar a sessão baseada no token
  async function getSession(): Promise<void> {
    if (localStorage.getItem('token') !== token.value) {
      setToken(token.value)
    }

    if (!token.value) {
      await router.push({ name: 'login', replace: true })
      return
    }

    try {
      const res = await apiClient.get('/auth/session')

      const { id, name, admin } = res.data
      user.value = {
        id: id,
        name: name,
        admin: admin,
      }

      if (admin) {
        await router.push({ name: 'admin', replace: true })
      } else if (id) {
        await router.push({ name: 'user', params: { id: id }, replace: true })
      } else {
        console.error('Falha ao identificar usuário')
        await router.push({ name: 'login', replace: true })
      }
    } catch (e) {
      const error = e as AxiosError
      if (error.response && error.response.status === 401) {
        console.error(error.message)
      } else {
        console.error('Erro ao obter sessão: ', e)
      }
      setToken(null)
      await router.push({ name: 'login', replace: true })
    }
  }

  // Função para logout
  async function logout() {
    setToken(null)
    user.value = null
    await router.push({ name: 'login' })
  }

  return {
    token,
    user,
    setToken,
    getSession,
    logout,
  }
})
