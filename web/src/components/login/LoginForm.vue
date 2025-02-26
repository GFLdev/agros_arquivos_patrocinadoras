<script setup lang="ts">
import { type Ref, ref } from 'vue'
import { PhCircleNotch, PhIdentificationCard, PhPassword, PhSignIn } from '@phosphor-icons/vue'
import InputText from '@/components/generic/InputText.vue'
import InputPassword from '@/components/generic/InputPassword.vue'
import type { CredentialsRequest } from '@/@types/Requests.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosError, AxiosResponse } from 'axios'
import { useAuthStore } from '@/stores/authStore.ts'
import { validatePassword, validateUsername } from '@/utils/validate.ts'
import type { LoginResponse } from '@/@types/Responses.ts'
import { AlertType } from '@/@types/Enumerations.ts'
import PopupAlert from '@/components/generic/PopupAlert.vue'
import { Alert } from '@/utils/modals.ts'

// Pinia Store
const authStore: ReturnType<typeof useAuthStore> = useAuthStore()

// Formulário
const formUsername: Ref<string> = ref<string>('')
const formPasswd: Ref<string> = ref<string>('')

// Status de carregamento
const isLoading: Ref<boolean> = ref<boolean>(false)

// Alerta
const alert: Ref<Alert> = ref<Alert>(new Alert())

/**
 * Lida com o processo de login, validando os campos de entrada, enviando as credenciais de acesso para o servidor
 * e processando a resposta do servidor. Atualiza o estado de autenticação.
 *
 * @return {Promise<void>} Uma Promise que é resolvida sem retorno quando o processo de login é concluído.
 */
async function handleSignIn(): Promise<void> {
  isLoading.value = true
  if (!formUsername.value || !formPasswd.value) {
    alert.value.handleAlert('Campos necessários não preenchidos', AlertType.Warning)
    return
  }

  const body: CredentialsRequest = {
    username: formUsername.value,
    password: formPasswd.value,
  }

  try {
    const res: AxiosResponse<LoginResponse> = await apiClient.post('/login', JSON.stringify(body))
    const token: string = res.data?.token

    if (token) {
      alert.value.handleAlert('Login realizado com sucesso', AlertType.Success)

      // Armazenar token e parâmetros de login para atualizar token
      authStore.setToken(token)
      authStore.loginBody = body
      await authStore.getSession()
    } else {
      alert.value.handleAlert('Token não recebido', AlertType.Error)
    }
  } catch (e: unknown) {
    const error: AxiosError = e as AxiosError
    if (error.response && error.response.status === 401) {
      alert.value.handleAlert('Credenciais inválidas', AlertType.Warning)
    } else {
      alert.value.handleAlert(`Erro ao fazer login: ${error.message || error}`, AlertType.Error)
    }
  } finally {
    isLoading.value = false
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
        v-model="formUsername"
        :disabled="isLoading"
        :left-inner-icon="PhIdentificationCard"
        required
      />
      <InputPassword
        label="Senha"
        placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
        v-model="formPasswd"
        :disabled="isLoading"
        :left-inner-icon="PhPassword"
        showable
        required
      />
      <div
        v-if="!validateUsername(formUsername) || !validatePassword(formPasswd)"
        class="w-full text-center text-sm font-light"
      >
        <p v-if="!validateUsername(formUsername)">O usuário deve ter entre 4 e 16 caracteres</p>
        <p v-if="!validatePassword(formPasswd)">A senha deve ter pelo menos 4 caracteres</p>
      </div>
    </section>
    <button
      type="submit"
      class="focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md px-3 py-1.5 text-base text-dark shadow-sm transition duration-200 ease-in-out focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2 focus-visible:outline-dark enabled:hover:bg-dark enabled:hover:text-white disabled:bg-dark disabled:text-white"
      :disabled="!validateUsername(formUsername) || !validatePassword(formPasswd) || isLoading"
      @click="handleSignIn"
    >
      <PhSignIn v-if="!isLoading" class="size-5" aria-hidden="true" />
      <PhCircleNotch v-else class="size-5 animate-spin" aria-hidden="true" />
      <span v-if="!isLoading">Entrar</span>
      <span v-else>Entrando</span>
    </button>
  </form>
  <PopupAlert :text="alert.text" :type="alert.type" :duration="alert.duration" v-model="alert.show" />
</template>

<style scoped>
*:disabled {
  @apply opacity-50;
}

*:enabled {
  @apply opacity-100;
}
</style>
