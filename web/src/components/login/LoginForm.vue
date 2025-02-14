<script setup lang="ts">
import { ref } from 'vue'
import { PhSignIn, PhCircleNotch, PhIdentificationCard, PhPassword } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import InputPassword from '@/components/generic/InputPassword.vue'
import type { LoginRequest } from '@/@types/Requests.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosError } from 'axios'
import { useAuthStore } from '@/stores/authStore.ts'
import { validatePassword, validateUsername } from '@/utils/validate.ts'

const authStore = useAuthStore()

const visible = ref<boolean>(false)
const username = ref<string | null>()
const passwd = ref<string | null>()
const loading = ref<boolean>(false)

// Função para gerenciar o envio do formulário para autenticação do usuário
async function handleSignIn(): Promise<void> {
  if (!username.value || !passwd.value) {
    // TODO: Manipular erro
    return
  }

  visible.value = false
  loading.value = true

  const body: LoginRequest = {
    username: username.value,
    password: passwd.value,
  }

  try {
    const res = await apiClient.post('/login', JSON.stringify(body))
    const token = res.data?.token

    if (token) {
      authStore.setToken(token)
      await authStore.getSession()
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
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form
    class="flex w-full flex-col items-center justify-center gap-6 rounded-none bg-white px-8 py-4 drop-shadow-2xl sm:mx-16 sm:w-fit sm:rounded-lg sm:px-16 sm:py-8"
  >
    <p class="text-agros-gray-dark text-center text-lg font-light">Por favor, preencha os campos para fazer o login.</p>
    <section class="flex w-full flex-col justify-center gap-6 px-8">
      <InputText
        label="Usuário"
        placeholder="Nome de usuário"
        v-model="username"
        :disabled="loading"
        :left-inner-icon="PhIdentificationCard"
        required
      />
      <InputPassword
        label="Senha"
        placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
        v-model="passwd"
        :disabled="loading"
        :left-inner-icon="PhPassword"
        showable
        required
      />
      <div
        v-if="!validateUsername(username) || !validatePassword(passwd)"
        class="w-full text-center text-sm font-light"
      >
        <p v-if="!validateUsername(username)">O usuário deve ter entre 4 e 16 caracteres</p>
        <p v-if="!validatePassword(passwd)">A senha deve ter pelo menos 4 caracteres</p>
      </div>
    </section>
    <button
      type="submit"
      class="focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md px-3 py-1.5 text-base text-dark shadow-sm transition duration-200 ease-in-out focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-dark enabled:hover:bg-dark enabled:hover:text-white disabled:bg-dark disabled:text-white"
      :disabled="!validateUsername(username) || !validatePassword(passwd) || loading"
      @click="handleSignIn"
    >
      <PhSignIn v-if="!loading" class="size-5" aria-hidden="true" />
      <PhCircleNotch v-else class="size-5 animate-spin" aria-hidden="true" />
      <span v-if="!loading">Entrar</span>
      <span v-else>Entrando</span>
    </button>
  </form>
</template>

<style scoped>
*:disabled {
  @apply opacity-50;
}

*:enabled {
  @apply opacity-100;
}
</style>
