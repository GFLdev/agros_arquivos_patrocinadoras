<script setup lang="ts">
import { PhIdentificationBadge, PhPassword, PhUser, PhUserPlus, PhXCircle } from '@phosphor-icons/vue'
import InputPassword from '@/components/generic/InputPassword.vue'
import InputText from '@/components/generic/InputText.vue'
import SubmitButton from '@/components/generic/SubmitButton.vue'
import CancelButton from '@/components/generic/CancelButton.vue'
import PopupWindow from '@/components/generic/PopupWindow.vue'
import type { UserRequest } from '@/@types/Requests.ts'
import apiClient from '@/services/axios.ts'
import type { AxiosError } from 'axios'
import { ref, watchEffect } from 'vue'
import { validatePassword, validateUsername } from '@/utils/validate.ts'
import router from '@/router'

// Formulário
const username = ref<string | null>(null)
const name = ref<string | null>(null)
const passwd = ref<string | null>(null)
const confirmPasswd = ref<string | null>(null)

// Validações
const loading = ref<boolean>(false)
const filled = ref<boolean>(false)
const matched = ref<boolean>(false)
const formValid = ref<boolean>(false)

// Estado de visibilidade
const showModel = defineModel<boolean>()

// Função para abrir janela de criação de usuário
async function handleCreateUser() {
  if (!username.value || !name.value || !passwd.value) {
    // TODO: Manipular erro
    return
  }
  loading.value = true

  const body: UserRequest = {
    username: username.value,
    name: name.value,
    password: passwd.value,
  }

  try {
    const res = await apiClient.post('/auth/user', JSON.stringify(body))
    const message = res.data?.message
    console.log(message)
    router.go(0)
  } catch (e: unknown) {
    const error = e as AxiosError
    if (error.response && error.response.status === 401) {
      // TODO: Exiba uma mensagem amigável para o usuário
      console.error('Credenciais inválidas')
    } else {
      // TODO: Tratar erros de login e exibir mensagens relevantes ao usuário
      console.error('Erro ao criar usuário:', error.message || error)
    }
  }
}

// Função para verificar se as senhas são iguais
function checkPasswords(): boolean {
  return passwd.value === confirmPasswd.value
}

watchEffect(() => {
  const u = !!username.value && username.value.length > 0
  const n = !!name.value && name.value.length > 0
  const p = !!passwd.value && passwd.value.length > 0
  const cp = !!confirmPasswd.value && confirmPasswd.value.length > 0
  const validLengths = validateUsername(username.value) && validatePassword(passwd.value)

  filled.value = u && n && p && cp
  matched.value = checkPasswords()
  formValid.value = filled.value && matched.value && validLengths
})
</script>

<template>
  <PopupWindow :title="`Criar novo usuário`" v-model="showModel">
    <form class="flex flex-col gap-4 space-y-4 overflow-scroll px-8 py-4">
      <div class="flex w-full flex-col gap-4">
        <!-- Campos do formulário -->
        <InputText
          placeholder="Nome de usuário"
          label="Usuário"
          v-model="username"
          :left-inner-icon="PhIdentificationBadge"
          :required="true"
        />
        <InputText
          placeholder="Nome Completo"
          label="Nome Completo"
          v-model="name"
          :left-inner-icon="PhUser"
          :required="true"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Senha"
          v-model="passwd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="true"
        />
        <InputPassword
          placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;"
          label="Confirmar Senha"
          v-model="confirmPasswd"
          :left-inner-icon="PhPassword"
          :showable="true"
          :required="true"
        />
        <!-- Avisos -->
        <div v-if="!formValid" class="w-full text-center text-sm font-light">
          <p v-if="!formValid">Preencha todos os campos</p>
          <p v-if="!validateUsername(username)">O usuário deve ter entre 4 e 16 caracteres</p>
          <p v-if="!validatePassword(passwd)">A senha deve ter pelo menos 4 caracteres</p>
          <p v-else-if="!matched">As senhas devem sem iguais</p>
        </div>
      </div>
      <div class="flex w-full flex-row items-center justify-end gap-4">
        <!-- Botões -->
        <CancelButton
          text="Cancelar"
          :on-click="() => (showModel = false)"
          :disabled="loading"
          :left-inner-icon="PhXCircle"
        />
        <SubmitButton
          text="Criar"
          loading-text="Criando"
          :on-click="handleCreateUser"
          :loading="loading"
          :disabled="loading || !formValid"
          :left-inner-icon="PhUserPlus"
        />
      </div>
    </form>
  </PopupWindow>
</template>

<style scoped></style>
