<script setup lang="ts">
import { ref } from 'vue'
import { PhSignIn, PhCircleNotch, PhIdentificationCard, PhPassword } from '@phosphor-icons/vue'
import InputText from '@/components/InputText.vue'
import InputPassword from '@/components/InputPassword.vue'
import type { LoginRequest } from '@/@types/Requests.ts'

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

  const res = await fetch('https://localhost:8080/login', {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Headers': 'Content-Type',
    },
    body: JSON.stringify(body),
  })
    .then((res) => res.json())
    .finally(() => (loading.value = false))

  console.log(res)
}
</script>

<template>
  <div
    class="flex w-full flex-col items-center justify-center gap-6 rounded-none bg-white px-8 py-4 drop-shadow-2xl sm:mx-16 sm:w-fit sm:rounded-lg sm:px-16 sm:py-8"
  >
    <p class="text-agros-gray-dark text-center text-lg font-light">
      Por favor, preencha os campos para fazer o login.
    </p>
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
        placeholder="&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;&#9679;"
        v-model="passwd"
        :disabled="loading"
        :left-inner-icon="PhPassword"
        showable
        required
      />
    </section>
    <button
      type="button"
      :class="`text-dark focus-visible:outline-dark ${loading ? 'disabled:bg-dark disabled:text-white' : ''} focus-visible:outline-offset inline-flex items-center gap-x-2 rounded-md px-3 py-1.5 text-base shadow-sm transition duration-200 ease-in-out focus:-outline-offset-2 focus-visible:outline focus-visible:outline-2 enabled:hover:bg-dark enabled:hover:text-white`"
      :disabled="!username || !passwd || loading"
      @click="handleSignIn"
    >
      <PhSignIn v-if="!loading" class="size-5" aria-hidden="true" />
      <PhCircleNotch v-else class="size-5 animate-spin" aria-hidden="true" />
      <span v-if="!loading">Entrar</span>
      <span v-else>Entrar</span>
    </button>
  </div>
</template>

<style scoped>
*:disabled {
  @apply opacity-50;
}

*:enabled {
  @apply opacity-100;
}
</style>
