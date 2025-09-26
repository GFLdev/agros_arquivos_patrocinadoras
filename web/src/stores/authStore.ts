import { defineStore } from 'pinia'
import { type Ref, ref } from 'vue'
import type { UserData } from '@/@types/Entities.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosResponse } from 'axios'
import type { CredentialsRequest } from '@/@types/Requests.ts'
import type { LoginResponse } from '@/@types/Responses.ts'
import router from '@/router'

// Pinia Store para gerenciar o estado de autenticação e o gerenciamento da sessão do usuário.
export const useAuthStore = defineStore('auth', () => {
  // Token JWT armazenado no localStorage.
  const token: Ref<string | null> = ref<string | null>(localStorage.getItem('token'))

  // Dados do usuário da sessão.
  const user: Ref<UserData | undefined> = ref<UserData>()

  // Parâmetros de login para atualizar token.
  const loginBody: Ref<CredentialsRequest | undefined> = ref<CredentialsRequest>()

  /**
   * Define o token de autenticação e atualiza o local storage de forma correspondente.
   *
   * @param {string | null} newToken - O novo token a ser definido. Se for null, o token será removido.
   * @return {void} Este método não retorna valor.
   */
  function setToken(newToken: string | null): void {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  /**
   * Envia uma solicitação de login para atualizar o token usando o estado atual do corpo de login.
   * Atualiza o token armazenado se um novo token for recebido, caso contrário, limpa o token e redefine o estado do
   * corpo de login em caso de falha.
   *
   * @return {Promise<void>} Uma promise que é resolvida quando a operação de atualização do token é concluída.
   */
  async function refreshToken(): Promise<void> {
    try {
      const res: AxiosResponse<LoginResponse> = await apiClient.post('/login', JSON.stringify(loginBody.value))
      const token: string = res.data?.token

      if (token) {
        setToken(token)
        return
      }

      setToken(null)
      loginBody.value = undefined
    } catch {
      setToken(null)
      loginBody.value = undefined
    }
  }

  /**
   * Gerencia a sessão do usuário verificando o token JWT, atualizando-o se necessário
   * e navegando para as páginas apropriadas com base no status do usuário (admin, usuário ou página de login).
   *
   * @return {Promise<void>} Uma promise que é resolvida quando o gerenciamento da sessão e a navegação são concluídos.
   */
  async function getSession(): Promise<void> {
    // Atualizar Cookie JWT
    if (localStorage.getItem('token') !== token.value) {
      setToken(token.value)
    }

    try {
      // Verificar JWT
      if (!token.value && loginBody.value) {
        await refreshToken()
      }

      // Tentar retomar sessão
      const res: AxiosResponse<LoginResponse> = await apiClient.get('/auth/session')
      user.value = res.data

      if (user.value.admin) {
        // Página de administrador
        await router.push({ name: 'admin', replace: true })
      } else if (user.value.id) {
        // Página de usuário
        await router.push({ name: 'user', params: { id: user.value.id }, replace: true })
      } else {
        // Erro
        await router.push({ name: 'login', replace: true })
      }
    } catch {
      setToken(null)
      await router.push({ name: 'login', replace: true })
    }
  }

  /**
   * Faz o logout do usuário da aplicação, reiniciando os tokens de autenticação, limpando os dados do usuário e
   * navegando para a página de login.
   *
   * @return {Promise<void>} Uma promise que é resolvida quando o processo de logout é concluído.
   */
  async function logout(): Promise<void> {
    setToken(null)
    loginBody.value = undefined
    user.value = undefined
    await router.push({ name: 'login' })
  }

  return {
    token,
    user,
    loginBody,
    setToken,
    getSession,
    logout,
  }
})
