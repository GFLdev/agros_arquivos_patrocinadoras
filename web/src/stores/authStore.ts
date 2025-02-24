import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import type { UserData } from '@/@types/Entities.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosError, AxiosResponse } from 'axios'
import type { LoginRequest } from '@/@types/Requests.ts'
import type { LoginResponse } from '@/@types/Responses.ts'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()

  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<UserData | null>(null)

  const loggedIn = ref<boolean>(false)
  const loginBody = ref<LoginRequest | null>(null)

  // Função para salvar o token no estado e no localStorage
  function setToken(newToken: string | null) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  async function refreshToken(): Promise<void> {
    try {
      const res: AxiosResponse<LoginResponse> = await apiClient.post('/login', JSON.stringify(loginBody.value))
      const token = res.data?.token

      if (token) {
        setToken(token)
      } else {
        console.error('Token não recebido')
      }
    } catch (e: unknown) {
      const error = e as AxiosError
      if (error.response && error.response.status === 401) {
        // TODO: Exiba uma mensagem amigável para o usuário
        console.error('Credenciais inválidas')
      } else {
        // TODO: Tratar erros de login e exibir mensagens relevantes ao usuário
        console.error('Erro ao fazer login:', error.message || error)
      }
    }
  }

  // Função para restaurar a sessão baseada no token
  async function getSession(): Promise<void> {
    if (localStorage.getItem('token') !== token.value) {
      setToken(token.value)
    }

    try {
      if (!token.value && !loggedIn.value) {
        if (!loggedIn.value) {
          await router.push({ name: 'login', replace: true })
          return
        }
        await refreshToken()
      }

      const res: AxiosResponse<LoginResponse> = await apiClient.get('/auth/session')

      const { id, name, admin } = res.data
      user.value = {
        id: id,
        name: name,
        admin: admin,
      }

      if (admin) {
        await router.push({ name: 'admin', replace: true })
        loggedIn.value = true
      } else if (id) {
        await router.push({ name: 'user', params: { id: id }, replace: true })
        loggedIn.value = true
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
    loggedIn,
    loginBody,
    setToken,
    getSession,
    logout,
  }
})
